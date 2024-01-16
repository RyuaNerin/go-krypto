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
	init(key []byte)
	process(rk *[rkSize]byte, dst, src []byte)
} = (*ariaContextAsm)(nil)

func newCipherAsm(key []byte) (cipher.Block, error) {
	ctx := new(ariaContextAsm)
	ctx.ctx.rounds = (len(key)*8 + 256) / 32

	ctx.init(key)
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

	ctx.process(&ctx.ctx.ek, dst, src)
}

func (ctx *ariaContextAsm) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/aria: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/aria: invalid block size %d (dst)", len(dst)))
	}

	ctx.process(&ctx.ctx.dk, dst, src)
}