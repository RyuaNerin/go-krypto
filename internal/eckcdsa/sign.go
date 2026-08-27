package eckcdsa

import (
	"crypto/elliptic"
	"hash"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

func SignUsingK(
	k, r, s *big.Int,
	curve elliptic.Curve, d *big.Int, xQ, yQ *big.Int,
	h hash.Hash,
	data []byte, buf []byte,
	tmp *big.Int,
) (bufOut []byte, ok bool) {
	curveParams := curve.Params()
	n := curveParams.N

	curveSize := internal.BitsToBytes(curveParams.BitSize)
	w := internal.DigitLen(curveParams.N)
	Lh := h.Size()
	L := h.BlockSize()

	LhIsBiggerThanW := Lh > w
	cQLen := max(L, 2*curveSize)

	buf = internal.Grow(buf, cQLen)

	////////////////////////////////////////////////////////////////////////////////////////////////////

	// 2: kG = (x1, y1) 계산
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("2: kG   = (x1, y1) 계산")
	kb := k.Bytes()
	defer zeroize.Zeroize(kb)
	x1, _ := curve.ScalarBaseMult(kb)
	x1Bytes := buf[:curveSize]
	x1.FillBytes(x1Bytes)
	zeroize.Zeroize(x1)
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
	clear(cQ)
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
	if t.Sign() <= 0 {
		return buf, false
	}

	// 9: t를 길이 w의 바이트 열 s로 변환
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("r = 0x" + hex.EncodeToString(r.Bytes()))
	// fmt.Println("s = 0x" + hex.EncodeToString(s.Bytes()))

	// 10: Σ = (r, s)를 반환
	return buf, true
}

func Verify(
	curve elliptic.Curve, xQ, yQ *big.Int,
	h hash.Hash,
	data []byte,
	r, s *big.Int,
) bool {
	curveParams := curve.Params()
	n := curveParams.N

	curveSize := internal.BitsToBytes(curveParams.BitSize)
	w := internal.DigitLen(curveParams.N)
	Lh := h.Size()
	L := h.BlockSize()

	LhIsBiggerThanW := Lh > w
	cQLen := max(L, 2*curveSize)

	tmp := make([]byte, cQLen)
	defer clear(tmp)

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
	defer zeroize.Zeroize(v)

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
	tQx, tQy := curve.ScalarMult(xQ, yQ, t.Bytes())
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
	defer zeroize.Zeroize(r2)
	// fmt.Println("r2 = 0x" + hex.EncodeToString(r2.Bytes()))
	// fmt.Println("r  = 0x" + hex.EncodeToString(r.Bytes()))

	return internal.BigEqual(r, r2)
}
