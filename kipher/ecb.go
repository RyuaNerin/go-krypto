package kipher

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/kipher"
)

// NewECBEncryptor returns a BlockMode which encrypts in electronic codebook
// mode, using the given Block.
//
// Zeroize() on the returned value also zeroizes the cipher b you passed in (in
// place), so do not reuse b afterward. Further use of the mode panics with
// krypto.ErrStateCleared.
func NewECBEncryptor(b cipher.Block) cipher.BlockMode {
	return wrapBlockMode(b, kipher.NewECBEncryptor(internal.ToBlock8(b)))
}

// NewECBDecryptor returns a BlockMode which decrypts in electronic codebook
// mode, using the given Block.
//
// Zeroize() on the returned value also zeroizes the cipher b you passed in (in
// place), so do not reuse b afterward. Further use of the mode panics with
// krypto.ErrStateCleared.
func NewECBDecryptor(b cipher.Block) cipher.BlockMode {
	return wrapBlockMode(b, kipher.NewECBDecryptor(internal.ToBlock8(b)))
}
