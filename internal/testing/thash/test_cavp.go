package thash

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
)

func (ht *HashTest) TestSMT(t *testing.T, skip bool) {
	t.Helper()
	if skip {
		t.Skip("skipping test due to missing CPU feature")
	}

	for _, sz := range ht.Sizes {
		t.Run(sz.Name, func(t *testing.T) {
			newHash := ht.getNew(t, sz.DigestSize)
			cavp.TestHashSMT(t, sz.KATFilename, newHash, sz.DigestSize)
		})
	}
}

func (ht *HashTest) TestSMTSum(t *testing.T, skip bool) {
	t.Helper()
	if skip {
		t.Skip("skipping test due to missing CPU feature")
	}

	for _, sz := range ht.Sizes {
		t.Run(sz.Name, func(t *testing.T) {
			sum := ht.getSum(t, sz.DigestSize)
			cavp.TastHashSumSMT(t, sz.KATFilename, sum, sz.DigestSize)
		})
	}
}
