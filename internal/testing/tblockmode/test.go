package tblockmode

import (
	"bytes"
	"crypto/cipher"
	"io"
	"slices"
	"strconv"
	"testing"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
	"github.com/RyuaNerin/go-krypto/internal/testing/tzeroize"

	"github.com/stretchr/testify/require"
)

const testGuardSize = 32

func skipTest(t testing.TB, skip bool) {
	t.Helper()
	if skip {
		t.Skip("skipping test due to missing CPU feature")
	}
}

func randomBytes(t testing.TB, r io.Reader, n int) []byte {
	t.Helper()
	b := make([]byte, n)
	_, err := io.ReadFull(r, b)
	require.NoError(t, err)
	return b
}

func normalizedBytes(src []byte, n int) []byte {
	if n == 0 {
		return nil
	}
	if len(src) == 0 {
		src = []byte{0}
	}
	if len(src) < n {
		src = bytes.Repeat(src, n/len(src)+1)
	}
	return slices.Clone(src[:n])
}

func guardedBuffer(n int) (buf []byte, check func(testing.TB, string)) {
	raw := make([]byte, testGuardSize+n+testGuardSize)
	zero := make([]byte, testGuardSize)
	return raw[testGuardSize : testGuardSize+n], func(t testing.TB, operation string) {
		t.Helper()
		require.Equalf(t, zero, raw[:testGuardSize], "%s wrote before dst", operation)
		require.Equalf(t, zero, raw[testGuardSize+n:], "%s wrote past dst", operation)
	}
}

func mustNewBlockMode(t testing.TB, newMode NewBlockMode, iv []byte, blockSize int) cipher.BlockMode {
	t.Helper()
	require.NotNil(t, newMode, "nil block mode constructor")

	var mode cipher.BlockMode
	require.NotPanics(t, func() { mode = newMode(iv) }, "constructor panicked for a valid IV")
	require.NotNil(t, mode, "constructor returned a nil BlockMode")
	require.Equal(t, blockSize, mode.BlockSize(), "unexpected BlockSize")
	return mode
}

func mustNewStream(t testing.TB, newStream NewStream, iv []byte) cipher.Stream {
	t.Helper()
	require.NotNil(t, newStream, "nil stream constructor")

	var stream cipher.Stream
	require.NotPanics(t, func() { stream = newStream(iv) }, "constructor panicked for a valid IV")
	require.NotNil(t, stream, "constructor returned a nil Stream")
	return stream
}

func mustNewAEAD(t testing.TB, info *AEADInfo, newAEAD NewAEAD) cipher.AEAD {
	t.Helper()
	require.NotNil(t, info)
	require.NotNil(t, newAEAD, "nil AEAD constructor")

	var aead cipher.AEAD
	var err error
	require.NotPanics(t, func() { aead, err = newAEAD() }, "constructor panicked")
	require.NoError(t, err, "constructor returned an error")
	require.NotNil(t, aead, "constructor returned a nil AEAD")
	require.Equal(t, info.NonceSize, aead.NonceSize(), "unexpected NonceSize")
	require.Equal(t, info.Overhead, aead.Overhead(), "unexpected Overhead")
	return aead
}

func blockModeOutput(mode cipher.BlockMode, src []byte) []byte {
	dst := make([]byte, len(src))
	mode.CryptBlocks(dst, src)
	return dst
}

func streamOutput(stream cipher.Stream, src []byte) []byte {
	dst := make([]byte, len(src))
	stream.XORKeyStream(dst, src)
	return dst
}

func (info *BlockModeInfo) TestNew(t *testing.T, newEncrypter, newDecrypter NewBlockMode, skip bool) {
	t.Helper()
	skipTest(t, skip)
	require.NotNil(t, info)
	require.Positive(t, info.BlockSize)
	require.GreaterOrEqual(t, info.IVSize, 0)

	r, _ := randutil.NewTestRand()
	for _, tc := range []struct {
		name string
		new  NewBlockMode
	}{
		{"Encrypt", newEncrypter},
		{"Decrypt", newDecrypter},
	} {
		t.Run(tc.name, func(t *testing.T) {
			iv := randomBytes(t, r, info.IVSize)

			t.Run("zeroize", func(t *testing.T) {
				mode := mustNewBlockMode(t, tc.new, slices.Clone(iv), info.BlockSize)
				_, ok := mode.(krypto.BlockMode)
				require.Truef(t, ok, "%T does not implement krypto.BlockMode", mode)
				tzeroize.TestZeroize(t, internal.UnwrapBlockMode(mode))
			})

			t.Run("stable BlockSize", func(t *testing.T) {
				mode := mustNewBlockMode(t, tc.new, slices.Clone(iv), info.BlockSize)
				src := randomBytes(t, r, 2*info.BlockSize)
				mode.CryptBlocks(make([]byte, len(src)), src)
				require.Equal(t, info.BlockSize, mode.BlockSize(), "BlockSize changed after use")
				require.Equal(t, info.BlockSize, mustNewBlockMode(t, tc.new, slices.Clone(iv), info.BlockSize).BlockSize(), "BlockSize changed between instances")
			})

			t.Run("IV ownership", func(t *testing.T) {
				callerIV := slices.Clone(iv)
				mode := mustNewBlockMode(t, tc.new, callerIV, info.BlockSize)
				reference := mustNewBlockMode(t, tc.new, slices.Clone(iv), info.BlockSize)
				for i := range callerIV {
					callerIV[i] ^= 0xff
				}

				src := randomBytes(t, r, 2*info.BlockSize)
				require.Equal(t, blockModeOutput(reference, src), blockModeOutput(mode, src), "mode retained the caller's mutable IV")
			})

			if info.IVSize > 0 {
				t.Run("rejects invalid IV sizes", func(t *testing.T) {
					for _, size := range []int{0, info.IVSize - 1, info.IVSize + 1, 2 * info.IVSize} {
						if size == info.IVSize {
							continue
						}
						badIV := randomBytes(t, r, size)
						original := slices.Clone(badIV)
						require.Panicsf(t, func() { tc.new(badIV) }, "constructor accepted a %d-byte IV", size)
						require.Equal(t, original, badIV, "constructor modified an invalid IV")
					}
				})
			}
		})
	}
}

func (info *BlockModeInfo) TestBlockMode(t *testing.T, newEncrypter, newDecrypter NewBlockMode, skip bool) {
	t.Helper()
	skipTest(t, skip)
	require.NotNil(t, info)
	require.Positive(t, info.BlockSize)
	require.GreaterOrEqual(t, info.IVSize, 0)

	r, _ := randutil.NewTestRand()
	iv := randomBytes(t, r, info.IVSize)
	blockCounts := []int{0, 1, 2, 3, 4, 7, 8, 13}

	t.Run("round trip", func(t *testing.T) {
		for _, blocks := range blockCounts {
			t.Run(testSizeName(blocks*info.BlockSize), func(t *testing.T) {
				src := randomBytes(t, r, blocks*info.BlockSize)
				original := slices.Clone(src)

				enc := mustNewBlockMode(t, newEncrypter, slices.Clone(iv), info.BlockSize)
				ct, checkCT := guardedBuffer(len(src))
				require.NotPanics(t, func() { enc.CryptBlocks(ct, src) })
				require.Equal(t, original, src, "encryption modified src")
				checkCT(t, "encryption")

				dec := mustNewBlockMode(t, newDecrypter, slices.Clone(iv), info.BlockSize)
				plain, checkPlain := guardedBuffer(len(src))
				require.NotPanics(t, func() { dec.CryptBlocks(plain, ct) })
				require.True(t, bytes.Equal(original, plain), "decrypt(encrypt(plain)) did not recover plain")
				checkPlain(t, "decryption")
			})
		}
	})

	t.Run("streaming state", func(t *testing.T) {
		src := randomBytes(t, r, 13*info.BlockSize)
		want := blockModeOutput(mustNewBlockMode(t, newEncrypter, slices.Clone(iv), info.BlockSize), src)

		for _, splitBlocks := range []int{0, 1, 4, 8, 13} {
			split := splitBlocks * info.BlockSize
			got := make([]byte, len(src))
			enc := mustNewBlockMode(t, newEncrypter, slices.Clone(iv), info.BlockSize)
			enc.CryptBlocks(got[:split], src[:split])
			enc.CryptBlocks(got[split:], src[split:])
			require.Equalf(t, want, got, "split encryption differs at block %d", splitBlocks)

			plain := make([]byte, len(src))
			dec := mustNewBlockMode(t, newDecrypter, slices.Clone(iv), info.BlockSize)
			dec.CryptBlocks(plain[:split], got[:split])
			dec.CryptBlocks(plain[split:], got[split:])
			require.Equalf(t, src, plain, "split decryption differs at block %d", splitBlocks)
		}
	})

	t.Run("in place", func(t *testing.T) {
		src := randomBytes(t, r, 8*info.BlockSize)
		want := blockModeOutput(mustNewBlockMode(t, newEncrypter, slices.Clone(iv), info.BlockSize), src)
		buf := slices.Clone(src)
		mustNewBlockMode(t, newEncrypter, slices.Clone(iv), info.BlockSize).CryptBlocks(buf, buf)
		require.Equal(t, want, buf)
		mustNewBlockMode(t, newDecrypter, slices.Clone(iv), info.BlockSize).CryptBlocks(buf, buf)
		require.Equal(t, src, buf)
	})

	t.Run("buffer contract", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			new  NewBlockMode
		}{
			{"Encrypt", newEncrypter},
			{"Decrypt", newDecrypter},
		} {
			t.Run(tc.name, func(t *testing.T) {
				newMode := func() cipher.BlockMode { return mustNewBlockMode(t, tc.new, slices.Clone(iv), info.BlockSize) }
				full := randomBytes(t, r, 2*info.BlockSize)

				for n := 1; n < info.BlockSize; n++ {
					require.Panicsf(t, func() { newMode().CryptBlocks(make([]byte, n), full[:n]) }, "accepted non-full input of length %d", n)
				}
				require.Panics(t, func() { newMode().CryptBlocks(make([]byte, len(full)-1), full) }, "accepted a short dst")

				for shift := 1; shift < info.BlockSize; shift++ {
					buf := randomBytes(t, r, info.BlockSize+shift)
					left, right := buf[:info.BlockSize], buf[shift:shift+info.BlockSize]
					require.Panicsf(t, func() { newMode().CryptBlocks(left, right) }, "accepted dst-before-src overlap, shift=%d", shift)
					require.Panicsf(t, func() { newMode().CryptBlocks(right, left) }, "accepted dst-after-src overlap, shift=%d", shift)
				}

				src := randomBytes(t, r, info.BlockSize)
				dst := randomBytes(t, r, info.BlockSize+testGuardSize)
				tail := slices.Clone(dst[info.BlockSize:])
				newMode().CryptBlocks(dst, src)
				require.Equal(t, tail, dst[info.BlockSize:], "wrote past src length")
			})
		}
	})
}

func (info *StreamInfo) TestNew(t *testing.T, newEncrypter, newDecrypter NewStream, skip bool) {
	t.Helper()
	skipTest(t, skip)
	require.NotNil(t, info)
	require.GreaterOrEqual(t, info.IVSize, 0)

	r, _ := randutil.NewTestRand()
	for _, tc := range []struct {
		name string
		new  NewStream
	}{
		{"Encrypt", newEncrypter},
		{"Decrypt", newDecrypter},
	} {
		t.Run(tc.name, func(t *testing.T) {
			iv := randomBytes(t, r, info.IVSize)

			t.Run("zeroize", func(t *testing.T) {
				stream := mustNewStream(t, tc.new, slices.Clone(iv))
				_, ok := stream.(krypto.Stream)
				require.Truef(t, ok, "%T does not implement krypto.Stream", stream)
				tzeroize.TestZeroize(t, internal.UnwrapStream(stream))
			})

			t.Run("IV ownership", func(t *testing.T) {
				callerIV := slices.Clone(iv)
				stream := mustNewStream(t, tc.new, callerIV)
				reference := mustNewStream(t, tc.new, slices.Clone(iv))
				for i := range callerIV {
					callerIV[i] ^= 0xff
				}
				src := randomBytes(t, r, 257)
				require.Equal(t, streamOutput(reference, src), streamOutput(stream, src), "stream retained the caller's mutable IV")
			})

			if info.IVSize > 0 {
				t.Run("rejects invalid IV sizes", func(t *testing.T) {
					for _, size := range []int{0, info.IVSize - 1, info.IVSize + 1, 2 * info.IVSize} {
						if size == info.IVSize {
							continue
						}
						badIV := randomBytes(t, r, size)
						original := slices.Clone(badIV)
						require.Panicsf(t, func() { tc.new(badIV) }, "constructor accepted a %d-byte IV", size)
						require.Equal(t, original, badIV, "constructor modified an invalid IV")
					}
				})
			}
		})
	}
}

func (info *StreamInfo) TestStream(t *testing.T, newEncrypter, newDecrypter NewStream, skip bool) {
	t.Helper()
	skipTest(t, skip)
	require.NotNil(t, info)
	require.GreaterOrEqual(t, info.IVSize, 0)

	r, _ := randutil.NewTestRand()
	iv := randomBytes(t, r, info.IVSize)
	base := max(info.IVSize, 16)
	sizes := []int{0, 1, base - 1, base, base + 1, 4*base - 1, 8*base + 5, 1024}

	t.Run("round trip", func(t *testing.T) {
		for _, size := range sizes {
			t.Run(testSizeName(size), func(t *testing.T) {
				src := randomBytes(t, r, size)
				original := slices.Clone(src)
				ct, checkCT := guardedBuffer(size)
				mustNewStream(t, newEncrypter, slices.Clone(iv)).XORKeyStream(ct, src)
				require.Equal(t, original, src, "encryption modified src")
				checkCT(t, "encryption")

				plain, checkPlain := guardedBuffer(size)
				mustNewStream(t, newDecrypter, slices.Clone(iv)).XORKeyStream(plain, ct)
				require.True(t, bytes.Equal(original, plain), "decrypt(encrypt(plain)) did not recover plain")
				checkPlain(t, "decryption")
			})
		}
	})

	t.Run("streaming state", func(t *testing.T) {
		src := randomBytes(t, r, 8*base+5)
		want := streamOutput(mustNewStream(t, newEncrypter, slices.Clone(iv)), src)
		for _, split := range []int{0, 1, base - 1, base, base + 1, len(src)} {
			got := make([]byte, len(src))
			enc := mustNewStream(t, newEncrypter, slices.Clone(iv))
			enc.XORKeyStream(got[:split], src[:split])
			enc.XORKeyStream(nil, nil)
			enc.XORKeyStream(got[split:], src[split:])
			require.Equalf(t, want, got, "split encryption differs at byte %d", split)

			plain := make([]byte, len(src))
			dec := mustNewStream(t, newDecrypter, slices.Clone(iv))
			dec.XORKeyStream(plain[:split], got[:split])
			dec.XORKeyStream(nil, nil)
			dec.XORKeyStream(plain[split:], got[split:])
			require.Equalf(t, src, plain, "split decryption differs at byte %d", split)
		}
	})

	t.Run("in place", func(t *testing.T) {
		src := randomBytes(t, r, 8*base+5)
		want := streamOutput(mustNewStream(t, newEncrypter, slices.Clone(iv)), src)
		buf := slices.Clone(src)
		mustNewStream(t, newEncrypter, slices.Clone(iv)).XORKeyStream(buf, buf)
		require.Equal(t, want, buf)
		mustNewStream(t, newDecrypter, slices.Clone(iv)).XORKeyStream(buf, buf)
		require.Equal(t, src, buf)
	})

	t.Run("buffer contract", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			new  NewStream
		}{
			{"Encrypt", newEncrypter},
			{"Decrypt", newDecrypter},
		} {
			t.Run(tc.name, func(t *testing.T) {
				newStream := func() cipher.Stream { return mustNewStream(t, tc.new, slices.Clone(iv)) }
				src := randomBytes(t, r, 2*base+1)
				require.Panics(t, func() { newStream().XORKeyStream(make([]byte, len(src)-1), src) }, "accepted a short dst")

				for shift := 1; shift < min(len(src), testGuardSize); shift++ {
					buf := randomBytes(t, r, len(src)+shift)
					left, right := buf[:len(src)], buf[shift:shift+len(src)]
					require.Panicsf(t, func() { newStream().XORKeyStream(left, right) }, "accepted dst-before-src overlap, shift=%d", shift)
					require.Panicsf(t, func() { newStream().XORKeyStream(right, left) }, "accepted dst-after-src overlap, shift=%d", shift)
				}

				dst := randomBytes(t, r, len(src)+testGuardSize)
				tail := slices.Clone(dst[len(src):])
				newStream().XORKeyStream(dst, src)
				require.Equal(t, tail, dst[len(src):], "wrote past src length")
			})
		}
	})
}

func (info *AEADInfo) TestNew(t *testing.T, newAEAD NewAEAD, skip bool) {
	t.Helper()
	skipTest(t, skip)
	require.NotNil(t, info)
	require.Positive(t, info.NonceSize)
	require.Positive(t, info.Overhead)

	t.Run("zeroize", func(t *testing.T) {
		aead := mustNewAEAD(t, info, newAEAD)
		_, ok := aead.(krypto.AEAD)
		require.Truef(t, ok, "%T does not implement krypto.AEAD", aead)
		tzeroize.TestZeroize(t, internal.UnwrapAEAD(aead))
	})

	t.Run("stable geometry", func(t *testing.T) {
		aead := mustNewAEAD(t, info, newAEAD)
		r, _ := randutil.NewTestRand()
		nonce := randomBytes(t, r, info.NonceSize)
		sealed := aead.Seal(nil, nonce, []byte("plaintext"), []byte("additional data"))
		require.Equal(t, info.NonceSize, aead.NonceSize(), "NonceSize changed after Seal")
		require.Equal(t, info.Overhead, aead.Overhead(), "Overhead changed after Seal")
		_, err := aead.Open(nil, nonce, sealed, []byte("additional data"))
		require.NoError(t, err)
		require.Equal(t, info.NonceSize, aead.NonceSize(), "NonceSize changed after Open")
		require.Equal(t, info.Overhead, aead.Overhead(), "Overhead changed after Open")
		mustNewAEAD(t, info, newAEAD)
	})
}

func (info *AEADInfo) TestAEAD(t *testing.T, newAEAD NewAEAD, skip bool) {
	t.Helper()
	skipTest(t, skip)
	require.NotNil(t, info)
	require.Positive(t, info.NonceSize)
	require.Positive(t, info.Overhead)

	r, _ := randutil.NewTestRand()
	nonce := randomBytes(t, r, info.NonceSize)
	sizes := []int{0, 1, 15, 16, 17, 127, 1024}

	t.Run("round trip", func(t *testing.T) {
		for _, size := range sizes {
			t.Run(testSizeName(size), func(t *testing.T) {
				aead := mustNewAEAD(t, info, newAEAD)
				plain := randomBytes(t, r, size)
				aad := randomBytes(t, r, size/2+1)
				plainCopy, aadCopy, nonceCopy := slices.Clone(plain), slices.Clone(aad), slices.Clone(nonce)

				sealed := aead.Seal(nil, nonce, plain, aad)
				require.Len(t, sealed, len(plain)+info.Overhead)
				require.Equal(t, plainCopy, plain, "Seal modified plaintext")
				require.Equal(t, aadCopy, aad, "Seal modified additional data")
				require.Equal(t, nonceCopy, nonce, "Seal modified nonce")

				sealedCopy := slices.Clone(sealed)
				opened, err := aead.Open(nil, nonce, sealed, aad)
				require.NoError(t, err)
				require.True(t, bytes.Equal(plain, opened), "Open did not recover plaintext")
				require.Equal(t, sealedCopy, sealed, "Open modified ciphertext")
				require.Equal(t, aadCopy, aad, "Open modified additional data")
				require.Equal(t, nonceCopy, nonce, "Open modified nonce")
			})
		}
	})

	t.Run("append semantics", func(t *testing.T) {
		aead := mustNewAEAD(t, info, newAEAD)
		plain := randomBytes(t, r, 257)
		aad := randomBytes(t, r, 31)
		prefix := randomBytes(t, r, 19)

		wantSealed := aead.Seal(slices.Clone(prefix), nonce, plain, aad)
		dst := make([]byte, len(prefix), len(prefix)+len(plain)+info.Overhead+testGuardSize)
		copy(dst, prefix)
		sealed := aead.Seal(dst, nonce, plain, aad)
		require.Equal(t, wantSealed, sealed, "Seal does not append to dst")
		require.Equal(t, prefix, sealed[:len(prefix)], "Seal modified the dst prefix")

		wantOpened := append(slices.Clone(prefix), plain...)
		opened, err := aead.Open(slices.Clone(prefix), nonce, sealed[len(prefix):], aad)
		require.NoError(t, err)
		require.Equal(t, wantOpened, opened, "Open does not append to dst")
	})

	t.Run("in place", func(t *testing.T) {
		aead := mustNewAEAD(t, info, newAEAD)
		plain := randomBytes(t, r, 257)
		aad := randomBytes(t, r, 31)
		want := aead.Seal(nil, nonce, plain, aad)

		buf := make([]byte, len(plain), len(plain)+info.Overhead)
		copy(buf, plain)
		sealed := aead.Seal(buf[:0], nonce, buf, aad)
		require.Equal(t, want, sealed, "in-place Seal differs")
		opened, err := aead.Open(sealed[:0], nonce, sealed, aad)
		require.NoError(t, err)
		require.Equal(t, plain, opened, "in-place Open differs")
	})

	t.Run("authentication", func(t *testing.T) {
		plain := randomBytes(t, r, 257)
		aad := randomBytes(t, r, 31)
		sealed := mustNewAEAD(t, info, newAEAD).Seal(nil, nonce, plain, aad)

		for _, index := range []int{0, len(sealed) / 2, len(sealed) - 1} {
			bad := slices.Clone(sealed)
			bad[index] ^= 0x80
			out, err := mustNewAEAD(t, info, newAEAD).Open(nil, nonce, bad, aad)
			require.Error(t, err, "Open accepted modified ciphertext at byte %d", index)
			require.Nil(t, out, "Open returned plaintext after authentication failure")
		}

		badAAD := slices.Clone(aad)
		badAAD[0] ^= 0x80
		out, err := mustNewAEAD(t, info, newAEAD).Open(nil, nonce, sealed, badAAD)
		require.Error(t, err, "Open accepted modified additional data")
		require.Nil(t, out)

		badNonce := slices.Clone(nonce)
		badNonce[0] ^= 0x80
		out, err = mustNewAEAD(t, info, newAEAD).Open(nil, badNonce, sealed, aad)
		require.Error(t, err, "Open accepted a different nonce")
		require.Nil(t, out)

		out, err = mustNewAEAD(t, info, newAEAD).Open(nil, nonce, sealed[:info.Overhead-1], aad)
		require.Error(t, err, "Open accepted ciphertext shorter than Overhead")
		require.Nil(t, out)
	})

	t.Run("nonce length", func(t *testing.T) {
		for _, size := range []int{0, info.NonceSize - 1, info.NonceSize + 1} {
			if size == info.NonceSize || size < 0 {
				continue
			}
			badNonce := make([]byte, size)
			require.Panicsf(t, func() { mustNewAEAD(t, info, newAEAD).Seal(nil, badNonce, nil, nil) }, "Seal accepted a %d-byte nonce", size)
			require.Panicsf(t, func() { _, _ = mustNewAEAD(t, info, newAEAD).Open(nil, badNonce, make([]byte, info.Overhead), nil) }, "Open accepted a %d-byte nonce", size)
		}
	})

	t.Run("inexact overlap", func(t *testing.T) {
		plain := randomBytes(t, r, 64)
		for shift := 1; shift < 8; shift++ {
			buf := make([]byte, shift+len(plain)+info.Overhead)
			copy(buf[shift:], plain)
			require.Panicsf(t, func() {
				mustNewAEAD(t, info, newAEAD).Seal(buf[:0], nonce, buf[shift:shift+len(plain)], nil)
			}, "Seal accepted inexact overlap, shift=%d", shift)

			sealed := mustNewAEAD(t, info, newAEAD).Seal(nil, nonce, plain, nil)
			buf = make([]byte, shift+len(sealed))
			copy(buf[shift:], sealed)
			require.Panicsf(t, func() {
				_, _ = mustNewAEAD(t, info, newAEAD).Open(buf[:0], nonce, buf[shift:], nil)
			}, "Open accepted inexact overlap, shift=%d", shift)
		}
	})
}

func testSizeName(size int) string {
	if size == 1024 {
		return "1K"
	}
	return fmtInt(size)
}

func fmtInt(n int) string {
	return strconv.Itoa(n)
}
