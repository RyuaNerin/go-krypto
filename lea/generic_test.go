package lea

import (
	"testing"
)

func Test_LEA_New_Go(t *testing.T)   { ci.TestNewSIMD(t, newCipherGo, false) }
func Test_LEA_Block_Go(t *testing.T) { ci.TestBlock(t, newCipherGo, false) }

func Fuzz_LEA_Block_Go(f *testing.F)       { ci.FuzzBlock(f, newCipherGo, false) }
func Fuzz_LEA_Wide_Go(f *testing.F)        { ci.FuzzWide(f, newCipherGo, false) }
func Fuzz_LEA_ShortBuffer_Go(f *testing.F) { ci.FuzzBlockShortBuffer(f, newCipherGo, false) }

func Benchmark_LEA_New_Go(b *testing.B)   { ci.BenchNew(b, newCipherGo, false) }
func Benchmark_LEA_Block_Go(b *testing.B) { ci.BenchBlock1(b, newCipherGo, false) }
