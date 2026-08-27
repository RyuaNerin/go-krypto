package lea

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var (
	_ krypto.CBCEncAble = (*leaContext)(nil)
	_ krypto.CBCDecAble = (*leaContext)(nil)
	_ krypto.CTRAble    = (*leaContext)(nil)
	_ krypto.GCMAble    = (*leaContext)(nil)
)

// for crypto/cipher
func (ctx *leaContext) NewCBCEncrypter(iv []byte) cipher.BlockMode { return newCBCEncrypter(ctx, iv) }

// for crypto/cipher
func (ctx *leaContext) NewCBCDecrypter(iv []byte) cipher.BlockMode { return newCBCDecrypter(ctx, iv) }

// for crypto/cipher
func (ctx *leaContext) NewCTR(iv []byte) cipher.Stream { return newCTR(ctx, iv) }

// for crypto/cipher
func (ctx *leaContext) NewGCM(nonceSize, tagSize int) (cipher.AEAD, error) {
	return newGCM(ctx, nonceSize, tagSize)
}

func newCBCEncrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	return kipher.NewCBCEncrypter(internal.ToBlock8(b), iv)
}

func newCBCDecrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	return kipher.NewCBCDecrypter(internal.ToBlock8(b), iv)
}

func newCTR(b cipher.Block, iv []byte) cipher.Stream {
	return kipher.NewCTR(internal.ToBlock8(b), iv)
}

func newGCM(b cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	return kipher.NewGCMWithSize(internal.ToBlock8(b), nonceSize, tagSize)
}
