//go:build amd64 && !purego
// +build amd64,!purego

package lea

import (
	"kryptosimd/internal/ptr"

	"golang.org/x/sys/cpu"
)

var (
	hasAVX2 = cpu.X86.HasAVX2 && cpu.X86.HasAVX && cpu.X86.HasSSSE3 && cpu.X86.HasSSE2
)

func init() {
	leaEnc4 = leaEnc4SSE2
	leaDec4 = leaDec4SSE2

	if hasAVX2 {
		leaEnc8 = leaEnc8AVX2
		leaDec8 = leaDec8AVX2
	}
}

func leaEnc4SSE2(ctx *leaContext, dst, src []byte) {
	__lea_encrypt_4block(ctx, ptr.BytePtr(dst), ptr.BytePtr(src))
}
func leaDec4SSE2(ctx *leaContext, dst, src []byte) {
	__lea_decrypt_4block(ctx, ptr.BytePtr(dst), ptr.BytePtr(src))
}

func leaEnc8AVX2(ctx *leaContext, dst, src []byte) {
	__lea_encrypt_8block(ctx, ptr.BytePtr(dst), ptr.BytePtr(src))
}
func leaDec8AVX2(ctx *leaContext, dst, src []byte) {
	__lea_decrypt_8block(ctx, ptr.BytePtr(dst), ptr.BytePtr(src))
}

//go:noescape
func __lea_encrypt_4block(ctx *leaContext, dst, src *byte)

//go:noescape
func __lea_decrypt_4block(ctx *leaContext, dst, src *byte)

//go:noescape
func __lea_encrypt_8block(ctx *leaContext, dst, src *byte)

//go:noescape
func __lea_decrypt_8block(ctx *leaContext, dst, src *byte)
