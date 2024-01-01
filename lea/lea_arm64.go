//go:build arm64 && !purego

package lea

func init() {
	leaEnc4 = __lea_encrypt_4block
	leaDec4 = __lea_decrypt_4block
}

//go:noescape
func __lea_encrypt_4block(ctx *leaContext, ct, pt []byte)

//go:noescape
func __lea_decrypt_4block(ctx *leaContext, pt, ct []byte)
