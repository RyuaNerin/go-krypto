package has160

import (
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

func Test_HAS160_ShortWrite(t *testing.T) { TestShortWrite(t, New()) }

func Benchmark_Hash_8_Go(b *testing.B)  { HB(b, New(), 8) }
func Benchmark_Hash_1K_Go(b *testing.B) { HB(b, New(), 1024) }
func Benchmark_Hash_8K_Go(b *testing.B) { HB(b, New(), 8192) }
