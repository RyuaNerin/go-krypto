//go:build gc && (arm64 || amd64) && !purego && go1.18
// +build gc
// +build arm64 amd64
// +build !purego
// +build go1.18

package lea

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal/kipherimpl"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var (
	_ kipherimpl.CBCDecAble = (*leaContextAsm)(nil)
	_ kipherimpl.CTRAble    = (*leaContextAsm)(nil)
	_ kipherimpl.GCMAble    = (*leaContextAsm)(nil)
)

// for crypto/cipher
func (ctx *leaContextAsm) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	return kipher.NewCBCDecrypter(ctx, iv)
}

// for crypto/cipher
func (ctx *leaContextAsm) NewCTR(iv []byte) cipher.Stream {
	return kipher.NewCTR(ctx, iv)
}

// for crypto/cipher
func (ctx *leaContextAsm) NewGCM(nonceSize, tagSize int) (cipher.AEAD, error) {
	return kipher.NewGCMWithSize(ctx, nonceSize, tagSize)
}
