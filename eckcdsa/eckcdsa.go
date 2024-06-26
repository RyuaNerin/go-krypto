// Package eckcdsa implements the EC-KCDSA(Korean Certificate-based Digital Signature Algorithm using Elliptic Curves) as defined in TTAK.KO-12.0015/R3
package eckcdsa

import (
	"crypto/elliptic"
	"errors"
	"hash"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	eckcdsainternal "github.com/RyuaNerin/go-krypto/internal/eckcdsa"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
)

// Generate the parameters
func GenerateKey(c elliptic.Curve, randReader io.Reader) (*PrivateKey, error) {
	randutil.MaybeReadByte(randReader)

	// https://cs.opensource.google/go/go/+/refs/tags/go1.20.7:src/crypto/ecdsa/ecdsa_legacy.go;l=20-31
	D := new(big.Int)
	err := randFieldElement(D, c, randReader)
	if err != nil {
		return nil, err
	}

	priv := &PrivateKey{
		D: D,
		PublicKey: PublicKey{
			Curve: c,
		},
	}
	priv.PublicKey.X, priv.PublicKey.Y = eckcdsainternal.XY(D, c)

	return priv, nil
}

// Sign data using K generated randomly like in crypto/ecdsa packages.
func Sign(rand io.Reader, priv *PrivateKey, h hash.Hash, data []byte) (r, s *big.Int, err error) {
	randutil.MaybeReadByte(rand)

	if priv == nil || priv.Curve == nil || priv.X == nil || priv.Y == nil || priv.D == nil || !priv.Curve.IsOnCurve(priv.X, priv.Y) {
		return nil, nil, errors.New(msgParametersNotSetUp)
	}

	h.Reset()
	h.Write(data)
	digest := h.Sum(nil)

	csprng, err := randutil.MixedCSPRNG(rand, priv.D.Bytes(), digest)
	if err != nil {
		return nil, nil, err
	}

	var ok bool
	k := new(big.Int)
	r, s = new(big.Int), new(big.Int)

	var buf []byte
	tmp := new(big.Int)
	for {
		// https://cs.opensource.google/go/go/+/refs/tags/go1.20.7:src/crypto/ecdsa/ecdsa_legacy.go;l=77-113
		err = randFieldElement(k, priv.Curve, csprng)
		if err != nil {
			return
		}

		buf, ok = signUsingK(k, r, s, priv, h, data, buf, tmp)
		if !ok {
			continue
		}

		return
	}
}

func signUsingK(k, r, s *big.Int, priv *PrivateKey, h hash.Hash, data []byte, buf []byte, tmp *big.Int) (bufOut []byte, ok bool) {
	curve := priv.Curve
	curveParams := curve.Params()
	n := curveParams.N

	curveSize := internal.BitsToBytes(curveParams.BitSize)
	w := internal.BigCeilLog2(curveParams.N, 8)
	Lh := h.Size()
	L := h.BlockSize()
	d := priv.D
	xQ := priv.X
	yQ := priv.Y

	LhIsBiggerThanW := Lh > w
	cQLen := L
	if cQLen < 2*curveSize {
		cQLen = 2 * curveSize
	}

	buf = internal.Grow(buf, cQLen)

	////////////////////////////////////////////////////////////////////////////////////////////////////

	// 2: kG = (x1, y1) 계산
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("2: kG   = (x1, y1) 계산")
	x1, _ := curve.ScalarBaseMult(k.Bytes())
	x1Bytes := buf[:curveSize]
	x1.FillBytes(x1Bytes)
	// fmt.Println("kGx, x1 = 0x" + hex.EncodeToString(x1Bytes))

	// 3: r ← Hash(x1)
	// 해시 코드의 바이트 길이 LH가 w( = log256n)보다 큰 경우 r ← Hash(x1) 연산을 r ← Hash(x1) mod 2^8w로
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("3: r ← Hash(x1)")
	// fmt.Println("kGx, x1 = 0x" + hex.EncodeToString(x1Bytes))
	h.Reset()
	h.Write(x1Bytes)
	rBytes := h.Sum(buf[:0])
	// fmt.Println("r       = 0x" + hex.EncodeToString(rBytes))
	if LhIsBiggerThanW {
		rBytes = rBytes[len(rBytes)-w:]
	}
	r.SetBytes(rBytes)
	// fmt.Println("r       = 0x" + hex.EncodeToString(r.Bytes()))

	// 4: cQ ← MSB(xQ ‖ yQ, L)
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("4: cQ ← MSB(xQ ‖ yQ, L)")
	// fmt.Println("xQ = 0x" + hex.EncodeToString(xQ.Bytes()))
	// fmt.Println("yQ = 0x" + hex.EncodeToString(yQ.Bytes()))
	cQ := buf[:cQLen]
	xQ.FillBytes(cQ[:curveSize])
	yQ.FillBytes(cQ[curveSize : curveSize*2])
	cQ = cQ[:L]
	// fmt.Println("cQ = 0x" + hex.EncodeToString(cQ))

	// 5: v ← Hash(cQ ‖ M)
	// 해시 코드의 바이트 길이 LH가 w( = log256n)보다 큰 경우 v ← Hash(cQ ‖ M) 연산을 v ← Hash(cQ ‖ M) mod 2^(8w)로 대체한다
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("5: v ← Hash(cQ ‖ M)")
	// fmt.Println("cQ = 0x" + hex.EncodeToString(cQ))
	// fmt.Println("M  = 0x" + hex.EncodeToString(M))
	h.Reset()
	h.Write(cQ)
	h.Write(data)
	vBytes := h.Sum(buf[:0])
	if LhIsBiggerThanW {
		vBytes = vBytes[len(vBytes)-w:]
	}
	v := tmp.SetBytes(vBytes)
	// fmt.Println("v  = 0x" + hex.EncodeToString(v.Bytes()))

	// 6: e ← (r ⊕ v) mod n
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("6: e ← (r ⊕ v) mod n")
	// fmt.Println("r = 0x" + hex.EncodeToString(r.Bytes()))
	// fmt.Println("v = 0x" + hex.EncodeToString(v.Bytes()))
	// fmt.Println("n = 0x" + hex.EncodeToString(n.Bytes()))
	e := tmp
	e.Mod(e.Xor(r, v), n)
	// fmt.Println("e = 0x" + hex.EncodeToString(e.Bytes()))

	// 7: t ← d(k - e) mod n
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("7: t ← d(k - e) mod n")
	// fmt.Println("d = 0x" + hex.EncodeToString(d.Bytes()))
	// fmt.Println("v = 0x" + hex.EncodeToString(v.Bytes()))
	// fmt.Println("n = 0x" + hex.EncodeToString(n.Bytes()))
	t := s
	t.Mod(t.Sub(k, e), n)
	t.Mod(t.Mul(d, t), n)
	// fmt.Println("t = 0x" + hex.EncodeToString(t.Bytes()))

	// 8: 만약 t = 0이면 단계 1로 간다.
	if t.Sign() <= 0 { //nolint:gosimple
		return buf, false
	}

	// 9: t를 길이 w의 바이트 열 s로 변환
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("r = 0x" + hex.EncodeToString(r.Bytes()))
	// fmt.Println("s = 0x" + hex.EncodeToString(s.Bytes()))

	// 10: Σ = (r, s)를 반환
	return buf, true
}

func Verify(pub *PublicKey, h hash.Hash, data []byte, r, s *big.Int) bool {
	if pub == nil || pub.Curve == nil || pub.X == nil || pub.Y == nil || !pub.Curve.IsOnCurve(pub.X, pub.Y) {
		return false
	}
	if r.Sign() <= 0 || s.Sign() <= 0 {
		return false
	}

	curve := pub.Curve
	curveParams := pub.Curve.Params()
	n := curveParams.N

	curveSize := internal.BitsToBytes(curveParams.BitSize)
	w := internal.BigCeilLog2(curveParams.N, 8)
	Lh := h.Size()
	L := h.BlockSize()
	xQ := pub.X
	yQ := pub.Y

	LhIsBiggerThanW := Lh > w
	cQLen := L
	if cQLen < 2*curveSize {
		cQLen = 2 * curveSize
	}

	tmp := make([]byte, cQLen)

	////////////////////////////////////////////////////////////////////////////////////////////////////

	// 0: (선택 사항) 서명자의 인증서를 확인하고, 서명 검증에 필요한 도메인 변수와 검증키 Q 추출
	// 1: Σ′ = (r′, s′)에 대해 |r′|가 해시 코드의 비트 길이와 일치하는지 여부와 s′을 정수로 변환한 t′에 대해 0 < t′ < n임을 확인
	// 해시 코드의 바이트 길이 LH가 w( = log256n)보다 큰 경우 단계 1의 r′의 바이트 길이와 LH의 비교를 r′의 바이트 길이와 w의 비교
	t := s

	if r.Sign() <= 0 {
		return false
	}
	if LhIsBiggerThanW {
		if internal.BitsToBytes(r.BitLen()) > curveSize {
			return false
		}
	} else {
		if internal.BitsToBytes(r.BitLen()) > Lh {
			return false
		}
	}
	if t.Sign() <= 0 || t.Cmp(n) >= 0 {
		return false
	}

	// 2: cQ ← MSB(xQ ‖ yQ, L)
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("2: cQ ← MSB(xQ ‖ yQ, L)")
	// fmt.Println("xQ = 0x" + hex.EncodeToString(xQ.Bytes()))
	// fmt.Println("yQ = 0x" + hex.EncodeToString(yQ.Bytes()))
	cQ := tmp[:cQLen]
	xQ.FillBytes(cQ[:curveSize])
	yQ.FillBytes(cQ[curveSize : curveSize*2])
	cQ = cQ[:L]
	// fmt.Println("cQ = 0x" + hex.EncodeToString(cQ))

	// 3: v′ ← Hash(cQ ‖ M′)
	// 해시 코드의 바이트 길이 LH가 w( = log256n)보다 큰 경우 단계 3의 v′ ← Hash(cQ ‖ M′) 연산을 v′ ← Hash(cQ ‖ M′) mod 2^(8w) 으로
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("3: v′ ← Hash(cQ ‖ M′)")
	// fmt.Println("cQ = 0x" + hex.EncodeToString(cQ))
	// fmt.Println("M  = 0x" + hex.EncodeToString(M))
	h.Reset()
	h.Write(cQ)
	h.Write(data)
	vBytes := h.Sum(tmp[:0])
	// fmt.Println("v  = 0x" + hex.EncodeToString(vBytes))
	if LhIsBiggerThanW {
		vBytes = vBytes[len(vBytes)-w:]
	}
	// fmt.Println("v  = 0x" + hex.EncodeToString(vBytes))
	v := new(big.Int).SetBytes(vBytes)

	// 4: e′ ← (r′ ⊕ v′) mod n
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("4: e′ ← (r′ ⊕ v′) mod n")
	// fmt.Println("r  = 0x" + hex.EncodeToString(r.Bytes()))
	// fmt.Println("v  = 0x" + hex.EncodeToString(v.Bytes()))
	// fmt.Println("n  = 0x" + hex.EncodeToString(n.Bytes()))
	e := v.Xor(r, v)
	// fmt.Println("e  = 0x" + hex.EncodeToString(e.Bytes()))
	e.Mod(e, n)
	// fmt.Println("e% = 0x" + hex.EncodeToString(e.Bytes()))

	// 5: (x2, y2) ← t′Q + e′G
	//		Q : 서명자의 검증키
	//		G : EC-KCDSA 도메인 변수의 하나로, EC-KCDSA는 기본점 G에 의해 생성되는 타원 곡선 부분군에서 정의
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("5: (x2, y2) ← t′Q + e′G")
	tQx, tQy := curve.ScalarMult(pub.X, pub.Y, t.Bytes())
	eGx, eGy := curve.ScalarBaseMult(e.Bytes())
	x2, _ := curve.Add(tQx, tQy, eGx, eGy)
	x2Bytes := tmp[:curveSize]
	x2.FillBytes(x2Bytes)

	// fmt.Println("tQ.x  = 0x" + hex.EncodeToString(tQx.Bytes()))
	// fmt.Println("tQ.y  = 0x" + hex.EncodeToString(tQy.Bytes()))
	// fmt.Println("eG.x  = 0x" + hex.EncodeToString(eGx.Bytes()))
	// fmt.Println("eG.y  = 0x" + hex.EncodeToString(eGy.Bytes()))
	// fmt.Println("x2    = 0x" + hex.EncodeToString(x2Bytes))

	// 6: Hash(x2′) = r′ 여부 확인
	// 해시 코드의 바이트 길이 LH가 w( = log256n)보다 큰 경우 단계 6의 Hash(x2′) = r′ 연산을 Hash(x2′) mod 2^(8w) = r′로 대체한다.
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("6: Hash(x2′) = r′ 여부 확인")
	// fmt.Println("x2 = 0x" + hex.EncodeToString(x2Bytes))
	h.Reset()
	h.Write(x2Bytes)
	rBytes := h.Sum(tmp[:0])
	if LhIsBiggerThanW {
		rBytes = rBytes[len(rBytes)-w:]
	}
	r2 := new(big.Int).SetBytes(rBytes)
	// fmt.Println("r2 = 0x" + hex.EncodeToString(r2.Bytes()))
	// fmt.Println("r  = 0x" + hex.EncodeToString(r.Bytes()))

	return internal.BigEqual(r, r2)
}

// https://cs.opensource.google/go/go/+/refs/tags/go1.20.7:src/crypto/ecdsa/ecdsa_legacy.go;l=168-188
// randFieldElement returns a random element of the order of the given
// curve using the procedure given in FIPS 186-4, Appendix B.5.2.
func randFieldElement(k *big.Int, c elliptic.Curve, rand io.Reader) (err error) {
	// See randomPoint for notes on the algorithm. This has to match, or s390x
	// signatures will come out different from other architectures, which will
	// break TLS recorded tests.

	N := c.Params().N                   // 1. N = len(n)
	b := make([]byte, (N.BitLen()+7)/8) // 2. If N is invalid, then return an ERROR indication, Invalid_d, and Invalid_Q.

	for {
		if _, err = io.ReadFull(rand, b); err != nil { // 3
			return
		}
		if excess := len(b)*8 - N.BitLen(); excess > 0 {
			b[0] >>= excess
		}
		// 6. If (c > n-2), then go to step 4.
		// 7. d = c + 1.
		//
		// d > n-1 ===> d >= n
		k.SetBytes(b)
		if k.Sign() > 0 && k.Cmp(N) < 0 {
			return
		}
	}

	/**
	B.4.2 Key Pair Generation by Testing Candidates

	In this method, a random number is obtained and tested to determine that it will produce a value
		of d in the correct range. If d is out-of-range, another random number is obtained (i.e., the
		process is iterated until an acceptable value of d is obtained.
	The following process or its equivalent may be used to generate an ECC key pair.
	Input:
		1. (q, FR, a, b {, domain_parameter_seed}, G, n, h)
			The domain parameters that are used for this process. n is a prime number,
			and G is a point on the elliptic curve.
	Output:
		1. status The status returned from the key pair generation procedure. The status will
			indicate SUCCESS or an ERROR.
		2. (d, Q) The generated private and public keys. If an error is encountered during
			the generation process, invalid values for d and Q should be returned, as
			represented by Invalid_d and Invalid_Q in the following specification. d is
			an integer, and Q is an elliptic curve point. The generated private key d is
			in the range [1, n-1].

	Process:
		1. N = len(n)
		2. If N is invalid, then return an ERROR indication, Invalid_d, and Invalid_Q.
		3. requested_security_strength = the security strength associated with N; see SP 800-57, Part 1.
		4. Obtain a string of N returned_bits from an RBG with a security strength of
			requested_security_strength or more. If an ERROR indication is returned, then
			return an ERROR indication, Invalid_d, and Invalid_Q.
		5. Convert returned_bits to the (non-negative) integer c (see Appendix C.2.1)
		6. If (c > n-2), then go to step 4.
		7. d = c + 1.
	*/
}
