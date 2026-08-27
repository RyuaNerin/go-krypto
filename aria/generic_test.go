package aria

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

var cg1 = internal.JustBlockFn(newCipher)

func Test_ARIA_New_Go(t *testing.T)   { ci.TestNewSIMD(t, cg1, false) }
func Test_ARIA_Block_Go(t *testing.T) { ci.TestBlock(t, cg1, false) }

func Test_ARIA_KAT_Go(t *testing.T) { ci.TestEcbKAT(t, cg1, false) }

func Fuzz_ARIA_Block_Go(f *testing.F)       { ci.FuzzBlock(f, cg1, false) }
func Fuzz_ARIA_Wide_Go(f *testing.F)        { ci.FuzzWide(f, cg1, false) }
func Fuzz_ARIA_ShortBuffer_Go(f *testing.F) { ci.FuzzBlockShortBuffer(f, cg1, false) }

func Benchmark_ARIA_New_Go(b *testing.B)   { ci.BenchNew(b, cg1, false) }
func Benchmark_ARIA_Block_Go(b *testing.B) { ci.BenchBlock(b, cg1, false) }
