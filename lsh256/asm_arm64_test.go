//go:build arm64 && !purego && (!gccgo || go1.18)
// +build arm64
// +build !purego
// +build !gccgo go1.18

package lsh256

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/aligned"

	. "github.com/RyuaNerin/testingutil"
)

var newNEON = simdSetNEON.NewContext

func Test_Aligned(t *testing.T) {
	var ctx lsh256ContextAsm
	aligned.TestIsAligned(t, 16, &ctx, "algType")
	aligned.TestIsAligned(t, 16, &ctx, "remainDataByteLen")
	aligned.TestIsAligned(t, 16, &ctx, "cvL")
	aligned.TestIsAligned(t, 16, &ctx, "cvR")
	aligned.TestIsAligned(t, 16, &ctx, "lastBlock")
}

func Test_ShortWrite_NEON(t *testing.T) { HTSWA(t, as, newNEON, !hasNEON) }

func Test_LSH224_NEON(t *testing.T) { HT(t, newNEON(Size224), testCases224, !hasNEON) }
func Test_LSH256_NEON(t *testing.T) { HT(t, newNEON(Size), testCases256, !hasNEON) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Hash_8_NEON(b *testing.B)  { HBA(b, as, newNEON, 8, !hasNEON) }
func Benchmark_Hash_1K_NEON(b *testing.B) { HBA(b, as, newNEON, 1024, !hasNEON) }
func Benchmark_Hash_8K_NEON(b *testing.B) { HBA(b, as, newNEON, 8192, !hasNEON) }
