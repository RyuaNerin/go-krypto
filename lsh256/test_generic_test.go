package lsh256

import (
	"bytes"
	"testing"
)

type testCase struct {
	M  []byte
	MD []byte
}

func Test_LSH224_Go(t *testing.T) { testGo(t, testCases224, lshType256H224) }
func Test_LSH256_Go(t *testing.T) { testGo(t, testCases256, lshType256H256) }

func testGo(t *testing.T, testCases []testCase, algType algType) {
	h := newContextGo(algType)

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
