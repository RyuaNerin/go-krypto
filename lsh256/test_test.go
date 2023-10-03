package lsh256

import (
	"hash"
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

var as = []CipherSize{
	{Name: "256", Size: Size},
	{Name: "224", Size: Size224},
}

func Test_ShortWrite(t *testing.T) {
	TA(t, as, func(t *testing.T, size int) {
		h := newContextGo(size)
		TestShortWrite(t, h)
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Hash_8_Go(b *testing.B)  { BA(b, as, bh(b, newContextGo, 8, true)) }
func Benchmark_Hash_1K_Go(b *testing.B) { BA(b, as, bh(b, newContextGo, 1024, true)) }
func Benchmark_Hash_8K_Go(b *testing.B) { BA(b, as, bh(b, newContextGo, 8196, true)) }

func bh(b *testing.B, newHash func(size int) hash.Hash, dataLen int, do bool) func(b *testing.B, size int) {
	return func(b *testing.B, size int) {
		if !do {
			b.Skip()
			return
		}

		HB(b, newHash(size), dataLen)
	}
}
