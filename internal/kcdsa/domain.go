package internal

import (
	"crypto/sha256"
	"hash"
)

const (
	L2048N224SHA224 = iota
	L2048N224SHA256
	L2048N256SHA256
	L3072N256SHA256
)

type Domain struct {
	A, B int // 소수 p와 q의 비트 길이를 각각 α와 β라 할 때, 두 값의 순서 쌍
	LH   int // 해시 코드의 비트 길이
	L    int // ℓ 해시 함수의 입력 블록 비트 길이

	NewHash func() hash.Hash
}

var (
	paramValuesMap = map[int]Domain{
		L2048N224SHA224: {
			A:       2048,
			B:       224,
			LH:      28,
			NewHash: sha256.New224,
			L:       512,
		},
		L2048N224SHA256: {
			A:       2048,
			B:       224,
			LH:      32,
			NewHash: sha256.New,
			L:       512,
		},
		L2048N256SHA256: {
			A:       2048,
			B:       256,
			LH:      32,
			NewHash: sha256.New,
			L:       512,
		},
		L3072N256SHA256: {
			A:       3072,
			B:       256,
			LH:      32,
			NewHash: sha256.New,
			L:       512,
		},
	}
)

func GetDomain(sizes int) (Domain, bool) {
	p, ok := paramValuesMap[sizes]
	return p, ok
}
