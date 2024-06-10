package kcdsa

import (
	"crypto"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
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
		params.J != nil &&
		params.J.Sign() > 0
}

// Equal reports whether p, q, g and sizes have the same value.
func (params *GenerationParameters) Equal(xx GenerationParameters) bool {
	return internal.BigEqual(params.J, xx.J) &&
		subtle.ConstantTimeEq(int32(params.Count), int32(xx.Count)) == 1 &&
		subtle.ConstantTimeCompare(params.Seed, xx.Seed) == 1
}

// crypto.Signer
func (priv *PrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	if _, ok := opts.(*SignerOpts); !ok {
		panic(msgInvalidSignerOpts)
	}

	return SignASN1(rand, priv, opts.(*SignerOpts).Sizes, digest)
}

// SignerOpts contains options for creating and verifying EC-KCDSA signatures.
type SignerOpts struct {
	Sizes ParameterSizes
}

// HashFunc returns opts.Hash so that [SignerOpts] implements [crypto.SignerOpts].
func (opts *SignerOpts) HashFunc() crypto.Hash {
	return crypto.SHA256
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
