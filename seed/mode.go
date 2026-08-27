package seed

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var (
	_ krypto.CBCEncAble = (*seedContext)(nil)
	_ krypto.CBCDecAble = (*seedContext)(nil)
	_ krypto.CTRAble    = (*seedContext)(nil)
	_ krypto.GCMAble    = (*seedContext)(nil)
)

// for crypto/cipher
func (ctx *seedContext) NewCBCEncrypter(iv []byte) cipher.BlockMode {
	return kipher.NewCBCEncrypter(internal.ToBlock8(ctx), iv)
}

// for crypto/cipher
func (ctx *seedContext) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	return kipher.NewCBCDecrypter(internal.ToBlock8(ctx), iv)
}

// for crypto/cipher
func (ctx *seedContext) NewCTR(iv []byte) cipher.Stream {
	return kipher.NewCTR(internal.ToBlock8(ctx), iv)
}

// for crypto/cipher
func (ctx *seedContext) NewGCM(nonceSize, tagSize int) (cipher.AEAD, error) {
	return kipher.NewGCMWithSize(internal.ToBlock8(ctx), nonceSize, tagSize)
}
