// Package eckcdsa implements the EC-KCDSA(Korean Certificate-based Digital Signature Algorithm using Elliptic Curves) as defined in TTAK.KO-12.0015/R3
package eckcdsa

import (
	"crypto/elliptic"
	"hash"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal/eckcdsa"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

const maxAttempts = 10

// GenerateKey generates an EC-KCDSA private key using randReader.
// The caller is responsible for providing a cryptographically secure random
// source, such as crypto/rand.Reader.
func GenerateKey(c elliptic.Curve, randReader io.Reader) (*PrivateKey, error) {
	randutil.MaybeReadByte(randReader)

	// https://cs.opensource.google/go/go/+/refs/tags/go1.20.7:src/crypto/ecdsa/ecdsa_legacy.go;l=20-31
	D := new(big.Int)
	err := randutil.RandFieldElement(D, c, randReader)
	if err != nil {
		return nil, err
	}

	priv := &PrivateKey{
		D: D,
		PublicKey: PublicKey{
			Curve: c,
		},
	}
	priv.PublicKey.X, priv.PublicKey.Y = eckcdsa.GenerateXY(D, c)

	return priv, nil
}

// Sign data using K generated randomly like in crypto/ecdsa packages.
func Sign(rand io.Reader, priv *PrivateKey, h hash.Hash, data []byte) (r, s *big.Int, err error) {
	randutil.MaybeReadByte(rand)

	if priv == nil || priv.Curve == nil || priv.X == nil || priv.Y == nil || priv.D == nil || !priv.Curve.IsOnCurve(priv.X, priv.Y) {
		return nil, nil, ErrParametersNotSetUp
	}
	if priv.D.Sign() <= 0 || priv.D.Cmp(priv.Curve.Params().N) >= 0 {
		return nil, nil, ErrInvalidPrivateKey
	}

	h.Reset()
	h.Write(data)
	digest := h.Sum(nil)
	defer zeroize.Zeroize(digest)

	dBytes := priv.D.Bytes()
	defer zeroize.Zeroize(dBytes)
	csprng, err := randutil.MixedCSPRNG(rand, dBytes, digest)
	if err != nil {
		return nil, nil, err
	}
	defer zeroize.Zeroize(csprng)

	var ok bool
	k := new(big.Int)
	r, s = new(big.Int), new(big.Int)

	var buf []byte
	tmp := new(big.Int)

	defer func() {
		zeroize.Zeroize(k)
		clear(buf)
	}()

	for attempts := maxAttempts; attempts > 0; attempts-- {
		// https://cs.opensource.google/go/go/+/refs/tags/go1.20.7:src/crypto/ecdsa/ecdsa_legacy.go;l=77-113
		err = randutil.RandFieldElement(k, priv.Curve, csprng)
		if err != nil {
			return nil, nil, err
		}

		buf, ok = eckcdsa.SignUsingK(
			k, r, s,
			priv.Curve, priv.D, priv.X, priv.Y,
			h,
			data, buf,
			tmp,
		)
		zeroize.Zeroize(tmp)
		if ok {
			return r, s, nil
		}
	}

	return nil, nil, ErrInvalidPrivateKey
}

func Verify(pub *PublicKey, h hash.Hash, data []byte, r, s *big.Int) bool {
	if pub == nil || pub.Curve == nil || pub.X == nil || pub.Y == nil || !pub.Curve.IsOnCurve(pub.X, pub.Y) {
		return false
	}
	if r == nil || s == nil || r.Sign() <= 0 || s.Sign() <= 0 {
		return false
	}

	return eckcdsa.Verify(pub.Curve, pub.X, pub.Y, h, data, r, s)
}
