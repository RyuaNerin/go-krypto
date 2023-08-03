//go:build amd64

package lsh256

import (
	"testing"

	"golang.org/x/sys/cpu"
)

func Benchmark_LSH256_Reset_AVX2(b *testing.B) {
	benchReset(b, newContextAsm(lshType256H256, simdSetAVX2), cpu.X86.HasAVX2)
}

func Benchmark_LSH256_Write_AVX2(b *testing.B) {
	benchWrite(b, newContextAsm(lshType256H256, simdSetAVX2), cpu.X86.HasAVX2)
}

func Benchmark_LSH256_WriteSum_AVX2(b *testing.B) {
	benchSum(b, newContextAsm(lshType256H256, simdSetAVX2), cpu.X86.HasAVX2)
}
