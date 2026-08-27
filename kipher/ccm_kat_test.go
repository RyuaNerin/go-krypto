package kipher_test

import (
	"crypto/aes"
	"encoding/hex"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/cavp"
	"github.com/RyuaNerin/go-krypto/kipher"
	"github.com/RyuaNerin/go-krypto/lea"

	"github.com/stretchr/testify/require"
)

func Test_CCM_LEA128_KAT(t *testing.T) {
	// Nonces of 7, 10, 12 and 13 bytes.
	cavp.TestCCMGe(t, "03. LEA/CCM_LEA-128_GE", lea.NewCipher)
}

// NIST SP 800-38C, Appendix C
func Test_CCM_SP800_38C(t *testing.T) {
	key := internal.HB("404142434445464748494a4b4c4d4e4f")
	testCases := []struct {
		nonce, adata, plaintext, ciphertext string
		tagSize                             int
	}{
		{"10111213141516", "0001020304050607", "20212223", "7162015b4dac255d", 4},
		{"1011121314151617", "000102030405060708090a0b0c0d0e0f", "202122232425262728292a2b2c2d2e2f", "d2a1f0e051ea5f62081a7792073d593d1fc64fbfaccd", 6},
		{"101112131415161718191a1b", "000102030405060708090a0b0c0d0e0f10111213", "202122232425262728292a2b2c2d2e2f3031323334353637", "e3b201a9f5b71a7a9b1ceaeccd97e70b6176aad9a4428aa5484392fbc1b09951", 8},
	}

	for _, tc := range testCases {
		b, err := aes.NewCipher(key)
		require.NoError(t, err)
		nonce := internal.HB(tc.nonce)

		aead, err := kipher.NewCCM(b, len(nonce), tc.tagSize)
		require.NoError(t, err)

		ct := aead.Seal(nil, nonce, internal.HB(tc.plaintext), internal.HB(tc.adata))
		require.Equal(t, tc.ciphertext, hex.EncodeToString(ct))

		pt, err := aead.Open(nil, nonce, ct, internal.HB(tc.adata))
		require.NoError(t, err)
		require.Equal(t, tc.plaintext, hex.EncodeToString(pt))
	}
}

func Test_CCM_MaxLength(t *testing.T) {
	// A 13-byte nonce leaves a 2-byte length field: at most 2^16 - 1 bytes.
	aead, err := kipher.NewCCM(newLEABlock(), 13, 16)
	require.NoError(t, err)
	nonce := make([]byte, 13)

	ct := aead.Seal(nil, nonce, make([]byte, 1<<16-1), nil)
	_, err = aead.Open(nil, nonce, ct, nil)
	require.NoError(t, err)

	require.Panics(t, func() { aead.Seal(nil, nonce, make([]byte, 1<<16), nil) })
	_, err = aead.Open(nil, nonce, make([]byte, 1<<16+16), nil)
	require.ErrorIs(t, err, kipher.ErrOpenFailed)
}
