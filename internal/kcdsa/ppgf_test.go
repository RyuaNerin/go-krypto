package kcdsa

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"testing"
)

func BenchmarkPPGF(b *testing.B) {
	rnd := bufio.NewReaderSize(rand.Reader, 1<<15)

	bench := func(h hash.Hash) func(b *testing.B) {
		return func(b *testing.B) {
			dstBits := h.BlockSize()*h.BlockSize() + 3

			src := make([]byte, 1024)
			dst := make([]byte, dstBits)

			b.ReportAllocs()
			b.ResetTimer()
			b.SetBytes(int64(dstBits))
			for i := 0; i < b.N; i++ {
				rnd.Read(src)
				dst = ppgf(dst[:0], dstBits, h, src)
			}
		}
	}

	b.Run("Hash", bench(&hashWrap{h: sha256.New()}))
	b.Run("MarshalableHash", bench(sha256.New()))
}

type hashWrap struct {
	h hash.Hash
}

func (h *hashWrap) Write(p []byte) (int, error) { return h.h.Write(p) }
func (h *hashWrap) Sum(b []byte) []byte         { return h.h.Sum(b) }
func (h *hashWrap) Reset()                      { h.h.Reset() }
func (h *hashWrap) Size() int                   { return h.h.Size() }
func (h *hashWrap) BlockSize() int              { return h.h.BlockSize() }
