package eckcdsa

import (
	"crypto"
	"crypto/elliptic"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

// PublicKey represents a EC-KCDSA public key.
type PublicKey struct {
	elliptic.Curve

	X *big.Int
	Y *big.Int
}

// PrivateKey represents a EC-KCDSA private key.
type PrivateKey struct {
	PublicKey
	D *big.Int
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
	return priv.PublicKey.Equal(&xx.PublicKey) && internal.BigIntEqual(priv.D, xx.D)
}

// Equal reports whether pub and x have the same value.
func (pub *PublicKey) Equal(x crypto.PublicKey) bool {
	xx, ok := x.(*PublicKey)
	if !ok {
		return false
	}
	return internal.BigIntEqual(pub.X, xx.X) && internal.BigIntEqual(pub.Y, xx.Y) &&
		pub.Curve == xx.Curve
}
