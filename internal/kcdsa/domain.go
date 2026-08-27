package kcdsa

import (
	"crypto"
	"hash"
	"math/big"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

const (
	A2048B224SHA224 = iota // 권고
	A2048B224SHA256        // A2048B224SHA224 와 강도 동일
	A2048B256SHA256        // 효율성이 떨어지지만, 필요에 따라 이용 가능하다.
	A3072B256SHA256        // 권고
	A1024B160HAS160        // 레거시
)

type Domain struct {
	A, B int // 소수 p와 q의 비트 길이를 각각 α와 β라 할 때, 두 값의 순서 쌍
	LH   int // 해시 코드의 비트 길이
	L    int // ℓ 해시 함수의 입력 블록 비트 길이

	NewHash func() hash.Hash
}

var paramValuesMap = map[int]Domain{
	A2048B224SHA224: {
		A:       2048,
		B:       224,
		NewHash: crypto.SHA224.New,
		LH:      224,
		L:       64,
	},
	A2048B224SHA256: {
		A:       2048,
		B:       224,
		NewHash: crypto.SHA256.New,
		LH:      256,
		L:       64,
	},
	A2048B256SHA256: {
		A:       2048,
		B:       256,
		NewHash: crypto.SHA256.New,
		LH:      256,
		L:       64,
	},
	A3072B256SHA256: {
		A:       3072,
		B:       256,
		NewHash: crypto.SHA256.New,
		LH:      256,
		L:       64,
	},
	A1024B160HAS160: {
		A:       1024,
		B:       160,
		NewHash: krypto.HAS160.New,
		LH:      krypto.HAS160.Size() * 8,
		L:       krypto.HAS160.BlockSize(),
	},
}

func GetDomain(sizes int) (Domain, bool) {
	p, ok := paramValuesMap[sizes]
	return p, ok
}

func IsSupportedSizes(P, Q *big.Int) bool {
	if P == nil || Q == nil || P.Sign() <= 0 || Q.Sign() <= 0 {
		return false
	}
	for _, domain := range paramValuesMap {
		if P.BitLen() == domain.A && Q.BitLen() == domain.B {
			return true
		}
	}
	return false
}

func IsValidElement(x, P *big.Int) bool {
	if x == nil || P == nil || x.Cmp(internal.One) <= 0 {
		return false
	}

	var pMinusOne big.Int
	defer zeroize.Zeroize(&pMinusOne)
	pMinusOne.Sub(P, internal.One)
	return x.Cmp(&pMinusOne) < 0
}

func IsValidPublicKey(P, Q, G, Y *big.Int, domain Domain) bool {
	return P != nil &&
		Q != nil &&
		P.Sign() > 0 && Q.Sign() > 0 &&
		P.BitLen() == domain.A &&
		Q.BitLen() == domain.B &&
		IsValidElement(G, P) &&
		IsValidElement(Y, P)
}

func IsValidParameters(P, Q, G *big.Int) bool {
	if !IsSupportedSizes(P, Q) || !IsValidElement(G, P) {
		return false
	}
	if !Q.ProbablyPrime(internal.NumMRTests) || !P.ProbablyPrime(internal.NumMRTests) {
		return false
	}

	var remainder big.Int
	defer zeroize.Zeroize(&remainder)

	remainder.Sub(P, internal.One)
	if remainder.Mod(&remainder, Q).Sign() != 0 {
		return false
	}

	return IsSubgroupElement(G, P, Q)
}

func IsSubgroupElement(x, P, Q *big.Int) bool {
	if !IsValidElement(x, P) || Q == nil || Q.Sign() <= 0 {
		return false
	}
	var power big.Int
	defer zeroize.Zeroize(&power)

	return power.Exp(x, Q, P).Cmp(internal.One) == 0
}
