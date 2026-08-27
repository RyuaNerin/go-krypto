package lea

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/tblock"
)

var ci = tblock.CipherInfo{
	BlockSize: BlockSize,
	Sizes: []tblock.CipherSize{
		{Name: "LEA-128", KeySize: 128 / 8, KATFilename: "03. LEA/LEA-128_(ECB)_KAT"},
		{Name: "LEA-192", KeySize: 192 / 8, KATFilename: "03. LEA/LEA-192_(ECB)_KAT"},
		{Name: "LEA-256", KeySize: 256 / 8, KATFilename: "03. LEA/LEA-256_(ECB)_KAT"},
	},
}

func Test_LEA_New(t *testing.T)   { ci.TestNew(t, NewCipher, KeySizeError(0)) }
func Test_LEA_Block(t *testing.T) { ci.TestBlock(t, NewCipher, false) }

func Test_LEA_KAT(t *testing.T) { ci.TestEcbKAT(t, NewCipher, false) }

func Fuzz_LEA_New(f *testing.F)         { ci.FuzzNew(f, NewCipher, false) }
func Fuzz_LEA_Block(f *testing.F)       { ci.FuzzBlock(f, NewCipher, false) }
func Fuzz_LEA_Wide(f *testing.F)        { ci.FuzzWide(f, NewCipher, false) }
func Fuzz_LEA_ShortBuffer(f *testing.F) { ci.FuzzBlockShortBuffer(f, NewCipher, false) }

func Benchmark_LEA_New(b *testing.B)   { ci.BenchNew(b, NewCipher, false) }
func Benchmark_LEA_Block(b *testing.B) { ci.BenchBlock(b, NewCipher, false) }
