//go:build (amd64 || amd64p32) && !purego
// +build amd64 amd64p32
// +build !purego

package lea

import "github.com/RyuaNerin/go-krypto/internal/golang.org/x/sys/cpu"

var (
	hasAVX2 = cpu.X86.HasAVX2 && cpu.X86.HasAVX && cpu.X86.HasSSSE3 && cpu.X86.HasSSE2
)

//go:noescape
func __lea_encrypt_4block(ctx *leaContext, dst, src *byte)

//go:noescape
func __lea_decrypt_4block(ctx *leaContext, dst, src *byte)

//go:noescape
func __lea_encrypt_8block(ctx *leaContext, dst, src *byte)

//go:noescape
func __lea_decrypt_8block(ctx *leaContext, dst, src *byte)
