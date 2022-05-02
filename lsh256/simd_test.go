//go:build amd64

package lsh256

import (
	"bytes"
	"io"
	"math/rand"
	"sync"
	"testing"

	"github.com/RyuaNerin/go-krypto/test"
)

const (
	testIter      = 512
	testMaxLength = 13
)

func testSimd(t *testing.T, algtype algType, simd simdSet) {
	r := rand.New(rand.NewSource(0))
	rGo := rand.New(rand.NewSource(0))
	rAsm := rand.New(rand.NewSource(0))

	hashBufGo := make([]byte, 0, Size)
	hashBufAsm := make([]byte, 0, Size)

	var ctxGo lsh256ContextGo
	var ctxAsm lsh256ContextAsm

	var w sync.WaitGroup
	for iter := 0; iter < testIter; iter++ {
		initContextGo(&ctxGo, algtype)
		initContextAsm(&ctxAsm, algtype, simd)

		n := r.Int63n(testMaxLength-1) + 1

		w.Add(2)
		go func() {
			defer w.Done()
			io.CopyN(&ctxGo, rGo, n)
		}()
		go func() {
			defer w.Done()
			io.CopyN(&ctxAsm, rAsm, n)
		}()
		w.Wait()

		hashBufGo = ctxGo.Sum(hashBufGo[:0])
		hashBufAsm = ctxAsm.Sum(hashBufAsm[:0])

		if !bytes.Equal(hashBufGo, hashBufAsm) {
			t.Error(test.DumpByteArray("Error", hashBufGo, hashBufAsm))
			return
		}
	}
}

func TestLSH256SSE2Sample(t *testing.T) {
	ctx := newContextAsm(lshType256H256, simdSetSSSE3)

	//ctx := newContextGo(lshType256H256)

	ctx.Write([]byte{1, 2, 3, 4, 5, 6, 7, 8})
	t.Error(ctx.Sum(nil)[0])
}

func TestLSH256SSE2(t *testing.T) {
	testSimd(t, lshType256H256, simdSetSSE2)
	testSimd(t, lshType256H224, simdSetSSE2)
}
func TestLSH256SSSE3(t *testing.T) {
	testSimd(t, lshType256H256, simdSetSSSE3)
	testSimd(t, lshType256H224, simdSetSSSE3)
}
func TestLSH256AVX2(t *testing.T) {
	testSimd(t, lshType256H256, simdSetAVX2)
	testSimd(t, lshType256H224, simdSetAVX2)
}

func BenchmarkLSH256SSE2Reset(b *testing.B) {
	benchReset(b, newContextAsm(lshType256H256, simdSetSSE2))
}
func BenchmarkLSH256SSE2Write(b *testing.B) {
	benchWrite(b, newContextAsm(lshType256H256, simdSetSSE2))
}
func BenchmarkLSH256SSE2Sum(b *testing.B) {
	benchSum(b, newContextAsm(lshType256H256, simdSetSSE2))
}

func BenchmarkLSH256SSSE3Reset(b *testing.B) {
	benchReset(b, newContextAsm(lshType256H256, simdSetSSSE3))
}
func BenchmarkLSH256SSSE3Write(b *testing.B) {
	benchWrite(b, newContextAsm(lshType256H256, simdSetSSSE3))
}
func BenchmarkLSH256SSSE3Sum(b *testing.B) {
	benchSum(b, newContextAsm(lshType256H256, simdSetSSSE3))
}

func BenchmarkLSH256AVX2Reset(b *testing.B) {
	benchReset(b, newContextAsm(lshType256H256, simdSetAVX2))
}
func BenchmarkLSH256AVX2Write(b *testing.B) {
	benchWrite(b, newContextAsm(lshType256H256, simdSetAVX2))
}
func BenchmarkLSH256AVX2Sum(b *testing.B) {
	benchSum(b, newContextAsm(lshType256H256, simdSetAVX2))
}
