//go:build !(arm64 || amd64) || purego || (gccgo && !go1.18)
// +build !arm64,!amd64 purego gccgo,!go1.18

package lsh256

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_ShortWrite_Go(t *testing.T) { HTSWA(t, as, newContextGo, false) }

func Test_LSH224_Go(t *testing.T) { HT(t, newContextGo(Size224), testCases224, false) }
func Test_LSH256_Go(t *testing.T) { HT(t, newContextGo(Size), testCases256, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Hash_8_Go(b *testing.B)  { HBA(b, as, newContextGo, 8, false) }
func Benchmark_Hash_1K_Go(b *testing.B) { HBA(b, as, newContextGo, 1024, false) }
func Benchmark_Hash_8K_Go(b *testing.B) { HBA(b, as, newContextGo, 8196, false) }
