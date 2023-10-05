package lsh512

import (
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

var as = []CipherSize{
	{Name: "512", Size: Size},
	{Name: "384", Size: Size384},
	{Name: "256", Size: Size256},
	{Name: "224", Size: Size224},
}

func Test_ShortWrite(t *testing.T) { HTSWA(t, as, newContextGo, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Hash_8_Go(b *testing.B)  { HBA(b, as, newContextGo, 8, false) }
func Benchmark_Hash_1K_Go(b *testing.B) { HBA(b, as, newContextGo, 1024, false) }
func Benchmark_Hash_8K_Go(b *testing.B) { HBA(b, as, newContextGo, 8196, false) }
