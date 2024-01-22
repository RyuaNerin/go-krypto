package has160

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_HAS160_ShortWrite(t *testing.T) { HTSW(t, New(), false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Hash_8(b *testing.B)  { HB(b, New(), 8, false) }
func Benchmark_Hash_1K(b *testing.B) { HB(b, New(), 1024, false) }
func Benchmark_Hash_8K(b *testing.B) { HB(b, New(), 8192, false) }
