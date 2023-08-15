//go:build amd64 && gc && !purego

package lsh512

import (
	"hash"
	"testing"
)

func Benchmark_Hash_8_SSE2(b *testing.B)  { benchmarkSize(b, newSSE2, 8, true) }
func Benchmark_Hash_1K_SSE2(b *testing.B) { benchmarkSize(b, newSSE2, 1024, true) }
func Benchmark_Hash_8K_SSE2(b *testing.B) { benchmarkSize(b, newSSE2, 8192, true) }

func newSSE2(size int) hash.Hash {
	return newContextAsm(size, simdSetSSE2)
}
