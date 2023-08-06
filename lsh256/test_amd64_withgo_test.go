//go:build amd64

package lsh256

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"testing"
)

const (
	testBlocks = 16 * 1024
)

func Test_LSH224_SSE2_CONTINUS(t *testing.T) { testAsmContinus(t, lshType256H224, simdSetSSE2, true) }
func Test_LSH256_SSE2_CONTINUS(t *testing.T) { testAsmContinus(t, lshType256H256, simdSetSSE2, true) }

func Test_LSH224_SSSE3_CONTINUS(t *testing.T) {
	testAsmContinus(t, lshType256H224, simdSetSSSE3, hasSSSE3)
}
func Test_LSH256_SSSE3_CONTINUS(t *testing.T) {
	testAsmContinus(t, lshType256H256, simdSetSSSE3, hasSSSE3)
}

func Test_LSH224_AVX2_CONTINUS(t *testing.T) {
	testAsmContinus(t, lshType256H224, simdSetAVX2, hasAVX2)
}
func Test_LSH256_AVX2_CONTINUS(t *testing.T) {
	testAsmContinus(t, lshType256H256, simdSetAVX2, hasAVX2)
}

func testAsmContinus(t *testing.T, algType algType, simd simdSet, do bool) {
	if !do {
		t.Skip()
		return
	}
	rnd := bufio.NewReaderSize(rand.Reader, 1<<15)

	var hGo lsh256ContextGo
	initContextGo(&hGo, algType)

	var hAsm lsh256ContextAsm
	initContextAsm(&hAsm, algType, simd)

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
}
