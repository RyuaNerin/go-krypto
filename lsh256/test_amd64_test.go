//go:build amd64 && gc && !purego

package lsh256

import (
	"hash"
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

func newSSE2(size int) hash.Hash  { return newContextAsm(size, simdSetSSE2) }
func newSSSE3(size int) hash.Hash { return newContextAsm(size, simdSetSSSE3) }
func newAVX2(size int) hash.Hash  { return newContextAsm(size, simdSetAVX2) }

func Test_ShortWrite_SSE2(t *testing.T)  { HTSWA(t, as, newSSE2, false) }
func Test_ShortWrite_SSSE3(t *testing.T) { HTSWA(t, as, newSSSE3, !hasSSSE3) }
func Test_ShortWrite_AVX2(t *testing.T)  { HTSWA(t, as, newAVX2, !hasAVX2) }

func Test_WITH_GO_SSE2(t *testing.T)  { HTSA(t, as, newContextGo, newSSE2, false) }
func Test_WITH_GO_SSSE3(t *testing.T) { HTSA(t, as, newContextGo, newSSE2, !hasSSSE3) }
func Test_WITH_GO_AVX2(t *testing.T)  { HTSA(t, as, newContextGo, newSSE2, !hasAVX2) }

func Test_LSH224_SSE2(t *testing.T) { HT(t, newSSE2(Size224), testCases224, false) }
func Test_LSH256_SSE2(t *testing.T) { HT(t, newSSE2(Size), testCases256, false) }

func Test_LSH224_SSSE3(t *testing.T) { HT(t, newSSSE3(Size224), testCases224, !hasSSSE3) }
func Test_LSH256_SSSE3(t *testing.T) { HT(t, newSSSE3(Size), testCases256, !hasSSSE3) }

func Test_LSH224_AVX2(t *testing.T) { HT(t, newAVX2(Size224), testCases224, !hasAVX2) }
func Test_LSH256_AVX2(t *testing.T) { HT(t, newAVX2(Size), testCases256, !hasAVX2) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Hash_8_SSE2(b *testing.B)  { HBA(b, as, newSSE2, 8, false) }
func Benchmark_Hash_1K_SSE2(b *testing.B) { HBA(b, as, newSSE2, 1024, false) }
func Benchmark_Hash_8K_SSE2(b *testing.B) { HBA(b, as, newSSE2, 8192, false) }

func Benchmark_Hash_8_SSSE3(b *testing.B)  { HBA(b, as, newSSSE3, 8, !hasSSSE3) }
func Benchmark_Hash_1K_SSSE3(b *testing.B) { HBA(b, as, newSSSE3, 1024, !hasSSSE3) }
func Benchmark_Hash_8K_SSSE3(b *testing.B) { HBA(b, as, newSSSE3, 8192, !hasSSSE3) }

func Benchmark_Hash_8_AVX2(b *testing.B)  { HBA(b, as, newAVX2, 8, !hasAVX2) }
func Benchmark_Hash_1K_AVX2(b *testing.B) { HBA(b, as, newAVX2, 1024, !hasAVX2) }
func Benchmark_Hash_8K_AVX2(b *testing.B) { HBA(b, as, newAVX2, 8192, !hasAVX2) }
