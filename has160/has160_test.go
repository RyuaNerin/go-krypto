package has160

import (
	"bytes"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"

	. "github.com/RyuaNerin/testingutil"
)

func Test_HAS160_ShortWrite(t *testing.T) { HTSW(t, New(), false) }

func TestHAS160(t *testing.T) {
	var sum []byte
	for tcIdx, tc := range testCases {
		expect := internal.HB(tc.MD)

		h := New()

		written := 0
		for n := 1; written < len(tc.MsgBytes); n++ {
			k := written + n
			if k > len(tc.MsgBytes) {
				k = len(tc.MsgBytes)
			}

			part := tc.MsgBytes[written:k]
			h.Write(part)
			written += n
		}

		sum = h.Sum(sum[:0])
		if !bytes.Equal(sum, expect) {
			t.Errorf("failed idx: %d\nexpect: %x\nactual: %x", tcIdx, expect, sum)
		}

		sum := Sum(tc.MsgBytes)
		if !bytes.Equal(sum[:], expect) {
			t.Errorf("failed idx: %d\nexpect: %x\nactual: %x", tcIdx, expect, sum)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	almost1K = 1024 - 5
	almost8K = 8*1024 - 5
)

func Benchmark_Hash_8(b *testing.B)  { HB(b, New(), 8, false) }
func Benchmark_Hash_1K(b *testing.B) { HB(b, New(), almost1K, false) }
func Benchmark_Hash_8K(b *testing.B) { HB(b, New(), almost8K, false) }
