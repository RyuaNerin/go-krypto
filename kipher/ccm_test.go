package kipher_test

import (
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/tblockmode"
	"github.com/RyuaNerin/go-krypto/kipher"
)

var ccmInfo = tblockmode.AEADInfo{
	NonceSize: 12,
	Overhead:  16,
}

func newCCM() (cipher.AEAD, error) {
	return kipher.NewCCM(newLEABlock(), ccmInfo.NonceSize, ccmInfo.Overhead)
}

func Test_CCM_New(t *testing.T)  { ccmInfo.TestNew(t, newCCM, false) }
func Test_CCM_AEAD(t *testing.T) { ccmInfo.TestAEAD(t, newCCM, false) }

func Fuzz_CCM_New(f *testing.F)      { ccmInfo.FuzzNew(f, newCCM, false) }
func Fuzz_CCM_AEAD(f *testing.F)     { ccmInfo.FuzzAEAD(f, newCCM, false) }
func Fuzz_CCM_AEADOpen(f *testing.F) { ccmInfo.FuzzAEADOpen(f, newCCM, false) }

func Benchmark_CCM_New(b *testing.B)  { ccmInfo.BenchNew(b, newCCM, false) }
func Benchmark_CCM_AEAD(b *testing.B) { ccmInfo.BenchAEAD(b, newCCM, false) }
