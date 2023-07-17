//go:build amd64

package lea

import (
	"testing"
)

func Benchmark_LEA128_Encrypt_4Blocks_Go(b *testing.B)   { block(b, 128, 4, leaEnc4Go, true) }
func Benchmark_LEA128_Encrypt_4Blocks_SSE2(b *testing.B) { block(b, 128, 4, leaEnc4SSE2, true) }

func Benchmark_LEA196_Encrypt_4Blocks_Go(b *testing.B)   { block(b, 196, 4, leaEnc4Go, true) }
func Benchmark_LEA196_Encrypt_4Blocks_SSE2(b *testing.B) { block(b, 196, 4, leaEnc4SSE2, true) }

func Benchmark_LEA256_Encrypt_4Blocks_Go(b *testing.B)   { block(b, 256, 4, leaEnc4Go, true) }
func Benchmark_LEA256_Encrypt_4Blocks_SSE2(b *testing.B) { block(b, 256, 4, leaEnc4SSE2, true) }

func Benchmark_LEA128_Encrypt_8Blocks_Go(b *testing.B)   { block(b, 128, 8, leaEnc8Go, true) }
func Benchmark_LEA128_Encrypt_8Blocks_SSE2(b *testing.B) { block(b, 128, 8, leaEnc8SSE2, true) }
func Benchmark_LEA128_Encrypt_8Blocks_AVX2(b *testing.B) { block(b, 128, 8, leaEnc8AVX2, true) }

func Benchmark_LEA196_Encrypt_8Blocks_Go(b *testing.B)   { block(b, 196, 8, leaEnc8Go, true) }
func Benchmark_LEA196_Encrypt_8Blocks_SSE2(b *testing.B) { block(b, 196, 8, leaEnc8SSE2, true) }
func Benchmark_LEA196_Encrypt_8Blocks_AVX2(b *testing.B) { block(b, 196, 8, leaEnc8AVX2, true) }

func Benchmark_LEA256_Encrypt_8Blocks_Go(b *testing.B)   { block(b, 256, 8, leaEnc8Go, true) }
func Benchmark_LEA256_Encrypt_8Blocks_SSE2(b *testing.B) { block(b, 256, 8, leaEnc8SSE2, true) }
func Benchmark_LEA256_Encrypt_8Blocks_AVX2(b *testing.B) { block(b, 256, 8, leaEnc8AVX2, true) }

func Benchmark_LEA128_Decrypt_4Blocks_Go(b *testing.B)   { block(b, 128, 4, leaDec4Go, true) }
func Benchmark_LEA128_Decrypt_4Blocks_SSE2(b *testing.B) { block(b, 128, 4, leaDec4SSE2, true) }

func Benchmark_LEA196_Decrypt_4Blocks_Go(b *testing.B)   { block(b, 196, 4, leaDec4Go, true) }
func Benchmark_LEA196_Decrypt_4Blocks_SSE2(b *testing.B) { block(b, 196, 4, leaDec4SSE2, true) }

func Benchmark_LEA256_Decrypt_4Blocks_Go(b *testing.B)   { block(b, 256, 4, leaDec4Go, true) }
func Benchmark_LEA256_Decrypt_4Blocks_SSE2(b *testing.B) { block(b, 256, 4, leaDec4SSE2, true) }

func Benchmark_LEA128_Decrypt_8Blocks_Go(b *testing.B)   { block(b, 128, 8, leaDec8Go, true) }
func Benchmark_LEA128_Decrypt_8Blocks_SSE2(b *testing.B) { block(b, 128, 8, leaDec8SSE2, true) }
func Benchmark_LEA128_Decrypt_8Blocks_AVX2(b *testing.B) { block(b, 128, 8, leaDec8AVX2, true) }

func Benchmark_LEA196_Decrypt_8Blocks_Go(b *testing.B)   { block(b, 196, 8, leaDec8Go, true) }
func Benchmark_LEA196_Decrypt_8Blocks_SSE2(b *testing.B) { block(b, 196, 8, leaDec8SSE2, true) }
func Benchmark_LEA196_Decrypt_8Blocks_AVX2(b *testing.B) { block(b, 196, 8, leaDec8AVX2, true) }

func Benchmark_LEA256_Decrypt_8Blocks_Go(b *testing.B)   { block(b, 256, 8, leaDec8Go, true) }
func Benchmark_LEA256_Decrypt_8Blocks_SSE2(b *testing.B) { block(b, 256, 8, leaDec8SSE2, true) }
func Benchmark_LEA256_Decrypt_8Blocks_AVX2(b *testing.B) { block(b, 256, 8, leaDec8AVX2, true) }
