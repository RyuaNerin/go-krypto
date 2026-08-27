package hight

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var (
	_ krypto.CBCEncAble = (*hightContext)(nil)
	_ krypto.CBCDecAble = (*hightContext)(nil)
	_ krypto.CTRAble    = (*hightContext)(nil)
)

// for crypto/cipher
func (ctx *hightContext) NewCBCEncrypter(iv []byte) cipher.BlockMode {
	return kipher.NewCBCEncrypter(internal.ToBlock8(ctx), iv)
}

// for crypto/cipher
func (ctx *hightContext) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	return kipher.NewCBCDecrypter(internal.ToBlock8(ctx), iv)
}

// for crypto/cipher
func (ctx *hightContext) NewCTR(iv []byte) cipher.Stream {
	return kipher.NewCTR(internal.ToBlock8(ctx), iv)
}
