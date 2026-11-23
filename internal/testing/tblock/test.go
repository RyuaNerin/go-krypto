package tblock

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"io"
	"math/rand/v2"
	"reflect"
	"slices"
	"sync"
	"testing"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
	"github.com/RyuaNerin/go-krypto/internal/testing/tzeroize"
	"github.com/RyuaNerin/go-krypto/internal/testing/unaligned"

	"github.com/stretchr/testify/require"
)

const (
	randomKeysPerSize  = 8
	randomBlocksPerKey = 32
	concurrentWorkers  = 16
	concurrentRounds   = 64
)

func (ci *CipherInfo) TestNew(t *testing.T, newCipher func(key []byte) (cipher.Block, error), keyError error) {
	t.Helper()

	r, _ := randutil.NewTestRand()

	t.Run("zeroize", func(t *testing.T) {
		for _, sz := range ci.Sizes {
			key := mustRandomBytes(t, r, sz.KeySize)
			b := mustNewBlock(t, newCipher, key)
			tzeroize.TestZeroize(t, internal.UnwrapBlock(b))
		}
	})

	t.Run("rejects invalid key sizes", ci.newInvalidKeySizesTest(r, newCipher, keyError))

	t.Run("same BlockSize", ci.newSameBlockSizeTest(r, newCipher))

	t.Run("key ownership", func(t *testing.T) {
		for _, sz := range ci.Sizes {
			testKeyOwnershipTest(t, r, sz.KeySize, newCipher)
		}
	})
}

// no rejects invalid key sizes test
func (ci *CipherInfo) TestNewSIMD(t *testing.T, newCipher func(key []byte) (cipher.Block, error), skip bool) {
	t.Helper()
	if skip {
		t.Skip("skipping test due to missing CPU feature")
	}

	r, _ := randutil.NewTestRand()

	t.Run("same BlockSize", ci.newSameBlockSizeTest(r, newCipher))

	t.Run("key ownership", func(t *testing.T) {
		for _, sz := range ci.Sizes {
			testKeyOwnershipTest(t, r, sz.KeySize, newCipher)
		}
	})
}

func (ci *CipherInfo) newInvalidKeySizesTest(r io.Reader, newCipher func([]byte) (cipher.Block, error), keyError error) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		t.Run("nil key", func(t *testing.T) {
			var err error
			var b cipher.Block
			require.NotPanics(t, func() { b, err = newCipher(nil) }, "constructor panicked for nil key")
			require.Nil(t, b, "constructor returned a non-nil block for nil key")
			require.ErrorIs(t, err, keyError, "constructor returned unexpected error for nil key")
		})

		maxKeySize := slices.MaxFunc(ci.Sizes, func(a, b CipherSize) int { return a.KeySize - b.KeySize }).KeySize

		key := mustRandomBytes(t, r, maxKeySize+256)
		for keySize := range key {
			var err error
			var b cipher.Block
			require.NotPanicsf(t, func() { b, err = newCipher(key[:keySize]) }, "constructor panicked for undeclared key size %d", keySize)
			if slices.ContainsFunc(ci.Sizes, func(sz CipherSize) bool { return sz.KeySize == keySize }) {
				require.NotNilf(t, b, "constructor returned a nil block for valid key size %d", keySize)
				require.NoErrorf(t, err, "constructor rejected valid key size %d", keySize)
			} else {
				require.Nilf(t, b, "constructor returned a non-nil block for undeclared key size %d", keySize)
				target := reflect.New(reflect.TypeOf(keyError)).Interface()
				require.ErrorAsf(t, err, target, "constructor returned unexpected error for invalid key size %d", keySize)
			}
		}
	}
}

func (ci *CipherInfo) newSameBlockSizeTest(r io.Reader, newCipher func([]byte) (cipher.Block, error)) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		{
			minKeySize := slices.MaxFunc(ci.Sizes, func(a, b CipherSize) int { return b.KeySize - a.KeySize }).KeySize

			b := mustNewBlock(t, newCipher, mustRandomBytes(t, r, minKeySize))

			require.Equal(t, ci.BlockSize, b.BlockSize(), "BlockSize is not stable across calls")
		}

		for _, sz := range ci.Sizes {
			for range randomKeysPerSize {
				b := mustNewBlock(t, newCipher, mustRandomBytes(t, r, sz.KeySize))
				require.Equalf(t, ci.BlockSize, b.BlockSize(), "BlockSize changed for key size %d", sz.KeySize)
			}
		}
	}
}

func testKeyOwnershipTest(t *testing.T, reader io.Reader, keySize int, newCipher func([]byte) (cipher.Block, error)) {
	t.Helper()

	key := mustRandomBytes(t, reader, keySize)
	originalKey := append([]byte(nil), key...)

	bGot := mustNewBlock(t, newCipher, key)
	bWant := mustNewBlock(t, newCipher, originalKey)

	plain := mustRandomBytes(t, reader, bGot.BlockSize())

	want := make([]byte, bGot.BlockSize())
	bWant.Encrypt(want, plain)

	for idx := range key {
		key[idx] ^= 0xff
	}
	got := make([]byte, bGot.BlockSize())
	bGot.Encrypt(got, plain)
	require.Equal(t, want, got, "cipher retained the caller's mutable key slice")
}

func (ci *CipherInfo) TestBlock(t *testing.T, newCipher func(key []byte) (cipher.Block, error), skip bool) {
	t.Helper()
	if skip {
		t.Skip("skipping test due to missing CPU feature")
	}

	r, rand := randutil.NewTestRand()

	require.NotNil(t, ci)
	require.NotNil(t, newCipher)

	for _, sz := range ci.Sizes {
		t.Run(sz.Name, newBlockTest(r, rand, sz.KeySize, newCipher))
	}

	t.Run("unaligned context", ci.newUnalignedContextTest(r, newCipher))
}

func newBlockTest(r io.Reader, rand *rand.Rand, keySize int, newCipher func([]byte) (cipher.Block, error)) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		key := mustRandomBytes(t, r, keySize)

		block := mustNewBlock(t, newCipher, key)
		blockSize := block.BlockSize()
		require.Positive(t, blockSize, "BlockSize must be positive")

		t.Run("independent", func(t *testing.T) {
			keyA := mustRandomBytes(t, r, keySize)
			keyB := mustRandomBytes(t, r, keySize)

			blockA := mustNewBlock(t, newCipher, keyA)
			blockB := mustNewBlock(t, newCipher, keyB)

			plain := mustRandomBytes(t, r, blockSize)
			wantA := make([]byte, blockSize)
			blockA.Encrypt(wantA, plain)

			scratch := make([]byte, blockSize)
			for range 8 {
				blockB.Encrypt(scratch, plain)
			}

			gotA := make([]byte, blockSize)
			blockA.Encrypt(gotA, plain)
			require.Equal(t, wantA, gotA, "using blockB affected blockA's output")
		})
		t.Run("short", testShortBuffers(block))
		t.Run("overlap", newOverlapTest(r, block))
		t.Run("round trip", func(t *testing.T) {
			for range randomKeysPerSize {
				key := mustRandomBytes(t, r, keySize)
				block := mustNewBlock(t, newCipher, key)
				require.Equal(t, blockSize, block.BlockSize(), "BlockSize changed between keys")

				for range randomBlocksPerKey {
					testRoundTrip(t, r, rand.IntN(blockSize+1), rand.IntN(blockSize+1), block)
				}
			}
		})
		t.Run("concurrent", newConcurrentUseTest(r, block, blockSize))
	}
}

func mustNewBlock(t *testing.T, newCipher func([]byte) (cipher.Block, error), key []byte) krypto.Block {
	t.Helper()

	var block cipher.Block
	var err error
	require.NotPanicsf(t, func() { block, err = newCipher(key) }, "constructor panicked for valid key size %d", len(key))
	require.NoErrorf(t, err, "constructor rejected valid key size %d", len(key))
	require.NotNil(t, block, "constructor returned a nil block")

	kryptoBlock, ok := block.(krypto.Block)
	require.True(t, ok, "constructor returned a krypto.Block that does not implement krypto.Block")
	return kryptoBlock
}

func testShortBuffers(block krypto.Block) func(t *testing.T) {
	blockSize := block.BlockSize()
	return func(t *testing.T) {
		t.Helper()

		full := make([]byte, krypto.MaxSupportedBlocks*blockSize)

		for size := range blockSize {
			short := make([]byte, size)
			require.Panics(t, func() { block.Encrypt(full, short) }, "Encrypt accepted short src")
			require.Panics(t, func() { block.Encrypt(short, full) }, "Encrypt accepted short dst")
			require.Panics(t, func() { block.Decrypt(full, short) }, "Decrypt accepted short src")
			require.Panics(t, func() { block.Decrypt(short, full) }, "Decrypt accepted short dst")
		}

		if block, ok := block.(krypto.Block4); ok {
			for size := range 4 * blockSize {
				short := make([]byte, size)
				require.Panics(t, func() { block.Encrypt4(full, short) }, "Encrypt accepted short src")
				require.Panics(t, func() { block.Encrypt4(short, full) }, "Encrypt accepted short dst")
				require.Panics(t, func() { block.Decrypt4(full, short) }, "Decrypt accepted short src")
				require.Panics(t, func() { block.Decrypt4(short, full) }, "Decrypt accepted short dst")
			}
		}
		if block, ok := block.(krypto.Block8); ok {
			for size := range 8 * blockSize {
				short := make([]byte, size)
				require.Panics(t, func() { block.Encrypt8(full, short) }, "Encrypt accepted short src")
				require.Panics(t, func() { block.Encrypt8(short, full) }, "Encrypt accepted short dst")
				require.Panics(t, func() { block.Decrypt8(full, short) }, "Decrypt accepted short src")
				require.Panics(t, func() { block.Decrypt8(short, full) }, "Decrypt accepted short dst")
			}
		}
	}
}

func newOverlapTest(reader io.Reader, block krypto.Block) func(*testing.T) {
	blockSize := block.BlockSize()

	return func(t *testing.T) {
		t.Helper()
		src := mustRandomBytes(t, reader, krypto.MaxSupportedBlocks*blockSize)
		want := make([]byte, krypto.MaxSupportedBlocks*blockSize)
		for idxBlock := range internal.MaxECBBlockSize {
			dst := want[idxBlock*blockSize : (idxBlock+1)*blockSize]
			src := src[idxBlock*blockSize : (idxBlock+1)*blockSize]
			block.Encrypt(dst, src)
		}

		t.Run("multiblock", func(t *testing.T) {
			if b4, ok := block.(krypto.Block4); ok {
				t.Run("4Blocks", newMultiblockTest(blockSize, b4.Encrypt4, b4.Decrypt4, 4, src, want))
			}
			if b8, ok := block.(krypto.Block8); ok {
				t.Run("8Blocks", newMultiblockTest(blockSize, b8.Encrypt8, b8.Decrypt8, 8, src, want))
			}
			if b16, ok := block.(krypto.Block16); ok {
				t.Run("16Blocks", newMultiblockTest(blockSize, b16.Encrypt16, b16.Decrypt16, 16, src, want))
			}
		})

		t.Run("exact", func(t *testing.T) {
			t.Run("1Block", newExactTest(blockSize, block.Encrypt, block.Decrypt, 1, src, want))
			if b4, ok := block.(krypto.Block4); ok {
				t.Run("4Blocks", newExactTest(blockSize, b4.Encrypt4, b4.Decrypt4, 4, src, want))
			}
			if b8, ok := block.(krypto.Block8); ok {
				t.Run("8Blocks", newExactTest(blockSize, b8.Encrypt8, b8.Decrypt8, 8, src, want))
			}
			if b16, ok := block.(krypto.Block16); ok {
				t.Run("16Blocks", newExactTest(blockSize, b16.Encrypt16, b16.Decrypt16, 16, src, want))
			}
		})

		t.Run("adjacent", func(t *testing.T) {
			t.Run("1Block", newAdjacentTest(blockSize, block.Encrypt, block.Decrypt, 1, src, want))
			if b4, ok := block.(krypto.Block4); ok {
				t.Run("4Blocks", newAdjacentTest(blockSize, b4.Encrypt4, b4.Decrypt4, 4, src, want))
			}
			if b8, ok := block.(krypto.Block8); ok {
				t.Run("8Blocks", newAdjacentTest(blockSize, b8.Encrypt8, b8.Decrypt8, 8, src, want))
			}
			if b16, ok := block.(krypto.Block16); ok {
				t.Run("16Blocks", newAdjacentTest(blockSize, b16.Encrypt16, b16.Decrypt16, 16, src, want))
			}
		})

		t.Run("overlap", func(t *testing.T) {
			t.Run("1Block", newOverlapTest2(blockSize, block.Encrypt, block.Decrypt, 1, reader))
			if b4, ok := block.(krypto.Block4); ok {
				t.Run("4Blocks", newOverlapTest2(blockSize, b4.Encrypt4, b4.Decrypt4, 4, reader))
			}
			if b8, ok := block.(krypto.Block8); ok {
				t.Run("8Blocks", newOverlapTest2(blockSize, b8.Encrypt8, b8.Decrypt8, 8, reader))
			}
			if b16, ok := block.(krypto.Block16); ok {
				t.Run("16Blocks", newOverlapTest2(blockSize, b16.Encrypt16, b16.Decrypt16, 16, reader))
			}
		})

		t.Run("no_overlap", func(t *testing.T) {
			t.Run("1Block", newNoOverlapTest(blockSize, block.Encrypt, block.Decrypt, 1, reader))
			if b4, ok := block.(krypto.Block4); ok {
				t.Run("4Blocks", newNoOverlapTest(blockSize, b4.Encrypt4, b4.Decrypt4, 4, reader))
			}
			if b8, ok := block.(krypto.Block8); ok {
				t.Run("8Blocks", newNoOverlapTest(blockSize, b8.Encrypt8, b8.Decrypt8, 8, reader))
			}
			if b16, ok := block.(krypto.Block16); ok {
				t.Run("16Blocks", newNoOverlapTest(blockSize, b16.Encrypt16, b16.Decrypt16, 16, reader))
			}
		})
	}
}

func newMultiblockTest(
	blockSize int,
	encrypt func([]byte, []byte),
	decrypt func([]byte, []byte),
	blocks int,
	src []byte,
	want []byte,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		size := blocks * blockSize
		src = src[:size]

		gotEnc := make([]byte, size)
		require.NotPanics(t, func() { encrypt(gotEnc, src) })
		require.Equal(t, want[:size], gotEnc, "Encrypt did not produce expected output")

		gotDec := make([]byte, size)
		require.NotPanics(t, func() { decrypt(gotDec, want[:size]) })
		require.Equal(t, src, gotDec, "Decrypt did not produce expected output")
	}
}

func newExactTest(
	blockSize int,
	encrypt func([]byte, []byte),
	decrypt func([]byte, []byte),
	blocks int,
	src []byte,
	want []byte,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		size := blocks * blockSize
		src = src[:size]
		want = want[:size]

		buf := slices.Clone(src)

		require.NotPanics(t, func() { encrypt(buf, buf) })
		require.Equal(t, want, buf)

		require.NotPanics(t, func() { decrypt(buf, buf) })
		require.Equal(t, src, buf)
	}
}

func newAdjacentTest(
	blockSize int,
	encrypt func([]byte, []byte),
	decrypt func([]byte, []byte),
	blocks int,
	src []byte,
	want []byte,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		srcSize := blocks * blockSize
		src := src[:srcSize]
		want := want[:srcSize]

		t.Run("encrypt/dst after src", func(t *testing.T) {
			buf := make([]byte, 2*srcSize)
			copy(buf[:srcSize], src)

			require.NotPanics(t, func() { encrypt(buf[srcSize:], buf[:srcSize]) })
			require.Equal(t, want, buf[srcSize:])
		})

		t.Run("encrypt/dst before src", func(t *testing.T) {
			buf := make([]byte, 2*srcSize)
			copy(buf[srcSize:], src)

			require.NotPanics(t, func() { encrypt(buf[:srcSize], buf[srcSize:]) })
			require.Equal(t, want, buf[:srcSize])
		})

		t.Run("decrypt/dst after src", func(t *testing.T) {
			buf := make([]byte, 2*srcSize)
			copy(buf[:srcSize], want)

			require.NotPanics(t, func() { decrypt(buf[srcSize:], buf[:srcSize]) })
			require.Equal(t, src, buf[srcSize:])
		})

		t.Run("decrypt/dst before src", func(t *testing.T) {
			buf := make([]byte, 2*srcSize)
			copy(buf[srcSize:], want)

			require.NotPanics(t, func() { decrypt(buf[:srcSize], buf[srcSize:]) })
			require.Equal(t, src, buf[:srcSize])
		})
	}
}

func newOverlapTest2(
	blockSize int,
	encrypt func([]byte, []byte),
	decrypt func([]byte, []byte),
	blocks int,
	reader io.Reader,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		size := blocks * blockSize

		for shift := 1; shift < size; shift++ {
			buf := mustRandomBytes(t, reader, size+shift)
			left := buf[:size]
			right := buf[shift : shift+size]

			require.Panics(t, func() { encrypt(left, right) }, "accepted dst-before-src overlap, shift=%d", shift)
			require.Panics(t, func() { encrypt(right, left) }, "accepted dst-after-src overlap, shift=%d", shift)
			require.Panics(t, func() { decrypt(left, right) }, "accepted dst-before-src overlap, shift=%d", shift)
			require.Panics(t, func() { decrypt(right, left) }, "accepted dst-after-src overlap, shift=%d", shift)
		}
	}
}

func newNoOverlapTest(
	blockSize int,
	encrypt func([]byte, []byte),
	decrypt func([]byte, []byte),
	blocks int,
	reader io.Reader,
) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		size := blocks * blockSize
		for _, shift := range []int{size, size + 1, 2 * size} {
			buf := mustRandomBytes(t, reader, size+shift)
			left := buf[:size]
			right := buf[shift : shift+size]

			srcCopy := slices.Clone(right)
			dstA := make([]byte, size)
			require.NotPanics(t, func() { encrypt(dstA, right) })
			require.NotPanics(t, func() { encrypt(left, srcCopy) }, "panicked on non-overlapping shift=%d", shift)
			require.Equal(t, dstA, left)

			dstB := make([]byte, size)
			require.NotPanics(t, func() { decrypt(dstB, dstA) })
			require.Equal(t, srcCopy, dstB)
		}
	}
}

func testRoundTrip(
	t *testing.T,
	reader io.Reader,
	prefix int,
	suffix int,
	block krypto.Block,
) {
	t.Helper()
	blockSize := block.BlockSize()

	srcGot := mustRandomBytes(t, reader, prefix+blockSize+suffix)
	srcWant := slices.Clone(srcGot)
	src := srcGot[prefix:]

	dstGot := mustRandomBytes(t, reader, prefix+blockSize+suffix)
	dstWant := slices.Clone(dstGot)
	dst := dstGot[prefix:]
	block.Encrypt(dst, src)

	require.Equal(t, srcWant, srcGot, "Encrypt modified src")
	require.Equal(t, dstWant[:prefix], dstGot[:prefix], "Encrypt wrote before dst")
	require.Equal(t, dstWant[prefix+blockSize:], dstGot[prefix+blockSize:], "Encrypt wrote after dst")

	repeated := make([]byte, blockSize)
	block.Encrypt(repeated, src)
	require.Equal(t, dst[:blockSize], repeated, "Encrypt is not deterministic")

	ciphertextBefore := append([]byte(nil), dstGot...)
	plainBuf := mustRandomBytes(t, reader, suffix+blockSize+prefix)
	plainBefore := append([]byte(nil), plainBuf...)
	plainDst := plainBuf[suffix:]
	block.Decrypt(plainDst, dst)

	require.Equal(t, ciphertextBefore, dstGot, "Decrypt modified src")
	require.Equal(t, src[:blockSize], plainDst[:blockSize], "Decrypt(Encrypt(plain)) did not recover plain")
	require.Equal(t, plainBefore[:suffix], plainBuf[:suffix], "Decrypt wrote before dst")
	require.Equal(t, plainBefore[suffix+blockSize:], plainBuf[suffix+blockSize:], "Decrypt wrote after dst")
}

func newConcurrentUseTest(reader io.Reader, block krypto.Block, blockSize int) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		input := mustRandomBytes(t, reader, blockSize)
		wantCipher := make([]byte, blockSize)
		block.Encrypt(wantCipher, input)

		errs := make(chan error, concurrentWorkers)

		var wg sync.WaitGroup
		for worker := range concurrentWorkers {
			wg.Go(func() {
				for round := range concurrentRounds {
					gotCipher := make([]byte, blockSize)
					block.Encrypt(gotCipher, input)
					if !bytes.Equal(wantCipher, gotCipher) {
						errs <- fmt.Errorf("worker %d round %d: concurrent Encrypt mismatch", worker, round)
						return
					}

					gotPlain := make([]byte, blockSize)
					block.Decrypt(gotPlain, gotCipher)
					if !bytes.Equal(input, gotPlain) {
						errs <- fmt.Errorf("worker %d round %d: concurrent Decrypt mismatch", worker, round)
						return
					}
				}
			})
		}
		wg.Wait()
		close(errs)

		for err := range errs {
			require.NoError(t, err)
		}
	}
}

func mustRandomBytes(t *testing.T, reader io.Reader, size int) []byte {
	t.Helper()
	buf := make([]byte, size)
	_, err := io.ReadFull(reader, buf)
	require.NoError(t, err)
	return buf
}

func (ci *CipherInfo) newUnalignedContextTest(r io.Reader, newCipher func([]byte) (cipher.Block, error)) func(t *testing.T) {
	return func(t *testing.T) {
		for _, sz := range ci.Sizes {
			t.Run(sz.Name, func(t *testing.T) {
				key := mustRandomBytes(t, r, sz.KeySize)
				plain := mustRandomBytes(t, r, krypto.MaxSupportedBlocks*ci.BlockSize)

				reference := mustNewBlock(t, newCipher, key)
				want := make([]byte, len(plain))
				for offset := 0; offset < len(plain); offset += ci.BlockSize {
					reference.Encrypt(want[offset:offset+ci.BlockSize], plain[offset:offset+ci.BlockSize])
				}

				unaligned.Test(t, 32, reference, func(residue uintptr) {
					context := unaligned.Copy(mustNewBlock(t, newCipher, key), 32, residue)
					block := context.Value

					got := make([]byte, ci.BlockSize)
					block.Encrypt(got, plain[:ci.BlockSize])
					require.Equalf(t, want[:ci.BlockSize], got, "Encrypt failed with an unaligned context (address mod 32 = %d)", residue)

					decrypted := make([]byte, ci.BlockSize)
					block.Decrypt(decrypted, got)
					require.Equalf(t, plain[:ci.BlockSize], decrypted, "Decrypt failed with an unaligned context (address mod 32 = %d)", residue)

					if block4, ok := block.(krypto.Block4); ok {
						got4 := make([]byte, 4*ci.BlockSize)
						block4.Encrypt4(got4, plain[:4*ci.BlockSize])
						require.Equalf(t, want[:4*ci.BlockSize], got4, "Encrypt4 failed with an unaligned context (address mod 32 = %d)", residue)

						decrypted4 := make([]byte, 4*ci.BlockSize)
						block4.Decrypt4(decrypted4, got4)
						require.Equalf(t, plain[:4*ci.BlockSize], decrypted4, "Decrypt4 failed with an unaligned context (address mod 32 = %d)", residue)
					}

					if block8, ok := block.(krypto.Block8); ok {
						got8 := make([]byte, 8*ci.BlockSize)
						block8.Encrypt8(got8, plain[:8*ci.BlockSize])
						require.Equalf(t, want[:8*ci.BlockSize], got8, "Encrypt8 failed with an unaligned context (address mod 32 = %d)", residue)

						decrypted8 := make([]byte, 8*ci.BlockSize)
						block8.Decrypt8(decrypted8, got8)
						require.Equalf(t, plain[:8*ci.BlockSize], decrypted8, "Decrypt8 failed with an unaligned context (address mod 32 = %d)", residue)
					}

					if block16, ok := block.(krypto.Block16); ok {
						got16 := make([]byte, 16*ci.BlockSize)
						block16.Encrypt16(got16, plain[:16*ci.BlockSize])
						require.Equalf(t, want[:16*ci.BlockSize], got16, "Encrypt16 failed with an unaligned context (address mod 32 = %d)", residue)

						decrypted16 := make([]byte, 16*ci.BlockSize)
						block16.Decrypt16(decrypted16, got16)
						require.Equalf(t, plain[:16*ci.BlockSize], decrypted16, "Decrypt16 failed with an unaligned context (address mod 32 = %d)", residue)
					}

					context.KeepAlive()
				})
			})
		}
	}
}
