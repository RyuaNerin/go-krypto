//go:build !noasm && !noasm_krypto && !noasm_krypto_lea && arm64

package lea

import (
	"testing"
)

var cNEON = newCipherAsm

func Test_LEA_Block_NEON(t *testing.T) { ci.TestBlock(t, cNEON, false) }

func Fuzz_LEA_Block_NEON(f *testing.F)       { ci.FuzzBlock(f, cNEON, false) }
func Fuzz_LEA_Wide_NEON(f *testing.F)        { ci.FuzzWide(f, cNEON, false) }
func Fuzz_LEA_ShortBuffer_NEON(f *testing.F) { ci.FuzzBlockShortBuffer(f, cNEON, false) }

func Benchmark_LEA_Block_NEON(b *testing.B) { ci.BenchBlock4(b, cNEON, false) }
