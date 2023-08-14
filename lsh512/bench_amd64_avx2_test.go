//go:build amd64

package lsh512

import (
	"hash"
	"testing"
)

func Benchmark_Hash_8_AVX2(b *testing.B)  { benchmarkSize(b, newAVX2, 8, true) }
func Benchmark_Hash_1K_AVX2(b *testing.B) { benchmarkSize(b, newAVX2, 1024, true) }
func Benchmark_Hash_8K_AVX2(b *testing.B) { benchmarkSize(b, newAVX2, 8192, true) }

func newAVX2(size int) hash.Hash {
	return newContextAsm(size, simdSetAVX2)
}
