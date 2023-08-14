//go:build amd64

package lsh512

import (
	"bytes"
	"log"
	"testing"
)

var (
	input = make([]byte, 256)
	a512  = []byte{230, 166, 58, 163, 40, 218, 17, 142, 7, 60, 35, 188, 131, 19, 32, 171, 175, 219, 140, 251, 19, 16, 97, 75, 94, 75, 66, 111, 48, 81, 68, 85}
	a384  = []byte{230, 166, 58, 163, 40, 218, 17, 142, 7, 60, 35, 188, 131, 19, 32, 171, 175, 219, 140, 251, 19, 16, 97, 75, 94, 75, 66, 111, 48, 81, 68, 85}
	a256  = []byte{230, 166, 58, 163, 40, 218, 17, 142, 7, 60, 35, 188, 131, 19, 32, 171, 175, 219, 140, 251, 19, 16, 97, 75, 94, 75, 66, 111, 48, 81, 68, 85}
	a224  = []byte{22, 120, 202, 175, 245, 115, 11, 8, 161, 36, 224, 164, 169, 146, 50, 189, 139, 212, 48, 102, 89, 55, 192, 53, 148, 136, 101, 121, 3, 245, 226, 105}
)

func init() {
	for idx := range input {
		input[idx] = byte(idx)
	}
}

func Test_LSH512_SSE2_DEV(t *testing.T) { testAsmDev(t, Size, a512, simdSetSSE2) }

/**
func Test_LSH512_SSSE3_DEV(t *testing.T) { testAsmDev(t, Size, a512, simdSetSSSE3) }

func Test_LSH512_AVX2_DEV(t *testing.T) { testAsmDev(t, Size, a512, simdSetAVX2) }
*/

func testAsmDev(t *testing.T, algType int, answer []byte, s simdSet) {
	ctx := lsh512ContextAsmData{
		algtype: uint32(algType),
	}
	s.init(&ctx)

	dst := make([]byte, Size)

	s.update(&ctx, input)
	for idx, v := range ctx.cv_l {
		log.Println("cv_l", idx, v)
	}
	for idx, v := range ctx.cv_r {
		log.Println("cv_r", idx, v)
	}
	s.final(&ctx, dst)

	if !bytes.Equal(answer, dst) {
		t.Fail()
	}
}
