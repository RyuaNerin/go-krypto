//go:build amd64 && !purego && (!gccgo || go1.18)
// +build amd64
// +build !purego
// +build !gccgo go1.18

package lsh256

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/aligned"

	. "github.com/RyuaNerin/testingutil"
)

var (
	newSSE2  = simdSetSSE2.NewContext
	newSSSE3 = simdSetSSSE3.NewContext
	newAVX2  = simdSetAVX2.NewContext
)

func Test_Aligned(t *testing.T) {
	var ctx lsh256ContextAsm
	aligned.TestIsAligned(t, 16, &ctx, "algType")
	aligned.TestIsAligned(t, 16, &ctx, "remainDataByteLen")
	aligned.TestIsAligned(t, 16, &ctx, "cvL")
	aligned.TestIsAligned(t, 16, &ctx, "cvR")
	aligned.TestIsAligned(t, 16, &ctx, "lastBlock")
}

func Test_ShortWrite_SSE2(t *testing.T)  { HTSWA(t, as, newSSE2, !hasSSE2) }
func Test_ShortWrite_SSSE3(t *testing.T) { HTSWA(t, as, newSSSE3, !hasSSSE3) }
func Test_ShortWrite_AVX2(t *testing.T)  { HTSWA(t, as, newAVX2, !hasAVX2) }

func Test_LSH224_SSE2(t *testing.T) { HT(t, newSSE2(Size224), testCases224, !hasSSE2) }
func Test_LSH256_SSE2(t *testing.T) { HT(t, newSSE2(Size), testCases256, !hasSSE2) }

func Test_LSH224_SSSE3(t *testing.T) { HT(t, newSSSE3(Size224), testCases224, !hasSSSE3) }
func Test_LSH256_SSSE3(t *testing.T) { HT(t, newSSSE3(Size), testCases256, !hasSSSE3) }

func Test_LSH224_AVX2(t *testing.T) { HT(t, newAVX2(Size224), testCases224, !hasAVX2) }
func Test_LSH256_AVX2(t *testing.T) { HT(t, newAVX2(Size), testCases256, !hasAVX2) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Hash_8_SSE2(b *testing.B)  { HBA(b, as, newSSE2, 8, !hasSSE2) }
func Benchmark_Hash_1K_SSE2(b *testing.B) { HBA(b, as, newSSE2, 1024, !hasSSE2) }
func Benchmark_Hash_8K_SSE2(b *testing.B) { HBA(b, as, newSSE2, 8192, !hasSSE2) }

func Benchmark_Hash_8_SSSE3(b *testing.B)  { HBA(b, as, newSSSE3, 8, !hasSSSE3) }
func Benchmark_Hash_1K_SSSE3(b *testing.B) { HBA(b, as, newSSSE3, 1024, !hasSSSE3) }
func Benchmark_Hash_8K_SSSE3(b *testing.B) { HBA(b, as, newSSSE3, 8192, !hasSSSE3) }

func Benchmark_Hash_8_AVX2(b *testing.B)  { HBA(b, as, newAVX2, 8, !hasAVX2) }
func Benchmark_Hash_1K_AVX2(b *testing.B) { HBA(b, as, newAVX2, 1024, !hasAVX2) }
func Benchmark_Hash_8K_AVX2(b *testing.B) { HBA(b, as, newAVX2, 8192, !hasAVX2) }
