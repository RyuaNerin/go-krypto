//go:build amd64

package lsh256

import (
	"bytes"
	"testing"
)

var (
	input = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	a224  = []byte{22, 120, 202, 175, 245, 115, 11, 8, 161, 36, 224, 164, 169, 146, 50, 189, 139, 212, 48, 102, 89, 55, 192, 53, 148, 136, 101, 121, 3, 245, 226, 105}
	a256  = []byte{230, 166, 58, 163, 40, 218, 17, 142, 7, 60, 35, 188, 131, 19, 32, 171, 175, 219, 140, 251, 19, 16, 97, 75, 94, 75, 66, 111, 48, 81, 68, 85}
)

func Test_LSH224_SSE2_DEV(t *testing.T) { testAsmDev(t, Size224, a224, SimdSetSSE2) }
func Test_LSH256_SSE2_DEV(t *testing.T) { testAsmDev(t, Size, a256, SimdSetSSE2) }

func Test_LSH224_SSSE3_DEV(t *testing.T) { testAsmDev(t, Size224, a224, SimdSetSSSE3) }
func Test_LSH256_SSSE3_DEV(t *testing.T) { testAsmDev(t, Size, a256, SimdSetSSSE3) }

func Test_LSH224_AVX2_DEV(t *testing.T) { testAsmDev(t, Size224, a224, SimdSetAVX2) }
func Test_LSH256_AVX2_DEV(t *testing.T) { testAsmDev(t, Size, a256, SimdSetAVX2) }

func testAsmDev(t *testing.T, algType int, answer []byte, s simdSet) {
	ctx := lsh256ContextAsmData{
		algtype:           uint32(algType),
		remain_databitlen: 0,
	}
	s.init(&ctx)

	dst := make([]byte, Size)

	s.update(&ctx, input)
	s.final(&ctx, dst)

	if !bytes.Equal(answer, dst) {
		t.Fail()
	}
}
