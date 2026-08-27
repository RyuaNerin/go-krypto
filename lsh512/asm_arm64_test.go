//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh512 && arm64

package lsh512

import (
	"testing"
)

var htNEON = newHashTest(simdNEON{})

func Test_LSH512_Hash_NEON(t *testing.T) { htNEON.TestHash(t, false) }

func Test_LSH512_SMT_NEON(t *testing.T)     { htNEON.TestSMT(t, false) }
func Test_LSH512_SMT_Sum_NEON(t *testing.T) { htNEON.TestSMTSum(t, false) }

func Fuzz_LSH512_Hash_NEON(f *testing.F) { htNEON.FuzzHash(f, false) }
func Fuzz_LSH512_Sum_NEON(f *testing.F)  { htNEON.FuzzSum(f, false) }

func Benchmark_LSH512_Hash_NEON(b *testing.B) { htNEON.BenchHash(b, false) }
func Benchmark_LSH512_Sum_NEON(b *testing.B)  { htNEON.BenchSum(b, false) }
