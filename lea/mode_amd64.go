//go:build !noasm && !noasm_krypto && !noasm_krypto_lea && amd64

package lea

import "crypto/cipher"

// for crypto/cipher
func (ctx *leaContextAsmAVX2) NewCBCEncrypter(iv []byte) cipher.BlockMode {
	return newCBCEncrypter(ctx, iv)
}

// for crypto/cipher
func (ctx *leaContextAsmAVX2) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	return newCBCDecrypter(ctx, iv)
}

// for crypto/cipher
func (ctx *leaContextAsmAVX2) NewCTR(iv []byte) cipher.Stream { return newCTR(ctx, iv) }

// for crypto/cipher
func (ctx *leaContextAsmAVX2) NewGCM(nonceSize, tagSize int) (cipher.AEAD, error) {
	return newGCM(ctx, nonceSize, tagSize)
}

// for crypto/cipher
func (ctx *leaContextAsmAVX512) NewCBCEncrypter(iv []byte) cipher.BlockMode {
	return newCBCEncrypter(ctx, iv)
}

// for crypto/cipher
func (ctx *leaContextAsmAVX512) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	return newCBCDecrypter(ctx, iv)
}

// for crypto/cipher
func (ctx *leaContextAsmAVX512) NewCTR(iv []byte) cipher.Stream { return newCTR(ctx, iv) }

// for crypto/cipher
func (ctx *leaContextAsmAVX512) NewGCM(nonceSize, tagSize int) (cipher.AEAD, error) {
	return newGCM(ctx, nonceSize, tagSize)
}
