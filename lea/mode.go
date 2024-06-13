//go:build (amd64 || arm64) && !purego && (!gccgo || go1.18)
// +build amd64 arm64
// +build !purego
// +build !gccgo go1.18

package lea

import (
	"crypto/cipher"

	ikipher "github.com/RyuaNerin/go-krypto/internal/kipher"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var (
	_ ikipher.CBCDecAble = (*leaContext)(nil)
	_ ikipher.CTRAble    = (*leaContext)(nil)
	_ ikipher.GCMAble    = (*leaContext)(nil)
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
