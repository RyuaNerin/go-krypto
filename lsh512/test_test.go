package lsh512

import (
	"bytes"
	"testing"
)

type testCase struct {
	Msg []byte
	MD  []byte
}

func testGo(t *testing.T, testCases []testCase, size int) {
	h := newHash(size)

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
