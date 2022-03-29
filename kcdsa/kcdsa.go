// Package kcdsa implements the KCDSA(Korean Certificate-based Digital Signature Algorithm) as defined in TTAK.KO-12.0001/R4
package kcdsa

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"hash"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal/randutil"
)

type Parameters struct {
	P, Q, G *big.Int
	Sizes   ParameterSizes
}

type PublicKey struct {
	Parameters
	Y *big.Int
}

type PrivateKey struct {
	PublicKey
	X *big.Int
}

var (
	ErrInvalidPublicKey      = errors.New("krypto/kcdsa: invalid public key")
	ErrInvalidParameterSizes = errors.New("krypto/kcdsa: invalid ParameterSizes")
	ErrParametersNotSetUp    = errors.New("krypto/kcdsa: parameters not set up before generating key")
)

type ParameterSizes int

const (
	L2048N224WithSHA224 ParameterSizes = iota
	L2048N224WithSHA256
	L2048N256WithSHA256
	L3072N256WithSHA256
)

const numMRTests = 64

func GenerateParameters(params *Parameters, rand io.Reader, sizes ParameterSizes) (err error) {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/crypto/dsa/dsa.go;l=65-155
	var L, N int
	switch sizes {
	case L2048N224WithSHA224:
		L = 2048
		N = 224
	case L2048N224WithSHA256:
		L = 2048
		N = 224
	case L2048N256WithSHA256:
		L = 2048
		N = 256
	case L3072N256WithSHA256:
		L = 3072
		N = 256
	default:
		return ErrInvalidParameterSizes
	}
	params.Sizes = sizes

	qBytes := make([]byte, N/8)
	pBytes := make([]byte, L/8)

	q := new(big.Int)
	p := new(big.Int)
	rem := new(big.Int)
	one := new(big.Int)
	one.SetInt64(1)

GeneratePrimes:
	for {
		if _, err := io.ReadFull(rand, qBytes); err != nil {
			return err
		}

		qBytes[len(qBytes)-1] |= 1
		qBytes[0] |= 0x80
		q.SetBytes(qBytes)

		if !q.ProbablyPrime(numMRTests) {
			continue
		}

		for i := 0; i < 4*L; i++ {
			if _, err := io.ReadFull(rand, pBytes); err != nil {
				return err
			}

			pBytes[len(pBytes)-1] |= 1
			pBytes[0] |= 0x80

			p.SetBytes(pBytes)
			rem.Mod(p, q)
			rem.Sub(rem, one)
			p.Sub(p, rem)
			if p.BitLen() < L {
				continue
			}

			if !p.ProbablyPrime(numMRTests) {
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
		xInv = fermatInverse(x, priv.Q)
		if xInv != nil {
			break
		}
	}

	// 전자서명 검증키 y 생성(Y = G^{X^{-1} mod Q} mod P)
	priv.Y = new(big.Int).Exp(priv.G, xInv, priv.P)
	priv.X = x

	return nil
}

var (
	biPow2l = new(big.Int).Exp(big.NewInt(2), big.NewInt(64*8), nil)
)

func Sign(randReader io.Reader, priv *PrivateKey, data io.Reader) (r, s *big.Int, err error) {
	randutil.MaybeReadByte(randReader)

	var h hash.Hash
	switch priv.Sizes {
	case L2048N224WithSHA224:
		h = sha256.New224()
	case L2048N224WithSHA256:
		h = sha256.New()
	case L2048N256WithSHA256:
		h = sha256.New()
	case L3072N256WithSHA256:
		h = sha256.New()
	default:
		err = ErrInvalidParameterSizes
		return
	}

	if priv.Q.Sign() <= 0 || priv.P.Sign() <= 0 || priv.G.Sign() <= 0 || priv.X.Sign() <= 0 || priv.Q.BitLen()%8 != 0 {
		err = ErrInvalidPublicKey
		return
	}

	one := new(big.Int)
	one.SetInt64(1)

	privQMinus1 := new(big.Int).Set(priv.Q)
	privQMinus1.Sub(privQMinus1, one)

	// step 1. 난수 k를 [1, Q-1]에서 임의로 선택한다.
	var K *big.Int
	for {
		K, err = rand.Int(randReader, privQMinus1)
		if err != nil {
			return
		}
		K.Add(K, one)

		if K.Sign() > 0 && K.Cmp(priv.Q) < 0 {
			break
		}
	}

	// step 2. W=G^K mod P를 계산한다.
	W := new(big.Int).Exp(priv.G, K, priv.P)

	//	step 3. 서명의 첫 부분 R=h(W)를 계산한다.
	h.Reset()
	h.Write(W.Bytes())
	R := new(big.Int).SetBytes(h.Sum(nil))

	// step 4. Z = Y mod 2^l
	// step 5. h = Hash(Z||M)을 계산한다.
	Z := new(big.Int).Set(priv.Y)
	Z.Mod(Z, biPow2l)
	ZBytes := make([]byte, 64)
	Z.FillBytes(ZBytes)

	h.Reset()
	h.Write(ZBytes)
	_, err = io.Copy(h, data)
	if err != nil {
		return
	}
	HBytes := h.Sum(nil)
	H := new(big.Int).SetBytes(HBytes)

	// step 6. E = (R ^ H) mod Q를 계산한다.
	E := new(big.Int).Exp(R, H, priv.Q)

	//step 7. S = X(K-E) mod Q를 계산한다.
	S := new(big.Int).Mul(priv.X, K.Sub(K, E))
	S.Mod(S, priv.Q)

	r = R
	s = S

	return
}

func Verify(pub *PublicKey, data io.Reader, R, S *big.Int) (ok bool, err error) {
	var h hash.Hash
	switch pub.Sizes {
	case L2048N224WithSHA224:
		h = sha256.New224()
	case L2048N224WithSHA256:
		h = sha256.New()
	case L2048N256WithSHA256:
		h = sha256.New()
	case L3072N256WithSHA256:
		h = sha256.New()
	default:
		return false, ErrInvalidParameterSizes
	}

	// step 1. 수신된 서명 {R', S'}에 대해 |R'|=LH, 0 < S' < Q 임을 확인한다.
	if pub.P.Sign() == 0 {
		return false, ErrInvalidPublicKey
	}

	if S.Sign() < 1 || S.Cmp(pub.Q) >= 0 {
		return false, ErrInvalidPublicKey
	}

	// step 2. Z = Y mod 2^l
	// step 3. h = Hash(Z||M)을 계산한다.
	Z := new(big.Int).Set(pub.Y)
	Z.Mod(Z, biPow2l)
	ZBytes := make([]byte, 64)
	Z.FillBytes(ZBytes)

	h.Reset()
	h.Write(ZBytes)
	_, err = io.Copy(h, data)
	if err != nil {
		return
	}
	HBytes := h.Sum(nil)
	H := new(big.Int).SetBytes(HBytes)

	// step 4. E' = (R' ^ H') mod Q을 계산한다.
	E := new(big.Int).Exp(R, H, pub.Q)

	// step 5. W' = Y ^ {S'} G ^ {E'} mod P를 계산한다.
	W := new(big.Int).Mul(
		new(big.Int).Exp(pub.Y, S, pub.P),
		new(big.Int).Exp(pub.G, E, pub.P),
	)
	W.Mod(W, pub.P)

	// step 6. h(W') = R'이 성립하는지 확인한다.
	h.Reset()
	h.Write(W.Bytes())
	r := new(big.Int).SetBytes(h.Sum(nil))

	return R.Cmp(r) == 0, nil
}

func fermatInverse(a, N *big.Int) *big.Int {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/crypto/dsa/dsa.go;l=188-192
	two := big.NewInt(2)
	nMinus2 := new(big.Int).Sub(N, two)
	return new(big.Int).Exp(a, nMinus2, N)
}
