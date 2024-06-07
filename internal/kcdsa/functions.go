package kcdsa

import (
	"hash"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

func GenerateJ(
	J *big.Int, buf []byte,
	///////////////
	seed []byte,
	h hash.Hash,
	d Domain,
) (bufNew []byte, ok bool) {
	// 2: Seed를 일방향 함수 PPGF의 입력으로 하여 비트 길이가 n = (α - β - 4)인 난수 U를 생성한다.
	// (U ← PPGF(Seed, n))
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("U ← PPGF(Seed, n)")
	bufNew = ppgf(buf, d.A-d.B-4, h, seed)
	// fmt.Println("U = 0x" + hex.EncodeToString(UBytes2))f

	// 3: U의 상위에 4 비트 '1000'을 붙이고 최하위 비트는 1로 만들어 이를 J로 둔다.
	// (J ← 2^(α-β-1) ∨ U ∨ 1)
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("J ← 2^(α-β-1) ∨ U ∨ 1")
	bufNew[len(bufNew)-1] |= 1
	J.SetBytes(bufNew)
	J.SetBit(J, d.A-d.B-1, 1)
	// fmt.Println("J = 0x" + hex.EncodeToString(J.Bytes()))

	// 4: 강한 소수 판정 알고리즘으로 J를 판정하여 소수가 아니면 단계 1로 간다.
	if !J.ProbablyPrime(internal.NumMRTests) {
		return
	}

	ok = true
	return
}

func GeneratePQ(
	P, Q *big.Int, buf []byte,
	///////////////
	J *big.Int,
	seed []byte,
	h hash.Hash,
	d Domain,
) (bufNew []byte, count int, ok bool) {
	// 5: Count를 0으로 둔다. (Count ← 0)
	count = 0

	var countB [4]byte

	bufNew = internal.Grow(buf, internal.BitsToBytes(d.B))

	// krypto: 성능 향상을 위해서 PPGF(Seed)의 State를 먼저 계산해둔다.
	ppgf := newPPGF(h, seed)

	// 7: Count > 2^24이면 단계 1로 간다.
	for count <= (1 << 24) {
		// 6: Count를 1 증가시킨다. (Count ← (Count + 1))
		internal.IncCtr(countB[:])
		count += 1

		// 8: Seed에 Count를 연접한 것을 일방향 함수 PPGF의 입력으로 하여 비트 길이가
		// β인 난수 UBytes를 생성한다. (UBytes ← PPGF(Seed ‖ Count, β))
		bufNew = ppgf.Read(bufNew, d.B, countB[:])

		// 9: U의 최상위 및 최하위 비트를 1로 만들어 이를 q로 둔다.
		// (q ← 2^(β-1) ∨ U ∨ 1)
		bufNew[len(bufNew)-1] |= 1
		Q.SetBytes(bufNew)
		Q.SetBit(Q, d.B-1, 1)

		// 10: p ← (2Jq + 1)의 비트 길이가 α보다 길면 단계 6으로 간다.
		P.Add(P.Lsh(P.Mul(J, Q), 1), internal.One)
		if P.BitLen() > d.A {
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

		ok = true
		return
	}

	return
}

func GenerateHG(
	H, G *big.Int,
	///////////////
	buf []byte,
	rand io.Reader,
	P, J *big.Int,
) (bufOut []byte, err error) {
	for {
		// 1: p보다 작은 임의의 수 h를 생성한다.
		//     1 < h < (p - 1)
		//
		//     0 < h < p 로 생성 한 다음에
		//
		bufOut, err = internal.ReadBigInt(H, rand, buf, P)
		if err != nil {
			return
		}
		H.Add(H, internal.Two)

		ok := GenerateG(G, P, J, H)
		if !ok {
			continue
		}

		return
	}
}

func GenerateG(
	G *big.Int,
	///////////////
	P, J, H *big.Int,
) (ok bool) {
	// 2: g ← h^(2J) mod p를 계산한다.
	G.Set(J)
	G.Lsh(G, 1)
	G.Exp(H, G, P)

	// 3: g = 1이면 단계 1로 간다.
	return G.Cmp(internal.One) != 0
}

func GenerateX(
	X *big.Int,
	Q *big.Int, upri, xkey []byte, h hash.Hash, d Domain,
) {
	// 3: XSEEDj ← PPGF(user_provided_random_input, b)
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("3: XSEEDj ← PPGF(user_provided_random_input, b)")
	xseed := ppgf(nil, d.B, h, upri)
	// fmt.Println("xseed = 0x" + hex.EncodeToString(xseed))

	// 4: XVAL ← (XKEY + XSEEDj) mod 2^b
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("4: XVAL ← (XKEY + XSEEDj) mod 2^b")
	var carry int
	xval := make([]byte, internal.BitsToBytes(d.B))
	for i := 0; i < len(xseed); i++ {
		idx := len(xseed) - i - 1
		sum := int(xseed[idx]) + carry

		if i < len(xkey) {
			sum += int(xkey[len(xkey)-i-1])
		}

		xval[idx] = byte(sum)
		carry = sum >> 8
	}
	xval = internal.RightMost(xval, d.B)
	// fmt.Println("xval = 0x" + hex.EncodeToString(xval.Bytes()))

	// 5: xj ← PPGF(XVAL, b) mod q
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("5: xj ← PPGF(XVAL, b) mod q")
	X.SetBytes(ppgf(xseed, d.B, h, xval))
	X.Mod(X, Q)
	// fmt.Println("X = 0x" + hex.EncodeToString(X.Bytes()))
}

func GenerateY(
	Y *big.Int,
	P, Q, G, X *big.Int,
) {
	// x의 역원 생성
	xInv := internal.FermatInverse(X, Q)

	// 전자서명 검증키 y 생성(Y = G^{X^{-1} mod Q} mod P)
	Y.Exp(G, xInv, P)
}
