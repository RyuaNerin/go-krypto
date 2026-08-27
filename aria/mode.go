package aria

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var (
	_ krypto.CBCEncAble = (*ariaContext)(nil)
	_ krypto.CBCDecAble = (*ariaContext)(nil)
	_ krypto.CTRAble    = (*ariaContext)(nil)
	_ krypto.GCMAble    = (*ariaContext)(nil)
)

// for crypto/cipher
func (ctx *ariaContext) NewCBCEncrypter(iv []byte) cipher.BlockMode {
	return kipher.NewCBCEncrypter(internal.ToBlock8(ctx), iv)
}

// for crypto/cipher
func (ctx *ariaContext) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	return kipher.NewCBCDecrypter(internal.ToBlock8(ctx), iv)
}

// for crypto/cipher
func (ctx *ariaContext) NewCTR(iv []byte) cipher.Stream {
	return kipher.NewCTR(internal.ToBlock8(ctx), iv)
}

// for crypto/cipher
func (ctx *ariaContext) NewGCM(nonceSize, tagSize int) (cipher.AEAD, error) {
	return kipher.NewGCMWithSize(internal.ToBlock8(ctx), nonceSize, tagSize)
}
