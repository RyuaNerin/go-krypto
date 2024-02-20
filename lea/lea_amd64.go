//go:build (amd64 || amd64p32) && !purego
// +build amd64 amd64p32
// +build !purego

package lea

func init() {
	leaEnc4 = toAsmFunc(__lea_encrypt_4block)
	leaDec4 = toAsmFunc(__lea_decrypt_4block)

	if hasAVX2 {
		leaEnc8 = toAsmFunc(__lea_encrypt_8block)
		leaDec8 = toAsmFunc(__lea_decrypt_8block)
	}
}
