package lsh256

import (
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/thash"
)

var (
	hi = thash.HashInfo{
		BlockSize: BlockSize,
		Sizes: []thash.CipherSize{
			{Name: "LSH-256-224", DigestSize: Size224, KATFilename: "06. LSH/LSH-256-224_(Byte)_SMT"},
			{Name: "LSH-256-256", DigestSize: Size, KATFilename: "06. LSH/LSH-256-256_(Byte)_SMT"},
		},
	}

	ht = thash.HashTest{
		HashInfo: hi,
		New: map[int]func() hash.Hash{
			Size224: New224,
			Size:    New,
		},
		Sum: map[int]any{
			Size224: Sum224,
			Size:    Sum256,
		},
	}
)

func Test_LSH256_New(t *testing.T)     { ht.TestNew(t, false) }
func Test_LSH256_Hash(t *testing.T)    { ht.TestHash(t, false) }
func Test_LSH256_SMT(t *testing.T)     { ht.TestSMT(t, false) }
func Test_LSH256_SMT_Sum(t *testing.T) { ht.TestSMTSum(t, false) }

func Fuzz_LSH256_New(f *testing.F)       { ht.FuzzNew(f, false) }
func Fuzz_LSH256_Hash(f *testing.F)      { ht.FuzzHash(f, false) }
func Fuzz_LSH256_Sum(f *testing.F)       { ht.FuzzSum(f, false) }
func Fuzz_LSH256_Marshal(f *testing.F)   { ht.FuzzMarshal(f, false) }
func Fuzz_LSH256_Unmarshal(f *testing.F) { ht.FuzzUnmarshal(f, false) }

func Benchmark_LSH256_New(b *testing.B)  { ht.BenchNew(b, false) }
func Benchmark_LSH256_Hash(b *testing.B) { ht.BenchHash(b, false) }
func Benchmark_LSH256_Sum(b *testing.B)  { ht.BenchSum(b, false) }
