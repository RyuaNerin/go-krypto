package kbkdf

import (
	"crypto/cipher"
	"crypto/hmac"
	"hash"

	"github.com/RyuaNerin/go-krypto/kipher"
)

// implements of Pseudo-Random Functions
type PRF interface {
	Sum(dst []byte, K []byte, src ...[]byte) []byte
}

////////////////////////////////////////////////////////////

type prfHMAC struct {
	h func() hash.Hash
}

func (hrf *prfHMAC) Sum(dst []byte, K []byte, src ...[]byte) []byte {
	h := hmac.New(hrf.h, K)
	for _, v := range src {
		h.Write(v)
	}
	return h.Sum(dst)
}

// New HMAC-based Pseudo-Random Functions
func NewHMACPRF(h func() hash.Hash) PRF {
	return &prfHMAC{
		h: h,
	}
}

////////////////////////////////////////////////////////////

type prfCMAC struct {
	cipher func(key []byte) (cipher.Block, error)
}

func (hrf *prfCMAC) Sum(dst []byte, K []byte, src ...[]byte) []byte {
	b, err := hrf.cipher(K)
	if err != nil {
		panic(err)
	}
	h := kipher.NewCMAC(b)
	for _, v := range src {
		h.Write(v)
	}
	return h.Sum(dst)
}

// New CMAC-based Pseudo-Random Functions
func NewCMACPRF(cipher func(key []byte) (cipher.Block, error)) PRF {
	return &prfCMAC{
		cipher: cipher,
	}
}
