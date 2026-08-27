package has160

import (
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/thash"
)

var (
	hi = thash.HashInfo{
		BlockSize: BlockSize,
		Sizes: []thash.CipherSize{
			{Name: "HAS-160", DigestSize: Size},
		},
	}

	ht = thash.HashTest{
		HashInfo: hi,
		New: map[int]func() hash.Hash{
			Size: New,
		},
		Sum: map[int]any{
			Size: Sum,
		},
	}
)

func Test_HAS160_New(t *testing.T)  { ht.TestNew(t, false) }
func Test_HAS160_Hash(t *testing.T) { ht.TestHash(t, false) }

func Fuzz_HAS160_New(f *testing.F)       { ht.FuzzNew(f, false) }
func Fuzz_HAS160_Hash(f *testing.F)      { ht.FuzzHash(f, false) }
func Fuzz_HAS160_Sum(f *testing.F)       { ht.FuzzSum(f, false) }
func Fuzz_HAS160_Marshal(f *testing.F)   { ht.FuzzMarshal(f, false) }
func Fuzz_HAS160_Unmarshal(f *testing.F) { ht.FuzzUnmarshal(f, false) }

func Benchmark_HAS160_New(b *testing.B)  { ht.BenchNew(b, false) }
func Benchmark_HAS160_Hash(b *testing.B) { ht.BenchHash(b, false) }
func Benchmark_HAS160_Sum(b *testing.B)  { ht.BenchSum(b, false) }
