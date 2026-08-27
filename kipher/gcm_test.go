package kipher_test

import (
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/tblockmode"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var gcmInfo = tblockmode.AEADInfo{
	NonceSize: kipher.GCMStandardNonceSize,
	Overhead:  kipher.GCMTagSize,
}

func newGCM() (cipher.AEAD, error)       { return kipher.NewGCM(newLEABlock()) }
func newCipherGCM() (cipher.AEAD, error) { return cipher.NewGCM(newLEABlock()) }

func Test_GCM_New(t *testing.T)  { gcmInfo.TestNew(t, newGCM, false) }
func Test_GCM_AEAD(t *testing.T) { gcmInfo.TestAEAD(t, newGCM, false) }
func Test_GCM_Compare(t *testing.T) {
	gcmInfo.TestCompare(t, newGCM, newCipherGCM, false)
}

func Fuzz_GCM_New(f *testing.F)      { gcmInfo.FuzzNew(f, newGCM, false) }
func Fuzz_GCM_AEAD(f *testing.F)     { gcmInfo.FuzzAEAD(f, newGCM, false) }
func Fuzz_GCM_AEADOpen(f *testing.F) { gcmInfo.FuzzAEADOpen(f, newGCM, false) }

func Benchmark_GCM_New(b *testing.B)  { gcmInfo.BenchNew(b, newGCM, false) }
func Benchmark_GCM_AEAD(b *testing.B) { gcmInfo.BenchAEADCompare(b, newGCM, newCipherGCM, false) }
