//go:build amd64 && gc && !purego

package lsh256

import (
	"bytes"
	"testing"
)

const (
	testBlocks = 16 * 1024
)

func Test_SSE2_CONTINUS(t *testing.T)  { testAsmContinus(t, Size, simdSetSSE2, true) }
func Test_SSSE3_CONTINUS(t *testing.T) { testAsmContinus(t, Size224, simdSetSSSE3, hasSSSE3) }
func Test_AVX2_CONTINUS(t *testing.T)  { testAsmContinus(t, Size224, simdSetAVX2, hasAVX2) }

func testAsmContinus(t *testing.T, size int, simd simdSet, do bool) {
	if !do {
		t.Skip()
		return
	}

	testSize(
		t,
		func(t *testing.T, size int) {

			var hGo lsh256ContextGo
			initContextGo(&hGo, size)

			var hAsm lsh256ContextAsm
			initContextAsm(&hAsm, size, simd)

			src := make([]byte, BlockSize)
			rnd.Read(src)

			dstGo := make([]byte, BlockSize)
			dstAsm := make([]byte, BlockSize)

			for i := 0; i < testBlocks; i++ {
				hGo.Reset()
				hAsm.Reset()

				hGo.Write(src)
				hAsm.Write(src)

				dstGo = hGo.Sum(dstGo[:0])
				dstAsm = hAsm.Sum(dstAsm[:0])

				if !bytes.Equal(dstGo, dstAsm) {
					t.Fail()
				}

				copy(src, dstGo)
			}
		},
	)
}
