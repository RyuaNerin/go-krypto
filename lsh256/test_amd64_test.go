package lsh256

import (
	"bytes"
	"testing"
)

type testCase struct {
	M  []byte
	MD []byte
}

func Test_LSH224_SSE2(t *testing.T) { testAsm(t, testCases224, lshType256H224, simdSetSSE2) }
func Test_LSH256_SSE2(t *testing.T) { testAsm(t, testCases256, lshType256H256, simdSetSSE2) }

func Test_LSH224_SSSE3(t *testing.T) { testAsm(t, testCases224, lshType256H224, simdSetSSSE3) }
func Test_LSH256_SSSE3(t *testing.T) { testAsm(t, testCases256, lshType256H256, simdSetSSSE3) }

func Test_LSH224_AVX2(t *testing.T) { testAsm(t, testCases224, lshType256H224, simdSetAVX2) }
func Test_LSH256_AVX2(t *testing.T) { testAsm(t, testCases256, lshType256H256, simdSetAVX2) }

func testAsm(t *testing.T, testCases []testCase, algType algType, simd simdSet) {
	h := newContextAsm(algType, simd)

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
