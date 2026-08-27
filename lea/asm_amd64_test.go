//go:build !noasm && !noasm_krypto && !noasm_krypto_lea && amd64

package lea

import (
	"testing"
)

var (
	cSSE2   = newCipherAsmSSE2
	cAVX2   = newCipherAsmAVX2
	cAVX512 = newCipherAsmAVX512
)

func Test_LEA_Block_SSE2(t *testing.T)   { ci.TestBlock(t, cSSE2, false) }
func Test_LEA_Block_AVX2(t *testing.T)   { ci.TestBlock(t, cAVX2, !hasAVX2) }
func Test_LEA_Block_AVX512(t *testing.T) { ci.TestBlock(t, cAVX512, !hasAVX512) }

func Fuzz_LEA_Block_SSE2(f *testing.F)       { ci.FuzzBlock(f, cSSE2, false) }
func Fuzz_LEA_Wide_SSE2(f *testing.F)        { ci.FuzzWide(f, cSSE2, false) }
func Fuzz_LEA_ShortBuffer_SSE2(f *testing.F) { ci.FuzzBlockShortBuffer(f, cSSE2, false) }

func Fuzz_LEA_Block_AVX2(f *testing.F)       { ci.FuzzBlock(f, cAVX2, !hasAVX2) }
func Fuzz_LEA_Wide_AVX2(f *testing.F)        { ci.FuzzWide(f, cAVX2, !hasAVX2) }
func Fuzz_LEA_ShortBuffer_AVX2(f *testing.F) { ci.FuzzBlockShortBuffer(f, cAVX2, !hasAVX2) }

func Fuzz_LEA_Block_AVX512(f *testing.F)       { ci.FuzzBlock(f, cAVX512, !hasAVX512) }
func Fuzz_LEA_Wide_AVX512(f *testing.F)        { ci.FuzzWide(f, cAVX512, !hasAVX512) }
func Fuzz_LEA_ShortBuffer_AVX512(f *testing.F) { ci.FuzzBlockShortBuffer(f, cAVX512, !hasAVX512) }

func Benchmark_LEA_Block_SSE2(b *testing.B)   { ci.BenchBlock4(b, cSSE2, false) }
func Benchmark_LEA_Block_AVX2(b *testing.B)   { ci.BenchBlock8(b, cAVX2, !hasAVX2) }
func Benchmark_LEA_Block_AVX512(b *testing.B) { ci.BenchBlock16(b, cAVX512, !hasAVX512) }
