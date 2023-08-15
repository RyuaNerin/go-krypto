package lsh256

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"testing"
)

type testCase struct {
	M  []byte
	MD []byte
}

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

func testGo(t *testing.T, testCases []testCase, size int) {
	h := newContextGo(size)

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

func testSize(t *testing.T, f func(t *testing.T, size int)) {
	tests := []struct {
		name string
		size int
	}{
		{"256", Size},
		{"224", Size224},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			f(t, test.size)
		})
	}
}
