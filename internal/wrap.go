package internal

import (
	"crypto/cipher"
	"hash"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

// ////////////////////////////////////////////////
// MARK: krypto.BlockMode

func WrapBlockMode(b cipher.BlockMode) krypto.BlockMode {
	if bm, ok := b.(krypto.BlockMode); ok {
		return bm
	}
	return &blockModeWrap{b}
}

type blockModeWrap struct{ cipher.BlockMode }

var _ krypto.BlockMode = (*blockModeWrap)(nil)

func (bm *blockModeWrap) Zeroize()                    { zeroize.Zeroize(bm.BlockMode) }
func (bm *blockModeWrap) BlockSize() int              { return bm.BlockMode.BlockSize() }
func (bm *blockModeWrap) CryptBlocks(dst, src []byte) { bm.BlockMode.CryptBlocks(dst, src) }

//////////////////////////////////////////////////
// MARK: krypto.Hash

func WrapHash(h hash.Hash) krypto.Hash {
	if kh, ok := h.(krypto.Hash); ok {
		return kh
	}
	return &hashWrap{h: h}
}

type hashWrap struct{ h hash.Hash }

var _ krypto.Hash = (*hashWrap)(nil)

func (h *hashWrap) Zeroize()                    { zeroize.Zeroize(h.h) }
func (h *hashWrap) Write(p []byte) (int, error) { return h.h.Write(p) }
func (h *hashWrap) Sum(b []byte) []byte         { return h.h.Sum(b) }
func (h *hashWrap) Reset()                      { h.h.Reset() }
func (h *hashWrap) Size() int                   { return h.h.Size() }
func (h *hashWrap) BlockSize() int              { return h.h.BlockSize() }

//////////////////////////////////////////////////
// MARK: krypto.Stream

func WrapStream(s cipher.Stream) krypto.Stream {
	if st, ok := s.(krypto.Stream); ok {
		return st
	}
	return &streamWrap{Stream: s}
}

type streamWrap struct{ cipher.Stream }

var _ krypto.Stream = (*streamWrap)(nil)

func (sw *streamWrap) Zeroize()                     { zeroize.Zeroize(sw.Stream) }
func (sw *streamWrap) XORKeyStream(dst, src []byte) { sw.Stream.XORKeyStream(dst, src) }

// ////////////////////////////////////////////////
// MARK: krypto.AEAD

func WrapAEAD(a cipher.AEAD) krypto.AEAD {
	if aa, ok := a.(krypto.AEAD); ok {
		return aa
	}
	return &aeadWrap{AEAD: a}
}

type aeadWrap struct{ cipher.AEAD }

var _ krypto.AEAD = (*aeadWrap)(nil)

func (aw *aeadWrap) Zeroize()       { zeroize.Zeroize(aw.AEAD) }
func (aw *aeadWrap) NonceSize() int { return aw.AEAD.NonceSize() }
func (aw *aeadWrap) Overhead() int  { return aw.AEAD.Overhead() }
func (aw *aeadWrap) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	return aw.AEAD.Seal(dst, nonce, plaintext, additionalData)
}

func (aw *aeadWrap) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	return aw.AEAD.Open(dst, nonce, ciphertext, additionalData)
}

// ////////////////////////////////////////////////
// MARK: krypto.Block

func ToBlock8Fn(fn func(key []byte) (cipher.Block, error)) func(key []byte) (krypto.Block8, error) {
	return func(key []byte) (krypto.Block8, error) {
		b, err := fn(key)
		if err != nil {
			return nil, err
		}
		return ToBlock8(b), nil
	}
}

func ToBlock8(b cipher.Block) krypto.Block8 {
	switch b := b.(type) {
	case krypto.Block8:
		return b
	case krypto.Block4:
		return &toBlock8{b, b.BlockSize()}
	default:
		return &toBlockN{b, b.BlockSize()}
	}
}

type (
	toBlockN struct {
		cipher.Block
		bs int
	}
	toBlock8 struct {
		krypto.Block4
		bs int
	}
)

var (
	_ krypto.Block8 = (*toBlockN)(nil)
	_ krypto.Block8 = (*toBlock8)(nil)
)

func (b *toBlockN) Zeroize() { zeroize.Zeroize(b.Block) }

func (b *toBlockN) Encrypt4(dst, src []byte) { b.encryptN(dst, src, 4) }
func (b *toBlockN) Decrypt4(dst, src []byte) { b.decryptN(dst, src, 4) }
func (b *toBlockN) Encrypt8(dst, src []byte) { b.encryptN(dst, src, 8) }
func (b *toBlockN) Decrypt8(dst, src []byte) { b.decryptN(dst, src, 8) }

func (b *toBlockN) encryptN(dst, src []byte, n int) {
	bs, end := b.bs, b.bs*n

	_ = dst[end-1]
	_ = src[end-1]

	for idx := 0; idx < end; idx += bs {
		b.Encrypt(dst[idx:idx+bs:idx+bs], src[idx:idx+bs:idx+bs])
	}
}

func (b *toBlockN) decryptN(dst, src []byte, n int) {
	bs, end := b.bs, b.bs*n

	_ = dst[end-1]
	_ = src[end-1]

	for idx := 0; idx < end; idx += bs {
		b.Decrypt(dst[idx:idx+bs:idx+bs], src[idx:idx+bs:idx+bs])
	}
}

func (b *toBlock8) Encrypt8(dst, src []byte) {
	n := b.bs * 4
	b.Encrypt4(dst[0*n:1*n:1*n], src[0*n:1*n:1*n])
	b.Encrypt4(dst[1*n:2*n:2*n], src[1*n:2*n:2*n])
}

func (b *toBlock8) Decrypt8(dst, src []byte) {
	n := b.bs * 4
	b.Decrypt4(dst[0*n:1*n:1*n], src[0*n:1*n:1*n])
	b.Decrypt4(dst[1*n:2*n:2*n], src[1*n:2*n:2*n])
}

// ////////////////////////////////////////////////
// MARK: mask Block

func JustBlockFn(fn func(key []byte) (cipher.Block, error)) func(key []byte) (cipher.Block, error) {
	return func(key []byte) (cipher.Block, error) {
		b, err := fn(key)
		if err != nil {
			return nil, err
		}
		return JustBlock(b), nil
	}
}

func JustBlock4Fn(fn func(key []byte) (cipher.Block, error)) func(key []byte) (cipher.Block, error) {
	return func(key []byte) (cipher.Block, error) {
		b, err := fn(key)
		if err != nil {
			return nil, err
		}
		return JustBlock4(b), nil
	}
}

// Hide 4, 8, 16 Blocks.
func JustBlock(b cipher.Block) cipher.Block {
	if b, ok := b.(krypto.Block4); ok {
		return &just1Block{b: b}
	}
	return b
}

// Hide 8, 16 Blocks.
func JustBlock4(b cipher.Block) cipher.Block {
	if b, ok := b.(krypto.Block8); ok {
		return &just4Block{b: b}
	}
	return b
}

type (
	just1Block struct{ b krypto.Block }
	just4Block struct{ b krypto.Block4 }
)

var (
	_ krypto.Block  = (*just1Block)(nil)
	_ krypto.Block4 = (*just4Block)(nil)
)

func (b *just1Block) Zeroize()                { b.b.Zeroize() }
func (b *just1Block) BlockSize() int          { return b.b.BlockSize() }
func (b *just1Block) Encrypt(dst, src []byte) { b.b.Encrypt(dst, src) }
func (b *just1Block) Decrypt(dst, src []byte) { b.b.Decrypt(dst, src) }

func (b *just4Block) Zeroize()                 { b.b.Zeroize() }
func (b *just4Block) BlockSize() int           { return b.b.BlockSize() }
func (b *just4Block) Encrypt(dst, src []byte)  { b.b.Encrypt(dst, src) }
func (b *just4Block) Decrypt(dst, src []byte)  { b.b.Decrypt(dst, src) }
func (b *just4Block) Encrypt4(dst, src []byte) { b.b.Encrypt4(dst, src) }
func (b *just4Block) Decrypt4(dst, src []byte) { b.b.Decrypt4(dst, src) }

// ////////////////////////////////////////////////
// MARK: unwrap

func UnwrapBlockMode(b cipher.BlockMode) cipher.BlockMode {
	if bm, ok := b.(*blockModeWrap); ok {
		return UnwrapBlockMode(bm.BlockMode)
	}
	return b
}

func UnwrapHash(h hash.Hash) hash.Hash {
	if hw, ok := h.(*hashWrap); ok {
		return UnwrapHash(hw.h)
	}
	return h
}

func UnwrapStream(s cipher.Stream) cipher.Stream {
	if sw, ok := s.(*streamWrap); ok {
		return UnwrapStream(sw.Stream)
	}
	return s
}

func UnwrapAEAD(a cipher.AEAD) cipher.AEAD {
	if aw, ok := a.(*aeadWrap); ok {
		return UnwrapAEAD(aw.AEAD)
	}
	return a
}

func UnwrapBlock(b cipher.Block) cipher.Block {
	switch bb := b.(type) {
	case *just1Block:
		return UnwrapBlock(bb.b)
	case *just4Block:
		return UnwrapBlock(bb.b)
	case *toBlockN:
		return UnwrapBlock(bb.Block)
	case *toBlock8:
		return UnwrapBlock(bb.Block4)
	}
	return b
}
