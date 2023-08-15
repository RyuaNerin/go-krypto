//go:build amd64

package lsh256

import (
	"bytes"
	"testing"
)

func Test_LSH224_SSE2(t *testing.T) { testAsm(t, testCases224, Size224, simdSetSSE2, true) }
func Test_LSH256_SSE2(t *testing.T) { testAsm(t, testCases256, Size, simdSetSSE2, true) }

func Test_LSH224_SSSE3(t *testing.T) { testAsm(t, testCases224, Size224, simdSetSSSE3, hasSSSE3) }
func Test_LSH256_SSSE3(t *testing.T) { testAsm(t, testCases256, Size, simdSetSSSE3, hasSSSE3) }

func Test_LSH224_AVX2(t *testing.T) { testAsm(t, testCases224, Size224, simdSetAVX2, hasAVX2) }
func Test_LSH256_AVX2(t *testing.T) { testAsm(t, testCases256, Size, simdSetAVX2, hasAVX2) }

func testAsm(t *testing.T, testCases []testCase, size int, simd simdSet, nonskip bool) {
	if !nonskip {
		t.Skip()
		return
	}

	h := newContextAsm(size, simd)

	out := make([]byte, BlockSize)

	for _, tc := range testCases {
		h.Reset()
		h.Write(tc.M)
		out = h.Sum(out[:0])

		if !bytes.Equal(out, tc.MD) {
			t.Fail()
		}
	}
}
