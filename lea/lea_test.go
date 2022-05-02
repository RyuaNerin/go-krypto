package lea

import (
	"fmt"
	"io"
	"math/rand"
	"testing"

	"github.com/RyuaNerin/go-krypto/test"
)

const (
	testBlockContextIter = 256
	testBlockBlockIter   = 1024
)

func testECB(t *testing.T, keySize int, encMode bool) {
	const blocks = 8

	key := make([]byte, keySize)

	src1 := make([]byte, BlockSize*blocks)
	src2 := make([]byte, BlockSize*blocks)

	dst1 := make([]byte, BlockSize*blocks)
	dst2 := make([]byte, BlockSize*blocks)

	r := rand.New(rand.NewSource(0))

	var leaKey1, leaKey2 leaContext
	leaKey2.ecb = true

	for keyIter := 0; keyIter < testBlockContextIter; keyIter++ {
		io.ReadFull(r, key)

		initContext(&leaKey1, key)
		initContext(&leaKey2, key)

		io.ReadFull(r, src1)
		copy(src2, src1)

		for blockIter := 0; blockIter < testBlockBlockIter; blockIter++ {
			if encMode {
				for i := 0; i < blocks; i++ {
					leaKey1.Encrypt(dst1[BlockSize*i:], src1[BlockSize*i:])
				}
				leaKey2.Encrypt(dst2, src2)
			} else {
				for i := 0; i < blocks; i++ {
					leaKey1.Decrypt(dst1[BlockSize*i:], src1[BlockSize*i:])
				}
				leaKey2.Decrypt(dst2, src2)
			}

			for i := 0; i < blocks; i++ {
				if dst1[i] != dst2[i] {
					t.Errorf(test.DumpByteArray(fmt.Sprintf("Error KeySize=%d / encMode=%t", keySize, encMode), dst1, dst2))
					return
				}
			}

			copy(src1, dst1)
			copy(src2, dst2)
		}
	}
}

func TestECBEnc(t *testing.T) {
	testECB(t, 16, true)
	testECB(t, 24, true)
	testECB(t, 32, true)
}
func TestECBDec(t *testing.T) {
	testECB(t, 16, false)
	testECB(t, 24, false)
	testECB(t, 32, false)
}

func benchNewCipher(b *testing.B, keySize int) {
	key := make([]byte, keySize)

	r := rand.New(rand.NewSource(0))
	io.ReadFull(r, key)

	var leaCtx leaContext

	b.SetBytes(int64(len(leaCtx.rk) * 4))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		leaSetKeyGo(leaCtx.rk, key)
	}
}

func BenchmarkNewLEA128(b *testing.B) {
	benchNewCipher(b, 16)
}
func BenchmarkNewLEA192(b *testing.B) {
	benchNewCipher(b, 24)
}
func BenchmarkNewLEA256(b *testing.B) {
	benchNewCipher(b, 32)
}

func benchBlock(b *testing.B, isAVX2 bool, keySize int, blocks int, f funcBlock) {
	if isAVX2 && !hasAVX2 {
		b.SkipNow()
		return
	}

	key := make([]byte, keySize)
	src := make([]byte, BlockSize*blocks)
	dst := make([]byte, BlockSize*blocks)

	r := rand.New(rand.NewSource(0))
	io.ReadFull(r, key)

	var leaCtx leaContext
	initContext(&leaCtx, key)

	b.SetBytes(int64(len(dst)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f(leaCtx.round, leaCtx.rk, dst, src)
	}
}

func BenchmarkLEA128Enc1Go(b *testing.B) {
	benchBlock(b, false, 16, 1, leaEnc1Go)
}
func BenchmarkLEA128Dec1Go(b *testing.B) {
	benchBlock(b, false, 16, 1, leaDec1Go)
}

func BenchmarkLEA192Enc1Go(b *testing.B) {
	benchBlock(b, false, 24, 1, leaEnc1Go)
}
func BenchmarkLEA192Dec1Go(b *testing.B) {
	benchBlock(b, false, 24, 1, leaDec1Go)
}

func BenchmarkLEA256Enc1Go(b *testing.B) {
	benchBlock(b, false, 32, 1, leaEnc1Go)
}
func BenchmarkLEA256Dec1Go(b *testing.B) {
	benchBlock(b, false, 32, 1, leaDec1Go)
}
