//go:build amd64

package lea

import (
	"crypto/cipher"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"testing"
)

const (
	testCtrRefillIvIter   = 256
	testCtrRefillStepIter = 1024

	testCipherKeySize   = 16
	testCipherKeyIter   = 256
	testCipherStepIter  = 1024
	testCipherMaxBlocks = 32 // Blocks
)

type leaContextGo struct {
	raw leaContext
}

func (leaCtx *leaContextGo) BlockSize() int {
	return leaCtx.raw.BlockSize()
}

func (leaCtx *leaContextGo) Encrypt(dst, src []byte) {
	leaCtx.raw.Encrypt(dst, src)
}

func (leaCtx *leaContextGo) Decrypt(dst, src []byte) {
	leaCtx.raw.Decrypt(dst, src)
}

func testCipherStream(
	t *testing.T,
	newBlockMode func(cipher.Block, []byte) cipher.Stream,
	typeCheck func(v interface{}) bool,
) {
	key := make([]byte, testCipherKeySize)

	iv := make([]byte, BlockSize)

	srcGo := make([]byte, testCipherMaxBlocks)
	srcAsm := make([]byte, testCipherMaxBlocks)

	dstGo := make([]byte, testCipherMaxBlocks)
	dstAsm := make([]byte, testCipherMaxBlocks)

	r := rand.New(rand.NewSource(0))

	var leaCtxGo leaContextGo
	var leaCtxAsm leaContext

	for keyIter := 0; keyIter < testCipherKeyIter; keyIter++ {
		io.ReadFull(r, key)
		initContext(&leaCtxGo.raw, key)
		initContext(&leaCtxAsm, key)

		io.ReadFull(r, iv)
		cipherGo := newBlockMode(&leaCtxGo, iv)
		cipherAsm := newBlockMode(&leaCtxAsm, iv)

		if !typeCheck(cipherAsm) {
			t.Errorf("invaild type : %s", reflect.TypeOf(cipherAsm))
			return
		}

		for blockIter := 0; blockIter < testCipherStepIter; blockIter++ {
			srcSize := 1 + r.Intn(testCipherMaxBlocks-1)

			io.ReadFull(r, srcGo[:srcSize])
			copy(srcAsm[:srcSize], srcGo[:srcSize])

			cipherGo.XORKeyStream(dstGo, srcGo[:srcSize])
			cipherAsm.XORKeyStream(dstAsm, srcAsm[:srcSize])

			for i := 0; i < srcSize; i++ {
				if dstGo[i] != dstAsm[i] {
					t.Error(dumpByteArray(fmt.Sprintf("Error / keyIter = %d / blockIter = %d / srcSize = %d", keyIter, blockIter, srcSize), dstGo[:srcSize], dstAsm[:srcSize]))
					return
				}
			}
		}
	}
}

func TestCtr(t *testing.T) {
	testCipherStream(
		t,
		cipher.NewCTR,
		func(v interface{}) bool {
			_, ok := v.(*leaCtrContext)
			return ok
		},
	)
}

func testCipherBlockMode(
	t *testing.T,
	newStream func(cipher.Block, []byte) cipher.BlockMode,
	typeCheck func(v interface{}) bool,
) {
	key := make([]byte, testCipherKeySize)

	iv := make([]byte, BlockSize)

	srcGo := make([]byte, BlockSize*testCipherMaxBlocks)
	srcAsm := make([]byte, BlockSize*testCipherMaxBlocks)

	dstGo := make([]byte, BlockSize*testCipherMaxBlocks)
	dstAsm := make([]byte, BlockSize*testCipherMaxBlocks)

	r := rand.New(rand.NewSource(0))

	var leaCtxGo leaContextGo
	var leaCtxAsm leaContext

	for keyIter := 0; keyIter < testCipherKeyIter; keyIter++ {
		io.ReadFull(r, key)
		initContext(&leaCtxGo.raw, key)
		initContext(&leaCtxAsm, key)

		io.ReadFull(r, iv)
		cipherGo := newStream(&leaCtxGo, iv)
		cipherAsm := newStream(&leaCtxAsm, iv)

		if !typeCheck(cipherAsm) {
			t.Errorf("invaild type : %s", reflect.TypeOf(cipherAsm))
			return
		}

		for blockIter := 0; blockIter < testCipherStepIter; blockIter++ {
			sz := BlockSize * (1 + r.Intn(testCipherMaxBlocks-1))

			io.ReadFull(r, srcGo[:sz])
			copy(srcAsm, srcGo[:sz])

			cipherGo.CryptBlocks(dstGo, srcGo[:sz])
			cipherAsm.CryptBlocks(dstAsm, srcAsm[:sz])

			for i := 0; i < sz; i++ {
				if dstGo[i] != dstAsm[i] {
					t.Error(dumpByteArray(fmt.Sprintf("Error, Blocks=%d", sz/8), dstGo, dstAsm))
					return
				}
			}
		}
	}
}

func TestCBCDecrypter(t *testing.T) {
	testCipherBlockMode(
		t,
		cipher.NewCBCDecrypter,
		func(v interface{}) bool {
			_, ok := v.(*leaCbcContext)
			return ok
		},
	)
}

func benchCipherStream(b *testing.B, useAsm bool, f func(cipher.Block, []byte) cipher.Stream) {
	key := make([]byte, testCipherKeySize)
	iv := make([]byte, BlockSize)

	src := make([]byte, BlockSize)
	dst := make([]byte, BlockSize)

	r := rand.New(rand.NewSource(0))
	io.ReadFull(r, key)
	io.ReadFull(r, iv)
	io.ReadFull(r, src)

	var leaCtx leaContext
	initContext(&leaCtx, key)

	var ctr cipher.Stream
	if useAsm {
		ctr = f(&leaCtx, iv)
	} else {
		var leaCtxGo leaContextGo
		leaCtxGo.raw = leaCtx

		ctr = f(&leaCtxGo, iv)
	}

	b.SetBytes(int64(len(dst)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctr.XORKeyStream(dst, src)
	}
}

func BenchmarkCTRGo(b *testing.B) {
	benchCipherStream(b, false, cipher.NewCTR)
}

func BenchmarkCTRAsm(b *testing.B) {
	benchCipherStream(b, true, cipher.NewCTR)
}

func benchCipherBlockMode(b *testing.B, blocks int, useAsm bool, f func(cipher.Block, []byte) cipher.BlockMode) {
	key := make([]byte, testCipherKeySize)
	iv := make([]byte, BlockSize)

	src := make([]byte, BlockSize*blocks)
	dst := make([]byte, BlockSize*blocks)

	r := rand.New(rand.NewSource(0))
	io.ReadFull(r, key)
	io.ReadFull(r, iv)
	io.ReadFull(r, src)

	var leaCtx leaContext
	initContext(&leaCtx, key)

	var blockMode cipher.BlockMode
	if useAsm {
		blockMode = f(&leaCtx, iv)
	} else {
		var leaCtxGo leaContextGo
		leaCtxGo.raw = leaCtx

		blockMode = f(&leaCtxGo, iv)
	}

	b.SetBytes(int64(len(dst)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blockMode.CryptBlocks(dst, src)
	}
}

func BenchmarkCBCDecrypterGo(b *testing.B) {
	benchCipherBlockMode(b, 1, false, cipher.NewCBCDecrypter)
}

func BenchmarkCBCDecrypterAsm1Block(b *testing.B) {
	benchCipherBlockMode(b, 1, true, cipher.NewCBCDecrypter)
}
func BenchmarkCBCDecrypterAsm4Blocks(b *testing.B) {
	benchCipherBlockMode(b, 4, true, cipher.NewCBCDecrypter)
}
func BenchmarkCBCDecrypterAsm8Blocks(b *testing.B) {
	benchCipherBlockMode(b, 8, true, cipher.NewCBCDecrypter)
}
