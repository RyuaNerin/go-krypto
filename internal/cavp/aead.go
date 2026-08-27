package cavp

import (
	"crypto/cipher"
	"testing"

	"github.com/stretchr/testify/require"
)

func testAEAD(
	t testing.TB,
	iv, input, add,
	wantCipher, wantTag []byte,
	newAEAD func() (cipher.AEAD, error),
) {
	t.Helper()

	h, err := newAEAD()
	require.NoError(t, err)

	got := h.Seal(nil, iv, input, add)
	require.Len(t, got, len(wantCipher)+len(wantTag))
	require.Equal(t, wantCipher, got[:len(input)], "ciphertext mismatch")
	require.Equal(t, wantTag, got[len(input):], "tag mismatch")
}
