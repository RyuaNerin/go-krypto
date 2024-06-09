//go:build arm64 && !purego && (!gccgo || go1.18)
// +build arm64
// +build !purego
// +build !gccgo go1.18

package lea

func init() {
	leaEnc4 = toAsmFunc(__lea_encrypt_4block)
	leaDec4 = toAsmFunc(__lea_decrypt_4block)
}
