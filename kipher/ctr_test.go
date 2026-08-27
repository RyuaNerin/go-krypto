//nolint:dupl
package kipher_test

import (
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/tblockmode"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var ctrInfo = tblockmode.StreamInfo{IVSize: blockSize}

func newCTR(iv []byte) cipher.Stream       { return kipher.NewCTR(newLEABlock(), iv) }
func newCipherCTR(iv []byte) cipher.Stream { return cipher.NewCTR(newLEABlock(), iv) }

func Test_CTR_New(t *testing.T)    { ctrInfo.TestNew(t, newCTR, newCTR, false) }
func Test_CTR_Stream(t *testing.T) { ctrInfo.TestStream(t, newCTR, newCTR, false) }
func Test_CTR_Compare(t *testing.T) {
	ctrInfo.TestCompare(t, newCTR, newCTR, newCipherCTR, newCipherCTR, false)
}

func Fuzz_CTR_New(f *testing.F)         { ctrInfo.FuzzNew(f, newCTR, newCTR, false) }
func Fuzz_CTR_Stream(f *testing.F)      { ctrInfo.FuzzStream(f, newCTR, newCTR, false) }
func Fuzz_CTR_ShortBuffer(f *testing.F) { ctrInfo.FuzzShortBuffer(f, newCTR, newCTR, false) }

func Benchmark_CTR_New(b *testing.B) { ctrInfo.BenchNew(b, newCTR, newCTR, false) }
func Benchmark_CTR_Stream(b *testing.B) {
	ctrInfo.BenchStreamCompare(b, newCTR, newCTR, newCipherCTR, newCipherCTR, false)
}
