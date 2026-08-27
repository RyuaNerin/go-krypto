package kcdsa

import (
	"crypto"
	"crypto/subtle"
	"io"
	"math/big"
	"slices"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	kcdsa1 "github.com/RyuaNerin/go-krypto/internal/kcdsa"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

// GenerationParameters holds the extra values required to regenerate a KCDSA key
// pair.
type GenerationParameters struct {
	J     *big.Int
	Seed  []byte
	Count int
}

// Parameters holds the domain parameters of a KCDSA key pair.
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

var (
	_ krypto.PrivateKey = (*PrivateKey)(nil)
	_ krypto.PublicKey  = (*PublicKey)(nil)
)

func (params *GenerationParameters) Zeroize() {
	zeroize.Zeroize(params.J)
	params.J = nil
	zeroize.Zeroize(params.Seed)
	params.Seed = nil
	params.Count = 0
}

func (params *Parameters) Zeroize() {
	params.GenParameters.Zeroize()
	zeroize.Zeroize(params.P)
	params.P = nil
	zeroize.Zeroize(params.Q)
	params.Q = nil
	zeroize.Zeroize(params.G)
	params.G = nil
}

func (pub *PublicKey) Zeroize() {
	pub.Parameters.Zeroize()
	zeroize.Zeroize(pub.Y)
	pub.Y = nil
}

func (priv *PrivateKey) Zeroize() {
	priv.PublicKey.Zeroize()
	zeroize.Zeroize(priv.X)
	priv.X = nil
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

// Clone returns a deep copy of params.
func (params GenerationParameters) Clone() GenerationParameters {
	clone := GenerationParameters{
		Seed:  slices.Clone(params.Seed),
		Count: params.Count,
	}
	if params.J != nil {
		clone.J = new(big.Int).Set(params.J)
	}
	return clone
}

// Clone returns a deep copy of params.
func (params Parameters) Clone() Parameters {
	clone := Parameters{
		GenParameters: params.GenParameters.Clone(),
	}
	if params.P != nil {
		clone.P = new(big.Int).Set(params.P)
	}
	if params.Q != nil {
		clone.Q = new(big.Int).Set(params.Q)
	}
	if params.G != nil {
		clone.G = new(big.Int).Set(params.G)
	}
	return clone
}

// Clone returns a deep copy of pub.
func (pub PublicKey) Clone() PublicKey {
	clone := PublicKey{
		Parameters: pub.Parameters.Clone(),
	}
	if pub.Y != nil {
		clone.Y = new(big.Int).Set(pub.Y)
	}
	return clone
}

// Clone returns a deep copy of priv.
func (priv PrivateKey) Clone() PrivateKey {
	clone := PrivateKey{
		PublicKey: priv.PublicKey.Clone(),
	}
	if priv.X != nil {
		clone.X = new(big.Int).Set(priv.X)
	}
	return clone
}

// IsValid reports whether params describes a valid generation state.
func (params *GenerationParameters) IsValid() bool {
	return params.Count > 0 &&
		params.Count <= 1<<24 &&
		len(params.Seed) > 0 &&
		params.J != nil &&
		params.J.Sign() > 0
}

// Validate checks the validity of the KCDSA parameters.
func (params *Parameters) Validate() error {
	if params == nil || !kcdsa1.IsValidParameters(params.P, params.Q, params.G) {
		return ErrInvalidParameters
	}
	return nil
}

// Validate Y belongs to its subgroup of order Q.
func (pub *PublicKey) Validate() error {
	if pub == nil {
		return ErrInvalidPublicKey
	}
	if err := pub.Parameters.Validate(); err != nil {
		return err
	}
	if !kcdsa1.IsSubgroupElement(pub.Y, pub.P, pub.Q) {
		return ErrInvalidPublicKey
	}
	return nil
}

// Equal reports whether p, q, g and sizes have the same value.
func (params *GenerationParameters) Equal(xx GenerationParameters) bool {
	return internal.BigEqual(params.J, xx.J) &&
		subtle.ConstantTimeEq(int32(params.Count), int32(xx.Count)) == 1 &&
		subtle.ConstantTimeCompare(params.Seed, xx.Seed) == 1
}

// crypto.Signer
//
// KCDSA hashes the message itself (together with the public key), so digest is
// the message, not its hash. opts must be a *SignerOpts.
func (priv *PrivateKey) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	kopts, ok := opts.(*SignerOpts)
	if !ok || kopts == nil {
		return nil, ErrInvalidSignerOpts
	}

	return SignASN1(rand, priv, kopts.Sizes, digest)
}

// SignerOpts contains options for creating and verifying KCDSA signatures.
type SignerOpts struct {
	Sizes ParameterSizes
}

// HashFunc returns the hash function of opts.Sizes, so that [SignerOpts]
// implements [crypto.SignerOpts]. It returns 0 for A1024B160HAS160, as HAS-160
// is not a crypto.Hash.
func (opts *SignerOpts) HashFunc() crypto.Hash {
	switch opts.Sizes {
	case A2048B224SHA224:
		return crypto.SHA224
	case A2048B224SHA256, A2048B256SHA256, A3072B256SHA256:
		return crypto.SHA256
	default:
		return 0
	}
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
