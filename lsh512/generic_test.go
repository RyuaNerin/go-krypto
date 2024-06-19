package lsh512

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_ShortWrite_Go(t *testing.T) { HTSWA(t, as, newContextGo, false) }

func Test_LSH512_224_Go(t *testing.T) { HT(t, newContextGo(Size224), testCases224, false) }
func Test_LSH512_256_Go(t *testing.T) { HT(t, newContextGo(Size256), testCases256, false) }
func Test_LSH512_384_Go(t *testing.T) { HT(t, newContextGo(Size384), testCases384, false) }
func Test_LSH512_512_Go(t *testing.T) { HT(t, newContextGo(Size), testCases512, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Hash_8_Go(b *testing.B)  { HBA(b, as, newContextGo, 8, false) }
func Benchmark_Hash_1K_Go(b *testing.B) { HBA(b, as, newContextGo, 1024, false) }
func Benchmark_Hash_8K_Go(b *testing.B) { HBA(b, as, newContextGo, 8196, false) }
