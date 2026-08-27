package tblockmode

import (
	"bytes"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

const fuzzMaxSize = 4096

func clampFuzzSize(size int) int {
	return min(max(size, 0), fuzzMaxSize)
}

func addIVSeeds(f *testing.F, ivSize int) {
	f.Add(make([]byte, ivSize), []byte(nil), 0)
	f.Add(bytes.Repeat([]byte{0xff}, ivSize), []byte("x"), 1)
	f.Add(bytes.Repeat([]byte{0xa5}, ivSize), []byte("The quick brown fox jumps over the lazy dog"), 17)
	f.Add(bytes.Repeat([]byte{0x5a}, ivSize), bytes.Repeat([]byte{0xff}, 1024), 513)
}

func (info *BlockModeInfo) FuzzNew(f *testing.F, newEncrypter, newDecrypter NewBlockMode, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}
	require.NotNil(f, info)
	require.Positive(f, info.BlockSize)
	require.GreaterOrEqual(f, info.IVSize, 0)
	require.NotNil(f, newEncrypter)
	require.NotNil(f, newDecrypter)

	f.Add(make([]byte, info.IVSize))
	f.Add(bytes.Repeat([]byte{0xff}, info.IVSize))
	f.Add([]byte("arbitrary IV material"))

	f.Fuzz(func(t *testing.T, input []byte) {
		for _, newMode := range []NewBlockMode{newEncrypter, newDecrypter} {
			if info.IVSize > 0 && len(input) != info.IVSize {
				original := slices.Clone(input)
				require.Panics(t, func() { newMode(input) }, "constructor accepted an invalid IV size")
				require.Equal(t, original, input, "constructor modified an invalid IV")
				continue
			}

			iv := normalizedBytes(input, info.IVSize)
			callerIV := slices.Clone(iv)
			mode := mustNewBlockMode(t, newMode, callerIV, info.BlockSize)
			reference := mustNewBlockMode(t, newMode, slices.Clone(iv), info.BlockSize)
			for i := range callerIV {
				callerIV[i] ^= 0xff
			}

			src := normalizedBytes(input, 2*info.BlockSize)
			require.Equal(t, blockModeOutput(reference, src), blockModeOutput(mode, src), "mode retained the caller's mutable IV")
		}
	})
}

func (info *BlockModeInfo) FuzzBlockMode(f *testing.F, newEncrypter, newDecrypter NewBlockMode, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}
	require.NotNil(f, info)
	require.Positive(f, info.BlockSize)
	require.GreaterOrEqual(f, info.IVSize, 0)
	require.NotNil(f, newEncrypter)
	require.NotNil(f, newDecrypter)
	addIVSeeds(f, info.IVSize)

	f.Fuzz(func(t *testing.T, ivInput, input []byte, split int) {
		iv := normalizedBytes(ivInput, info.IVSize)
		blocks := clampFuzzSize(len(input)) / info.BlockSize
		if len(input) > 0 && blocks == 0 {
			blocks = 1
		}
		src := normalizedBytes(input, blocks*info.BlockSize)
		original := slices.Clone(src)

		oneShot := blockModeOutput(mustNewBlockMode(t, newEncrypter, slices.Clone(iv), info.BlockSize), src)
		point := 0
		if blocks > 0 {
			point = (split%(blocks+1) + blocks + 1) % (blocks + 1)
		}
		point *= info.BlockSize

		ct, checkCT := guardedBuffer(len(src))
		enc := mustNewBlockMode(t, newEncrypter, slices.Clone(iv), info.BlockSize)
		enc.CryptBlocks(ct[:point], src[:point])
		enc.CryptBlocks(ct[point:], src[point:])
		require.Equal(t, oneShot, ct, "chunked encryption differs")
		require.Equal(t, original, src, "encryption modified src")
		checkCT(t, "encryption")

		plain, checkPlain := guardedBuffer(len(src))
		dec := mustNewBlockMode(t, newDecrypter, slices.Clone(iv), info.BlockSize)
		dec.CryptBlocks(plain[:point], ct[:point])
		dec.CryptBlocks(plain[point:], ct[point:])
		require.True(t, bytes.Equal(original, plain), "block mode round trip failed")
		checkPlain(t, "decryption")

		inPlace := slices.Clone(original)
		mustNewBlockMode(t, newEncrypter, slices.Clone(iv), info.BlockSize).CryptBlocks(inPlace, inPlace)
		require.True(t, bytes.Equal(oneShot, inPlace), "in-place encryption differs")
		mustNewBlockMode(t, newDecrypter, slices.Clone(iv), info.BlockSize).CryptBlocks(inPlace, inPlace)
		require.True(t, bytes.Equal(original, inPlace), "in-place decryption differs")
	})
}

func (info *BlockModeInfo) FuzzBlockModeShortBuffer(f *testing.F, newEncrypter, newDecrypter NewBlockMode, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}
	require.NotNil(f, info)
	require.Positive(f, info.BlockSize)
	require.GreaterOrEqual(f, info.IVSize, 0)
	require.NotNil(f, newEncrypter)
	require.NotNil(f, newDecrypter)

	f.Add(0, 0)
	f.Add(info.BlockSize, info.BlockSize)
	f.Add(info.BlockSize-1, info.BlockSize)
	f.Add(info.BlockSize, info.BlockSize-1)
	f.Add(2*info.BlockSize, 2*info.BlockSize)

	f.Fuzz(func(t *testing.T, dstLen, srcLen int) {
		dstLen, srcLen = clampFuzzSize(dstLen), clampFuzzSize(srcLen)
		iv := make([]byte, info.IVSize)
		dst, src := make([]byte, dstLen), make([]byte, srcLen)

		for _, newMode := range []NewBlockMode{newEncrypter, newDecrypter} {
			call := func() { mustNewBlockMode(t, newMode, iv, info.BlockSize).CryptBlocks(dst, src) }
			valid := srcLen%info.BlockSize == 0 && dstLen >= srcLen
			if valid {
				require.NotPanics(t, call)
			} else {
				require.Panics(t, call)
			}
		}
	})
}

func (info *StreamInfo) FuzzNew(f *testing.F, newEncrypter, newDecrypter NewStream, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}
	require.NotNil(f, info)
	require.GreaterOrEqual(f, info.IVSize, 0)
	require.NotNil(f, newEncrypter)
	require.NotNil(f, newDecrypter)

	f.Add(make([]byte, info.IVSize))
	f.Add(bytes.Repeat([]byte{0xff}, info.IVSize))
	f.Add([]byte("arbitrary IV material"))

	f.Fuzz(func(t *testing.T, input []byte) {
		for _, newStream := range []NewStream{newEncrypter, newDecrypter} {
			if info.IVSize > 0 && len(input) != info.IVSize {
				original := slices.Clone(input)
				require.Panics(t, func() { newStream(input) }, "constructor accepted an invalid IV size")
				require.Equal(t, original, input, "constructor modified an invalid IV")
				continue
			}

			iv := normalizedBytes(input, info.IVSize)
			callerIV := slices.Clone(iv)
			stream := mustNewStream(t, newStream, callerIV)
			reference := mustNewStream(t, newStream, slices.Clone(iv))
			for i := range callerIV {
				callerIV[i] ^= 0xff
			}
			src := normalizedBytes(input, 257)
			require.Equal(t, streamOutput(reference, src), streamOutput(stream, src), "stream retained the caller's mutable IV")
		}
	})
}

func (info *StreamInfo) FuzzStream(f *testing.F, newEncrypter, newDecrypter NewStream, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}
	require.NotNil(f, info)
	require.GreaterOrEqual(f, info.IVSize, 0)
	require.NotNil(f, newEncrypter)
	require.NotNil(f, newDecrypter)
	addIVSeeds(f, info.IVSize)

	f.Fuzz(func(t *testing.T, ivInput, input []byte, split int) {
		iv := normalizedBytes(ivInput, info.IVSize)
		src := slices.Clone(input[:min(len(input), fuzzMaxSize)])
		original := slices.Clone(src)
		point := 0
		if len(src) > 0 {
			point = (split%(len(src)+1) + len(src) + 1) % (len(src) + 1)
		}

		oneShot := streamOutput(mustNewStream(t, newEncrypter, slices.Clone(iv)), src)
		ct, checkCT := guardedBuffer(len(src))
		enc := mustNewStream(t, newEncrypter, slices.Clone(iv))
		enc.XORKeyStream(ct[:point], src[:point])
		enc.XORKeyStream(nil, nil)
		enc.XORKeyStream(ct[point:], src[point:])
		require.Equal(t, oneShot, ct, "chunked encryption differs")
		require.Equal(t, original, src, "encryption modified src")
		checkCT(t, "encryption")

		plain, checkPlain := guardedBuffer(len(src))
		dec := mustNewStream(t, newDecrypter, slices.Clone(iv))
		dec.XORKeyStream(plain[:point], ct[:point])
		dec.XORKeyStream(nil, nil)
		dec.XORKeyStream(plain[point:], ct[point:])
		require.True(t, bytes.Equal(original, plain), "stream round trip failed")
		checkPlain(t, "decryption")

		inPlace := slices.Clone(original)
		mustNewStream(t, newEncrypter, slices.Clone(iv)).XORKeyStream(inPlace, inPlace)
		require.True(t, bytes.Equal(oneShot, inPlace), "in-place encryption differs")
		mustNewStream(t, newDecrypter, slices.Clone(iv)).XORKeyStream(inPlace, inPlace)
		require.True(t, bytes.Equal(original, inPlace), "in-place decryption differs")
	})
}

func (info *StreamInfo) FuzzShortBuffer(f *testing.F, newEncrypter, newDecrypter NewStream, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}
	require.NotNil(f, info)
	require.GreaterOrEqual(f, info.IVSize, 0)
	require.NotNil(f, newEncrypter)
	require.NotNil(f, newDecrypter)

	f.Add(0, 0)
	f.Add(1, 1)
	f.Add(16, 17)
	f.Add(17, 16)
	f.Add(fuzzMaxSize, fuzzMaxSize)

	f.Fuzz(func(t *testing.T, dstLen, srcLen int) {
		dstLen, srcLen = clampFuzzSize(dstLen), clampFuzzSize(srcLen)
		iv := make([]byte, info.IVSize)
		dst, src := make([]byte, dstLen), make([]byte, srcLen)

		for _, newStream := range []NewStream{newEncrypter, newDecrypter} {
			call := func() { mustNewStream(t, newStream, iv).XORKeyStream(dst, src) }
			if dstLen >= srcLen {
				require.NotPanics(t, call)
			} else {
				require.Panics(t, call)
			}
		}
	})
}

func (info *AEADInfo) FuzzNew(f *testing.F, newAEAD NewAEAD, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}
	require.NotNil(f, info)
	require.Positive(f, info.NonceSize)
	require.Positive(f, info.Overhead)
	require.NotNil(f, newAEAD)

	f.Add([]byte(nil))
	f.Add([]byte("x"))
	f.Add(bytes.Repeat([]byte{0xff}, 257))

	f.Fuzz(func(t *testing.T, input []byte) {
		aead := mustNewAEAD(t, info, newAEAD)
		nonce := normalizedBytes(input, info.NonceSize)
		plain := input[:min(len(input), fuzzMaxSize)]
		sealed := aead.Seal(nil, nonce, plain, nil)
		require.Len(t, sealed, len(plain)+info.Overhead)
		require.Equal(t, info.NonceSize, aead.NonceSize())
		require.Equal(t, info.Overhead, aead.Overhead())
	})
}

func (info *AEADInfo) FuzzAEAD(f *testing.F, newAEAD NewAEAD, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}
	require.NotNil(f, info)
	require.Positive(f, info.NonceSize)
	require.Positive(f, info.Overhead)
	require.NotNil(f, newAEAD)

	f.Add(make([]byte, info.NonceSize), []byte(nil), []byte(nil), 0)
	f.Add(bytes.Repeat([]byte{0xff}, info.NonceSize), []byte("plaintext"), []byte("additional data"), 1)
	f.Add(bytes.Repeat([]byte{0xa5}, info.NonceSize), bytes.Repeat([]byte{0x5a}, 1024), bytes.Repeat([]byte{0xff}, 31), 512)

	f.Fuzz(func(t *testing.T, nonceInput, plaintext, additionalData []byte, tamper int) {
		nonce := normalizedBytes(nonceInput, info.NonceSize)
		plaintext = slices.Clone(plaintext[:min(len(plaintext), fuzzMaxSize)])
		additionalData = slices.Clone(additionalData[:min(len(additionalData), fuzzMaxSize)])
		plainCopy, aadCopy, nonceCopy := slices.Clone(plaintext), slices.Clone(additionalData), slices.Clone(nonce)

		aead := mustNewAEAD(t, info, newAEAD)
		sealed := aead.Seal(nil, nonce, plaintext, additionalData)
		require.Len(t, sealed, len(plaintext)+info.Overhead)
		require.Equal(t, plainCopy, plaintext, "Seal modified plaintext")
		require.Equal(t, aadCopy, additionalData, "Seal modified additional data")
		require.Equal(t, nonceCopy, nonce, "Seal modified nonce")

		sealedCopy := slices.Clone(sealed)
		opened, err := aead.Open(nil, nonce, sealed, additionalData)
		require.NoError(t, err)
		require.True(t, bytes.Equal(plaintext, opened), "Open did not recover plaintext")
		require.Equal(t, sealedCopy, sealed, "Open modified ciphertext")

		inPlace := make([]byte, len(plaintext), len(plaintext)+info.Overhead)
		copy(inPlace, plaintext)
		inPlace = mustNewAEAD(t, info, newAEAD).Seal(inPlace[:0], nonce, inPlace, additionalData)
		opened, err = mustNewAEAD(t, info, newAEAD).Open(inPlace[:0], nonce, inPlace, additionalData)
		require.NoError(t, err)
		require.True(t, bytes.Equal(plaintext, opened), "in-place Open did not recover plaintext")

		bad := slices.Clone(sealed)
		index := (tamper%len(bad) + len(bad)) % len(bad)
		bad[index] ^= 1
		opened, err = mustNewAEAD(t, info, newAEAD).Open(nil, nonce, bad, additionalData)
		require.Error(t, err, "Open accepted modified ciphertext")
		require.Nil(t, opened)
	})
}

// FuzzAEADOpen feeds arbitrary, attacker-controlled ciphertext to Open. It is
// intentionally separate from FuzzAEAD so callers can dedicate fuzzing time to
// the authentication and parsing path.
func (info *AEADInfo) FuzzAEADOpen(f *testing.F, newAEAD NewAEAD, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}
	require.NotNil(f, info)
	require.Positive(f, info.NonceSize)
	require.Positive(f, info.Overhead)
	require.NotNil(f, newAEAD)

	f.Add(make([]byte, info.NonceSize), []byte(nil), []byte(nil))
	f.Add(bytes.Repeat([]byte{0xff}, info.NonceSize), make([]byte, info.Overhead), []byte("additional data"))
	f.Add(bytes.Repeat([]byte{0xa5}, info.NonceSize), bytes.Repeat([]byte{0x5a}, 257), bytes.Repeat([]byte{0xff}, 31))

	f.Fuzz(func(t *testing.T, nonceInput, ciphertext, additionalData []byte) {
		nonce := normalizedBytes(nonceInput, info.NonceSize)
		ciphertext = slices.Clone(ciphertext[:min(len(ciphertext), fuzzMaxSize)])
		additionalData = slices.Clone(additionalData[:min(len(additionalData), fuzzMaxSize)])
		cipherCopy, aadCopy, nonceCopy := slices.Clone(ciphertext), slices.Clone(additionalData), slices.Clone(nonce)

		var plaintext []byte
		var err error
		require.NotPanics(t, func() {
			plaintext, err = mustNewAEAD(t, info, newAEAD).Open(nil, nonce, ciphertext, additionalData)
		}, "Open panicked for a correctly sized nonce")
		require.Equal(t, cipherCopy, ciphertext, "Open modified ciphertext")
		require.Equal(t, aadCopy, additionalData, "Open modified additional data")
		require.Equal(t, nonceCopy, nonce, "Open modified nonce")

		if err != nil {
			return
		}
		require.Len(t, plaintext, len(ciphertext)-info.Overhead)
		resealed := mustNewAEAD(t, info, newAEAD).Seal(nil, nonce, plaintext, additionalData)
		require.Equal(t, ciphertext, resealed, "Open accepted ciphertext that does not re-seal identically")
	})
}
