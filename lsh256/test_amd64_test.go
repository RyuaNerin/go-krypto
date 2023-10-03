//go:build amd64 && gc && !purego

package lsh256

import (
	"hash"
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

func Test_ShortWrite_SSE2(t *testing.T)  { testShortWriteAsm(t, simdSetSSE2, true) }
func Test_ShortWrite_SSSE3(t *testing.T) { testShortWriteAsm(t, simdSetSSSE3, hasSSSE3) }
func Test_ShortWrite_AVX2(t *testing.T)  { testShortWriteAsm(t, simdSetAVX2, hasAVX2) }

func testShortWriteAsm(t *testing.T, simd simdSet, do bool) {
	if !do {
		t.Skip()
		return
	}

	TA(t, as, func(t *testing.T, size int) {
		h := newContextGo(size)
		TestShortWrite(t, h)
	})
}

func Test_WITH_GO_SSE2(t *testing.T)  { withGo(t, simdSetSSE2, true) }
func Test_WITH_GO_SSSE3(t *testing.T) { withGo(t, simdSetSSSE3, hasSSSE3) }
func Test_WITH_GO_AVX2(t *testing.T)  { withGo(t, simdSetAVX2, hasAVX2) }

func withGo(t *testing.T, simd simdSet, do bool) {
	if !do {
		t.Skip()
		return
	}

	TA(
		t,
		as,
		func(t *testing.T, size int) {
			hGo := newContextGo(size)
			hAsm := newContextAsm(size, simd)

			HTTC(
				t,
				size*8,
				func(dst, p []byte) []byte {
					hGo.Reset()
					hGo.Write(p)
					return hGo.Sum(dst)
				},
				func(dst, p []byte) []byte {
					hAsm.Reset()
					hAsm.Write(p)
					return hAsm.Sum(dst)
				},
			)
		},
	)
}

func Test_LSH224_SSE2(t *testing.T) { testAsm(t, testCases224, Size224, simdSetSSE2, true) }
func Test_LSH256_SSE2(t *testing.T) { testAsm(t, testCases256, Size, simdSetSSE2, true) }

func Test_LSH224_SSSE3(t *testing.T) { testAsm(t, testCases224, Size224, simdSetSSSE3, hasSSSE3) }
func Test_LSH256_SSSE3(t *testing.T) { testAsm(t, testCases256, Size, simdSetSSSE3, hasSSSE3) }

func Test_LSH224_AVX2(t *testing.T) { testAsm(t, testCases224, Size224, simdSetAVX2, hasAVX2) }
func Test_LSH256_AVX2(t *testing.T) { testAsm(t, testCases256, Size, simdSetAVX2, hasAVX2) }

func testAsm(t *testing.T, testCases []HashTestCase, algType int, simd simdSet, do bool) {
	if !do {
		t.Skip()
		return
	}

	HT(t, newContextAsm(algType, simd), testCases)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func newSSE2(size int) hash.Hash  { return newContextAsm(size, simdSetSSE2) }
func newSSSE3(size int) hash.Hash { return newContextAsm(size, simdSetSSSE3) }
func newAVX2(size int) hash.Hash  { return newContextAsm(size, simdSetAVX2) }

func Benchmark_Hash_8_SSE2(b *testing.B)  { bh(b, newSSE2, 8, true) }
func Benchmark_Hash_1K_SSE2(b *testing.B) { bh(b, newSSE2, 1024, true) }
func Benchmark_Hash_8K_SSE2(b *testing.B) { bh(b, newSSE2, 8192, true) }

func Benchmark_Hash_8_SSSE3(b *testing.B)  { bh(b, newSSSE3, 8, hasSSSE3) }
func Benchmark_Hash_1K_SSSE3(b *testing.B) { bh(b, newSSSE3, 1024, hasSSSE3) }
func Benchmark_Hash_8K_SSSE3(b *testing.B) { bh(b, newSSSE3, 8192, hasSSSE3) }

func Benchmark_Hash_8_AVX2(b *testing.B)  { bh(b, newAVX2, 8, hasAVX2) }
func Benchmark_Hash_1K_AVX2(b *testing.B) { bh(b, newAVX2, 1024, hasAVX2) }
func Benchmark_Hash_8K_AVX2(b *testing.B) { bh(b, newAVX2, 8192, hasAVX2) }
