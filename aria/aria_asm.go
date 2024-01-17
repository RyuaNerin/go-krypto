//go:build (amd64 || arm64) && !purego
// +build amd64 arm64
// +build !purego

package aria

import (
	"crypto/cipher"
	"fmt"
)

type ariaContextAsm struct {
	ctx ariaContext
}

var _ interface {
	initRoundKey(key []byte)
	process(rk []byte, dst, src []byte)
} = (*ariaContextAsm)(nil)

func newCipherAsm(key []byte) (cipher.Block, error) {
	ctx := new(ariaContextAsm)
	ctx.ctx.rounds = (len(key) + 32) / 4

	ctx.initRoundKey(key)
	return ctx, nil
}

func (ctx *ariaContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *ariaContextAsm) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/aria: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/aria: invalid block size %d (dst)", len(dst)))
	}

	ctx.process(dst, src, ctx.ctx.ek[:])
}

func (ctx *ariaContextAsm) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/aria: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/aria: invalid block size %d (dst)", len(dst)))
	}

	ctx.process(dst, src, ctx.ctx.dk[:])
}
