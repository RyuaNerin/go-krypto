package lsh256

import (
	"hash"
	"testing"
)

const (
	benchBlockSize = BlockSize/2 + 1
)

func Benchmark_LSH224_Reset_Go(b *testing.B) { benchReset(b, newContextGo(lshType256H224), true) }
func Benchmark_LSH256_Reset_Go(b *testing.B) { benchReset(b, newContextGo(lshType256H256), true) }

func Benchmark_LSH224_Write_Go(b *testing.B) { benchWrite(b, newContextGo(lshType256H224), true) }
func Benchmark_LSH256_Write_Go(b *testing.B) { benchWrite(b, newContextGo(lshType256H256), true) }

func Benchmark_LSH224_WriteSum_Go(b *testing.B) { benchSum(b, newContextGo(lshType256H224), true) }
func Benchmark_LSH256_WriteSum_Go(b *testing.B) { benchSum(b, newContextGo(lshType256H256), true) }

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
