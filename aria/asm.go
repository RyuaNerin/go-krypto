//go:build (arm64 || amd64) && !purego && (!gccgo || go1.18)
// +build arm64 amd64
// +build !purego
// +build !gccgo go1.18

package aria

import (
	"crypto/cipher"
	"fmt"
)

type ariaContextAsm struct {
	ariaContext
}

func newCipherAsm(key []byte) (cipher.Block, error) {
	ctx := new(ariaContextAsm)
	initRoundKeyAsm(ctx, key)
	return ctx, nil
}

func (ctx *ariaContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *ariaContextAsm) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf(msgFormatInvalidBlockSizeSrc, len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf(msgFormatInvalidBlockSizeDst, len(dst)))
	}

	processAsm(dst, src, ctx.ek[:], ctx.rounds)
}

func (ctx *ariaContextAsm) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf(msgFormatInvalidBlockSizeSrc, len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf(msgFormatInvalidBlockSizeDst, len(dst)))
	}

	processAsm(dst, src, ctx.dk[:], ctx.rounds)
}
