package tblockmode

import (
	"slices"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/randutil"

	"github.com/stretchr/testify/require"
)

// TestCompare compares a block mode implementation with a reference
// implementation for both encryption and decryption.
func (info *BlockModeInfo) TestCompare(
	t *testing.T,
	newEncrypter, newDecrypter NewBlockMode,
	newReferenceEncrypter, newReferenceDecrypter NewBlockMode,
	skip bool,
) {
	t.Helper()
	skipTest(t, skip)
	require.NotNil(t, info)
	require.Positive(t, info.BlockSize)
	require.GreaterOrEqual(t, info.IVSize, 0)
	require.NotNil(t, newEncrypter)
	require.NotNil(t, newDecrypter)
	require.NotNil(t, newReferenceEncrypter)
	require.NotNil(t, newReferenceDecrypter)

	r, _ := randutil.NewTestRand()
	iv := randomBytes(t, r, info.IVSize)

	for _, tc := range []struct {
		name         string
		newMode      NewBlockMode
		newReference NewBlockMode
	}{
		{"Encrypt", newEncrypter, newReferenceEncrypter},
		{"Decrypt", newDecrypter, newReferenceDecrypter},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, blocks := range []int{0, 1, 2, 3, 4, 7, 8, 13} {
				t.Run(testSizeName(blocks*info.BlockSize), func(t *testing.T) {
					src := randomBytes(t, r, blocks*info.BlockSize)
					got := blockModeOutput(
						mustNewBlockMode(t, tc.newMode, slices.Clone(iv), info.BlockSize),
						src,
					)
					want := blockModeOutput(
						mustNewBlockMode(t, tc.newReference, slices.Clone(iv), info.BlockSize),
						src,
					)
					require.Equal(t, want, got, "output differs from reference implementation")
				})
			}
		})
	}
}

// TestCompare compares a stream implementation with a reference
// implementation for both encryption and decryption.
func (info *StreamInfo) TestCompare(
	t *testing.T,
	newEncrypter, newDecrypter NewStream,
	newReferenceEncrypter, newReferenceDecrypter NewStream,
	skip bool,
) {
	t.Helper()
	skipTest(t, skip)
	require.NotNil(t, info)
	require.GreaterOrEqual(t, info.IVSize, 0)
	require.NotNil(t, newEncrypter)
	require.NotNil(t, newDecrypter)
	require.NotNil(t, newReferenceEncrypter)
	require.NotNil(t, newReferenceDecrypter)

	r, _ := randutil.NewTestRand()
	iv := randomBytes(t, r, info.IVSize)
	base := max(info.IVSize, 16)

	for _, tc := range []struct {
		name         string
		newStream    NewStream
		newReference NewStream
	}{
		{"Encrypt", newEncrypter, newReferenceEncrypter},
		{"Decrypt", newDecrypter, newReferenceDecrypter},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, size := range []int{0, 1, base - 1, base, base + 1, 4*base - 1, 8*base + 5, 1024} {
				t.Run(testSizeName(size), func(t *testing.T) {
					src := randomBytes(t, r, size)
					got := streamOutput(mustNewStream(t, tc.newStream, slices.Clone(iv)), src)
					want := streamOutput(mustNewStream(t, tc.newReference, slices.Clone(iv)), src)
					require.Equal(t, want, got, "output differs from reference implementation")
				})
			}
		})
	}
}

// TestCompare compares an AEAD implementation with a reference implementation.
func (info *AEADInfo) TestCompare(t *testing.T, newAEAD, newReference NewAEAD, skip bool) {
	t.Helper()
	skipTest(t, skip)
	require.NotNil(t, info)
	require.Positive(t, info.NonceSize)
	require.Positive(t, info.Overhead)
	require.NotNil(t, newAEAD)
	require.NotNil(t, newReference)

	r, _ := randutil.NewTestRand()
	nonce := randomBytes(t, r, info.NonceSize)

	for _, size := range []int{0, 1, 15, 16, 17, 127, 1024} {
		t.Run(testSizeName(size), func(t *testing.T) {
			plaintext := randomBytes(t, r, size)
			additionalData := randomBytes(t, r, size/2+1)

			got := mustNewAEAD(t, info, newAEAD).Seal(nil, nonce, plaintext, additionalData)
			want := mustNewAEAD(t, info, newReference).Seal(nil, nonce, plaintext, additionalData)
			require.Equal(t, want, got, "Seal output differs from reference implementation")

			gotPlaintext, gotErr := mustNewAEAD(t, info, newAEAD).Open(nil, nonce, want, additionalData)
			wantPlaintext, wantErr := mustNewAEAD(t, info, newReference).Open(nil, nonce, want, additionalData)
			require.NoError(t, gotErr, "Open rejected output from reference implementation")
			require.NoError(t, wantErr, "reference implementation rejected its own output")
			require.Equal(t, wantPlaintext, gotPlaintext, "Open output differs from reference implementation")
		})
	}
}
