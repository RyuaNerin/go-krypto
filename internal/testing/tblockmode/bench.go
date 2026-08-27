package tblockmode

import (
	"crypto/cipher"
	"io"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/randutil"

	"github.com/stretchr/testify/require"
)

type benchmarkSize struct {
	name string
	size int
}

var streamBenchmarkSizes = []benchmarkSize{
	{"8", 8},
	{"64", 64},
	{"1K", 1024 - 5},
	{"8K", 8*1024 - 5},
}

func blockModeBenchmarkSizes(blockSize int) []benchmarkSize {
	return []benchmarkSize{
		{testSizeName(blockSize), blockSize},
		{"1K", (1024 / blockSize) * blockSize},
		{"8K", (8 * 1024 / blockSize) * blockSize},
	}
}

func (info *BlockModeInfo) BenchNew(b *testing.B, newEncrypter, newDecrypter NewBlockMode, skip bool) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	require.NotNil(b, info)
	require.Positive(b, info.BlockSize)
	require.GreaterOrEqual(b, info.IVSize, 0)
	require.NotNil(b, newEncrypter)
	require.NotNil(b, newDecrypter)
	iv := make([]byte, info.IVSize)

	for _, tc := range []struct {
		name string
		new  NewBlockMode
	}{
		{"Encrypt", newEncrypter},
		{"Decrypt", newDecrypter},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				mode := tc.new(iv)
				require.NotNil(b, mode)
			}
		})
	}
}

func (info *BlockModeInfo) BenchBlockMode(b *testing.B, newEncrypter, newDecrypter NewBlockMode, skip bool) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	require.NotNil(b, info)
	require.Positive(b, info.BlockSize)
	require.GreaterOrEqual(b, info.IVSize, 0)
	require.NotNil(b, newEncrypter)
	require.NotNil(b, newDecrypter)
	r, _ := randutil.NewBenchRand()
	iv := make([]byte, info.IVSize)
	require.NoError(b, fillRandom(r, iv))

	for _, size := range blockModeBenchmarkSizes(info.BlockSize) {
		b.Run(size.name, func(b *testing.B) {
			src := make([]byte, size.size)
			dst := make([]byte, size.size)
			require.NoError(b, fillRandom(r, src))

			b.Run("Encrypt", benchBlockModeOperation(newEncrypter(iv), dst, src))
			b.Run("Decrypt", benchBlockModeOperation(newDecrypter(iv), dst, src))
		})
	}
}

// BenchCompare benchmarks a block mode implementation and a reference
// implementation for both encryption and decryption.
func (info *BlockModeInfo) BenchCompare(
	b *testing.B,
	newEncrypter, newDecrypter NewBlockMode,
	newReferenceEncrypter, newReferenceDecrypter NewBlockMode,
	skip bool,
) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	b.Run("Implementation", func(b *testing.B) {
		info.BenchBlockMode(b, newEncrypter, newDecrypter, false)
	})
	b.Run("Reference", func(b *testing.B) {
		info.BenchBlockMode(b, newReferenceEncrypter, newReferenceDecrypter, false)
	})
}

func benchBlockModeOperation(mode cipher.BlockMode, dst, src []byte) func(*testing.B) {
	return func(b *testing.B) {
		b.Helper()
		require.NotNil(b, mode)
		b.ReportAllocs()
		b.SetBytes(int64(len(src)))
		for b.Loop() {
			mode.CryptBlocks(dst, src)
		}
	}
}

func (info *StreamInfo) BenchNew(b *testing.B, newEncrypter, newDecrypter NewStream, skip bool) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	require.NotNil(b, info)
	require.GreaterOrEqual(b, info.IVSize, 0)
	require.NotNil(b, newEncrypter)
	require.NotNil(b, newDecrypter)
	iv := make([]byte, info.IVSize)

	for _, tc := range []struct {
		name string
		new  NewStream
	}{
		{"Encrypt", newEncrypter},
		{"Decrypt", newDecrypter},
	} {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				stream := tc.new(iv)
				require.NotNil(b, stream)
			}
		})
	}
}

func (info *StreamInfo) BenchStream(b *testing.B, newEncrypter, newDecrypter NewStream, skip bool) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	require.NotNil(b, info)
	require.GreaterOrEqual(b, info.IVSize, 0)
	require.NotNil(b, newEncrypter)
	require.NotNil(b, newDecrypter)
	r, _ := randutil.NewBenchRand()
	iv := make([]byte, info.IVSize)
	require.NoError(b, fillRandom(r, iv))

	for _, size := range streamBenchmarkSizes {
		b.Run(size.name, func(b *testing.B) {
			src := make([]byte, size.size)
			dst := make([]byte, size.size)
			require.NoError(b, fillRandom(r, src))

			b.Run("Encrypt", benchStreamOperation(newEncrypter(iv), dst, src))
			b.Run("Decrypt", benchStreamOperation(newDecrypter(iv), dst, src))
		})
	}
}

// BenchStreamCompare benchmarks a stream implementation and a reference
// implementation for both encryption and decryption.
func (info *StreamInfo) BenchStreamCompare(
	b *testing.B,
	newEncrypter, newDecrypter NewStream,
	newReferenceEncrypter, newReferenceDecrypter NewStream,
	skip bool,
) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	b.Run("Implementation", func(b *testing.B) {
		info.BenchStream(b, newEncrypter, newDecrypter, false)
	})
	b.Run("Reference", func(b *testing.B) {
		info.BenchStream(b, newReferenceEncrypter, newReferenceDecrypter, false)
	})
}

func benchStreamOperation(stream cipher.Stream, dst, src []byte) func(*testing.B) {
	return func(b *testing.B) {
		b.Helper()
		require.NotNil(b, stream)
		b.ReportAllocs()
		b.SetBytes(int64(len(src)))
		for b.Loop() {
			stream.XORKeyStream(dst, src)
		}
	}
}

func (info *AEADInfo) BenchNew(b *testing.B, newAEAD NewAEAD, skip bool) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	require.NotNil(b, info)
	require.Positive(b, info.NonceSize)
	require.Positive(b, info.Overhead)
	require.NotNil(b, newAEAD)

	b.ReportAllocs()
	for b.Loop() {
		aead, err := newAEAD()
		require.NoError(b, err)
		require.NotNil(b, aead)
	}
}

func (info *AEADInfo) BenchAEAD(b *testing.B, newAEAD NewAEAD, skip bool) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	require.NotNil(b, info)
	require.Positive(b, info.NonceSize)
	require.Positive(b, info.Overhead)
	require.NotNil(b, newAEAD)
	r, _ := randutil.NewBenchRand()
	nonce := make([]byte, info.NonceSize)
	additionalData := make([]byte, 13)
	require.NoError(b, fillRandom(r, nonce))
	require.NoError(b, fillRandom(r, additionalData))

	for _, size := range streamBenchmarkSizes {
		b.Run(size.name, func(b *testing.B) {
			aead := mustNewAEAD(b, info, newAEAD)
			plaintext := make([]byte, size.size)
			require.NoError(b, fillRandom(r, plaintext))
			ciphertext := aead.Seal(nil, nonce, plaintext, additionalData)

			b.Run("Seal", func(b *testing.B) {
				out := make([]byte, 0, len(ciphertext))
				b.ReportAllocs()
				b.SetBytes(int64(len(plaintext)))
				for b.Loop() {
					out = aead.Seal(out[:0], nonce, plaintext, additionalData)
				}
			})

			b.Run("Open", func(b *testing.B) {
				out := make([]byte, 0, len(plaintext))
				b.ReportAllocs()
				b.SetBytes(int64(len(plaintext)))
				for b.Loop() {
					var err error
					out, err = aead.Open(out[:0], nonce, ciphertext, additionalData)
					require.NoError(b, err)
				}
			})
		})
	}
}

// BenchAEADCompare benchmarks an AEAD implementation and a reference implementation.
func (info *AEADInfo) BenchAEADCompare(b *testing.B, newAEAD, newReference NewAEAD, skip bool) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	b.Run("Implementation", func(b *testing.B) {
		info.BenchAEAD(b, newAEAD, false)
	})
	b.Run("Reference", func(b *testing.B) {
		info.BenchAEAD(b, newReference, false)
	})
}

func fillRandom(r io.Reader, dst []byte) error {
	_, err := io.ReadFull(r, dst)
	return err
}
