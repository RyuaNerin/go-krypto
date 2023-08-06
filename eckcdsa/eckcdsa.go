// Package eckcdsa implements the EC-KCDSA(Korean Certificate-based Digital Signature Algorithm using Elliptic Curves) as defined in TTAK.KO-12.0015/R3
package eckcdsa

import (
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"hash"
	"io"
	"math"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
)

type PublicKey struct {
	elliptic.Curve

	X *big.Int
	Y *big.Int
}
type PrivateKey struct {
	PublicKey
	D *big.Int
}

var (
	ErrParametersNotSetUp = errors.New("krypto/eckcdsa: parameters not set up before generating key")
	ErrInvalidK           = errors.New("krypto/eckcdsa: invalid k. use other value")
)

type paramValues struct {
	/**
	MSB(S, L) 바이트 열 S의 맨 좌측 MSB_L 바이트
	여기에서 MSB_L은 해시 함수 메시지 입력 블록 바이트 길이
	*/
	MSB_L int
	K     int

	NewHash func() hash.Hash
}

var (
	one = big.NewInt(1)
	two = big.NewInt(2)
)

func GenerateKey(c elliptic.Curve, randReader io.Reader) (*PrivateKey, error) {
	randutil.MaybeReadByte(randReader)

	// https://cs.opensource.google/go/go/+/refs/tags/go1.20.7:src/crypto/ecdsa/ecdsa_legacy.go;l=20-31
	d, err := randFieldElement(c, randReader)
	if err != nil {
		return nil, err
	}

	dInv := internal.FermatInverse(d, c.Params().N)

	priv := new(PrivateKey)
	priv.PublicKey.Curve = c
	priv.D = d
	priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(dInv.Bytes())
	return priv, nil
}

// https://cs.opensource.google/go/go/+/refs/tags/go1.20.7:src/crypto/ecdsa/ecdsa_legacy.go;l=168-188
// randFieldElement returns a random element of the order of the given
// curve using the procedure given in FIPS 186-4, Appendix B.5.2.
func randFieldElement(c elliptic.Curve, rand io.Reader) (k *big.Int, err error) {
	// See randomPoint for notes on the algorithm. This has to match, or s390x
	// signatures will come out different from other architectures, which will
	// break TLS recorded tests.
	for {
		N := c.Params().N
		b := make([]byte, (N.BitLen()+7)/8)
		if _, err = io.ReadFull(rand, b); err != nil {
			return
		}
		if excess := len(b)*8 - N.BitLen(); excess > 0 {
			b[0] >>= excess
		}
		k = new(big.Int).SetBytes(b)
		if k.Sign() != 0 && k.Cmp(N) < 0 {
			return
		}
	}
}

func Sign(randReader io.Reader, priv *PrivateKey, h hash.Hash, M []byte) (r, s *big.Int, err error) {
	randutil.MaybeReadByte(randReader)

	curveParams := priv.Curve.Params()
	n := curveParams.N

	Nsub1 := new(big.Int).Sub(n, one)

	var k *big.Int
	for i := 0; i < 100; i++ {
		// 1: 난수 k를 [1, (n - 1)]에서 임의로 선택(8.2절 참조).
		k, err = rand.Int(randReader, Nsub1)
		if err != nil {
			return nil, nil, err
		}

		if k.Sign() != 1 {
			continue
		}

		r, s, err = SignWithK(k, priv, h, M)
		if err == ErrInvalidK {
			continue
		}
		return
	}
	return
}

func SignWithK(k *big.Int, priv *PrivateKey, h hash.Hash, M []byte) (r, s *big.Int, err error) {
	if priv == nil || priv.Curve == nil || priv.X == nil || priv.Y == nil || priv.D == nil || !priv.Curve.IsOnCurve(priv.X, priv.Y) {
		return nil, nil, ErrParametersNotSetUp
	}

	curveParams := priv.Curve.Params()
	n := curveParams.N

	w := int(math.Ceil(float64(curveParams.BitSize) / 8))
	K := w
	Lh := h.Size()
	L := h.BlockSize()
	d := priv.D
	xQ := priv.X
	yQ := priv.Y

	Lh_is_bigger_than_w := Lh > w

	var two_8w *big.Int
	if Lh_is_bigger_than_w {
		two_8w = big.NewInt(8)
		two_8w.Exp(two, two_8w.Mul(two_8w, big.NewInt(int64(w))), nil)
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////

	// 2: kG = (x1, y1) 계산
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("2: kG   = (x1, y1) 계산")
	x1, _ := priv.Curve.ScalarBaseMult(k.Bytes())
	//fmt.Println("kGx, x1 = 0x" + hex.EncodeToString(x1.Bytes()))

	// 3: r ← Hash(x1)
	//해시 코드의 바이트 길이 LH가 w( = log256n)보다 큰 경우 r ← Hash(x1) 연산을 r ← Hash(x1) mod 2^8w로
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("3: r ← Hash(x1)")
	//fmt.Println("kGx, x1 = 0x" + hex.EncodeToString(x1.Bytes()))
	h.Reset()
	h.Write(x1.Bytes())
	rBytes := h.Sum(nil)
	r = new(big.Int).SetBytes(rBytes)
	if Lh_is_bigger_than_w {
		r = r.Mod(r, two_8w)
	}
	//fmt.Println("r       = 0x" + hex.EncodeToString(r.Bytes()))

	// 4: cQ ← MSB(xQ ‖ yQ, L)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("4: cQ ← MSB(xQ ‖ yQ, L)")
	//fmt.Println("xQ = 0x" + hex.EncodeToString(xQ.Bytes()))
	//fmt.Println("yQ = 0x" + hex.EncodeToString(yQ.Bytes()))
	cQ := append(
		padLeft(xQ.Bytes(), K),
		padLeft(yQ.Bytes(), K)...,
	)
	cQ = padRight(cQ, L)
	//fmt.Println("cQ = 0x" + hex.EncodeToString(cQ))

	// 5: v ← Hash(cQ ‖ M)
	//해시 코드의 바이트 길이 LH가 w( = log256n)보다 큰 경우 v ← Hash(cQ ‖ M) 연산을 v ← Hash(cQ ‖ M) mod 2^(8w)로 대체한다
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("5: v ← Hash(cQ ‖ M)")
	//fmt.Println("cQ = 0x" + hex.EncodeToString(cQ))
	//fmt.Println("M  = 0x" + hex.EncodeToString(M))
	h.Reset()
	h.Write(cQ)
	h.Write(M)
	vBytes := h.Sum(nil)
	v := new(big.Int).SetBytes(vBytes)
	if Lh_is_bigger_than_w {
		v = v.Mod(v, two_8w)
	}
	//fmt.Println("v  = 0x" + hex.EncodeToString(v.Bytes()))

	// 6: e ← (r ⊕ v) mod n
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("6: e ← (r ⊕ v) mod n")
	//fmt.Println("r = 0x" + hex.EncodeToString(r.Bytes()))
	//fmt.Println("v = 0x" + hex.EncodeToString(v.Bytes()))
	//fmt.Println("n = 0x" + hex.EncodeToString(n.Bytes()))
	e := new(big.Int)
	e.Mod(e.Xor(r, v), n)
	//fmt.Println("e = 0x" + hex.EncodeToString(e.Bytes()))

	// 7: t ← d(k - e) mod n
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("7: t ← d(k - e) mod n")
	//fmt.Println("d = 0x" + hex.EncodeToString(d.Bytes()))
	//fmt.Println("v = 0x" + hex.EncodeToString(v.Bytes()))
	//fmt.Println("n = 0x" + hex.EncodeToString(n.Bytes()))
	t := new(big.Int)
	t.Mod(t.Sub(k, e), n)
	t.Mod(t.Mul(d, t), n)
	//fmt.Println("t = 0x" + hex.EncodeToString(t.Bytes()))

	//8: 만약 t = 0이면 단계 1로 간다.
	if t.Sign() <= 0 {
		return nil, nil, ErrInvalidK
	}

	//9: t를 길이 w의 바이트 열 s로 변환
	//fmt.Println("--------------------------------------------------")
	s = t
	//fmt.Println("r = 0x" + hex.EncodeToString(r.Bytes()))
	//fmt.Println("s = 0x" + hex.EncodeToString(s.Bytes()))

	//10: Σ = (r, s)를 반환
	return r, s, nil
}

func Verify(pub *PublicKey, h hash.Hash, M []byte, r, s *big.Int) bool {
	if pub == nil || pub.Curve == nil || pub.X == nil || pub.Y == nil || !pub.Curve.IsOnCurve(pub.X, pub.Y) {
		return false
	}

	curveParams := pub.Curve.Params()
	n := curveParams.N

	w := int(math.Ceil(float64(curveParams.BitSize) / 8))
	K := w
	L := h.BlockSize()
	Lh := h.Size()
	xQ := pub.X
	yQ := pub.Y

	Lh_is_bigger_than_w := Lh > w

	////////////////////////////////////////////////////////////////////////////////////////////////////

	// 0: (선택 사항) 서명자의 인증서를 확인하고, 서명 검증에 필요한 도메인 변수와 검증키 Q 추출
	// 1: Σ′ = (r′, s′)에 대해 |r′|가 해시 코드의 비트 길이와 일치하는지 여부와 s′을 정수로 변환한 t′에 대해 0 < t′ < n임을 확인
	// 해시 코드의 바이트 길이 LH가 w( = log256n)보다 큰 경우 단계 1의 r′의 바이트 길이와 LH의 비교를 r′의 바이트 길이와 w의 비교
	t := s

	// 사전 계산
	two_8w := big.NewInt(8)
	two_8w.Exp(two, two_8w.Mul(two_8w, big.NewInt(int64(w))), nil)

	if r.Sign() <= 0 {
		return false
	}
	if Lh_is_bigger_than_w {
		if r.BitLen()/8 > w {
			return false
		}
	} else {
		if r.BitLen()/8 > Lh {
			return false
		}
	}
	if t.Sign() <= 0 || t.Cmp(n) >= 0 {
		return false
	}

	// 2: cQ ← MSB(xQ ‖ yQ, L)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("2: cQ ← MSB(xQ ‖ yQ, L)")
	//fmt.Println("xQ = 0x" + hex.EncodeToString(xQ.Bytes()))
	//fmt.Println("yQ = 0x" + hex.EncodeToString(yQ.Bytes()))
	cQ := append(
		padLeft(xQ.Bytes(), K),
		padLeft(yQ.Bytes(), K)...,
	)
	cQ = padRight(cQ, L)
	//fmt.Println("cQ = 0x" + hex.EncodeToString(cQ))

	// 3: v′ ← Hash(cQ ‖ M′)
	// 해시 코드의 바이트 길이 LH가 w( = log256n)보다 큰 경우 단계 3의 v′ ← Hash(cQ ‖ M′) 연산을 v′ ← Hash(cQ ‖ M′) mod 2^(8w) 으로
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("3: v′ ← Hash(cQ ‖ M′)")
	//fmt.Println("cQ = 0x" + hex.EncodeToString(cQ))
	//fmt.Println("M  = 0x" + hex.EncodeToString(M))
	h.Reset()
	h.Write(cQ)
	h.Write(M)
	vBytes := h.Sum(nil)
	v := new(big.Int).SetBytes(vBytes)
	if Lh_is_bigger_than_w {
		v.Mod(v, two_8w)
	}

	// 4: e′ ← (r′ ⊕ v′) mod n
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("4: e′ ← (r′ ⊕ v′) mod n")
	//fmt.Println("r  = 0x" + hex.EncodeToString(r.Bytes()))
	//fmt.Println("v  = 0x" + hex.EncodeToString(v.Bytes()))
	//fmt.Println("n  = 0x" + hex.EncodeToString(n.Bytes()))
	e := new(big.Int).Xor(r, v)
	e.Mod(e, n)
	//fmt.Println("e  = 0x" + hex.EncodeToString(e.Bytes()))

	// 5: (x2, y2) ← t′Q + e′G
	//		Q : 서명자의 검증키
	//		G : EC-KCDSA 도메인 변수의 하나로, EC-KCDSA는 기본점 G에 의해 생성되는 타원 곡선 부분군에서 정의
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("5: (x2, y2) ← t′Q + e′G")
	x21, y21 := curveParams.ScalarMult(pub.X, pub.Y, t.Bytes())
	x22, y22 := curveParams.ScalarBaseMult(e.Bytes())
	x2, _ := curveParams.Add(x21, y21, x22, y22)
	//fmt.Println("x2  = 0x" + hex.EncodeToString(x2.Bytes()))

	// 6: Hash(x2′) = r′ 여부 확인
	// 해시 코드의 바이트 길이 LH가 w( = log256n)보다 큰 경우 단계 6의 Hash(x2′) = r′ 연산을 Hash(x2′) mod 2^(8w) = r′로 대체한다.
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("6: Hash(x2′) = r′ 여부 확인")
	//fmt.Println("x2 = 0x" + hex.EncodeToString(x2.Bytes()))
	h.Reset()
	h.Write(x2.Bytes())
	rBytes := h.Sum(nil)
	r2 := new(big.Int).SetBytes(rBytes)
	if Lh_is_bigger_than_w {
		r2.Mod(r2, two_8w)
	}
	//fmt.Println("r2 = 0x" + hex.EncodeToString(r2.Bytes()))
	//fmt.Println("r  = 0x" + hex.EncodeToString(r.Bytes()))

	return r.Cmp(r2) == 0
}

func padLeft(arr []byte, l int) []byte {
	if len(arr) >= l {
		return arr[:l]
	}

	n := make([]byte, l)
	copy(n[l-len(arr):], arr)

	return n
}

func padRight(arr []byte, l int) []byte {
	if len(arr) >= l {
		return arr[:l]
	}

	n := make([]byte, l)
	copy(n, arr)

	return n
}
