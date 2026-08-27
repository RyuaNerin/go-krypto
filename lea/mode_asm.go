//go:build !noasm && !noasm_krypto && !noasm_krypto_lea && (amd64 || arm64)

package lea

import "crypto/cipher"

// for crypto/cipher
func (ctx *leaContextAsm) NewCBCEncrypter(iv []byte) cipher.BlockMode {
	return newCBCEncrypter(ctx, iv)
}

// for crypto/cipher
func (ctx *leaContextAsm) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	return newCBCDecrypter(ctx, iv)
}

// for crypto/cipher
func (ctx *leaContextAsm) NewCTR(iv []byte) cipher.Stream { return newCTR(ctx, iv) }

// for crypto/cipher
func (ctx *leaContextAsm) NewGCM(nonceSize, tagSize int) (cipher.AEAD, error) {
	return newGCM(ctx, nonceSize, tagSize)
}
