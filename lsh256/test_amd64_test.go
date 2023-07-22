//go:build amd64

package lsh256

import (
	"log"
	"testing"
)

const (
	BlockSize = 128
)

type algType int

const (
	lshType256H256 algType = 32 // 256
	lshType256H224         = 28 // 224
)

type simdSet struct {
	init   func(ctx *lsh256ContextAsmData)
	update func(ctx *lsh256ContextAsmData, data []byte, databitlen uint32)
	final  func(ctx *lsh256ContextAsmData, hashval []byte)
}

var (
	simdSetSSE2 = simdSet{
		init:   lsh256InitSSE2,
		update: lsh256UpdateSSE2,
		final:  lsh256FinalSSE2,
	}
	simdSetAVX2 = simdSet{
		init:   lsh256InitAVX2,
		update: lsh256UpdateAVX2,
		final:  lsh256FinalAVX2,
	}
)

func Test_LSH224_SSE2(t *testing.T) { testAsm(t, lshType256H224, simdSetSSE2) }
func Test_LSH256_SSE2(t *testing.T) { testAsm(t, lshType256H256, simdSetSSE2) }

func Test_LSH224_AVX2(t *testing.T) { testAsm(t, lshType256H224, simdSetAVX2) }
func Test_LSH256_AVX2(t *testing.T) { testAsm(t, lshType256H256, simdSetAVX2) }

func testAsm(t *testing.T, algType algType, s simdSet) {
	ctx := lsh256ContextAsmData{
		algtype:           uint32(algType),
		remain_databitlen: 0,
		cv_l:              make([]uint32, 32/4),
		cv_r:              make([]uint32, 32/4),
		last_block:        make([]byte, 128),
	}
	s.init(&ctx)

	src := make([]byte, BlockSize)
	for idx := range src {
		src[idx] = byte(idx)
	}
	dst := make([]byte, BlockSize)

	s.update(&ctx, src, uint32(len(src)*8))
	s.final(&ctx, dst)

	log.Println(dst)
}
