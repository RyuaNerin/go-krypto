package lea

import (
	"log"
	"testing"
)

const (
	BlockSize = 16
)

var (
	round = 24
	rk    = make([]uint32, 192)
	src   = make([]byte, BlockSize*8)
)

func init() {
	for idx := range src {
		src[idx] = byte(idx)
	}
}

func Test_leaEnc4SSE2(t *testing.T) {
	dst := make([]byte, BlockSize*4)
	leaEnc4SSE2(round, rk, dst, src)
	log.Println(dst)
}

func Test_leaEnc8AVX2(t *testing.T) {
	dst := make([]byte, BlockSize*8)
	leaEnc8AVX2(round, rk, dst, src)
	log.Println(dst)
}

func Test_leaDec4SSE2(t *testing.T) {
	dst := make([]byte, BlockSize*4)
	leaDec4SSE2(round, rk, dst, src)
	log.Println(dst)
}

func Test_leaDec8AVX2(t *testing.T) {
	dst := make([]byte, BlockSize*8)
	leaDec8AVX2(round, rk, dst, src)
	log.Println(dst)
}
