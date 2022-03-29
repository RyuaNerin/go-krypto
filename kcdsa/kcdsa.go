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

const numMRTests = 64

// Generate the paramters
// Used the prime number generator used in crypto/dsa package.
func GenerateParameters(params *Parameters, rand io.Reader, sizes ParameterSizes) (err error) {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/crypto/dsa/dsa.go;l=65-155
	var L, N int
	switch sizes {
	case L2048N224SHA224:
		L = 2048
		N = 224
	case L2048N224SHA256:
		L = 2048
		N = 224
	case L2048N256SHA256:
		L = 2048
		N = 256
	case L3072N256SHA256:
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
	i2l_512 = new(big.Int).Exp(big.NewInt(2), big.NewInt(512), nil)
)

func Sign(randReader io.Reader, priv *PrivateKey, data io.Reader) (r, s *big.Int, err error) {
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

	return sign(K, priv, data)
}

func getParams(params *Parameters) (h hash.Hash, b int, i2l *big.Int, err error) {
	// b = n/8
	// i2l = 사전 계산된 2^l
	switch params.Sizes {
	case L2048N224SHA224:
		return sha256.New224(), 224 / 8, i2l_512, nil
	case L2048N224SHA256:
		return sha256.New(), 224 / 8, i2l_512, nil
	case L2048N256SHA256:
		return sha256.New(), 256 / 8, i2l_512, nil
	case L3072N256SHA256:
		return sha256.New(), 256 / 8, i2l_512, nil
	default:
		return nil, 0, nil, ErrInvalidParameterSizes
	}
}

func sign(K *big.Int, priv *PrivateKey, data io.Reader) (r, s *big.Int, err error) {
	h, b, i2l, err := getParams(&priv.Parameters)
	if err != nil {
		return nil, nil, err
	}

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
	RBytes = RBytes[len(RBytes)-b:]
	R := new(big.Int).SetBytes(RBytes)
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	// step 4. Z = Y mod 2^l
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 4. Z = Y mod 2^l")
	//fmt.Println("Y = 0x" + hex.EncodeToString(priv.Y.Bytes()))
	//fmt.Println("2l = 0x" + hex.EncodeToString(i2l.Bytes()))
	Z := new(big.Int).Mod(priv.Y, i2l)
	ZBytes := Z.Bytes()
	//fmt.Println("Z = 0x" + hex.EncodeToString(ZBytes))

	// step 5. h = trunc(Hash(Z||M), β)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 5. h = trunc(Hash(Z||M), β)")
	h.Reset()
	h.Write(ZBytes)
	_, err = io.Copy(h, data)
	if err != nil {
		return
	}
	HBytes := h.Sum(nil)
	HBytes = HBytes[len(HBytes)-b:]
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

func Verify(pub *PublicKey, data io.Reader, R, S *big.Int) (ok bool, err error) {
	h, b, i2l, err := getParams(&pub.Parameters)
	if err != nil {
		return false, err
	}

	// step 1. 수신된 서명 {R', S'}에 대해 |R'|=LH, 0 < S' < Q 임을 확인한다.
	if pub.P.Sign() == 0 {
		return false, ErrInvalidPublicKey
	}

	if S.Sign() < 1 || S.Cmp(pub.Q) >= 0 {
		return false, ErrInvalidPublicKey
	}

	// step 2. Z = Y mod 2^l
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 2. Z = Y mod 2^l")
	//fmt.Println("Y = 0x" + hex.EncodeToString(pub.Y.Bytes()))
	//fmt.Println("2l = 0x" + hex.EncodeToString(i2l.Bytes()))
	Z := new(big.Int).Mod(pub.Y, i2l)
	ZBytes := Z.Bytes()
	//fmt.Println("Z = 0x" + hex.EncodeToString(ZBytes))

	// step 3. h = trunc(Hash(Z||M), β)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 3. h = trunc(Hash(Z||M), β)")
	h.Reset()
	h.Write(ZBytes)
	_, err = io.Copy(h, data)
	if err != nil {
		return
	}
	HBytes := h.Sum(nil)
	HBytes = HBytes[len(HBytes)-b:]
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
	rBytes = rBytes[len(rBytes)-b:]
	r := new(big.Int).SetBytes(rBytes)
	//fmt.Println("r = 0x" + hex.EncodeToString(r.Bytes()))
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	return R.Cmp(r) == 0, nil
}

func fermatInverse(a, N *big.Int) *big.Int {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/crypto/dsa/dsa.go;l=188-192
	two := big.NewInt(2)
	nMinus2 := new(big.Int).Sub(N, two)
	return new(big.Int).Exp(a, nMinus2, N)
}
