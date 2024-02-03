package kipher

import (
	"bytes"
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

// NewCBCEncrypter returns a BlockMode which encrypts in cipher block chaining
// mode, using the given Block. The length of iv must be the same as the
// Block's block size.
func NewCBCEncrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	return cipher.NewCBCEncrypter(b, iv)
}

// NewCBCDecrypter returns a BlockMode which decrypts in cipher block chaining
// mode, using the given Block. The length of iv must be the same as the
// Block's block size and must match the iv used to encrypt the data.
func NewCBCDecrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	if kb, ok := b.(kryptoBlock); ok {
		return &cbc2{
			b:         kb,
			blockSize: b.BlockSize(),
			iv:        bytes.Clone(iv),
			tmp:       make([]byte, 8*b.BlockSize()),
		}
	}
	return cipher.NewCBCDecrypter(b, iv)
}

type cbc2 struct {
	b         kryptoBlock
	blockSize int
	iv        []byte
	tmp       []byte
}

func (b *cbc2) BlockSize() int {
	return b.BlockSize()
}

func (b *cbc2) CryptBlocks(dst, src []byte) {
	if len(src)%b.blockSize != 0 {
		panic("krypto/kipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("krypto/kipher: output smaller than input")
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic("krypto/kipher: invalid buffer overlap")
	}

	var (
		bs0 = 0 * b.blockSize
		bs1 = 1 * b.blockSize
		bs4 = 4 * b.blockSize
		bs8 = 8 * b.blockSize
	)

	remainBlock := len(src) / b.blockSize

	dstIdx := remainBlock * b.blockSize
	srcIdx := remainBlock * b.blockSize

	for remainBlock >= 8 {
		remainBlock -= 8
		dstIdx -= bs8
		srcIdx -= bs8

		b.b.Decrypt8(b.tmp[:bs8], src[srcIdx:])
		if remainBlock > 0 {
			subtle.XORBytes(dst[dstIdx+bs0:dstIdx+bs8], b.tmp[bs0:bs8], src[srcIdx-bs1:])
		} else {
			// Ignore the first block, must use iv.
			subtle.XORBytes(dst[dstIdx+bs1:dstIdx+bs8], b.tmp[bs1:bs8], src[srcIdx-bs0:])
		}
	}

	for remainBlock >= 4 {
		remainBlock -= 4
		dstIdx -= bs4
		srcIdx -= bs4

		b.b.Decrypt4(b.tmp[:bs4], src[srcIdx:])
		if remainBlock > 0 {
			subtle.XORBytes(dst[dstIdx+bs0:dstIdx+bs4], b.tmp[bs0:bs4], src[srcIdx-bs1:])
		} else {
			// Ignore the first block, must use iv.
			subtle.XORBytes(dst[dstIdx+bs1:dstIdx+bs4], b.tmp[bs1:bs4], src[srcIdx-bs0:])
		}
	}

	for remainBlock >= 1 {
		remainBlock -= 1
		dstIdx -= bs1
		srcIdx -= bs1

		b.b.Decrypt(b.tmp[:bs1], src[srcIdx:])

		// Ignore the first block, must use iv.
		if remainBlock > 0 {
			subtle.XORBytes(dst[dstIdx+bs0:dstIdx+bs1], b.tmp[bs0:bs1], src[srcIdx-bs1:])
		}
	}

	subtle.XORBytes(dst[:bs1], b.tmp[:bs1], b.iv)
	copy(b.iv, src[len(src)-bs1:])
}
