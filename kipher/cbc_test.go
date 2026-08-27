package kipher_test

import (
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/tblockmode"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var cbcInfo = tblockmode.BlockModeInfo{
	BlockSize: blockSize,
	IVSize:    blockSize,
}

func newCBCE(iv []byte) cipher.BlockMode       { return kipher.NewCBCEncrypter(newLEABlock(), iv) }
func newCBCD(iv []byte) cipher.BlockMode       { return kipher.NewCBCDecrypter(newLEABlock(), iv) }
func newCipherCBCE(iv []byte) cipher.BlockMode { return cipher.NewCBCEncrypter(newLEABlock(), iv) }
func newCipherCBCD(iv []byte) cipher.BlockMode { return cipher.NewCBCDecrypter(newLEABlock(), iv) }

func Test_CBC_New(t *testing.T)       { cbcInfo.TestNew(t, newCBCE, newCBCD, false) }
func Test_CBC_BlockMode(t *testing.T) { cbcInfo.TestBlockMode(t, newCBCE, newCBCD, false) }
func Test_CBC_Compare(t *testing.T) {
	cbcInfo.TestCompare(t, newCBCE, newCBCD, newCipherCBCE, newCipherCBCD, false)
}

func Fuzz_CBC_New(f *testing.F)         { cbcInfo.FuzzNew(f, newCBCE, newCBCD, false) }
func Fuzz_CBC_BlockMode(f *testing.F)   { cbcInfo.FuzzBlockMode(f, newCBCE, newCBCD, false) }
func Fuzz_CBC_ShortBuffer(f *testing.F) { cbcInfo.FuzzBlockModeShortBuffer(f, newCBCE, newCBCD, false) }

func Benchmark_CBC_New(b *testing.B) { cbcInfo.BenchNew(b, newCBCE, newCBCD, false) }
func Benchmark_CBC_BlockMode(b *testing.B) {
	cbcInfo.BenchCompare(b, newCBCE, newCBCD, newCipherCBCE, newCipherCBCD, false)
}
