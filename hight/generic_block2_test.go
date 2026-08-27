//go:build amd64 || arm64 || loong64 || mips || mips64 || mips64le || mipsle || ppc64 || ppc64le || riscv64

package hight

import (
	"testing"
)

var cg4 = NewCipher

func Test_HIGHT_Block_Extended(t *testing.T) { ci.TestBlock(t, cg4, false) }

func Fuzz_HIGHT_Block_GoExtended(f *testing.F)       { ci.FuzzBlock(f, cg4, false) }
func Fuzz_HIGHT_Wide_GoExtended(f *testing.F)        { ci.FuzzWide(f, cg4, false) }
func Fuzz_HIGHT_ShortBuffer_GoExtended(f *testing.F) { ci.FuzzBlockShortBuffer(f, cg4, false) }

func Benchmark_HIGHT_Extended(b *testing.B) { ci.BenchBlock4(b, cg4, false) }
