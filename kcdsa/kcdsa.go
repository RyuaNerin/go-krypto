// Package kcdsa implements the KCDSA(Korean Certificate-based Digital Signature Algorithm) as defined in TTAK.KO-12.0001/R4
package kcdsa

import (
	"errors"
	"hash"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	kcdsainternal "github.com/RyuaNerin/go-krypto/internal/kcdsa"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
)

var (
	ErrInvalidPublicKey      = errors.New("krypto/kcdsa: invalid public key")
	ErrInvalidTTAKParameters = errors.New("krypto/kcdsa: invalid ttak parameters")
	ErrInvalidParameterSizes = errors.New("krypto/kcdsa: invalid ParameterSizes")
	ErrParametersNotSetUp    = errors.New("krypto/kcdsa: parameters not set up before generating key")
)

type ParameterSizes int

const (
	L2048N224SHA224 ParameterSizes = kcdsainternal.L2048N224SHA224
	L2048N224SHA256 ParameterSizes = kcdsainternal.L2048N224SHA256
	L2048N256SHA256 ParameterSizes = kcdsainternal.L2048N256SHA256
	L3072N256SHA256 ParameterSizes = kcdsainternal.L3072N256SHA256
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

	tmp := make([]byte, internal.Bytes(domain.A))

	J := new(big.Int)

	P := new(big.Int)
	Q := new(big.Int)

GeneratePrimes:
	for {
		// 2: Seed를 일방향 함수 PPGF의 입력으로 하여 비트 길이가 n = (α - β - 4)인 난수 U를 생성한다.
		// (U ← PPGF(Seed, n))
		U, err := internal.ReadBits(tmp[:0], rand, domain.A-domain.B)
		if err != nil {
			return err
		}

		// 3: U의 상위에 4 비트 '1000'을 붙이고 최하위 비트는 1로 만들어 이를 J로 둔다.
		// (J ← 2^(α-β-1) ∨ U ∨ 1)
		U[0] = (U[0] & 0b0000_1111) | 0b1000_0000
		U[len(U)-1] |= 1
		J.SetBytes(U)

		// 4: 강한 소수 판정 알고리즘으로 J를 판정하여 소수가 아니면 단계 1로 간다.
		if !J.ProbablyPrime(internal.NumMRTests) {
			continue
		}

		for i := 0; i < 4*domain.A; i++ {
			// 8: Seed에 Count를 연접한 것을 일방향 함수 PPGF의 입력으로 하여 비트 길이가
			// β인 난수 U를 생성한다. (U ← PPGF(Seed ‖ Count, β))
			U, err := internal.ReadBits(tmp[:0], rand, domain.B)
			if err != nil {
				return err
			}

			// 9: U의 최상위 및 최하위 비트를 1로 만들어 이를 q로 둔다.
			// (q ← 2^(β-1) ∨ U ∨ 1)
			U[0] |= 0b1000_0000
			U[len(U)-1] |= 1
			Q.SetBytes(U)

			// 10: p ← (2Jq + 1)의 비트 길이가 α보다 길면 단계 6으로 간다.
			P.Add(P.Lsh(P.Mul(J, Q), 1), one)
			if P.BitLen() > domain.A {
				continue
			}

			// 11: 강한 소수 판정 알고리즘으로 q를 판정하여 소수가 아니면 단계 6으로 간다.
			if !Q.ProbablyPrime(internal.NumMRTests) {
				continue
			}

			// 12: 강한 소수 판정 알고리즘으로 p를 판정하여 소수가 아니면 단계 6으로 간다
			if !P.ProbablyPrime(internal.NumMRTests) {
				continue
			}

			params.P = P
			params.Q = Q
			break GeneratePrimes
		}
	}

	_, G, err := generateHG(rand, P, J)
	if err != nil {
		return err
	}
	params.G = G
	return nil
}

func GenerateKey(priv *PrivateKey, rand io.Reader) error {
	if priv.P == nil || priv.Q == nil || priv.G == nil {
		return ErrParametersNotSetUp
	}

	x := new(big.Int)
	xBytes := make([]byte, priv.Q.BitLen()/8)

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

	priv.X = x
	priv.Y = kcdsainternal.Y(priv.P, priv.Q, priv.G, priv.X)

	return nil
}

func Sign(randReader io.Reader, priv *PrivateKey, sizes ParameterSizes, data []byte) (r, s *big.Int, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, nil, ErrInvalidParameterSizes
	}

	randutil.MaybeReadByte(randReader)

	if priv.Q.Sign() <= 0 || priv.P.Sign() <= 0 || priv.G.Sign() <= 0 || priv.X.Sign() <= 0 || priv.Q.BitLen()%8 != 0 {
		err = ErrInvalidPublicKey
		return
	}

	qblen := priv.Q.BitLen()

	var K *big.Int
	var attempts int
	for attempts = 10; attempts > 0; attempts-- {
		// step 1. 난수 k를 [1, Q-1]에서 임의로 선택한다.

		for {
			b := make([]byte, internal.Bytes(qblen))
			if _, err = io.ReadFull(randReader, b); err != nil {
				return
			}
			if excess := len(b)*8 - qblen; excess > 0 {
				b[0] >>= excess
			}
			K = new(big.Int).SetBytes(b)
			K.Add(K, one)

			if K.Sign() > 0 && K.Cmp(priv.Q) < 0 {
				break
			}
		}

		r, s, err = sign(priv, domain, K, data)
		if err != nil {
			return nil, nil, err
		}
		if r.Sign() == 0 {
			continue
		}

		if s.Sign() != 0 {
			break
		}
	}

	// Only degenerate private keys will require more than a handful of
	// attempts.
	if attempts == 0 {
		return nil, nil, ErrInvalidPublicKey
	}

	return
}

func Verify(pub *PublicKey, sizes ParameterSizes, data []byte, R, S *big.Int) bool {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return false
	}

	// step 1. 수신된 서명 {R', S'}에 대해 |R'|=LH, 0 < S' < Q 임을 확인한다.
	if pub.P.Sign() <= 0 {
		return false
	}

	if R.Sign() < 1 {
		return false
	}
	if S.Sign() < 1 || S.Cmp(pub.Q) >= 0 {
		return false
	}

	return verify(pub, domain, data, R, S)
}

func sign(priv *PrivateKey, domain kcdsainternal.Domain, K *big.Int, data []byte) (r, s *big.Int, err error) {
	h := domain.NewHash()

	B := domain.B
	l := h.BlockSize()

	tmp := make([]byte, internal.Bytes(domain.A))

	// step 2. w = g^k mod p
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 2. w = g^k mod p")
	//fmt.Println("G = 0x" + hex.EncodeToString(priv.G.Bytes()))
	//fmt.Println("K = 0x" + hex.EncodeToString(K.Bytes()))
	//fmt.Println("P = 0x" + hex.EncodeToString(priv.P.Bytes()))

	W := new(big.Int).Exp(priv.G, K, priv.P)
	WBytes := tmp[:internal.Bytes(domain.A)]
	W.FillBytes(WBytes)

	//fmt.Println("W = 0x" + hex.EncodeToString(WBytes))
	//fmt.Println("W = 0x" + hex.EncodeToString(W.Bytes()))

	//	step 3. R = h(W) mod 2^β (w를 바이트 열로 변환 후 해시한 결과의 바이트 열에서 	β 비트만큼 절삭):
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 3. R = h(W) mod 2^β")
	h.Reset()
	h.Write(WBytes)
	RBytes := internal.TruncateLeft(h.Sum(tmp[:0]), B)
	R := new(big.Int).SetBytes(RBytes)
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	// step 4. Z = Y mod 2^l
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 4. Z = Y mod 2^l")
	//fmt.Println("Y = 0x" + hex.EncodeToString(priv.Y.Bytes()))
	ZBytesLen := internal.Bytes(priv.Y.BitLen())
	if ZBytesLen < l {
		ZBytesLen = l
	}
	ZBytes := tmp[:ZBytesLen]
	priv.Y.FillBytes(ZBytes)
	ZBytes = internal.TruncateLeft(ZBytes, l*8)
	//fmt.Println("Z = 0x" + hex.EncodeToString(ZBytes))

	// step 5. h = trunc(Hash(Z||M), β)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 5. h = trunc(Hash(Z||M), β)")
	h.Reset()
	h.Write(ZBytes)
	h.Write(data)
	HBytes := internal.TruncateLeft(h.Sum(tmp[:0]), B)
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

func verify(pub *PublicKey, domain kcdsainternal.Domain, data []byte, R, S *big.Int) bool {
	h := domain.NewHash()

	B := domain.B
	l := h.BlockSize()

	tmpSize := l
	YBytesLen := internal.Bytes(pub.Y.BitLen())
	PBytesLen := internal.Bytes(pub.P.BitLen())
	if tmpSize < YBytesLen {
		tmpSize = YBytesLen
	}
	if tmpSize < PBytesLen {
		tmpSize = PBytesLen
	}

	tmp := make([]byte, tmpSize)

	// step 2. Z = Y mod 2^l
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 2. Z = Y mod 2^l")
	//fmt.Println("Y = 0x" + hex.EncodeToString(pub.Y.Bytes()))
	if YBytesLen < l {
		YBytesLen = l
	}
	ZBytes := tmp[:YBytesLen]
	pub.Y.FillBytes(ZBytes)
	ZBytes = internal.TruncateLeft(ZBytes, l*8)
	//fmt.Println("Z = 0x" + hex.EncodeToString(ZBytes))

	// step 3. h = trunc(Hash(Z||M), β)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 3. h = trunc(Hash(Z||M), β)")
	h.Reset()
	h.Write(ZBytes)
	h.Write(data)
	HBytes := internal.TruncateLeft(h.Sum(tmp[:0]), B)
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

	WBytes := tmp[:PBytesLen]
	W.FillBytes(WBytes)
	//fmt.Println("W = 0x" + hex.EncodeToString(W.Bytes()))

	// step 6. trunc(Hash(W'), β) = R'이 성립하는지 확인한다.
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 6. trunc(Hash(W'), β) = R'")
	h.Reset()
	h.Write(WBytes)
	rBytes := internal.TruncateLeft(h.Sum(tmp[:0]), B)
	r := new(big.Int).SetBytes(rBytes)
	//fmt.Println("r = 0x" + hex.EncodeToString(r.Bytes()))
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	return internal.BigIntEqual(R, r)
}
