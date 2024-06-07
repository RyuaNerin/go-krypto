package kcdsa

import (
	"encoding/binary"
	"errors"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

var ErrInvalidGenerationParameters = errors.New("krypto/kcdsa: invalid generation parameters")

type GeneratedParameter struct {
	P     *big.Int
	Q     *big.Int
	G     *big.Int
	J     *big.Int
	Seed  []byte
	Count int
	H     *big.Int
}

func GenerateParameters(rand io.Reader, domain Domain) (
	generated GeneratedParameter,
	err error,
) {
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

func RegenerateParameters(
	rand io.Reader,
	domain Domain,
	J *big.Int,
	seed []byte,
	count int,
) (
	P, Q, G *big.Int,
	err error,
) {
	P = new(big.Int)
	Q = new(big.Int)
	G = new(big.Int)

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
		err = ErrInvalidGenerationParameters
		return
	}

	// 11: 강한 소수 판정 알고리즘으로 q를 판정하여 소수가 아니면 단계 6으로 간다.
	if !Q.ProbablyPrime(internal.NumMRTests) {
		err = ErrInvalidGenerationParameters
		return
	}

	// 12: 강한 소수 판정 알고리즘으로 p를 판정하여 소수가 아니면 단계 6으로 간다
	if !P.ProbablyPrime(internal.NumMRTests) {
		err = ErrInvalidGenerationParameters
		return
	}

	H := new(big.Int)
	_, err = GenerateHG(H, G, buf, rand, P, J)
	if err != nil {
		return
	}

	return
}
