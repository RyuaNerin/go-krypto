//go:build !(amd64 || arm64) || purego
// +build !amd64,!arm64 purego

package lsh256

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_ShortWrite(t *testing.T) { HTSWA(t, as, newContext, false) }

func Test_LSH224_Go(t *testing.T) { HT(t, newContext(Size224), testCases224, false) }
func Test_LSH256_Go(t *testing.T) { HT(t, newContext(Size), testCases256, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Hash_8_Go(b *testing.B)  { HBA(b, as, newContext, 8, false) }
func Benchmark_Hash_1K_Go(b *testing.B) { HBA(b, as, newContext, 1024, false) }
func Benchmark_Hash_8K_Go(b *testing.B) { HBA(b, as, newContext, 8196, false) }
