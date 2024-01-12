//go:build (amd64 || arm64) && !purego
// +build amd64 arm64
// +build !purego

package lea

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

type cbcDecAble interface {
	NewCBCDecrypter(iv []byte) cipher.BlockMode
}

// Assert that leaContext implements the cbcDecAble interfaces.
var _ cbcDecAble = (*leaContext)(nil)

type cbcContext struct {
	ctx *leaContext
	iv  [BlockSize]byte
}

func (ctx *leaContext) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	cbc := &cbcContext{
		ctx: ctx,
	}
	copy(cbc.iv[:], iv)
	return cbc
}

func (b *cbcContext) BlockSize() int {
	return b.BlockSize()
}

func (b *cbcContext) CryptBlocks(dst, src []byte) {
	if len(src)%BlockSize != 0 {
		panic("krypto/lea: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("krypto/lea: output smaller than input")
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic("krypto/lea: invalid buffer overlap")
	}

	const (
		bs0 = 0 * BlockSize
		bs1 = 1 * BlockSize
		bs4 = 4 * BlockSize
		bs8 = 8 * BlockSize
	)

	var tmp0 [bs8]byte
	tmp := tmp0[:]

	remainBlock := len(src) / BlockSize

	dstIdx := remainBlock * BlockSize
	srcIdx := remainBlock * BlockSize

	for remainBlock >= 8 {
		remainBlock -= 8
		dstIdx -= bs8
		srcIdx -= bs8

		leaDec8(b.ctx, tmp, src[srcIdx:])
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

		leaDec4(b.ctx, tmp, src[srcIdx:])
		if remainBlock > 0 {
			subtle.XORBytes(dst[dstIdx+bs0:dstIdx+bs4], tmp[bs0:bs4], src[srcIdx-bs1:])
		} else {
			// Ignore the first block, must use iv.
			subtle.XORBytes(dst[dstIdx+bs1:dstIdx+bs4], tmp[bs1:bs4], src[srcIdx-bs0:])
		}
	}

	for remainBlock >= 1 {
		remainBlock -= 1
		dstIdx -= BlockSize
		srcIdx -= BlockSize

		leaDec1(b.ctx, tmp, src[srcIdx:])

		// Ignore the first block, must use iv.
		if remainBlock > 0 {
			subtle.XORBytes(dst[dstIdx+bs0:dstIdx+bs1], tmp[bs0:bs1], src[srcIdx-bs1:])
		}
	}

	subtle.XORBytes(dst[:bs1], tmp[:bs1], b.iv[:])
	copy(b.iv[:], src[len(src)-BlockSize:])
}
