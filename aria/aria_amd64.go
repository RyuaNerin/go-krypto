//go:build amd64 && !purego
// +build amd64,!purego

package aria

import (
	"unsafe"

	"golang.org/x/sys/cpu"
)

var (
	hasSSSE3 = cpu.X86.HasSSSE3
)

//go:noescape
func __initEncKey_SSSE3(rk, key unsafe.Pointer, keyBits uint64)

//go:noescape
func __initDecKey_SSSE3(rk, key unsafe.Pointer, keyBits uint64)

//go:noescape
func __process_SSSE3(dst, src, rk unsafe.Pointer, rounds uint64)

func init() {
	if hasSSSE3 {
		newCipher = newCipherAsm
	}
}

func (ctx *ariaContextAsm) init(key []byte) {
	keyBits := uint64(len(key) * 8)
	__initEncKey_SSSE3(unsafe.Pointer(&ctx.ctx.ek), unsafe.Pointer(&key[0]), keyBits)
	__initDecKey_SSSE3(unsafe.Pointer(&ctx.ctx.dk), unsafe.Pointer(&key[0]), keyBits)
}

func (ctx *ariaContextAsm) process(rk *[rkSize]byte, dst, src []byte) {
	__process_SSSE3(
		unsafe.Pointer(&dst[0]),
		unsafe.Pointer(&dst[0]),
		unsafe.Pointer(rk),
		uint64(ctx.ctx.rounds),
	)
}
