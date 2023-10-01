//go:build amd64 && gc && !purego

package lea

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

type cbcDecAble interface {
	NewCBCDecrypter(iv []byte) cipher.BlockMode
}

// Assert that leaContext implements the cbcDecAble interfaces.
var _ cbcDecAble = (*leaContext)(nil)

type cbcContext struct {
	ctx *leaContext
	iv  []byte
	tmp []byte
}

func (ctx *leaContext) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	cbc := &cbcContext{
		ctx: ctx,
		iv:  make([]byte, BlockSize),
		tmp: make([]byte, BlockSize*8),
	}
	copy(cbc.iv, iv)
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
	if subtle.InexactOverlap(dst[:len(src)], src) {
		panic("krypto/lea: invalid buffer overlap")
	}

	remainBlock := len(src) / BlockSize

	dstIdx := remainBlock * BlockSize
	srcIdx := remainBlock * BlockSize

	for remainBlock >= 8 {
		remainBlock -= 8
		dstIdx -= BlockSize * 8
		srcIdx -= BlockSize * 8

		dstLocal := dst[dstIdx : dstIdx+BlockSize*8]
		leaDec8(b.ctx, dstLocal, src[srcIdx:])
		if remainBlock > 0 {
			xorBytes(dst[dstIdx:], dstLocal, src[srcIdx-BlockSize:])
		} else {
			// Ignore the first block, must use iv.
			xorBytes(dst[dstIdx+BlockSize:], dstLocal[BlockSize:], src[srcIdx:])
		}
	}

	for remainBlock >= 4 {
		remainBlock -= 4
		dstIdx -= BlockSize * 4
		srcIdx -= BlockSize * 4

		dstLocal := dst[dstIdx : dstIdx+BlockSize*4]
		leaDec4(b.ctx, dstLocal, src[srcIdx:])
		if remainBlock > 0 {
			xorBytes(dst[dstIdx:], dstLocal, src[srcIdx-BlockSize:])
		} else {
			// Ignore the first block, must use iv.
			xorBytes(dst[dstIdx+BlockSize:], dstLocal[BlockSize:], src[srcIdx:])
		}
	}

	for remainBlock >= 1 {
		remainBlock -= 1
		dstIdx -= BlockSize
		srcIdx -= BlockSize

		dstLocal := dst[dstIdx : dstIdx+BlockSize]
		leaDec1(b.ctx, dstLocal, src[srcIdx:])

		if remainBlock > 0 { // Ignore the first block, must use iv.
			xorBytes(dst[dstIdx:], dstLocal, src[srcIdx-BlockSize:])
		}
	}

	xorBytes(dst, dst[:BlockSize], b.iv)
	copy(b.iv, src[len(src)-BlockSize:])
}
