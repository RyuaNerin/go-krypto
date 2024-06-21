//go:build (arm64 || amd64) && !purego && (!gccgo || go1.18)
// +build arm64 amd64
// +build !purego
// +build !gccgo go1.18

package lea

import (
	"crypto/cipher"
	"fmt"

	"github.com/RyuaNerin/go-krypto/internal/memory"
)

type funcBlock func(ctx *leaContext, dst, src []byte)

func newCipherAsm(key []byte) (cipher.Block, error) {
	ctx := new(leaContextAsm)
	setKeyGo(&ctx.leaContext, key)

	return ctx, nil
}

func toAsmFunc(f func(ctx *leaContext, dst, src *byte)) func(ctx *leaContext, dst, src []byte) {
	return func(ctx *leaContext, dst, src []byte) {
		f(ctx, memory.P8(dst), memory.P8(src))
	}
}

var (
	leaEnc8 funcBlock = func(ctx *leaContext, dst, src []byte) {
		__lea_encrypt_4block(ctx, memory.P8(dst[:BlockSize*4]), memory.P8(src[:BlockSize*4]))
		__lea_encrypt_4block(ctx, memory.P8(dst[BlockSize*4:]), memory.P8(src[BlockSize*4:]))
	}
	leaDec8 funcBlock = func(ctx *leaContext, dst, src []byte) {
		__lea_decrypt_4block(ctx, memory.P8(dst[:BlockSize*4]), memory.P8(src[:BlockSize*4]))
		__lea_decrypt_4block(ctx, memory.P8(dst[BlockSize*4:]), memory.P8(src[BlockSize*4:]))
	}
)

type leaContextAsm struct {
	leaContext
}

func (ctx *leaContextAsm) Encrypt4(dst, src []byte) {
	if len(src) < BlockSize*4 {
		panic(fmt.Sprintf(msgInvalidBlockSizeSrcFormat, len(src)))
	}
	if len(dst) < BlockSize*4 {
		panic(fmt.Sprintf(msgInvalidBlockSizeDstFormat, len(dst)))
	}

	__lea_encrypt_4block(&ctx.leaContext, memory.P8(dst), memory.P8(src))
}

func (ctx *leaContextAsm) Decrypt4(dst, src []byte) {
	if len(src) < BlockSize*4 {
		panic(fmt.Sprintf(msgInvalidBlockSizeSrcFormat, len(src)))
	}
	if len(dst) < BlockSize*4 {
		panic(fmt.Sprintf(msgInvalidBlockSizeDstFormat, len(dst)))
	}

	__lea_decrypt_4block(&ctx.leaContext, memory.P8(dst), memory.P8(src))
}

func (ctx *leaContextAsm) Encrypt8(dst, src []byte) {
	if len(src) < BlockSize*8 {
		panic(fmt.Sprintf(msgInvalidBlockSizeSrcFormat, len(src)))
	}
	if len(dst) < BlockSize*8 {
		panic(fmt.Sprintf(msgInvalidBlockSizeDstFormat, len(dst)))
	}

	leaEnc8(&ctx.leaContext, dst, src)
}

func (ctx *leaContextAsm) Decrypt8(dst, src []byte) {
	if len(src) < BlockSize*8 {
		panic(fmt.Sprintf(msgInvalidBlockSizeSrcFormat, len(src)))
	}
	if len(dst) < BlockSize*8 {
		panic(fmt.Sprintf(msgInvalidBlockSizeDstFormat, len(dst)))
	}

	leaDec8(&ctx.leaContext, dst, src)
}
