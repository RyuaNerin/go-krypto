//go:build amd64 && gc && !purego

package lsh512

import (
	"hash"
	"testing"
)

func Benchmark_Hash_8_AVX2(b *testing.B)  { benchmarkSize(b, newAVX2, 8, hasAVX2) }
func Benchmark_Hash_1K_AVX2(b *testing.B) { benchmarkSize(b, newAVX2, 1024, hasAVX2) }
func Benchmark_Hash_8K_AVX2(b *testing.B) { benchmarkSize(b, newAVX2, 8192, hasAVX2) }

func newAVX2(size int) hash.Hash {
	return newContextAsm(size, simdSetAVX2)
}
