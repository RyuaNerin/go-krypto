package lsh512

import (
	"bytes"
	"testing"
)

type testCase struct {
	Msg []byte
	MD  []byte
}

func Test_LSH512_224_SSE2(t *testing.T) { testAsm(t, testCases224, Size224, simdSetSSE2) }
func Test_LSH512_256_SSE2(t *testing.T) { testAsm(t, testCases256, Size256, simdSetSSE2) }
func Test_LSH512_384_SSE2(t *testing.T) { testAsm(t, testCases384, Size384, simdSetSSE2) }
func Test_LSH512_512_SSE2(t *testing.T) { testAsm(t, testCases512, Size, simdSetSSE2) }

func Test_LSH512_224_SSSE3(t *testing.T) { testAsm(t, testCases224, Size224, simdSetSSSE3) }
func Test_LSH512_256_SSSE3(t *testing.T) { testAsm(t, testCases256, Size256, simdSetSSSE3) }
func Test_LSH512_384_SSSE3(t *testing.T) { testAsm(t, testCases384, Size384, simdSetSSSE3) }
func Test_LSH512_512_SSSE3(t *testing.T) { testAsm(t, testCases512, Size, simdSetSSSE3) }

func Test_LSH512_224_AVX2(t *testing.T) { testAsm(t, testCases224, Size224, simdSetAVX2) }
func Test_LSH512_256_AVX2(t *testing.T) { testAsm(t, testCases256, Size256, simdSetAVX2) }
func Test_LSH512_384_AVX2(t *testing.T) { testAsm(t, testCases384, Size384, simdSetAVX2) }
func Test_LSH512_512_AVX2(t *testing.T) { testAsm(t, testCases512, Size, simdSetAVX2) }

func testAsm(t *testing.T, testCases []testCase, algType int, simd simdSet) {
	h := newContextAsm(algType, simd)

	out := make([]byte, BlockSize)

	for _, tc := range testCases {
		h.Reset()
		h.Write(tc.Msg)
		out = h.Sum(out[:0])

		if !bytes.Equal(out, tc.MD) {
			t.Fail()
		}
	}
}