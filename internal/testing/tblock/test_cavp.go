package tblock

import (
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
)

func (ci *CipherInfo) TestEcbKAT(t *testing.T, newCipher func(key []byte) (cipher.Block, error), skip bool) {
	t.Helper()
	if skip {
		t.Skip("skipping test due to missing CPU feature")
	}

	for _, sz := range ci.Sizes {
		t.Run(sz.Name, func(t *testing.T) { cavp.TestBlockKAT(t, sz.KATFilename, cavp.ECB, newCipher) })
	}
}
