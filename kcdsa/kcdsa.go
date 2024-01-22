// Package kcdsa implements the KCDSA(Korean Certificate-based Digital Signature Algorithm) as defined in TTAK.KO-12.0001/R4
package kcdsa

import (
	"crypto/rand"
	"errors"
	"hash"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
	kcdsainternal "github.com/RyuaNerin/go-krypto/kcdsa/internal"
)

var (
	ErrInvalidPublicKey      = errors.New("krypto/kcdsa: invalid public key")
	ErrInvalidParameterSizes = errors.New("krypto/kcdsa: invalid ParameterSizes")
	ErrParametersNotSetUp    = errors.New("krypto/kcdsa: parameters not set up before generating key")
)

type ParameterSizes int

const (
	L2048N224SHA224 ParameterSizes = iota
	L2048N224SHA256
	L2048N256SHA256
	L3072N256SHA256
)

func (ps ParameterSizes) Hash() hash.Hash {
	domain, ok := kcdsainternal.GetDomain(int(ps))
	if !ok {
		panic(ErrInvalidParameterSizes.Error())
	}
	return domain.NewHash()
}

var (
	one = big.NewInt(1)
)

// Generate the paramters
// using the prime number generator used in crypto/dsa package.
func GenerateParameters(params *Parameters, rand io.Reader, sizes ParameterSizes) (err error) {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/crypto/dsa/dsa.go;l=65-155
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return ErrInvalidParameterSizes
	}

	qBytes := make([]byte, domain.B/8)
	pBytes := make([]byte, domain.A/8)

	q := new(big.Int)
	p := new(big.Int)
	rem := new(big.Int)

GeneratePrimes:
	for {
		if _, err := io.ReadFull(rand, qBytes); err != nil {
			return err
		}

		qBytes[len(qBytes)-1] |= 1
		qBytes[0] |= 0x80
		q.SetBytes(qBytes)

		if !q.ProbablyPrime(internal.NumMRTests) {
			continue
		}

		for i := 0; i < 4*domain.A; i++ {
			if _, err := io.ReadFull(rand, pBytes); err != nil {
				return err
			}

			pBytes[len(pBytes)-1] |= 1
			pBytes[0] |= 0x80

			p.SetBytes(pBytes)
			rem.Mod(p, q)
			rem.Sub(rem, one)
			p.Sub(p, rem)
			if p.BitLen() < domain.A {
				continue
			}

			if !p.ProbablyPrime(internal.NumMRTests) {
				continue
			}

			params.P = p
			params.Q = q
			break GeneratePrimes
		}
	}

	h := new(big.Int)
	h.SetInt64(2)
	g := new(big.Int)

	pm1 := new(big.Int).Sub(p, one)
	e := new(big.Int).Div(pm1, q)

	for {
		g.Exp(h, e, p)
		if g.Cmp(one) == 0 {
			h.Add(h, one)
			continue
		}

		params.G = g
		return nil
	}
}

func GenerateKey(priv *PrivateKey, rand io.Reader) error {
	if priv.P == nil || priv.Q == nil || priv.G == nil {
		return ErrParametersNotSetUp
	}

	x := new(big.Int)
	xBytes := make([]byte, priv.Q.BitLen()/8)
	xInv := new(big.Int)

	for {
		_, err := io.ReadFull(rand, xBytes)
		if err != nil {
			return err
		}
		x.SetBytes(xBytes)
		if x.Sign() != 0 && x.Cmp(priv.Q) < 0 {
			break
		}
	}

	// x의 역원 생성
	xInv = internal.FermatInverse(x, priv.Q)

	// 전자서명 검증키 y 생성(Y = G^{X^{-1} mod Q} mod P)
	priv.Y = new(big.Int).Exp(priv.G, xInv, priv.P)
	priv.X = x

	return nil
}

// Sign data using K generated randomly like in crypto/dsa packages.
func Sign(randReader io.Reader, priv *PrivateKey, h hash.Hash, data []byte) (r, s *big.Int, err error) {
	randutil.MaybeReadByte(randReader)

	if priv.Q.Sign() <= 0 || priv.P.Sign() <= 0 || priv.G.Sign() <= 0 || priv.X.Sign() <= 0 || priv.Q.BitLen()%8 != 0 {
		err = ErrInvalidPublicKey
		return
	}

	privQMinus1 := new(big.Int).Sub(priv.Q, one)

	// step 1. 난수 k를 [1, Q-1]에서 임의로 선택한다.
	var K *big.Int
	for {
		// K = [0 ~ q-2]
		K, err = rand.Int(randReader, privQMinus1)
		if err != nil {
			return
		}
		// k =  K + 1 -> [1 ~ q-1]
		K.Add(K, one)

		if K.Sign() > 0 && K.Cmp(priv.Q) < 0 {
			break
		}
	}

	return kcdsainternal.Sign(priv.P, priv.Q, priv.G, priv.Y, priv.X, K, h, data)
}

func Verify(pub *PublicKey, h hash.Hash, data []byte, R, S *big.Int) bool {
	return kcdsainternal.Verify(pub.P, pub.Q, pub.G, pub.Y, h, data, R, S)
}
