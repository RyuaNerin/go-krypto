//go:build amd64 && gc && !purego

package lsh512

import (
	"testing"
)

func Test_ShortWrite_SSE2(t *testing.T)  { testShortWrite(t, simdSetSSE2, true) }
func Test_ShortWrite_SSSE3(t *testing.T) { testShortWrite(t, simdSetSSSE3, hasSSSE3) }
func Test_ShortWrite_AVX2(t *testing.T)  { testShortWrite(t, simdSetAVX2, hasAVX2) }

func testShortWrite(t *testing.T, simd simdSet, do bool) {
	if !do {
		t.Skip()
		return
	}

	testSize(
		t,
		func(t *testing.T, size int) {
			h := newContextAsm(size, simd)

			buf := make([]byte, shortWriteSize)
			for i := 1; i < shortWriteSize; i++ {
				rnd.Read(buf[:i])
				n, err := h.Write(buf[:i])
				if err != nil {
					t.Error(err)
				}
				if n != i {
					t.Fail()
				}

				rnd.Read(buf[:1])
				switch buf[0] % 5 {
				case 0:
					h.Reset()
				case 1:
					h.Sum(buf[:0])
				}
			}
		},
	)
}
