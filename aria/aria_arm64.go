//go:build arm64 && !purego
// +build arm64,!purego

package aria

import (
	"crypto/cipher"

	"kryptosimd/internal/ptr"
)

//go:noescape
func __encKeySetup_NEON(rk *byte, mk *byte, keyBytes uint64)

//go:noescape
func __decKeySetup_NEON(rk *byte, rounds uint64)

//go:noescape
func __process_NEON(dst, src, rk *byte, round uint64)

func newCipher(key []byte) (cipher.Block, error) {
	return newCipherAsm(key)
}

func (ctx *ariaContextAsm) initRoundKey(key []byte) {
	__encKeySetup_NEON(ptr.BytePtr(ctx.ctx.ek[:]), ptr.BytePtr(key), uint64(len(key)))

	ctx.ctx.dk = ctx.ctx.ek
	__decKeySetup_NEON(ptr.BytePtr(ctx.ctx.dk[:]), uint64(ctx.ctx.rounds))
}

func (ctx *ariaContextAsm) process(dst, src, rk []byte) {
	__process_NEON(ptr.BytePtr(dst), ptr.BytePtr(src), ptr.BytePtr(rk), uint64(ctx.ctx.rounds))
}
