//nolint:dupl
package kipher_test

import (
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/tblockmode"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var ofbInfo = tblockmode.StreamInfo{IVSize: blockSize}

func newOFB(iv []byte) cipher.Stream       { return kipher.NewOFB(newLEABlock(), iv) }
func newCipherOFB(iv []byte) cipher.Stream { return cipher.NewOFB(newLEABlock(), iv) }

func Test_OFB_New(t *testing.T)    { ofbInfo.TestNew(t, newOFB, newOFB, false) }
func Test_OFB_Stream(t *testing.T) { ofbInfo.TestStream(t, newOFB, newOFB, false) }
func Test_OFB_Compare(t *testing.T) {
	ofbInfo.TestCompare(t, newOFB, newOFB, newCipherOFB, newCipherOFB, false)
}

func Fuzz_OFB_New(f *testing.F)    { ofbInfo.FuzzNew(f, newOFB, newOFB, false) }
func Fuzz_OFB_Stream(f *testing.F) { ofbInfo.FuzzStream(f, newOFB, newOFB, false) }

func Fuzz_OFB_ShortBuffer(f *testing.F) {
	ofbInfo.FuzzShortBuffer(f, newOFB, newOFB, false)
}

func Benchmark_OFB_New(b *testing.B) { ofbInfo.BenchNew(b, newOFB, newOFB, false) }
func Benchmark_OFB_Stream(b *testing.B) {
	ofbInfo.BenchStreamCompare(b, newOFB, newOFB, newCipherOFB, newCipherOFB, false)
}
