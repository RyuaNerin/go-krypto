package tblock

import (
	"crypto/cipher"
	"io"
	"testing"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/randutil"

	"github.com/stretchr/testify/require"
)

func (ci *CipherInfo) BenchNew(b *testing.B, newCipher func(key []byte) (cipher.Block, error), skip bool) {
	b.Helper()

	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}

	for _, sz := range ci.Sizes {
		b.Run(
			sz.Name,
			func(b *testing.B) {
				r, _ := randutil.NewBenchRand()

				key := make([]byte, sz.KeySize)

				for b.Loop() {
					_, err := io.ReadFull(r, key)
					require.NoError(b, err)

					_, err = newCipher(key)
					require.NoError(b, err)
				}
			},
		)
	}
}

func (ci *CipherInfo) BenchBlock(
	b *testing.B,
	newCipher func(key []byte) (cipher.Block, error),
	skip bool,
) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	for _, sz := range ci.Sizes {
		b.Run(sz.Name, newBenchBlock(newCipher, sz.KeySize, true, true, true, true))
	}
}

func (ci *CipherInfo) BenchBlock1(
	b *testing.B,
	newCipher func(key []byte) (cipher.Block, error),
	skip bool,
) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	for _, sz := range ci.Sizes {
		b.Run(sz.Name, newBenchBlock(newCipher, sz.KeySize, true, false, false, false))
	}
}

func (ci *CipherInfo) BenchBlock14(
	b *testing.B,
	newCipher func(key []byte) (cipher.Block, error),
	skip bool,
) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	for _, sz := range ci.Sizes {
		b.Run(sz.Name, newBenchBlock(newCipher, sz.KeySize, true, true, false, false))
	}
}

func (ci *CipherInfo) BenchBlock4(
	b *testing.B,
	newCipher func(key []byte) (cipher.Block, error),
	skip bool,
) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	for _, sz := range ci.Sizes {
		b.Run(sz.Name, newBenchBlock(newCipher, sz.KeySize, false, true, false, false))
	}
}

func (ci *CipherInfo) BenchBlock8(
	b *testing.B,
	newCipher func(key []byte) (cipher.Block, error),
	skip bool,
) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	for _, sz := range ci.Sizes {
		b.Run(sz.Name, newBenchBlock(newCipher, sz.KeySize, false, false, true, false))
	}
}

func (ci *CipherInfo) BenchBlock16(
	b *testing.B,
	newCipher func(key []byte) (cipher.Block, error),
	skip bool,
) {
	b.Helper()
	if skip {
		b.Skip("skipping benchmark due to missing CPU feature")
	}
	for _, sz := range ci.Sizes {
		b.Run(sz.Name, newBenchBlock(newCipher, sz.KeySize, false, false, false, true))
	}
}

func newBenchBlock(
	newCipher func(key []byte) (cipher.Block, error),
	keySize int,
	bench1Block bool,
	bench4Blocks bool,
	bench8Blocks bool,
	bench16Blocks bool,
) func(b *testing.B) {
	return func(b *testing.B) {
		b.Helper()

		r, _ := randutil.NewBenchRand()

		key := make([]byte, keySize)
		_, err := io.ReadFull(r, key)
		require.NoError(b, err)

		c, err := newCipher(key)
		require.NoError(b, err)

		bench := func(blocks int, fn func(dst, src []byte)) func(b *testing.B) {
			return func(b *testing.B) {
				b.Helper()
				src := make([]byte, c.BlockSize()*blocks)
				dst := make([]byte, c.BlockSize()*blocks)

				b.SetBytes(int64(len(src)))
				for b.Loop() {
					_, err := io.ReadFull(r, src)
					require.NoError(b, err)

					fn(dst, src)
				}
			}
		}

		if bench1Block {
			b.Run("Encrypt/1Block", bench(1, c.Encrypt))
			b.Run("Decrypt/1Block", bench(1, c.Decrypt))
		}

		if bench4Blocks {
			if c, ok := c.(krypto.Block4); ok {
				b.Run("Encrypt/4Blocks", bench(4, c.Encrypt4))
				b.Run("Decrypt/4Blocks", bench(4, c.Decrypt4))
			} else {
				b.Run("Encrypt/4Blocks", func(b *testing.B) { b.Skip("skipping benchmark due to missing Block4 implementation") })
				b.Run("Decrypt/4Blocks", func(b *testing.B) { b.Skip("skipping benchmark due to missing Block4 implementation") })
			}
		}

		if bench8Blocks {
			if c, ok := c.(krypto.Block8); ok {
				b.Run("Encrypt/8Blocks", bench(8, c.Encrypt8))
				b.Run("Decrypt/8Blocks", bench(8, c.Decrypt8))
			} else {
				b.Run("Encrypt/8Blocks", func(b *testing.B) { b.Skip("skipping benchmark due to missing Block8 implementation") })
				b.Run("Decrypt/8Blocks", func(b *testing.B) { b.Skip("skipping benchmark due to missing Block8 implementation") })
			}
		}

		if bench16Blocks {
			if c, ok := c.(krypto.Block16); ok {
				b.Run("Encrypt/16Blocks", bench(16, c.Encrypt16))
				b.Run("Decrypt/16Blocks", bench(16, c.Decrypt16))
			} else {
				b.Run("Encrypt/16Blocks", func(b *testing.B) { b.Skip("skipping benchmark due to missing Block16 implementation") })
				b.Run("Decrypt/16Blocks", func(b *testing.B) { b.Skip("skipping benchmark due to missing Block16 implementation") })
			}
		}
	}
}
