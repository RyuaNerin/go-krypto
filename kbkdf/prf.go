package kbkdf

import (
	"crypto/cipher"
	"hash"

	"github.com/RyuaNerin/go-krypto/cmac"
	"github.com/RyuaNerin/go-krypto/hmac"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

// PRF is a pseudo-random function used to derive key material in KBKDF.
type PRF interface {
	// Sum appends the pseudo-random output of the PRF to dst. The output is
	// derived from K as the key and the elements of src as additional input.
	Sum(dst []byte, K []byte, src ...[]byte) []byte
}

////////////////////////////////////////////////////////////

type prfHMAC struct {
	h func() hash.Hash
}

func (hrf *prfHMAC) Sum(dst, key []byte, src ...[]byte) []byte {
	h := hmac.New(hrf.h, key)
	defer zeroize.Zeroize(h)

	for _, v := range src {
		h.Write(v)
	}
	return h.Sum(dst)
}

// NewHMACPRF returns a pseudo-random function based on HMAC.
//
// The returned value zeroizes the underlying hash after each call, so secret
// material is wiped during and after derivation.
func NewHMACPRF(h func() hash.Hash) PRF {
	return &prfHMAC{
		h: h,
	}
}

////////////////////////////////////////////////////////////

type prfCMAC struct {
	cipher func(key []byte) (cipher.Block, error)
}

func (hrf *prfCMAC) Sum(dst, key []byte, src ...[]byte) []byte {
	b, err := hrf.cipher(key)
	if err != nil {
		panic(err)
	}
	defer zeroize.Zeroize(b)

	h := cmac.New(b)
	defer zeroize.Zeroize(h)

	for _, v := range src {
		h.Write(v)
	}
	return h.Sum(dst)
}

// NewCMACPRF returns a pseudo-random function based on CMAC.
//
// The returned value zeroizes the derived cipher and CMAC after each call, so
// secret material is wiped during and after derivation.
func NewCMACPRF(cipher func(key []byte) (cipher.Block, error)) PRF {
	return &prfCMAC{
		cipher: cipher,
	}
}
