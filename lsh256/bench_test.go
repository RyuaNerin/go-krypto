package lsh256

import (
	"hash"
	"testing"
)

const (
	benchBlockSize = BlockSize/2 + 1
)

func benchReset(b *testing.B, h hash.Hash, nonskip bool) {
	if !nonskip {
		b.Skip()
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Reset()
	}
}
func benchWrite(b *testing.B, h hash.Hash, nonskip bool) {
	if !nonskip {
		b.Skip()
		return
	}
	buf := make([]byte, benchBlockSize)
	for idx := range buf {
		buf[idx] = byte(idx % 256)
	}

	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(buf)
	}
}
func benchSum(b *testing.B, h hash.Hash, nonskip bool) {
	if !nonskip {
		b.Skip()
		return
	}
	buf := make([]byte, benchBlockSize)
	for idx := range buf {
		buf[idx] = byte(idx % 256)
	}

	o := make([]byte, h.Size())

	b.SetBytes(int64(len(buf)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(buf)
		h.Sum(o[:0])
	}
}
