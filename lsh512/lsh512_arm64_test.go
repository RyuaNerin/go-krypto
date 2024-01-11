//go:build arm64 && !purego
// +build arm64,!purego

package lsh512

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

var newNEON = simdSetNEON.NewContext

func Test_ShortWrite_NEON(t *testing.T) { HTSWA(t, as, newNEON, false) }

func Test_LSH512_224_NEON(t *testing.T) { HT(t, newNEON(Size224), testCases224, false) }
func Test_LSH512_256_NEON(t *testing.T) { HT(t, newNEON(Size256), testCases256, false) }
func Test_LSH512_384_NEON(t *testing.T) { HT(t, newNEON(Size384), testCases384, false) }
func Test_LSH512_512_NEON(t *testing.T) { HT(t, newNEON(Size), testCases512, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Hash_8_NEON(b *testing.B)  { HBA(b, as, newNEON, 8, false) }
func Benchmark_Hash_1K_NEON(b *testing.B) { HBA(b, as, newNEON, 1024, false) }
func Benchmark_Hash_8K_NEON(b *testing.B) { HBA(b, as, newNEON, 8192, false) }
