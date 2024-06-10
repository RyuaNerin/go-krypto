package kcdsa

import (
	"crypto/sha256"
	"hash"

	"github.com/RyuaNerin/go-krypto/has160" //nolint:staticcheck
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
		NewHash: sha256.New224,
		LH:      sha256.Size224 * 8,
		L:       sha256.BlockSize,
	},
	A2048B224SHA256: {
		A:       2048,
		B:       224,
		NewHash: sha256.New,
		LH:      sha256.Size * 8,
		L:       sha256.BlockSize,
	},
	A2048B256SHA256: {
		A:       2048,
		B:       256,
		NewHash: sha256.New,
		LH:      sha256.Size * 8,
		L:       sha256.BlockSize,
	},
	A3072B256SHA256: {
		A:       3072,
		B:       256,
		NewHash: sha256.New,
		LH:      sha256.Size * 8,
		L:       sha256.BlockSize,
	},
	A1024B160HAS160: {
		A:       1024,
		B:       160,
		NewHash: has160.New,
		LH:      has160.Size * 8,
		L:       has160.BlockSize,
	},
}

func GetDomain(sizes int) (Domain, bool) {
	p, ok := paramValuesMap[sizes]
	return p, ok
}
