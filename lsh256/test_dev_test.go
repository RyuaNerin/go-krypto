//go:build amd64

package lsh256

import (
	"bytes"
	"log"
	"testing"
)

var (
	input = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	a224  = []byte{22, 120, 202, 175, 245, 115, 11, 8, 161, 36, 224, 164, 169, 146, 50, 189, 139, 212, 48, 102, 89, 55, 192, 53, 148, 136, 101, 121, 3, 245, 226, 105}
	a256  = []byte{24, 235, 63, 13, 177, 177, 2, 51, 177, 247, 172, 35, 8, 10, 194, 85, 240, 18, 233, 33, 40, 20, 177, 182, 223, 125, 148, 192, 148, 189, 166, 255}
)

func Test_LSH224_SSE2_DEV(t *testing.T) { testAsmDev(t, lshType256H224, a224, simdSetSSE2) }
func Test_LSH256_SSE2_DEV(t *testing.T) { testAsmDev(t, lshType256H256, a256, simdSetSSE2) }

/**
func Test_LSH224_SSSE3(t *testing.T) { testAsm(t, lshType256H224, simdSetSSSE3) }
func Test_LSH256_SSSE3(t *testing.T) { testAsm(t, lshType256H256, simdSetSSSE3) }
*/

func Test_LSH224_AVX2_DEV(t *testing.T) { testAsmDev(t, lshType256H224, a224, simdSetAVX2) }
func Test_LSH256_AVX2_DEV(t *testing.T) { testAsmDev(t, lshType256H256, a256, simdSetAVX2) }

func testAsmDev(t *testing.T, algType algType, answer []byte, s simdSet) {
	ctx := lsh256ContextAsmData{
		algtype:           uint32(algType),
		remain_databitlen: 0,
	}
	s.init(&ctx)

	dst := make([]byte, Size)

	s.update(&ctx, input, uint32(len(input)*8))
	s.final(&ctx, dst)

	log.Println(dst)
	if !bytes.Equal(answer, dst) {
		t.Fail()
	}
}
