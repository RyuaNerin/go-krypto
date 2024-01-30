package kipher

import (
	"bytes"

	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

func newCBCDecrypter2(b kryptoBlock, iv []byte) *cbc2 {
	return &cbc2{
		b:         b,
		blockSize: b.BlockSize(),
		iv:        bytes.Clone(iv),
		tmp:       make([]byte, 8*b.BlockSize()),
	}
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
		panic("krypto/lea: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("krypto/lea: output smaller than input")
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic("krypto/lea: invalid buffer overlap")
	}

	var (
		bs0 = 0 * b.blockSize
		bs1 = 1 * b.blockSize
		bs4 = 4 * b.blockSize
		bs8 = 8 * b.blockSize
	)

	tmp := make([]byte, bs8)

	remainBlock := len(src) / b.blockSize

	dstIdx := remainBlock * b.blockSize
	srcIdx := remainBlock * b.blockSize

	for remainBlock >= 8 {
		remainBlock -= 8
		dstIdx -= bs8
		srcIdx -= bs8

		b.b.Decrypt8(tmp[:bs8], src[srcIdx:])
		if remainBlock > 0 {
			subtle.XORBytes(dst[dstIdx+bs0:dstIdx+bs8], tmp[bs0:bs8], src[srcIdx-bs1:])
		} else {
			// Ignore the first block, must use iv.
			subtle.XORBytes(dst[dstIdx+bs1:dstIdx+bs8], tmp[bs1:bs8], src[srcIdx-bs0:])
		}
	}

	for remainBlock >= 4 {
		remainBlock -= 4
		dstIdx -= bs4
		srcIdx -= bs4

		b.b.Decrypt4(tmp[:bs4], src[srcIdx:])
		if remainBlock > 0 {
			subtle.XORBytes(dst[dstIdx+bs0:dstIdx+bs4], tmp[bs0:bs4], src[srcIdx-bs1:])
		} else {
			// Ignore the first block, must use iv.
			subtle.XORBytes(dst[dstIdx+bs1:dstIdx+bs4], tmp[bs1:bs4], src[srcIdx-bs0:])
		}
	}

	for remainBlock >= 1 {
		remainBlock -= 1
		dstIdx -= bs1
		srcIdx -= bs1

		b.b.Decrypt(tmp[:bs1], src[srcIdx:])

		// Ignore the first block, must use iv.
		if remainBlock > 0 {
			subtle.XORBytes(dst[dstIdx+bs0:dstIdx+bs1], tmp[bs0:bs1], src[srcIdx-bs1:])
		}
	}

	subtle.XORBytes(dst[:bs1], tmp[:bs1], b.iv)
	copy(b.iv, src[len(src)-bs1:])
}

func (b *cbc2) IV() []byte { return bytes.Clone(b.iv) }
