//go:build arm64 && !purego

package aria

import (
	"unsafe"
)

//go:noescape
func __initEncKey_NEON(rk, key unsafe.Pointer, keyBits uint64)

//go:noescape
func __initDecKey_NEON(rk, key unsafe.Pointer, keyBits uint64)

//go:noescape
func __process_NEON(dst, src, rk unsafe.Pointer, rounds uint64)

func init() {
	newCipher = newCipherAsm
}

func (ctx *ariaContextAsm) init(key []byte) {
	keyBits := uint64(len(key) * 8)
	__initEncKey_NEON(unsafe.Pointer(&ctx.ctx.ek), unsafe.Pointer(&key[0]), keyBits)
	__initEncKey_NEON(unsafe.Pointer(&ctx.ctx.dk), unsafe.Pointer(&key[0]), keyBits)
}

func (ctx *ariaContextAsm) process(rk *[rkSize]byte, dst, src []byte) {
	__process_NEON(
		unsafe.Pointer(&dst[0]),
		unsafe.Pointer(&src[0]),
		unsafe.Pointer(rk),
		uint64(ctx.ctx.rounds),
	)
}
