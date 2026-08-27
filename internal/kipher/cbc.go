package kipher

import (
	"crypto/subtle"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/alias"
)

func NewCBCEncrypter(b krypto.Block8, iv []byte) *CBCEncrypter {
	return (*CBCEncrypter)(newCBC(b, iv))
}

func NewCBCDecrypter(b krypto.Block8, iv []byte) *CBCDecrypter {
	return (*CBCDecrypter)(newCBC(b, iv))
}

type (
	cbc struct {
		cipher    krypto.Block8
		blockSize int
		iv        []byte
		tmp       []byte
	}
	CBCEncrypter cbc
	CBCDecrypter cbc
)

var (
	_ krypto.BlockMode = (*CBCEncrypter)(nil)
	_ krypto.BlockMode = (*CBCDecrypter)(nil)
)

func newCBC(b krypto.Block8, iv []byte) *cbc {
	bs := b.BlockSize()
	if bs <= 0 {
		panic(msgInvalidBlockSize)
	}
	buf := make([]byte, 9*bs)
	copy(buf, iv)

	cbc := &cbc{
		cipher:    b,
		blockSize: bs,
		iv:        buf[:bs],
		tmp:       buf[bs:],
	}
	return cbc
}

func (ctx *cbc) Zeroize() {
	ctx.cipher = nil
	ctx.blockSize = 0
	clear(ctx.iv)
	ctx.iv = nil
	clear(ctx.tmp)
	ctx.tmp = nil
}

func (ctx *CBCEncrypter) Zeroize() { (*cbc)(ctx).Zeroize() }
func (ctx *CBCDecrypter) Zeroize() { (*cbc)(ctx).Zeroize() }

func (ctx *CBCEncrypter) BlockSize() int { return ctx.cipher.BlockSize() }
func (ctx *CBCDecrypter) BlockSize() int { return ctx.cipher.BlockSize() }

func (ctx *CBCEncrypter) CryptBlocks(dst, src []byte) {
	bs := ctx.blockSize
	if bs <= 0 {
		panic(msgInvalidBlockSize)
	}
	if len(src)%bs != 0 {
		panic(msgNotFullBlocks)
	}
	if len(dst) < len(src) {
		panic(msgSmallDst)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgBufferOverlap)
	}
	dst = dst[:len(src)]

	iv := ctx.iv

	for len(src) >= bs {
		// Write the xor to dst, then encrypt in place.
		subtle.XORBytes(dst[:bs], src[:bs], iv)
		ctx.cipher.Encrypt(dst[:bs], dst[:bs])

		// Move to the next block with this block as the next iv.
		iv = dst[:bs]
		src = src[bs:]
		dst = dst[bs:]
	}

	// Save the iv for the next CryptBlocks call.
	copy(ctx.iv, iv)
}

func (ctx *CBCDecrypter) CryptBlocks(dst, src []byte) {
	bs := ctx.blockSize
	if bs <= 0 {
		panic(msgInvalidBlockSize)
	}
	if len(src)%bs != 0 {
		panic(msgNotFullBlocks)
	}
	if len(dst) < len(src) {
		panic(msgSmallDst)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgBufferOverlap)
	}
	dst = dst[:len(src)]

	if len(src) == 0 {
		return
	}

	if &dst[0] == &src[0] {
		nextIV := ctx.tmp[bs : 2*bs]
		copy(nextIV, src[len(src)-bs:])

		for start := len(src) - bs; start > 0; start -= bs {
			ctx.cipher.Decrypt(ctx.tmp[:bs], src[start:start+bs])
			subtle.XORBytes(dst[start:start+bs], ctx.tmp[:bs], src[start-bs:start])
		}

		ctx.cipher.Decrypt(ctx.tmp[:bs], src[:bs])
		subtle.XORBytes(dst[:bs], ctx.tmp[:bs], ctx.iv)
		copy(ctx.iv, nextIV)
		return
	}

	var (
		bs0 = 0 * bs
		bs1 = 1 * bs
		bs4 = 4 * bs
		bs8 = 8 * bs
	)

	remainBlock := len(src) / bs

	dstIdx := remainBlock * bs
	srcIdx := remainBlock * bs

	for remainBlock >= 8 {
		remainBlock -= 8
		dstIdx -= bs8
		srcIdx -= bs8

		ctx.cipher.Decrypt8(ctx.tmp[:bs8], src[srcIdx:])
		if remainBlock > 0 {
			subtle.XORBytes(dst[dstIdx+bs0:dstIdx+bs8], ctx.tmp[bs0:bs8], src[srcIdx-bs1:])
		} else {
			// Ignore the first block, must use iv.
			subtle.XORBytes(dst[dstIdx+bs1:dstIdx+bs8], ctx.tmp[bs1:bs8], src[srcIdx-bs0:])
		}
	}

	for remainBlock >= 4 {
		remainBlock -= 4
		dstIdx -= bs4
		srcIdx -= bs4

		ctx.cipher.Decrypt4(ctx.tmp[:bs4], src[srcIdx:])
		if remainBlock > 0 {
			subtle.XORBytes(dst[dstIdx+bs0:dstIdx+bs4], ctx.tmp[bs0:bs4], src[srcIdx-bs1:])
		} else {
			// Ignore the first block, must use iv.
			subtle.XORBytes(dst[dstIdx+bs1:dstIdx+bs4], ctx.tmp[bs1:bs4], src[srcIdx-bs0:])
		}
	}

	for remainBlock >= 1 {
		remainBlock -= 1
		dstIdx -= bs1
		srcIdx -= bs1

		ctx.cipher.Decrypt(ctx.tmp[:bs1], src[srcIdx:])

		// Ignore the first block, must use iv.
		if remainBlock > 0 {
			subtle.XORBytes(dst[dstIdx+bs0:dstIdx+bs1], ctx.tmp[bs0:bs1], src[srcIdx-bs1:])
		}
	}

	subtle.XORBytes(dst[:bs1], ctx.tmp[:bs1], ctx.iv)
	copy(ctx.iv, src[len(src)-bs1:])
}
