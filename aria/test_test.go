package aria

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/testingutil"
)

var (
	allSizes = []testingutil.CipherSize{
		{Name: "128", Size: 128},
		{Name: "196", Size: 196},
		{Name: "256", Size: 256},
	}
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New(b *testing.B) {
	testingutil.BA(b, allSizes, func(b *testing.B, keySize int) {
		testingutil.BBNew(b, keySize, 0, testingutil.BIW(NewCipher))
	})
}

func Benchmark_Encrypt(b *testing.B) { bench(b, testingutil.CE) }
func Benchmark_Decrypt(b *testing.B) { bench(b, testingutil.CD) }

func bench(b *testing.B, do testingutil.BD) {
	testingutil.BA(b, allSizes, func(b *testing.B, keySize int) {
		testingutil.BBDo(
			b,
			keySize,
			0,
			BlockSize,
			testingutil.BIW(NewCipher),
			do,
		)
	})
}
