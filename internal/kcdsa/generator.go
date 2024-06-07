package kcdsa

import (
	"encoding/binary"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

type GeneratedParameter struct {
	P     *big.Int
	Q     *big.Int
	G     *big.Int
	J     *big.Int
	Seed  []byte
	Count int
	H     *big.Int
}

func GenerateParametersFast(rand io.Reader, d Domain) (generated GeneratedParameter, err error) {
	P, Q, G := new(big.Int), new(big.Int), new(big.Int)

	buf := make([]byte, internal.BitsToBytes(d.A))

	// TTAK
	// U: random len(U) = (α - β - 4)
	// J: 2^(a-b-1) ∨ U ∨ 1
	// q: 2^(b-1) ∨ U ∨ 1
	// p: 2Jq + 1, len(p) <= a
	// H: 1 < H < (p - 1)
	// g: h^(2J) mod p, g > 1
	//
	// Fast
	// p: prime, 2^(a-1)  < U < 2^a
	// q: p-1 is multiple of q, 2^(b-1)  < U < 2^b,
	// F: 1 < F < p -1, F ** ((p - 1)/Q) mod p > 1
	// G = F**((p-1)/Q) mod p, Z*p에 있는 위수 q의 요소

	tmp := new(big.Int)
	F := new(big.Int)

	// https://github.com/golang/go/blob/go1.22.4/src/crypto/dsa/dsa.go#L98-L136
GeneratePrimes:
	for {
		if buf, err = internal.ReadBits(buf, rand, d.B); err != nil {
			return
		}
		buf[len(buf)-1] |= 1
		Q.SetBytes(buf)
		Q.SetBit(Q, d.B-1, 1)

		if !Q.ProbablyPrime(internal.NumMRTests) {
			continue
		}

		for i := 0; i < 4*d.A; i++ {
			if buf, err = internal.ReadBits(buf, rand, d.A); err != nil {
				return
			}
			buf[len(buf)-1] |= 1
			P.SetBytes(buf)
			P.SetBit(P, d.A-1, 1)

			// P - (P % Q) - 1
			P.Sub(P, tmp.Sub(tmp.Mod(P, Q), internal.One))
			if P.BitLen() < d.A {
				continue
			}

			if !P.ProbablyPrime(internal.NumMRTests) {
				continue
			}

			break GeneratePrimes
		}
	}

	tmp.Div(tmp.Sub(P, internal.One), Q)

	for {
		if buf, err = internal.ReadBits(buf, rand, d.A); err != nil {
			return
		}
		F.SetBytes(buf)
		F.Add(F, internal.Two)
		if F.Cmp(P) >= 0 {
			continue
		}

		G.Exp(F, tmp, P)
		if G.Cmp(internal.One) <= 0 {
			continue
		}
		if G.Cmp(P) >= 0 {
			continue
		}

		break
	}

	return GeneratedParameter{
		P: P,
		Q: Q,
		G: G,
	}, nil
}

func GenerateParametersTTAK(rand io.Reader, domain Domain) (generated GeneratedParameter, err error) {
	h := domain.NewHash()

	generated.J = new(big.Int)
	generated.P = new(big.Int)
	generated.Q = new(big.Int)
	generated.H = new(big.Int)
	generated.G = new(big.Int)

	// p. 13
	generated.Seed = make([]byte, internal.BitsToBytes(domain.B))

	var ok bool
	var buf []byte
	for {
		_, err = io.ReadFull(rand, generated.Seed)
		if err != nil {
			return
		}

		// 2 ~ 4
		buf, ok = GenerateJ(generated.J, buf, generated.Seed, h, domain)
		if !ok {
			continue
		}

		// 5 ~ 12
		buf, generated.Count, ok = GeneratePQ(generated.P, generated.Q, buf, generated.J, generated.Seed, h, domain)
		if !ok {
			continue
		}

		_, err = GenerateHG(generated.H, generated.G, buf, rand, generated.P, generated.J)
		if err != nil {
			return
		}

		return
	}
}

func RegeneratePQ(
	domain Domain,
	J *big.Int,
	seed []byte,
	count int,
) (
	P, Q *big.Int,
	ok bool,
) {
	P = new(big.Int)
	Q = new(big.Int)

	var CountB [4]byte
	binary.BigEndian.PutUint32(CountB[:], uint32(count))

	buf := make([]byte, internal.BitsToBytes(domain.B))

	// 8: Seed에 Count를 연접한 것을 일방향 함수 PPGF의 입력으로 하여 비트 길이가
	// β인 난수 U를 생성한다. (U ← PPGF(Seed ‖ Count, β))
	U := ppgf(buf, domain.B, domain.NewHash(), seed, CountB[:])

	// 9: U의 최상위 및 최하위 비트를 1로 만들어 이를 q로 둔다.
	// (q ← 2^(β-1) ∨ U ∨ 1)
	U[len(U)-1] |= 1
	Q.SetBytes(U)
	Q.SetBit(Q, domain.B-1, 1)

	// 10: p ← (2Jq + 1)의 비트 길이가 α보다 길면 단계 6으로 간다.
	P.Add(P.Lsh(P.Mul(J, Q), 1), internal.One)
	if P.BitLen() > domain.A {
		return nil, nil, false
	}

	// 11: 강한 소수 판정 알고리즘으로 q를 판정하여 소수가 아니면 단계 6으로 간다.
	if !Q.ProbablyPrime(internal.NumMRTests) {
		return nil, nil, false
	}

	// 12: 강한 소수 판정 알고리즘으로 p를 판정하여 소수가 아니면 단계 6으로 간다
	if !P.ProbablyPrime(internal.NumMRTests) {
		return nil, nil, false
	}

	return P, Q, true
}
