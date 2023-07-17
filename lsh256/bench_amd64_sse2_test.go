//go:build amd64

package lsh256

import "testing"

func Benchmark_LSH224_Reset_SSE2(b *testing.B) {
	benchReset(b, newContextAsm(lshType256H224, simdSetSSE2), true)
}
func Benchmark_LSH256_Reset_SSE2(b *testing.B) {
	benchReset(b, newContextAsm(lshType256H256, simdSetSSE2), true)
}

func Benchmark_LSH224_Write_SSE2(b *testing.B) {
	benchWrite(b, newContextAsm(lshType256H224, simdSetSSE2), true)
}
func Benchmark_LSH256_Write_SSE2(b *testing.B) {
	benchWrite(b, newContextAsm(lshType256H256, simdSetSSE2), true)
}

func Benchmark_LSH224_WriteSum_SSE2(b *testing.B) {
	benchSum(b, newContextAsm(lshType256H224, simdSetSSE2), true)
}
func Benchmark_LSH256_WriteSum_SSE2(b *testing.B) {
	benchSum(b, newContextAsm(lshType256H256, simdSetSSE2), true)
}
