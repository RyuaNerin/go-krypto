//go:build amd64

package lsh256

import "testing"

func Benchmark_LSH224_Reset_SSSE3(b *testing.B) {
	benchReset(b, newContextAsm(lshType256H224, simdSetSSSE3), true)
}
func Benchmark_LSH256_Reset_SSSE3(b *testing.B) {
	benchReset(b, newContextAsm(lshType256H256, simdSetSSSE3), true)
}

func Benchmark_LSH224_Write_SSSE3(b *testing.B) {
	benchWrite(b, newContextAsm(lshType256H224, simdSetSSSE3), true)
}
func Benchmark_LSH256_Write_SSSE3(b *testing.B) {
	benchWrite(b, newContextAsm(lshType256H256, simdSetSSSE3), true)
}

func Benchmark_LSH224_WriteSum_SSSE3(b *testing.B) {
	benchSum(b, newContextAsm(lshType256H224, simdSetSSSE3), true)
}
func Benchmark_LSH256_WriteSum_SSSE3(b *testing.B) {
	benchSum(b, newContextAsm(lshType256H256, simdSetSSSE3), true)
}
