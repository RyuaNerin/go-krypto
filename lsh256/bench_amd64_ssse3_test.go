//go:build amd64

package lsh256

import (
	"testing"

	"golang.org/x/sys/cpu"
)

func Benchmark_LSH256_Reset_SSSE3(b *testing.B) {
	benchReset(b, newContextAsm(lshType256H256, simdSetSSSE3), cpu.X86.HasSSSE3)
}

func Benchmark_LSH256_Write_SSSE3(b *testing.B) {
	benchWrite(b, newContextAsm(lshType256H256, simdSetSSSE3), cpu.X86.HasSSSE3)
}

func Benchmark_LSH256_WriteSum_SSSE3(b *testing.B) {
	benchSum(b, newContextAsm(lshType256H256, simdSetSSSE3), cpu.X86.HasSSSE3)
}
