//go:build arm64 && gc && !purego

package lea

import (
	"unsafe"
)

func init() {
	leaEnc4 = leaEnc4NEON
	leaDec4 = leaDec4NEON

	leaEnc8 = leaEnc8NEON
	leaDec8 = leaDec8NEON
}

func leaEnc4NEON(ctx *leaContext, dst, src []byte) {
	lea_encrypt_4block(
		unsafe.Pointer(&dst[0]),
		unsafe.Pointer(&src[0]),
		unsafe.Pointer(&ctx.rk[0]),
		uint64(ctx.round),
	)
}
func leaDec4NEON(ctx *leaContext, dst, src []byte) {
	lea_decrypt_4block(
		unsafe.Pointer(&dst[0]),
		unsafe.Pointer(&src[0]),
		unsafe.Pointer(&ctx.rk[0]),
		uint64(ctx.round),
	)
}

func leaEnc8NEON(ctx *leaContext, dst, src []byte) {
	leaEnc4NEON(ctx, dst[0x00:], src[0x00:])
	leaEnc4NEON(ctx, dst[0x40:], src[0x40:])
}
func leaDec8NEON(ctx *leaContext, dst, src []byte) {
	leaDec4NEON(ctx, dst[0x00:], src[0x00:])
	leaDec4NEON(ctx, dst[0x40:], src[0x40:])
}
