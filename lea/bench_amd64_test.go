//go:build amd64

package lea

import (
	"testing"
)

func Benchmark_Encrypt_4Blocks_Go(b *testing.B)   { benchAll(b, block(b, 4, leaEnc4Go, true)) }
func Benchmark_Encrypt_4Blocks_SSE2(b *testing.B) { benchAll(b, block(b, 4, leaEnc4SSE2, true)) }

func Benchmark_Encrypt_8Blocks_Go(b *testing.B)   { benchAll(b, block(b, 4, leaEnc4Go, true)) }
func Benchmark_Encrypt_8Blocks_SSE2(b *testing.B) { benchAll(b, block(b, 4, leaEnc4SSE2, true)) }
func Benchmark_Encrypt_8Blocks_AVX2(b *testing.B) { benchAll(b, block(b, 4, leaEnc8AVX2, hasAVX2)) }
