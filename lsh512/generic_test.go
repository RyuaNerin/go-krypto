package lsh512

import (
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/thash"
)

var htGo = thash.HashTest{
	HashInfo: hi,
	New: map[int]func() hash.Hash{
		Size224: func() hash.Hash { return newContextGo(Size224) },
		Size256: func() hash.Hash { return newContextGo(Size256) },
		Size384: func() hash.Hash { return newContextGo(Size384) },
		Size:    func() hash.Hash { return newContextGo(Size) },
	},
	Sum: map[int]any{
		Size224: func(b []byte) (o [Size224]byte) {
			s := sumGo(Size224, b)
			copy(o[:], s[:Size224])
			return o
		},
		Size256: func(b []byte) (o [Size256]byte) {
			s := sumGo(Size256, b)
			copy(o[:], s[:Size256])
			return o
		},
		Size384: func(b []byte) (o [Size384]byte) {
			s := sumGo(Size384, b)
			copy(o[:], s[:Size384])
			return o
		},
		Size: sumGo,
	},
}

func Test_LSH512_Hash_Go(t *testing.T)    { htGo.TestHash(t, false) }
func Test_LSH512_SMT_Go(t *testing.T)     { htGo.TestSMT(t, false) }
func Test_LSH512_SMT_Sum_Go(t *testing.T) { htGo.TestSMTSum(t, false) }

func Fuzz_LSH512_Hash_Go(f *testing.F) { htGo.FuzzHash(f, false) }
func Fuzz_LSH512_Sum_Go(f *testing.F)  { htGo.FuzzSum(f, false) }

func Benchmark_LSH512_Hash_Go(b *testing.B) { htGo.BenchHash(b, false) }
func Benchmark_LSH512_Sum_Go(b *testing.B)  { htGo.BenchSum(b, false) }
