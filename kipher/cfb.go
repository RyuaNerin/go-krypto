package kipher

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/cfb.go

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/kipher"
)

// NewCFBEncrypter returns a Stream which encrypts with cipher feedback mode,
// using the given Block. The iv must be the same length as the Block's block
// size.
//
// Zeroize() on the returned value also zeroizes the cipher b you passed in (in
// place), so do not reuse b afterward. Further use of the mode panics with
// krypto.ErrStateCleared.
func NewCFBEncrypter(b cipher.Block, iv []byte) cipher.Stream {
	if len(iv) != b.BlockSize() {
		panic(msgInvalidIVLength)
	}
	return wrapStream(b, kipher.NewCFBEncrypter(internal.ToBlock8(b), iv, b.BlockSize()))
}

func NewCFBEncrypterWithBlockSize(b cipher.Block, iv []byte, cfbBlockByteSize int) cipher.Stream {
	if len(iv) != b.BlockSize() {
		panic(msgInvalidIVLength)
	}
	if cfbBlockByteSize < 1 || cfbBlockByteSize > b.BlockSize() {
		panic(msgInvalidCFBSize)
	}
	return wrapStream(b, kipher.NewCFBEncrypter(internal.ToBlock8(b), iv, cfbBlockByteSize))
}

// NewCFBDecrypter returns a Stream which decrypts with cipher feedback mode,
// using the given Block. The iv must be the same length as the Block's block
// size.
//
// Zeroize() on the returned value also zeroizes the cipher b you passed in (in
// place), so do not reuse b afterward. Further use of the mode panics with
// krypto.ErrStateCleared.
func NewCFBDecrypter(b cipher.Block, iv []byte) cipher.Stream {
	if len(iv) != b.BlockSize() {
		panic(msgInvalidIVLength)
	}

	return wrapStream(b, kipher.NewCFBDecrypter(internal.ToBlock8(b), iv, b.BlockSize()))
}

func NewCFBDecrypterWithBlockSize(b cipher.Block, iv []byte, cfbBlockByteSize int) cipher.Stream {
	if len(iv) != b.BlockSize() {
		panic(msgInvalidIVLength)
	}
	if cfbBlockByteSize < 1 || cfbBlockByteSize > b.BlockSize() {
		panic(msgInvalidCFBSize)
	}
	return wrapStream(b, kipher.NewCFBDecrypter(internal.ToBlock8(b), iv, cfbBlockByteSize))
}
