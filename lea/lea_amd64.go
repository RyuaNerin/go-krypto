//go:build amd64 && !purego
// +build amd64,!purego

package lea

import (
	"golang.org/x/sys/cpu"
)

var (
	hasAVX2 = cpu.X86.HasAVX2 && cpu.X86.HasAVX && cpu.X86.HasSSSE3 && cpu.X86.HasSSE2
)

func init() {
	leaEnc4 = __lea_encrypt_4block
	leaDec4 = __lea_decrypt_4block

	if hasAVX2 {
		leaEnc8 = __lea_encrypt_8block
		leaDec8 = __lea_decrypt_8block
	}
}

//go:noescape
func __lea_encrypt_4block(ctx *leaContext, dst, src []byte)

//go:noescape
func __lea_decrypt_4block(ctx *leaContext, dst, src []byte)

//go:noescape
func __lea_encrypt_8block(ctx *leaContext, dst, src []byte)

//go:noescape
func __lea_decrypt_8block(ctx *leaContext, dst, src []byte)
