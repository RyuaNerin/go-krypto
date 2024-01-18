// Package kcdsa implements the KCDSA(Korean Certificate-based Digital Signature Algorithm) as defined in TTAK.KO-12.0001/R4
package kcdsa

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
	"github.com/RyuaNerin/go-krypto/kcdsa/kcdsattak"
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

type paramValues struct {
	kcdsattak.Domain

	Int2L *big.Int
}

var (
	i2l_512 = new(big.Int).Exp(big.NewInt(2), big.NewInt(512), nil)

	paramValuesMap = map[ParameterSizes]paramValues{
		L2048N224SHA224: {
			Domain: kcdsattak.Domain{
				A:       2048,
				B:       224,
				LH:      28,
				NewHash: sha256.New224,
				L:       512,
			},
			Int2L: i2l_512,
		},
		L2048N224SHA256: {
			Domain: kcdsattak.Domain{
				A:       2048,
				B:       224,
				LH:      32,
				NewHash: sha256.New,
				L:       512,
			},
			Int2L: i2l_512,
		},
		L2048N256SHA256: {
			Domain: kcdsattak.Domain{
				A:       2048,
				B:       256,
				LH:      32,
				NewHash: sha256.New,
				L:       512,
			},
			Int2L: i2l_512,
		},
		L3072N256SHA256: {
			Domain: kcdsattak.Domain{
				A:       3072,
				B:       256,
				LH:      32,
				NewHash: sha256.New,
				L:       512,
			},
			Int2L: i2l_512,
		},
	}
)

func (ps ParameterSizes) domain() (paramValues, error) {
	p, ok := paramValuesMap[ps]
	if !ok {
		return p, ErrInvalidParameterSizes
	}
	return p, nil
}

// Generate the paramters
// using the prime number generator used in crypto/dsa package.
func GenerateParameters(params *Parameters, rand io.Reader, sizes ParameterSizes) (err error) {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/crypto/dsa/dsa.go;l=65-155
	domain, err := sizes.domain()
	if err != nil {
		return err
	}
	params.Sizes = sizes

	qBytes := make([]byte, domain.B/8)
	pBytes := make([]byte, domain.A/8)

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
		if xInv != nil {
			break
		}
	}

	// 전자서명 검증키 y 생성(Y = G^{X^{-1} mod Q} mod P)
	priv.Y = new(big.Int).Exp(priv.G, xInv, priv.P)
	priv.X = x

	return nil
}

// Sign data using K generated randomly like in crypto/dsa packages.
func Sign(randReader io.Reader, priv *PrivateKey, data []byte) (r, s *big.Int, err error) {
	switch priv.Sizes {
	case L2048N224SHA224:
	case L2048N224SHA256:
	case L2048N256SHA256:
	case L3072N256SHA256:
	default:
		err = ErrInvalidParameterSizes
		return
	}

	randutil.MaybeReadByte(randReader)

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

	return SignUsingK(K, priv, data)
}

// Sign data using K Specified
func SignUsingK(K *big.Int, priv *PrivateKey, data []byte) (r, s *big.Int, err error) {
	domain, err := priv.Sizes.domain()
	if err != nil {
		return nil, nil, err
	}
	h := domain.NewHash()

	// step 2. w = g^k mod p
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 2. w = g^k mod p")
	//fmt.Println("G = 0x" + hex.EncodeToString(priv.G.Bytes()))
	//fmt.Println("K = 0x" + hex.EncodeToString(K.Bytes()))
	//fmt.Println("P = 0x" + hex.EncodeToString(priv.P.Bytes()))
	W := new(big.Int).Exp(priv.G, K, priv.P)
	//fmt.Println("W = 0x" + hex.EncodeToString(W.Bytes()))

	//	step 3. R = h(W) mod 2^β (w를 바이트 열로 변환 후 해시한 결과의 바이트 열에서 	β 비트만큼 절삭):
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 3. R = h(W) mod 2^β")
	h.Reset()
	h.Write(W.Bytes())
	RBytes := h.Sum(nil)
	RBytes = RBytes[len(RBytes)-(domain.B/8):]
	R := new(big.Int).SetBytes(RBytes)
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	// step 4. Z = Y mod 2^l
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 4. Z = Y mod 2^l")
	//fmt.Println("Y = 0x" + hex.EncodeToString(priv.Y.Bytes()))
	//fmt.Println("2l = 0x" + hex.EncodeToString(i2l.Bytes()))
	Z := new(big.Int).Mod(priv.Y, domain.Int2L)
	ZBytes := Z.Bytes()
	//fmt.Println("Z = 0x" + hex.EncodeToString(ZBytes))

	// step 5. h = trunc(Hash(Z||M), β)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 5. h = trunc(Hash(Z||M), β)")
	h.Reset()
	h.Write(ZBytes)
	h.Write(data)
	HBytes := h.Sum(nil)
	HBytes = HBytes[len(HBytes)-(domain.B/8):]
	H := new(big.Int).SetBytes(HBytes)
	//fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))

	// step 6. E = (R xor H) mod Q
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 6. E = (R xor H) mod Q")
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))
	//fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))
	//fmt.Println("Q = 0x" + hex.EncodeToString(priv.Q.Bytes()))
	E := new(big.Int).Xor(R, H)
	E.Mod(E, priv.Q)
	//fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))

	//step 7. S = X(K-E) mod Q
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 7. S = X(K-E) mod Q")
	//fmt.Println("X = 0x" + hex.EncodeToString(priv.X.Bytes()))
	//fmt.Println("K = 0x" + hex.EncodeToString(K.Bytes()))
	//fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))
	//fmt.Println("Q = 0x" + hex.EncodeToString(priv.Q.Bytes()))
	K.Mod(K.Sub(K, E), priv.Q)
	S := new(big.Int).Mul(priv.X, K)
	S.Mod(S, priv.Q)
	//fmt.Println("S = 0x" + hex.EncodeToString(S.Bytes()))

	r = R
	s = S

	return
}

func Verify(pub *PublicKey, data []byte, R, S *big.Int) bool {
	domain, err := pub.Sizes.domain()
	if err != nil {
		return false
	}
	h := domain.NewHash()

	// step 1. 수신된 서명 {R', S'}에 대해 |R'|=LH, 0 < S' < Q 임을 확인한다.
	if pub.P.Sign() == 0 {
		return false
	}

	if S.Sign() < 1 || S.Cmp(pub.Q) >= 0 {
		return false
	}

	// step 2. Z = Y mod 2^l
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 2. Z = Y mod 2^l")
	//fmt.Println("Y = 0x" + hex.EncodeToString(pub.Y.Bytes()))
	//fmt.Println("2l = 0x" + hex.EncodeToString(i2l.Bytes()))
	Z := new(big.Int).Mod(pub.Y, domain.Int2L)
	ZBytes := Z.Bytes()
	//fmt.Println("Z = 0x" + hex.EncodeToString(ZBytes))

	// step 3. h = trunc(Hash(Z||M), β)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 3. h = trunc(Hash(Z||M), β)")
	h.Reset()
	h.Write(ZBytes)
	h.Write(data)
	HBytes := h.Sum(nil)
	HBytes = HBytes[len(HBytes)-(domain.B/8):]
	H := new(big.Int).SetBytes(HBytes)
	//fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))

	// step 4. E' = (R' xor H') mod Q
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 4. E' = (R' xor H') mod Q")
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))
	//fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))
	//fmt.Println("Q = 0x" + hex.EncodeToString(pub.Q.Bytes()))
	E := new(big.Int).Xor(R, H)
	E.Mod(E, pub.Q)
	//fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))

	// step 5. W' = Y ^ {S'} G ^ {E'} mod P
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 5. W' = Y ^ {S'} G ^ {E'} mod P")
	//fmt.Println("Y = 0x" + hex.EncodeToString(pub.Y.Bytes()))
	//fmt.Println("G = 0x" + hex.EncodeToString(pub.G.Bytes()))
	//fmt.Println("S = 0x" + hex.EncodeToString(S.Bytes()))
	//fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))
	//fmt.Println("P = 0x" + hex.EncodeToString(pub.P.Bytes()))
	W := new(big.Int).Exp(pub.Y, S, pub.P)
	E.Exp(pub.G, E, pub.P)
	W.Mul(W, E)
	W.Mod(W, pub.P)
	//fmt.Println("W = 0x" + hex.EncodeToString(W.Bytes()))

	// step 6. trunc(Hash(W'), β) = R'이 성립하는지 확인한다.
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 6. trunc(Hash(W'), β) = R'")
	h.Reset()
	h.Write(W.Bytes())
	rBytes := h.Sum(nil)
	rBytes = rBytes[len(rBytes)-(domain.B/8):]
	r := new(big.Int).SetBytes(rBytes)
	//fmt.Println("r = 0x" + hex.EncodeToString(r.Bytes()))
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	return bigIntEqual(R, r)
}
