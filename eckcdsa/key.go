package eckcdsa

import (
	"crypto"
	"crypto/elliptic"
	"math/big"
)

type PublicKey struct {
	elliptic.Curve

	X *big.Int
	Y *big.Int
}

type PrivateKey struct {
	PublicKey
	D *big.Int
}

// Public returns the public key corresponding to priv.
func (priv *PrivateKey) Public() crypto.PublicKey {
	return &priv.PublicKey
}

// Equal reports whether pub and x have the same value.
func (pub *PublicKey) Equal(x crypto.PublicKey) bool {
	xx, ok := x.(*PublicKey)
	if !ok {
		return false
	}
	return pub.X.Cmp(xx.X) == 0 && pub.Y.Cmp(xx.Y) == 0 && pub.Curve == xx.Curve
}
