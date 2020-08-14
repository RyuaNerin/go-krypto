package kipher

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"

	"kipher/lsh256"
	"kipher/lsh512"
)

func hash_1k(b *testing.B, h hash.Hash) {
	buf := make([]byte, 1024)

	for i := 0; i < b.N; i++ {
		h.Write(buf)
		h.Sum(nil)
	}
}

func Benchmark_HASH_SHA256_1K(b *testing.B) {
	hash_1k(b, sha256.New())
}
func Benchmark_HASH_SHA512_1K(b *testing.B) {
	hash_1k(b, sha512.New())
}

func Benchmark_HASH_LSH256_1K(b *testing.B) {
	hash_1k(b, lsh256.New())
}

func Benchmark_HASH_LSH512_1K(b *testing.B) {
	hash_1k(b, lsh512.New())
}
