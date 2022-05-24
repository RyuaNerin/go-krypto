package krypto

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/lsh256"
	"github.com/RyuaNerin/go-krypto/lsh512"
)

func Benchmark_HASH_SHA256_1K(b *testing.B) {
	benchHash1k(b, sha256.New())
}
func Benchmark_HASH_SHA512_1K(b *testing.B) {
	benchHash1k(b, sha512.New())
}

func Benchmark_HASH_LSH256_1K(b *testing.B) {
	benchHash1k(b, lsh256.New())
}
func Benchmark_HASH_LSH512_1K(b *testing.B) {
	benchHash1k(b, lsh512.New())
}

func benchHash1k(b *testing.B, h hash.Hash) {
	buf := make([]byte, 1024)

	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(buf)
		h.Sum(nil)
	}
}
