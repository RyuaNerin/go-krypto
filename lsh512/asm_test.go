//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh512 && (arm64 || amd64)

package lsh512

import (
	"hash"

	"github.com/RyuaNerin/go-krypto/internal/testing/thash"
)

func newHashTest(simd simdSet) thash.HashTest {
	return thash.HashTest{
		HashInfo: hi,
		New: map[int]func() hash.Hash{
			Size224: func() hash.Hash { return newContextWith(Size224, simd) },
			Size256: func() hash.Hash { return newContextWith(Size256, simd) },
			Size384: func() hash.Hash { return newContextWith(Size384, simd) },
			Size:    func() hash.Hash { return newContextWith(Size, simd) },
		},
		Sum: map[int]any{
			Size224: func(b []byte) (o [28]byte) {
				s := sumWith(Size224, b, simd)
				copy(o[:], s[:Size224])
				return o
			},
			Size256: func(b []byte) (o [32]byte) {
				s := sumWith(Size256, b, simd)
				copy(o[:], s[:Size256])
				return o
			},
			Size384: func(b []byte) (o [48]byte) {
				s := sumWith(Size384, b, simd)
				copy(o[:], s[:Size384])
				return o
			},
			Size: func(b []byte) [64]byte { return sumWith(Size, b, simd) },
		},
	}
}
