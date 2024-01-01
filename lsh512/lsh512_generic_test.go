//go:build !(amd64 || arm64) || purego

package lsh512

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_ShortWrite(t *testing.T) { HTSWA(t, as, newContext, false) }

func Test_LSH512_224(t *testing.T) { HT(t, newContext(Size224), testCases224, false) }
func Test_LSH512_256(t *testing.T) { HT(t, newContext(Size256), testCases256, false) }
func Test_LSH512_384(t *testing.T) { HT(t, newContext(Size384), testCases384, false) }
func Test_LSH512_512(t *testing.T) { HT(t, newContext(Size), testCases512, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Hash_8_Go(b *testing.B)  { HBA(b, as, newContext, 8, false) }
func Benchmark_Hash_1K_Go(b *testing.B) { HBA(b, as, newContext, 1024, false) }
func Benchmark_Hash_8K_Go(b *testing.B) { HBA(b, as, newContext, 8196, false) }
