//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh256 && arm64

package lsh256

import (
	"testing"
)

var htNEON = newHashTest(simdNEON{})

func Test_LSH256_Hash_NEON(t *testing.T) { htNEON.TestHash(t, false) }

func Test_LSH256_SMT_NEON(t *testing.T)     { htNEON.TestSMT(t, false) }
func Test_LSH256_SMT_Sum_NEON(t *testing.T) { htNEON.TestSMTSum(t, false) }

func Fuzz_LSH256_Hash_NEON(f *testing.F) { htNEON.FuzzHash(f, false) }
func Fuzz_LSH256_Sum_NEON(f *testing.F)  { htNEON.FuzzSum(f, false) }

func Benchmark_LSH256_Hash_NEON(b *testing.B) { htNEON.BenchHash(b, false) }
func Benchmark_LSH256_Sum_NEON(b *testing.B)  { htNEON.BenchSum(b, false) }
