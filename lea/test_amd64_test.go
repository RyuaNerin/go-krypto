//go:build amd64 && gc && !purego

package lea

import (
	"testing"
)

func Test_Encrypt_4Blocks_SSE2(t *testing.T) { testAll(t, tb(4, leaEnc4Go, leaEnc4SSE2, false)) }
func Test_Decrypt_4Blocks_SSE2(t *testing.T) { testAll(t, tb(4, leaDec4Go, leaDec4SSE2, false)) }

func Test_Encrypt_8Blocks_AVX2(t *testing.T) { testAll(t, tb(8, leaEnc8Go, leaEnc8AVX2, !hasAVX2)) }
func Test_Decrypt_8Blocks_AVX2(t *testing.T) { testAll(t, tb(8, leaDec8Go, leaDec8AVX2, !hasAVX2)) }

func Benchmark_Encrypt_4Blocks_SSE2(b *testing.B) { benchAll(b, bb(4, leaEnc4SSE2, false)) }
func Benchmark_Decrypt_4Blocks_SSE2(b *testing.B) { benchAll(b, bb(4, leaDec4SSE2, false)) }

func Benchmark_Encrypt_8Blocks_AVX2(b *testing.B) { benchAll(b, bb(8, leaEnc8AVX2, !hasAVX2)) }
func Benchmark_Decrypt_8Blocks_AVX2(b *testing.B) { benchAll(b, bb(8, leaDec8AVX2, !hasAVX2)) }
