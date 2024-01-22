package kcdsattak

import (
	"crypto"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

type Parameters struct {
	P, Q, G *big.Int

	J     *big.Int
	Seed  []byte
	Count int
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
	return priv.PublicKey.Equal(&xx.PublicKey) && internal.BigIntEqual(priv.X, xx.X)
}

// Equal reports whether pub and y have the same value.
func (pub *PublicKey) Equal(x crypto.PublicKey) bool {
	xx, ok := x.(*PublicKey)
	if !ok {
		return false
	}
	return pub.Parameters.Equal(xx.Parameters) && internal.BigIntEqual(pub.Y, xx.Y)
}

// Equal reports whether p, q, g and sizes have the same value.
func (params *Parameters) Equal(xx Parameters) bool {
	return internal.BigIntEqual(params.P, xx.P) &&
		internal.BigIntEqual(params.Q, xx.Q) &&
		internal.BigIntEqual(params.G, xx.G)
}
