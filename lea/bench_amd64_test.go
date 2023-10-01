//go:build amd64 && gc && !purego

package lea

import (
	"testing"
)

func Benchmark_Encrypt_4Blocks_SSE2(b *testing.B) { benchAll(b, block(b, 4, leaEnc4SSE2, true)) }
func Benchmark_Encrypt_8Blocks_SSE2(b *testing.B) { benchAll(b, block(b, 8, leaEnc8SSE2, true)) }
func Benchmark_Encrypt_8Blocks_AVX2(b *testing.B) { benchAll(b, block(b, 8, leaEnc8AVX2, hasAVX2)) }

func Benchmark_Decrypt_4Blocks_SSE2(b *testing.B) { benchAll(b, block(b, 4, leaDec4SSE2, true)) }
func Benchmark_Decrypt_8Blocks_SSE2(b *testing.B) { benchAll(b, block(b, 8, leaDec8SSE2, true)) }
func Benchmark_Decrypt_8Blocks_AVX2(b *testing.B) { benchAll(b, block(b, 8, leaDec8AVX2, hasAVX2)) }
