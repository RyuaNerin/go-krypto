//go:build gc && arm64 && !purego && go1.18
// +build gc,arm64,!purego,go1.18

package lea

//go:noescape
func __lea_encrypt_4block(ctx *leaContext, ct, pt *byte)

//go:noescape
func __lea_decrypt_4block(ctx *leaContext, pt, ct *byte)
