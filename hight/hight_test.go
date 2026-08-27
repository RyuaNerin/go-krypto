package hight

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/testing/tblock"
)

var ci = tblock.CipherInfo{
	BlockSize: BlockSize,
	Sizes: []tblock.CipherSize{
		{Name: "HIGHT", KeySize: 128 / 8, KATFilename: "04. HIGHT/HIGHT_(ECB)_KAT"},
	},
}

var cg = internal.JustBlockFn(NewCipher)

func Test_HIGHT_New(t *testing.T)   { ci.TestNew(t, cg, KeySizeError(0)) }
func Test_HIGHT_Block(t *testing.T) { ci.TestBlock(t, cg, false) }

func Test_HIGHT_KAT(t *testing.T) { ci.TestEcbKAT(t, cg, false) }

func Fuzz_HIGHT_New(f *testing.F)         { ci.FuzzNew(f, NewCipher, false) }
func Fuzz_HIGHT_Block(f *testing.F)       { ci.FuzzBlock(f, NewCipher, false) }
func Fuzz_HIGHT_Wide(f *testing.F)        { ci.FuzzWide(f, NewCipher, false) }
func Fuzz_HIGHT_ShortBuffer(f *testing.F) { ci.FuzzBlockShortBuffer(f, NewCipher, false) }

func Benchmark_HIGHT_New(b *testing.B)   { ci.BenchNew(b, cg, false) }
func Benchmark_HIGHT_Block(b *testing.B) { ci.BenchBlock1(b, cg, false) }
