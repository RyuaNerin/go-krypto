//go:build amd64

package lsh256

import (
	"bytes"
	"testing"
)

const (
	testBlocks = 16 * 1024
)

func Test_LSH224_SSE2(t *testing.T) { testAsm(t, lshType256H224, simdSetSSE2) }
func Test_LSH256_SSE2(t *testing.T) { testAsm(t, lshType256H256, simdSetSSE2) }

func Test_LSH224_AVX2(t *testing.T) { testAsm(t, lshType256H224, simdSetAVX2) }
func Test_LSH256_AVX2(t *testing.T) { testAsm(t, lshType256H256, simdSetAVX2) }

func testAsm(t *testing.T, algType algType, simd simdSet) {
	var hGo lsh256ContextGo
	initContextGo(&hGo, algType)

	var hAsm lsh256ContextAsm
	initContextAsm(&hAsm, algType, simd)

	src := make([]byte, BlockSize)
	for idx := range src {
		src[idx] = byte(idx)
	}
	dstGo := make([]byte, BlockSize)
	dstAsm := make([]byte, BlockSize)

	for i := 0; i < testBlocks; i++ {
		hGo.Reset()
		hAsm.Reset()

		// check
		for i := 0; i < len(hAsm.data_cv_l); i++ {
			if hGo.cv[i] != hAsm.data_cv_l[i] {
				t.Fail()
			}
		}
		for i := 0; i < len(hAsm.data_cv_r); i++ {
			if hGo.cv[len(hAsm.data_cv_l)+i] != hAsm.data_cv_r[i] {
				t.Fail()
			}
		}

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
