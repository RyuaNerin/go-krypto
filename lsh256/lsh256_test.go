package lsh256

import (
	"hash"
	"io"
	"math/rand"
	"testing"
)

const (
	benchBlockSize = 1024*BlockSize - 1
)

func benchReset(b *testing.B, h hash.Hash) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Reset()
	}
}
func benchWrite(b *testing.B, h hash.Hash) {
	r := rand.New(rand.NewSource(0))
	buf := make([]byte, benchBlockSize)
	io.ReadFull(r, buf)

	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(buf)
	}
}
func benchSum(b *testing.B, h hash.Hash) {
	r := rand.New(rand.NewSource(0))
	buf := make([]byte, benchBlockSize)
	io.ReadFull(r, buf)

	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(buf)
		h.Sum(nil)
	}
}

func BenchmarkLSH256GoReset(b *testing.B) {
	benchReset(b, newContextGo(lshType256H256))
}
func BenchmarkLSH256GoWrite(b *testing.B) {
	benchWrite(b, newContextGo(lshType256H256))
}
func BenchmarkLSH256GoSum(b *testing.B) {
	benchSum(b, newContextGo(lshType256H256))
}
