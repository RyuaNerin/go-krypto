//go:build (amd64 || arm64) && !purego
// +build amd64 arm64
// +build !purego

package lea

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var (
	_ internal.CBCDecAble = (*leaContext)(nil)
	_ internal.CTRAble    = (*leaContext)(nil)
	_ internal.GCMAble    = (*leaContext)(nil)
)

// for crypto/cipher
func (ctx *leaContext) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	return kipher.NewCBCDecrypter(ctx, iv)
}

// for crypto/cipher
func (ctx *leaContext) NewCTR(iv []byte) cipher.Stream {
	return kipher.NewCTR(ctx, iv)
}

// for crypto/cipher
func (ctx *leaContext) NewGCM(nonceSize, tagSize int) (cipher.AEAD, error) {
	return kipher.NewGCM(ctx, nonceSize, tagSize)
}
