//go:build arm64 && !purego
// +build arm64,!purego

package lea

import "kryptosimd/internal/ptr"

func init() {
	leaEnc4 = leaEnc4NEON
	leaDec4 = leaDec4NEON
}

func leaEnc4NEON(ctx *leaContext, dst, src []byte) {
	__lea_encrypt_4block(ctx, ptr.BytePtr(dst), ptr.BytePtr(src))
}
func leaDec4NEON(ctx *leaContext, dst, src []byte) {
	__lea_decrypt_4block(ctx, ptr.BytePtr(dst), ptr.BytePtr(src))
}

//go:noescape
func __lea_encrypt_4block(ctx *leaContext, ct, pt *byte)

//go:noescape
func __lea_decrypt_4block(ctx *leaContext, pt, ct *byte)
