package eckcdsa

import (
	"crypto"
	"crypto/elliptic"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
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

var (
	_ krypto.PrivateKey = (*PrivateKey)(nil)
	_ krypto.PublicKey  = (*PublicKey)(nil)

	_ crypto.Signer     = (*PrivateKey)(nil)
	_ crypto.SignerOpts = (*SignerOpts)(nil)
)

func (pub *PublicKey) Zeroize() {
	pub.Curve = nil
	zeroize.Zeroize(pub.X)
	pub.X = nil
	zeroize.Zeroize(pub.Y)
	pub.Y = nil
}

func (priv *PrivateKey) Zeroize() {
	priv.PublicKey.Zeroize()
	zeroize.Zeroize(priv.D)
	priv.D = nil
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

// Clone returns a deep copy of pub.
func (priv PublicKey) Clone() PublicKey {
	clone := PublicKey{
		Curve: priv.Curve,
	}
	if priv.X != nil {
		clone.X = new(big.Int).Set(priv.X)
	}
	if priv.Y != nil {
		clone.Y = new(big.Int).Set(priv.Y)
	}
	return clone
}

// Clone returns a deep copy of priv.
func (priv PrivateKey) Clone() PrivateKey {
	clone := PrivateKey{
		PublicKey: priv.PublicKey.Clone(),
	}
	if priv.D != nil {
		clone.D = new(big.Int).Set(priv.D)
	}
	return clone
}

// crypto.Signer
func (priv *PrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	if opts == nil || !opts.HashFunc().Available() {
		return nil, ErrInvalidSignerOpts
	}
	return SignASN1(rand, priv, opts.HashFunc().New(), digest)
}

// SignerOpts contains options for creating and verifying EC-KCDSA signatures.
//
// It wraps a crypto.Hash and implements [crypto.SignerOpts].
type SignerOpts struct {
	Hash crypto.Hash
}

// HashFunc returns opts.Hash, allowing [SignerOpts] to satisfy [crypto.SignerOpts].
func (opts *SignerOpts) HashFunc() crypto.Hash {
	return opts.Hash
}
