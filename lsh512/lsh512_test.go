package lsh512

import (
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/thash"
)

var (
	hi = thash.HashInfo{
		BlockSize: BlockSize,
		Sizes: []thash.CipherSize{
			{Name: "LSH-512-224", DigestSize: Size224, KATFilename: "06. LSH/LSH-512-224_(Byte)_SMT"},
			{Name: "LSH-512-256", DigestSize: Size256, KATFilename: "06. LSH/LSH-512-256_(Byte)_SMT"},
			{Name: "LSH-512-384", DigestSize: Size384, KATFilename: "06. LSH/LSH-512-384_(Byte)_SMT"},
			{Name: "LSH-512-512", DigestSize: Size, KATFilename: "06. LSH/LSH-512-512_(Byte)_SMT"},
		},
	}

	ht = thash.HashTest{
		HashInfo: hi,
		New: map[int]func() hash.Hash{
			Size224: New224,
			Size256: New256,
			Size384: New384,
			Size:    New,
		},
		Sum: map[int]any{
			Size224: Sum224,
			Size256: Sum256,
			Size384: Sum384,
			Size:    Sum512,
		},
	}
)

func Test_LSH512_New(t *testing.T)     { ht.TestNew(t, false) }
func Test_LSH512_Hash(t *testing.T)    { ht.TestHash(t, false) }
func Test_LSH512_SMT(t *testing.T)     { ht.TestSMT(t, false) }
func Test_LSH512_SMT_Sum(t *testing.T) { ht.TestSMTSum(t, false) }

func Fuzz_LSH512_New(f *testing.F)       { ht.FuzzNew(f, false) }
func Fuzz_LSH512_Hash(f *testing.F)      { ht.FuzzHash(f, false) }
func Fuzz_LSH512_Sum(f *testing.F)       { ht.FuzzSum(f, false) }
func Fuzz_LSH512_Marshal(f *testing.F)   { ht.FuzzMarshal(f, false) }
func Fuzz_LSH512_Unmarshal(f *testing.F) { ht.FuzzUnmarshal(f, false) }

func Benchmark_LSH512_New(b *testing.B)  { ht.BenchNew(b, false) }
func Benchmark_LSH512_Hash(b *testing.B) { ht.BenchHash(b, false) }
func Benchmark_LSH512_Sum(b *testing.B)  { ht.BenchSum(b, false) }
