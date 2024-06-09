//go:build amd64 && !purego && (!gccgo || go1.18)
// +build amd64
// +build !purego
// +build !gccgo go1.18

package lea

import "github.com/RyuaNerin/go-krypto/internal/golang.org/x/sys/cpu"

var hasAVX2 = cpu.X86.HasAVX2 && cpu.X86.HasAVX && cpu.X86.HasSSSE3 && cpu.X86.HasSSE2

func init() {
	leaEnc4 = toAsmFunc(__lea_encrypt_4block)
	leaDec4 = toAsmFunc(__lea_decrypt_4block)

	if hasAVX2 {
		leaEnc8 = toAsmFunc(__lea_encrypt_8block)
		leaDec8 = toAsmFunc(__lea_decrypt_8block)
	}
}
