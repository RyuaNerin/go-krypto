//go:build arm64 && !purego && (!gccgo || go1.18)
// +build arm64
// +build !purego
// +build !gccgo go1.18

package lea

//go:noescape
func __lea_encrypt_4block(ctx *leaContext, ct, pt *byte)

//go:noescape
func __lea_decrypt_4block(ctx *leaContext, pt, ct *byte)
