package eckcdsa

import (
	"crypto"
	"crypto/elliptic"
	"crypto/subtle"
	"math/big"
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
	return priv.PublicKey.Equal(&xx.PublicKey) && bigIntEqual(priv.D, xx.D)
}

// Equal reports whether pub and x have the same value.
func (pub *PublicKey) Equal(x crypto.PublicKey) bool {
	xx, ok := x.(*PublicKey)
	if !ok {
		return false
	}
	return bigIntEqual(pub.X, xx.X) && bigIntEqual(pub.Y, xx.Y) &&
		pub.Curve == xx.Curve
}

func bigIntEqual(a, b *big.Int) bool {
	return subtle.ConstantTimeCompare(a.Bytes(), b.Bytes()) == 1
}
