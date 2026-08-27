//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh512 && amd64

package lsh512

import (
	"testing"
)

var (
	htSSE2   = newHashTest(simdSSSE2{})
	htSSSE3  = newHashTest(simdSSSE3{})
	htAVX2   = newHashTest(simdAVX2{})
	htAVX512 = newHashTest(simdAVX512{})
)

func Test_LSH512_Hash_SSE2(t *testing.T)   { htSSE2.TestHash(t, false) }
func Test_LSH512_Hash_SSSE3(t *testing.T)  { htSSSE3.TestHash(t, !hasSSSE3) }
func Test_LSH512_Hash_AVX2(t *testing.T)   { htAVX2.TestHash(t, !hasAVX2) }
func Test_LSH512_Hash_AVX512(t *testing.T) { htAVX512.TestHash(t, !hasAVX512) }

func Test_LSH512_SMT_SSE2(t *testing.T)       { htSSE2.TestSMT(t, false) }
func Test_LSH512_SMT_Sum_SSE2(t *testing.T)   { htSSE2.TestSMTSum(t, false) }
func Test_LSH512_SMT_SSSE3(t *testing.T)      { htSSSE3.TestSMT(t, !hasSSSE3) }
func Test_LSH512_SMT_Sum_SSSE3(t *testing.T)  { htSSSE3.TestSMTSum(t, !hasSSSE3) }
func Test_LSH512_SMT_AVX2(t *testing.T)       { htAVX2.TestSMT(t, !hasAVX2) }
func Test_LSH512_SMT_Sum_AVX2(t *testing.T)   { htAVX2.TestSMTSum(t, !hasAVX2) }
func Test_LSH512_SMT_AVX512(t *testing.T)     { htAVX512.TestSMT(t, !hasAVX512) }
func Test_LSH512_SMT_Sum_AVX512(t *testing.T) { htAVX512.TestSMTSum(t, !hasAVX512) }

func Fuzz_LSH512_Hash_SSE2(f *testing.F)   { htSSE2.FuzzHash(f, false) }
func Fuzz_LSH512_Hash_SSSE3(f *testing.F)  { htSSSE3.FuzzHash(f, !hasSSSE3) }
func Fuzz_LSH512_Hash_AVX2(f *testing.F)   { htAVX2.FuzzHash(f, !hasAVX2) }
func Fuzz_LSH512_Hash_AVX512(f *testing.F) { htAVX512.FuzzHash(f, !hasAVX512) }

func Fuzz_LSH512_Sum_SSE2(f *testing.F)   { htSSE2.FuzzSum(f, false) }
func Fuzz_LSH512_Sum_SSSE3(f *testing.F)  { htSSSE3.FuzzSum(f, !hasSSSE3) }
func Fuzz_LSH512_Sum_AVX2(f *testing.F)   { htAVX2.FuzzSum(f, !hasAVX2) }
func Fuzz_LSH512_Sum_AVX512(f *testing.F) { htAVX512.FuzzSum(f, !hasAVX512) }

func Benchmark_LSH512_Hash_SSE2(b *testing.B)   { htSSE2.BenchHash(b, false) }
func Benchmark_LSH512_Hash_SSSE3(b *testing.B)  { htSSSE3.BenchHash(b, !hasSSSE3) }
func Benchmark_LSH512_Hash_AVX2(b *testing.B)   { htAVX2.BenchHash(b, !hasAVX2) }
func Benchmark_LSH512_Hash_AVX512(b *testing.B) { htAVX512.BenchHash(b, !hasAVX512) }

func Benchmark_LSH512_Sum_SSE2(b *testing.B)   { htSSE2.BenchSum(b, false) }
func Benchmark_LSH512_Sum_SSSE3(b *testing.B)  { htSSSE3.BenchSum(b, !hasSSSE3) }
func Benchmark_LSH512_Sum_AVX2(b *testing.B)   { htAVX2.BenchSum(b, !hasAVX2) }
func Benchmark_LSH512_Sum_AVX512(b *testing.B) { htAVX512.BenchSum(b, !hasAVX512) }
