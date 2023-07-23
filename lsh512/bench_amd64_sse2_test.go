//go:build amd64

package lsh512

import "testing"

func Benchmark_LSH512_Reset_SSE2(b *testing.B) { benchReset(b, newContextAsm(Size, simdSetSSE2), true) }
func Benchmark_LSH512_Write_SSE2(b *testing.B) {
	benchWrite(b, newContextAsm(Size, simdSetSSE2), true)
}
func Benchmark_LSH512_WriteSum_SSE2(b *testing.B) {
	benchSum(b, newContextAsm(Size, simdSetSSE2), true)
}
