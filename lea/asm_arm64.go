//go:build arm64 && !purego
// +build arm64,!purego

package lea

//go:noescape
func __lea_encrypt_4block(ctx *leaContext, ct, pt *byte)

//go:noescape
func __lea_decrypt_4block(ctx *leaContext, pt, ct *byte)
