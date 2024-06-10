package eckcdsa

import (
	"crypto"
	"crypto/elliptic"
	"io"
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
	return priv.PublicKey.Equal(&xx.PublicKey) && internal.BigEqual(priv.D, xx.D)
}

// Equal reports whether pub and x have the same value.
func (pub *PublicKey) Equal(x crypto.PublicKey) bool {
	xx, ok := x.(*PublicKey)
	if !ok {
		return false
	}
	return internal.BigEqual(pub.X, xx.X) && internal.BigEqual(pub.Y, xx.Y) &&
		pub.Curve == xx.Curve
}

// crypto.Signer
func (priv *PrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	return SignASN1(rand, priv, opts.HashFunc().New(), digest)
}

// SignerOpts contains options for creating and verifying EC-KCDSA signatures.
type SignerOpts struct {
	Hash crypto.Hash
}

// HashFunc returns opts.Hash so that [SignerOpts] implements [crypto.SignerOpts].
func (opts *SignerOpts) HashFunc() crypto.Hash {
	return opts.Hash
}

type (
	stdPublicKey interface {
		Equal(x crypto.PublicKey) bool
	}
	stdPrivateKey interface {
		Public() crypto.PublicKey
		Equal(x crypto.PrivateKey) bool
	}
)

var (
	_ stdPublicKey      = (*PublicKey)(nil)
	_ stdPrivateKey     = (*PrivateKey)(nil)
	_ crypto.Signer     = (*PrivateKey)(nil)
	_ crypto.SignerOpts = (*SignerOpts)(nil)
)
