//go:build gc && (arm64 || amd64) && !purego && go1.18
// +build gc
// +build arm64 amd64
// +build !purego
// +build go1.18

package lea

import (
	"crypto/cipher"

	ikipher "github.com/RyuaNerin/go-krypto/internal/kipher"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var (
	_ ikipher.CBCDecAble = (*leaContextAsm)(nil)
	_ ikipher.CTRAble    = (*leaContextAsm)(nil)
	_ ikipher.GCMAble    = (*leaContextAsm)(nil)
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
