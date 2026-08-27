package aria

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/tblock"
)

var ci = tblock.CipherInfo{
	BlockSize: BlockSize,
	Sizes: []tblock.CipherSize{
		{Name: "ARIA-128", KeySize: 128 / 8, KATFilename: "01. ARIA/ARIA-128_(ECB)_KAT"},
		{Name: "ARIA-192", KeySize: 192 / 8, KATFilename: "01. ARIA/ARIA-192_(ECB)_KAT"},
		{Name: "ARIA-256", KeySize: 256 / 8, KATFilename: "01. ARIA/ARIA-256_(ECB)_KAT"},
	},
}

func Test_ARIA_New(t *testing.T)   { ci.TestNew(t, NewCipher, KeySizeError(0)) }
func Test_ARIA_Block(t *testing.T) { ci.TestBlock(t, NewCipher, false) }

func Test_ARIA_KAT(t *testing.T) { ci.TestEcbKAT(t, NewCipher, false) }

func Fuzz_ARIA_New(f *testing.F)         { ci.FuzzNew(f, NewCipher, false) }
func Fuzz_ARIA_Block(f *testing.F)       { ci.FuzzBlock(f, NewCipher, false) }
func Fuzz_ARIA_Wide(f *testing.F)        { ci.FuzzWide(f, NewCipher, false) }
func Fuzz_ARIA_ShortBuffer(f *testing.F) { ci.FuzzBlockShortBuffer(f, NewCipher, false) }

func Benchmark_ARIA_New(b *testing.B)   { ci.BenchNew(b, NewCipher, false) }
func Benchmark_ARIA_Block(b *testing.B) { ci.BenchBlock(b, NewCipher, false) }
