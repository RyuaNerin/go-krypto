package kipher_test

import (
	"crypto/cipher"
	"strconv"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/tblockmode"
	"github.com/RyuaNerin/go-krypto/kipher"

	"github.com/stretchr/testify/require"
)

var (
	newCFBE = kipher.NewCFBEncrypterWithBlockSize
	newCFBD = kipher.NewCFBDecrypterWithBlockSize
)

func newCFB8E(iv []byte) cipher.Stream   { return newCFBE(newLEABlock(), iv, 1) }
func newCFB8D(iv []byte) cipher.Stream   { return newCFBD(newLEABlock(), iv, 1) }
func newCFB32E(iv []byte) cipher.Stream  { return newCFBE(newLEABlock(), iv, 4) }
func newCFB32D(iv []byte) cipher.Stream  { return newCFBD(newLEABlock(), iv, 4) }
func newCFB64E(iv []byte) cipher.Stream  { return newCFBE(newLEABlock(), iv, 8) }
func newCFB64D(iv []byte) cipher.Stream  { return newCFBD(newLEABlock(), iv, 8) }
func newCFB128E(iv []byte) cipher.Stream { return kipher.NewCFBEncrypter(newLEABlock(), iv) }
func newCFB128D(iv []byte) cipher.Stream { return kipher.NewCFBDecrypter(newLEABlock(), iv) }

func newCipherCFB128E(iv []byte) cipher.Stream { return cipher.NewCFBEncrypter(newLEABlock(), iv) }
func newCipherCFB128D(iv []byte) cipher.Stream { return cipher.NewCFBDecrypter(newLEABlock(), iv) }

var (
	cfbInfo   = tblockmode.StreamInfo{IVSize: blockSize}
	cfbBlocks = []struct {
		Name  string
		Size  int
		fnEnc func(iv []byte) cipher.Stream
		fnDec func(iv []byte) cipher.Stream
	}{
		{Name: "CFB8", Size: 1, fnEnc: newCFB8E, fnDec: newCFB8D},
		{Name: "CFB32", Size: 4, fnEnc: newCFB32E, fnDec: newCFB32D},
		{Name: "CFB64", Size: 8, fnEnc: newCFB64E, fnDec: newCFB64D},
		{Name: "CFB128", Size: 16, fnEnc: newCFB128E, fnDec: newCFB128D},
	}
)

func Test_CFB_New(t *testing.T) {
	for _, block := range cfbBlocks {
		t.Run(block.Name, func(t *testing.T) {
			cfbInfo.TestNew(t, block.fnEnc, block.fnDec, false)
		})
	}
}

func Test_CFB_InvalidSegmentSize(t *testing.T) {
	iv := make([]byte, blockSize)

	for _, size := range []int{-1, 0, blockSize + 1} {
		t.Run("encrypt/"+strconv.Itoa(size), func(t *testing.T) {
			require.Panics(t, func() {
				kipher.NewCFBEncrypterWithBlockSize(newLEABlock(), iv, size)
			})
		})
		t.Run("decrypt/"+strconv.Itoa(size), func(t *testing.T) {
			require.Panics(t, func() {
				kipher.NewCFBDecrypterWithBlockSize(newLEABlock(), iv, size)
			})
		})
	}
}

func Test_CFB_Stream(t *testing.T) {
	for _, block := range cfbBlocks {
		t.Run(block.Name, func(t *testing.T) {
			cfbInfo.TestStream(t, block.fnEnc, block.fnDec, false)
		})
	}
}

func Test_CFB_Compare(t *testing.T) {
	cfbInfo.TestCompare(t, newCFB128E, newCFB128D, newCipherCFB128E, newCipherCFB128D, false)
}

func Fuzz_CFB8_New(f *testing.F)          { cfbInfo.FuzzNew(f, newCFB8E, newCFB8D, false) }
func Fuzz_CFB64_New(f *testing.F)         { cfbInfo.FuzzNew(f, newCFB64E, newCFB64D, false) }
func Fuzz_CFB32_New(f *testing.F)         { cfbInfo.FuzzNew(f, newCFB32E, newCFB32D, false) }
func Fuzz_CFB128_New(f *testing.F)        { cfbInfo.FuzzNew(f, newCFB128E, newCFB128D, false) }
func Fuzz_CFB8_Stream(f *testing.F)       { cfbInfo.FuzzStream(f, newCFB8E, newCFB8D, false) }
func Fuzz_CFB64_Stream(f *testing.F)      { cfbInfo.FuzzStream(f, newCFB64E, newCFB64D, false) }
func Fuzz_CFB32_Stream(f *testing.F)      { cfbInfo.FuzzStream(f, newCFB32E, newCFB32D, false) }
func Fuzz_CFB128_Stream(f *testing.F)     { cfbInfo.FuzzStream(f, newCFB128E, newCFB128D, false) }
func Fuzz_CFB8_ShortBuffer(f *testing.F)  { cfbInfo.FuzzShortBuffer(f, newCFB8E, newCFB8D, false) }
func Fuzz_CFB64_ShortBuffer(f *testing.F) { cfbInfo.FuzzShortBuffer(f, newCFB64E, newCFB64D, false) }

func Fuzz_CFB32_ShortBuffer(f *testing.F) { cfbInfo.FuzzShortBuffer(f, newCFB32E, newCFB32D, false) }

func Fuzz_CFB128_ShortBuffer(f *testing.F) { cfbInfo.FuzzShortBuffer(f, newCFB128E, newCFB128D, false) }

func Benchmark_CFB_New(b *testing.B) {
	for _, block := range cfbBlocks {
		b.Run(block.Name, func(b *testing.B) {
			cfbInfo.BenchNew(b, block.fnEnc, block.fnDec, false)
		})
	}
}

func Benchmark_CFB_Stream(b *testing.B) {
	for _, block := range cfbBlocks {
		b.Run(block.Name, func(b *testing.B) {
			if block.Size == 16 {
				cfbInfo.BenchStreamCompare(b, newCFB128E, newCFB128D, newCipherCFB128E, newCipherCFB128D, false)
			} else {
				cfbInfo.BenchStream(b, block.fnEnc, block.fnDec, false)
			}
		})
	}
}
