package kipher

// Based on https://github.com/golang/go/blob/go1.21.6/src/krypto/kipher/cbc.go

import (
	"bytes"
	"crypto/cipher"
	"crypto/subtle"

	"github.com/RyuaNerin/go-krypto/internal/alias"
)

type cbc struct {
	b         cipher.Block
	blockSize int
	iv        []byte
	tmp       []byte
}

func newCBC(b cipher.Block, iv []byte) *cbc {
	return &cbc{
		b:         b,
		blockSize: b.BlockSize(),
		iv:        bytes.Clone(iv),
		tmp:       make([]byte, b.BlockSize()),
	}
}

type cbcEncrypter cbc

// NewCBCEncrypter returns a BlockMode which encrypts in cipher block chaining
// mode, using the given Block. The length of iv must be the same as the
// Block's block size.
func NewCBCEncrypter(b cipher.Block, iv []byte) BlockMode {
	if len(iv) != b.BlockSize() {
		panic("cipher.NewCBCEncrypter: IV length must equal block size")
	}
	if kb, ok := b.(kryptoBlock); ok {
		return newCBCDecrypter2(kb, iv)
	}
	return (*cbcEncrypter)(newCBC(b, iv))
}

// newCBCGenericEncrypter returns a BlockMode which encrypts in cipher block chaining
// mode, using the given Block. The length of iv must be the same as the
// Block's block size. This always returns the generic non-asm encrypter for use
// in fuzz testing.
func newCBCGenericEncrypter(b cipher.Block, iv []byte) BlockMode {
	if len(iv) != b.BlockSize() {
		panic("cipher.NewCBCEncrypter: IV length must equal block size")
	}
	return (*cbcEncrypter)(newCBC(b, iv))
}

func (x *cbcEncrypter) BlockSize() int { return x.blockSize }

func (x *cbcEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("krypto/kipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("krypto/kipher: output smaller than input")
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic("krypto/kipher: invalid buffer overlap")
	}

	iv := x.iv

	for len(src) > 0 {
		// Write the xor to dst, then encrypt in place.
		subtle.XORBytes(dst[:x.blockSize], src[:x.blockSize], iv)
		x.b.Encrypt(dst[:x.blockSize], dst[:x.blockSize])

		// Move to the next block with this block as the next iv.
		iv = dst[:x.blockSize]
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}

	// Save the iv for the next CryptBlocks call.
	copy(x.iv, iv)
}

func (x *cbcEncrypter) SetIV(iv []byte) {
	if len(iv) != len(x.iv) {
		panic("cipher: incorrect length IV")
	}
	copy(x.iv, iv)
}

type cbcDecrypter cbc

// NewCBCDecrypter returns a BlockMode which decrypts in cipher block chaining
// mode, using the given Block. The length of iv must be the same as the
// Block's block size and must match the iv used to encrypt the data.
func NewCBCDecrypter(b cipher.Block, iv []byte) BlockMode {
	if len(iv) != b.BlockSize() {
		panic("cipher.NewCBCDecrypter: IV length must equal block size")
	}
	return (*cbcDecrypter)(newCBC(b, iv))
}

// newCBCGenericDecrypter returns a BlockMode which encrypts in cipher block chaining
// mode, using the given Block. The length of iv must be the same as the
// Block's block size. This always returns the generic non-asm decrypter for use in
// fuzz testing.
func newCBCGenericDecrypter(b cipher.Block, iv []byte) BlockMode {
	if len(iv) != b.BlockSize() {
		panic("cipher.NewCBCDecrypter: IV length must equal block size")
	}
	return (*cbcDecrypter)(newCBC(b, iv))
}

func (x *cbcDecrypter) BlockSize() int { return x.blockSize }

func (x *cbcDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("krypto/kipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("krypto/kipher: output smaller than input")
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic("krypto/kipher: invalid buffer overlap")
	}
	if len(src) == 0 {
		return
	}

	// For each block, we need to xor the decrypted data with the previous block's ciphertext (the iv).
	// To avoid making a copy each time, we loop over the blocks BACKWARDS.
	end := len(src)
	start := end - x.blockSize
	prev := start - x.blockSize

	// Copy the last block of ciphertext in preparation as the new iv.
	copy(x.tmp, src[start:end])

	// Loop over all but the first block.
	for start > 0 {
		x.b.Decrypt(dst[start:end], src[start:end])
		subtle.XORBytes(dst[start:end], dst[start:end], src[prev:start])

		end = start
		start = prev
		prev -= x.blockSize
	}

	// The first block is special because it uses the saved iv.
	x.b.Decrypt(dst[start:end], src[start:end])
	subtle.XORBytes(dst[start:end], dst[start:end], x.iv)

	// Set the new iv to the first block we copied earlier.
	x.iv, x.tmp = x.tmp, x.iv
}

func (x *cbcDecrypter) SetIV(iv []byte) {
	if len(iv) != len(x.iv) {
		panic("cipher: incorrect length IV")
	}
	copy(x.iv, iv)
}

func (x *cbcEncrypter) IV() []byte { return bytes.Clone(x.iv) }
func (x *cbcDecrypter) IV() []byte { return bytes.Clone(x.iv) }
