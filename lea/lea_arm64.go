//go:build arm64 && !purego
// +build arm64,!purego

package lea

func init() {
	leaEnc4 = toAsmFunc(__lea_encrypt_4block)
	leaDec4 = toAsmFunc(__lea_decrypt_4block)
}
