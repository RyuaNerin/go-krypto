package kipher

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/cfb.go

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/kipher"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

// NewCFBEncrypter returns a Stream which encrypts with cipher feedback mode,
// using the given Block. The iv must be the same length as the Block's block
// size.
func NewCFBEncrypter(block cipher.Block, iv []byte) cipher.Stream {
	if len(iv) != block.BlockSize() {
		panic(msgInvalidIVLength)
	}
	return cipher.NewCFBEncrypter(block, iv)
}

func NewCFBEncrypterWithBlockSize(block cipher.Block, iv []byte, cfbBlockByteSize int) cipher.Stream {
	if len(iv) != block.BlockSize() {
		panic(msgInvalidIVLength)
	}
	if block.BlockSize() == cfbBlockByteSize {
		return cipher.NewCFBEncrypter(block, iv)
	}
	return newCFB(block, iv, false, cfbBlockByteSize)
}

// NewCFB8Decrypter returns a Stream which decrypts with cipher feedback mode,
// using the given Block. The iv must be the same length as the Block's block
// size.
func NewCFBDecrypter(block cipher.Block, iv []byte) cipher.Stream {
	if len(iv) != block.BlockSize() {
		panic(msgInvalidIVLength)
	}

	if kb, ok := block.(kipher.Block); ok {
		return newCFBFullDec(kb, iv)
	}
	return cipher.NewCFBDecrypter(block, iv)
}

func NewCFBDecrypterWithBlockSize(block cipher.Block, iv []byte, cfbBlockByteSize int) cipher.Stream {
	if len(iv) != block.BlockSize() {
		panic(msgInvalidIVLength)
	}
	if block.BlockSize() == cfbBlockByteSize {
		if kb, ok := block.(kipher.Block); ok {
			return newCFBFullDec(kb, iv)
		}
		return cipher.NewCFBDecrypter(block, iv)
	}
	return newCFB(block, iv, true, cfbBlockByteSize)
}

func newCFB(block cipher.Block, iv []byte, decrypt bool, cfbBlockByteSize int) cipher.Stream {
	blockSize := block.BlockSize()

	on := make([]byte, blockSize*2)
	x := &cfb{
		b:       block,
		out:     on[:blockSize],
		next:    on[blockSize:],
		outUsed: cfbBlockByteSize,
		bytes:   cfbBlockByteSize,
		decrypt: decrypt,
	}
	copy(x.next, iv)

	return x
}

func newCFBFullDec(block kipher.Block, iv []byte) cipher.Stream {
	blockSize := block.BlockSize()

	on := make([]byte, blockSize*(8+9))
	x := &cfbFullDec{
		b:       block,
		bs:      blockSize,
		out:     on[:8*blockSize],
		next:    on[8*blockSize:],
		outUsed: 0,
	}
	x.b.Encrypt(x.out, iv)

	return x
}

type cfb struct {
	b       cipher.Block
	next    []byte
	out     []byte
	outUsed int
	decrypt bool

	bytes int
}

func (x *cfb) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic(msgSmallDst)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgBufferOverlap)
	}
	for len(src) > 0 {
		if x.outUsed == x.bytes {
			x.b.Encrypt(x.out, x.next)
			copy(x.next, x.next[x.bytes:])
			x.outUsed = 0
		}

		if x.decrypt {
			// We can precompute a larger segment of the
			// keystream on decryption. This will allow
			// larger batches for xor, and we should be
			// able to match CTR/OFB performance.
			copy(x.next[len(x.next)-x.bytes+x.outUsed:], src)
		}
		n := subtle.XORBytes(dst, src, x.out[x.outUsed:x.bytes])
		if !x.decrypt {
			copy(x.next[len(x.next)-x.bytes+x.outUsed:], dst)
		}
		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}
}

type cfbFullDec struct {
	b  kipher.Block
	bs int

	next    []byte
	out     []byte
	outUsed int
}

func (x *cfbFullDec) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic(msgSmallDst)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgBufferOverlap)
	}

	var (
		bs1 = x.bs * 1
		bs4 = x.bs * 4
		bs5 = x.bs * 5
		bs8 = x.bs * 8
	)

	// Align
	if x.outUsed < bs1 {
		copy(x.next[bs8+x.outUsed:], src)
		n := subtle.XORBytes(dst, src, x.out[x.outUsed:bs1])
		dst = dst[n:]
		src = src[n:]
		x.outUsed += n

		if x.outUsed < bs1 {
			return
		}
	}

	// 8 Blocks
	for len(src) >= bs8 {
		copy(x.next, x.next[bs8:])
		copy(x.next[bs1:], src)
		x.b.Encrypt8(x.out, x.next)

		subtle.XORBytes(dst, src, x.out)
		dst = dst[bs8:]
		src = src[bs8:]
	}

	// 4 Blocks
	for len(src) >= bs4 {
		copy(x.next[bs4:], x.next[bs8:])
		copy(x.next[bs5:], src)
		x.b.Encrypt4(x.out[bs4:], x.next[bs4:])

		subtle.XORBytes(dst, src, x.out[bs4:])
		dst = dst[bs4:]
		src = src[bs4:]
	}

	// remains
	for len(src) > 0 {
		if x.outUsed == bs1 {
			copy(x.next, x.next[bs8:])
			x.b.Encrypt(x.out, x.next)
			x.outUsed = 0
		}

		copy(x.next[bs8+x.outUsed:], src)
		n := subtle.XORBytes(dst, src, x.out[x.outUsed:bs1])
		dst = dst[n:]
		src = src[n:]
		x.outUsed += n
	}
}
