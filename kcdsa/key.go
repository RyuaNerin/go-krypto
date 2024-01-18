package kcdsa

import (
	"crypto"
	"crypto/subtle"
	"math/big"
)

type Parameters struct {
	P, Q, G *big.Int
	Sizes   ParameterSizes
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
	return priv.PublicKey.Equal(&xx.PublicKey) && bigIntEqual(priv.X, xx.X)
}

// Equal reports whether pub and y have the same value.
func (pub *PublicKey) Equal(x crypto.PublicKey) bool {
	xx, ok := x.(*PublicKey)
	if !ok {
		return false
	}
	return pub.Parameters.Equal(xx.Parameters) && bigIntEqual(pub.Y, xx.Y)
}

// Equal reports whether p, q, g and sizes have the same value.
func (params *Parameters) Equal(xx Parameters) bool {
	return bigIntEqual(params.P, xx.P) &&
		bigIntEqual(params.Q, xx.Q) &&
		bigIntEqual(params.G, xx.G) &&
		params.Sizes == xx.Sizes
}

func bigIntEqual(a, b *big.Int) bool {
	return subtle.ConstantTimeCompare(a.Bytes(), b.Bytes()) == 1
}
