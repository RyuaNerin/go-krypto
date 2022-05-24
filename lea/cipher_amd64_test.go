//go:build amd64

package lea

import (
	"fmt"
	"io"
	"math/rand"
	"testing"

	"github.com/RyuaNerin/go-krypto/test"
)

func testAsmBlockWithKeySize(
	t *testing.T,
	keySize int,
	blockCount int,
	fgo func(round int, rk []uint32, dst []byte, src []byte),
	fas func(round int, rk []uint32, dst []byte, src []byte),
) {
	key := make([]byte, keySize)

	srcGo := make([]byte, BlockSize*blockCount)
	srcAs := make([]byte, BlockSize*blockCount)

	dstGo := make([]byte, BlockSize*blockCount)
	dstAs := make([]byte, BlockSize*blockCount)

	r := rand.New(rand.NewSource(0))

	var leaKeyGo leaContext
	var leaKeyAs leaContext

	for keyIter := 0; keyIter < testBlockContextIter; keyIter++ {
		io.ReadFull(r, key)

		initContext(&leaKeyGo, key)
		initContext(&leaKeyAs, key)

		io.ReadFull(r, srcGo)
		copy(srcAs, srcGo)

		for blockIter := 0; blockIter < testBlockBlockIter; blockIter++ {

			fgo(leaKeyGo.round, leaKeyGo.rk[:], dstGo, srcGo)
			fas(leaKeyAs.round, leaKeyAs.rk[:], dstAs, srcAs)

			for i := 0; i < BlockSize*blockCount; i++ {
				if dstGo[i] != dstAs[i] {
					t.Errorf(test.DumpByteArray(fmt.Sprintf("Error KeySize=%d", keySize), dstGo, dstAs))
					return
				}
			}

			copy(srcGo, dstGo)
			copy(srcAs, dstAs)
		}
	}
}
func testAsmBlock(
	t *testing.T,
	blockCount int,
	fgo func(round int, rk []uint32, dst []byte, src []byte),
	fasm func(round int, rk []uint32, dst []byte, src []byte),
) {
	testAsmBlockWithKeySize(t, 16, blockCount, fgo, fasm)
	testAsmBlockWithKeySize(t, 24, blockCount, fgo, fasm)
	testAsmBlockWithKeySize(t, 32, blockCount, fgo, fasm)
}

func Test_Asm_Encrypt_4Blocks_SSE2(t *testing.T) {
	testAsmBlock(t, 4, leaEnc4Go, leaEnc4SSE2)
}
func Test_Asm_Decrypt_4Blocks_SSE2(t *testing.T) {
	testAsmBlock(t, 4, leaDec4Go, leaDec4SSE2)
}

func Test_Asm_Encrypt_8Blocks_AVX2(t *testing.T) {
	if !hasAVX2 {
		t.SkipNow()
		return
	}

	testAsmBlock(t, 8, leaEnc8Go, leaEnc8AVX2)
}
func Test_Asm_Decrypt_8Blocks_AVX2(t *testing.T) {
	if !hasAVX2 {
		t.SkipNow()
		return
	}

	testAsmBlock(t, 8, leaDec8Go, leaDec8AVX2)
}

func Benchmark_LEA128_Encrypt_4Blocks_SSE2(b *testing.B) {
	benchBlock(b, false, 16, 4, leaEnc4SSE2)
}
func Benchmark_LEA128_Decrypt_4Blocks_SSE2(b *testing.B) {
	benchBlock(b, false, 16, 4, leaDec4SSE2)
}

func Benchmark_LEA192_Encrypt_4Blocks_SSE2(b *testing.B) {
	benchBlock(b, false, 24, 4, leaEnc4SSE2)
}
func Benchmark_LEA192_Decrypt_4Blocks_SSE2(b *testing.B) {
	benchBlock(b, false, 24, 4, leaDec4SSE2)
}

func Benchmark_LEA256_Encrypt_4Blocks_SSE2(b *testing.B) {
	benchBlock(b, false, 32, 4, leaEnc4SSE2)
}
func Benchmark_LEA256_Decrypt_4Blocks_SSE2(b *testing.B) {
	benchBlock(b, false, 32, 4, leaDec4SSE2)
}

func Benchmark_LEA128_Encrypt_8Blocks_AVX2(b *testing.B) {
	benchBlock(b, true, 16, 8, leaEnc8AVX2)
}
func Benchmark_LEA128_Decrypt_8Blocks_AVX2(b *testing.B) {
	benchBlock(b, true, 16, 8, leaDec8AVX2)
}

func Benchmark_LEA192_Encrypt_8Blocks_AVX2(b *testing.B) {
	benchBlock(b, true, 24, 8, leaEnc8AVX2)
}
func Benchmark_LEA192_Decrypt_8Blocks_AVX2(b *testing.B) {
	benchBlock(b, true, 24, 8, leaDec8AVX2)
}

func Benchmark_LEA256_Encrypt_8Blocks_AVX2(b *testing.B) {
	benchBlock(b, true, 32, 8, leaEnc8AVX2)
}
func Benchmark_LEA256_Decrypt_8Blocks_AVX2(b *testing.B) {
	benchBlock(b, true, 32, 8, leaDec8AVX2)
}
