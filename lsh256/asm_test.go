//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh256 && (arm64 || amd64)

package lsh256

import (
	"hash"

	"github.com/RyuaNerin/go-krypto/internal/testing/thash"
)

func newHashTest(simd simdSet) thash.HashTest {
	return thash.HashTest{
		HashInfo: hi,
		New: map[int]func() hash.Hash{
			Size224: func() hash.Hash { return newContextWith(Size224, simd) },
			Size:    func() hash.Hash { return newContextWith(Size, simd) },
		},
		Sum: map[int]any{
			Size224: func(b []byte) (o [28]byte) {
				s := sumWith(Size224, b, simd)
				copy(o[:], s[:Size224])
				return o
			},
			Size: func(b []byte) [32]byte { return sumWith(Size, b, simd) },
		},
	}
}
