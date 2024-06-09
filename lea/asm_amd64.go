//go:build amd64 && !purego && (!gccgo || go1.18)
// +build amd64
// +build !purego
// +build !gccgo go1.18

package lea

//go:noescape
func __lea_encrypt_4block(ctx *leaContext, dst, src *byte)

//go:noescape
func __lea_decrypt_4block(ctx *leaContext, dst, src *byte)

//go:noescape
func __lea_encrypt_8block(ctx *leaContext, dst, src *byte)

//go:noescape
func __lea_decrypt_8block(ctx *leaContext, dst, src *byte)
