package cavp_test

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
	"github.com/RyuaNerin/go-krypto/lsh256"
)

func TestHashSumSMTArray(t *testing.T) {
	cavp.TastHashSumSMT(t, "06. LSH/LSH-256-224_(Byte)_SMT", lsh256.Sum224, lsh256.Size224)
}

func TestHashSumSMTSlice(t *testing.T) {
	sum := func(input []byte) []byte {
		got := lsh256.Sum224(input)
		return got[:]
	}
	cavp.TastHashSumSMT(t, "06. LSH/LSH-256-224_(Byte)_SMT", sum, lsh256.Size224)
}
