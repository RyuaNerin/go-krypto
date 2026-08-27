package seed

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/testing/tblock"
)

var ci = tblock.CipherInfo{
	BlockSize: BlockSize,
	Sizes: []tblock.CipherSize{
		{Name: "SEED", KeySize: 16, KATFilename: "02. SEED/SEED-128_(ECB)_KAT"},
	},
}

var cg1 = internal.JustBlockFn(NewCipher)

func Test_SEED_New(t *testing.T)   { ci.TestNew(t, NewCipher, KeySizeError(0)) }
func Test_SEED_Block(t *testing.T) { ci.TestBlock(t, cg1, false) }

func Test_SEED_KAT(t *testing.T) { ci.TestEcbKAT(t, cg1, false) }

func Fuzz_SEED_New(f *testing.F)         { ci.FuzzNew(f, NewCipher, false) }
func Fuzz_SEED_Block(f *testing.F)       { ci.FuzzBlock(f, NewCipher, false) }
func Fuzz_SEED_Wide(f *testing.F)        { ci.FuzzWide(f, NewCipher, false) }
func Fuzz_SEED_ShortBuffer(f *testing.F) { ci.FuzzBlockShortBuffer(f, NewCipher, false) }

func Benchmark_SEED_New(b *testing.B)   { ci.BenchNew(b, NewCipher, false) }
func Benchmark_SEED_Block(b *testing.B) { ci.BenchBlock(b, cg1, false) }
