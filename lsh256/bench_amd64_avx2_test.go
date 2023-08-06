//go:build amd64

package lsh256

import (
	"hash"
	"testing"
)

func Benchmark_Hash_8_AVX2(b *testing.B)  { benchmarkSize(b, newAVX2, 8, hasAVX2) }
func Benchmark_Hash_1K_AVX2(b *testing.B) { benchmarkSize(b, newAVX2, 1024, hasAVX2) }
func Benchmark_Hash_8K_AVX2(b *testing.B) { benchmarkSize(b, newAVX2, 8192, hasAVX2) }

func newAVX2(algType algType) hash.Hash {
	return newContextAsm(algType, simdSetAVX2)
}
