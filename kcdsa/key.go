package kcdsa

import (
	"crypto"
	"crypto/subtle"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

type GenerationParameters struct {
	J     *big.Int
	Seed  []byte
	Count int
}

type Parameters struct {
	P, Q, G *big.Int

	GenParameters GenerationParameters
}

// PublicKey represents a KCDSA public key.
type PublicKey struct {
	Parameters
	Y *big.Int
}

// PrivateKey represents a KCDSA private key.
type PrivateKey struct {
	PublicKey
	X *big.Int
}

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() crypto.PublicKey {
	return &priv.PublicKey
}

// Equal reports whether priv and x have the same value.
func (priv *PrivateKey) Equal(x crypto.PrivateKey) bool {
	xx, ok := x.(*PrivateKey)
	if !ok {
		return false
	}
	return priv.PublicKey.Equal(&xx.PublicKey) && internal.BigEqual(priv.X, xx.X)
}

// Equal reports whether pub and y have the same value.
func (pub *PublicKey) Equal(x crypto.PublicKey) bool {
	xx, ok := x.(*PublicKey)
	if !ok {
		return false
	}
	return pub.Parameters.Equal(xx.Parameters) && internal.BigEqual(pub.Y, xx.Y)
}

// Equal reports whether p, q, g and sizes have the same value.
func (params Parameters) Equal(xx Parameters) bool {
	return internal.BigEqual(params.P, xx.P) &&
		internal.BigEqual(params.Q, xx.Q) &&
		internal.BigEqual(params.G, xx.G)
}

func (params *GenerationParameters) IsValid() bool {
	return params.Count > 0 &&
		len(params.Seed) > 0 &&
		params.J == nil &&
		params.J.Sign() > 0
}

// Equal reports whether p, q, g and sizes have the same value.
func (params *GenerationParameters) Equal(xx GenerationParameters) bool {
	return internal.BigEqual(params.J, xx.J) &&
		subtle.ConstantTimeEq(int32(params.Count), int32(xx.Count)) == 1 &&
		subtle.ConstantTimeCompare(params.Seed, xx.Seed) == 1
}
