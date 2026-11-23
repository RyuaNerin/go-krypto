package thash

import (
	"bytes"
	"encoding"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"math/rand/v2"
	"slices"
	"sync"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
	"github.com/RyuaNerin/go-krypto/internal/testing/tzeroize"
	"github.com/RyuaNerin/go-krypto/internal/testing/unaligned"

	"github.com/stretchr/testify/require"
)

const (
	randomInputsPerRun = 8
	randomChunkTrials  = 32
	largeInputSize     = 1 << 20
	concurrentWorkers  = 16
	concurrentRounds   = 64
)

func (ht *HashTest) getNew(t testing.TB, digestSize int) func() hash.Hash {
	t.Helper()

	newHash, ok := ht.New[digestSize]
	require.Truef(t, ok, "no constructor for digest size %d", digestSize)
	require.NotNilf(t, newHash, "nil constructor for digest size %d", digestSize)
	return newHash
}

func (ht *HashTest) getSum(t testing.TB, digestSize int) func([]byte) []byte {
	t.Helper()

	sumFunc, ok := ht.Sum[digestSize]
	require.Truef(t, ok, "no sum function for digest size %d", digestSize)
	require.NotNilf(t, sumFunc, "nil sum function for digest size %d", digestSize)
	return NormalizeSum(sumFunc, digestSize)
}

func (ht *HashTest) TestNew(t *testing.T, skip bool) {
	t.Helper()
	if skip {
		t.Skip("skipping test due to missing CPU feature")
	}

	r, _ := randutil.NewTestRand()

	for _, size := range ht.Sizes {
		t.Run(size.Name, ht.newNewTest(r, size.DigestSize))
	}
}

func (ht *HashTest) TestHash(t *testing.T, skip bool) {
	t.Helper()
	if skip {
		t.Skip("skipping test due to missing CPU feature")
	}

	r, rand := randutil.NewTestRand()

	for _, size := range ht.Sizes {
		t.Run(size.Name, ht.newHashTest(r, rand, size.DigestSize))
	}
}

func (ht *HashTest) TestSum(t *testing.T, skip bool) {
	t.Helper()
	if skip {
		t.Skip("skipping test due to missing CPU feature")
	}

	if len(ht.Sum) == 0 {
		t.Skip("no sum function")
	}

	r, rand := randutil.NewTestRand()

	for _, size := range ht.Sizes {
		if _, ok := ht.Sum[size.DigestSize]; !ok {
			continue
		}
		t.Run(size.Name, ht.newSumFuncTest(r, rand, size.DigestSize))
	}
}

func (ht *HashTest) newNewTest(r io.Reader, digestSize int) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		newHash := ht.getNew(t, digestSize)

		t.Run("zeroize", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			tzeroize.TestZeroize(t, internal.UnwrapHash(h))
		})

		t.Run("declared digest size", func(t *testing.T) {
			got := mustNewHash(t, newHash).Size()
			require.Equalf(t, digestSize, got, "Size() does not match the HashTest.New key")

			if len(ht.HashInfo.Sizes) == 0 {
				return
			}
			declared := slices.ContainsFunc(ht.HashInfo.Sizes, func(sz CipherSize) bool { return sz.DigestSize == got })
			require.Truef(t, declared, "Size() = %d is not declared in HashInfo.Sizes", got)
		})

		t.Run("same Size", ht.newSameSizeTest(r, newHash))

		t.Run("same BlockSize", ht.newSameBlockSizeTest(r, newHash))

		t.Run("input ownership", func(t *testing.T) {
			testInputOwnershipTest(t, r, ht.HashInfo.BlockSize, newHash)
		})

		if !ht.HashInfo.SkipUnalignedTest {
			t.Run("unaligned context", newUnalignedContextTest(r, newHash))
		}
	}
}

func newUnalignedContextTest(r io.Reader, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		reference := mustNewHash(t, newHash)
		input := mustRandomBytes(t, r, 3*reference.BlockSize()+7)
		want := mustSum(t, newHash, input)

		unaligned.Test(t, 32, reference, func(residue uintptr) {
			context := unaligned.Copy(mustNewHash(t, newHash), 32, residue)
			h := context.Value
			h.Reset()

			split := h.BlockSize() - 1
			n, err := h.Write(input[:split])
			require.NoErrorf(t, err, "Write failed at address mod 32 = %d", residue)
			require.Equalf(t, split, n, "Write returned a wrong count at address mod 32 = %d", residue)

			n, err = h.Write(input[split:])
			require.NoErrorf(t, err, "Write failed at address mod 32 = %d", residue)
			require.Equalf(t, len(input)-split, n, "Write returned a wrong count at address mod 32 = %d", residue)
			require.Equalf(t, want, h.Sum(nil), "hash failed with an unaligned context (address mod 32 = %d)", residue)

			context.KeepAlive()
		})
	}
}

func (ht *HashTest) newSameSizeTest(r io.Reader, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		h := mustNewHash(t, newHash)
		want := h.Size()
		require.Positive(t, want, "Size must be positive")

		h.Write(mustRandomBytes(t, r, 3*h.BlockSize()+7))
		require.Equal(t, want, h.Size(), "Size changed after Write")

		require.Len(t, h.Sum(nil), want, "digest length does not match Size")
		require.Equal(t, want, h.Size(), "Size changed after Sum")

		h.Reset()
		require.Equal(t, want, h.Size(), "Size changed after Reset")

		for range randomInputsPerRun {
			require.Equal(t, want, mustNewHash(t, newHash).Size(), "Size is not stable across instances")
		}
	}
}

func (ht *HashTest) newSameBlockSizeTest(r io.Reader, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		h := mustNewHash(t, newHash)
		want := h.BlockSize()
		require.Positive(t, want, "BlockSize must be positive")
		if ht.HashInfo.BlockSize != 0 {
			require.Equal(t, want, ht.HashInfo.BlockSize, "BlockSize does not match HashInfo.BlockSize")
		}

		h.Write(mustRandomBytes(t, r, 3*want))
		require.Equal(t, want, h.BlockSize(), "BlockSize changed after Write")

		h.Sum(nil)
		require.Equal(t, want, h.BlockSize(), "BlockSize changed after Sum")

		h.Reset()
		require.Equal(t, want, h.BlockSize(), "BlockSize changed after Reset")

		for range randomInputsPerRun {
			require.Equal(t, want, mustNewHash(t, newHash).BlockSize(), "BlockSize is not stable across instances")
		}
	}
}

func testInputOwnershipTest(t *testing.T, reader io.Reader, blockSize int, newHash func() hash.Hash) {
	t.Helper()

	src := mustRandomBytes(t, reader, 3*blockSize+7)
	originalSrc := slices.Clone(src)
	want := mustSum(t, newHash, originalSrc)

	h := mustNewHash(t, newHash)
	h.Write(src)
	require.Equal(t, originalSrc, src, "Write modified the caller's input slice")

	for idx := range src {
		src[idx] ^= 0xff
	}
	require.Equal(t, want, h.Sum(nil), "hash retained the caller's mutable input slice")
}

func (ht *HashTest) newHashTest(r io.Reader, rand *rand.Rand, digestSize int) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		newHash := ht.getNew(t, digestSize)

		h := mustNewHash(t, newHash)
		blockSize := h.BlockSize()
		require.Equal(t, digestSize, h.Size(), "Size does not match the HashTest.New key")

		t.Run("independent", ht.newIndependentTest(newHash))
		t.Run("empty", ht.newEmptyInputTest(newHash))
		t.Run("write", ht.newWriteTest(r, blockSize, newHash))
		t.Run("sum", ht.newSumTest(digestSize, newHash))
		t.Run("reset", ht.newResetTest(r, rand, blockSize, newHash))
		t.Run("streaming", ht.newStreamingTest(r, rand, blockSize, newHash))
		t.Run("boundary", ht.newBoundaryTest(r, blockSize, newHash))
		t.Run("sensitivity", ht.newSensitivityTest(r, rand, blockSize, newHash))
		t.Run("io", ht.newIOTest(newHash))
		t.Run("marshal", ht.newMarshalTest(r, rand, blockSize, newHash))
		t.Run("concurrent", ht.newConcurrentUseTest(r, blockSize, newHash))
	}
}

func (ht *HashTest) newIndependentTest(newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		hA := mustNewHash(t, newHash)
		hB := mustNewHash(t, newHash)

		hA.Write([]byte("only-A"))
		wantA := slices.Clone(hA.Sum(nil))

		require.Equal(t, mustNewHash(t, newHash).Sum(nil), hB.Sum(nil), "writing to hA affected hB")

		for range randomInputsPerRun {
			hB.Write([]byte("only-B"))
			hB.Sum(nil)
		}
		require.Equal(t, wantA, hA.Sum(nil), "using hB affected hA's output")
	}
}

func (ht *HashTest) newEmptyInputTest(newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		want := mustNewHash(t, newHash).Sum(nil)
		require.Len(t, want, mustNewHash(t, newHash).Size())
		require.Equal(t, want, mustNewHash(t, newHash).Sum(nil), "empty digest is not deterministic")

		h := mustNewHash(t, newHash)
		h.Write(nil)
		h.Write([]byte{})
		require.Equal(t, want, h.Sum(nil), "empty Write changed the state")

		require.NotEqual(t, want, mustSum(t, newHash, []byte{0}), `H("") equals H("\x00")`)
	}
}

func (ht *HashTest) newWriteTest(r io.Reader, blockSize int, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		t.Run("return values", func(t *testing.T) {
			for _, size := range []int{0, 1, 5, blockSize - 1, blockSize, blockSize + 1, 8 * blockSize} {
				if size < 0 {
					continue
				}
				h := mustNewHash(t, newHash)

				var n int
				var err error
				require.NotPanicsf(t, func() { n, err = h.Write(mustRandomBytes(t, r, size)) }, "Write panicked for %d bytes", size)
				require.NoErrorf(t, err, "Write returned an error for %d bytes", size)
				require.Equalf(t, size, n, "Write returned a wrong count for %d bytes", size)
			}
		})

		t.Run("nil and empty", func(t *testing.T) {
			h := mustNewHash(t, newHash)

			n, err := h.Write(nil)
			require.NoError(t, err)
			require.Zero(t, n)

			n, err = h.Write([]byte{})
			require.NoError(t, err)
			require.Zero(t, n)
		})

		t.Run("does not modify src", func(t *testing.T) {
			for _, size := range []int{1, blockSize - 1, blockSize, blockSize + 1, 8 * blockSize} {
				if size < 0 {
					continue
				}
				src := mustRandomBytes(t, r, size)
				want := slices.Clone(src)

				h := mustNewHash(t, newHash)
				h.Write(src)
				h.Sum(nil)

				require.Equalf(t, want, src, "Write modified src, size=%d", size)
			}
		})

		t.Run("after Sum", func(t *testing.T) {
			head := mustRandomBytes(t, r, blockSize+3)
			tail := mustRandomBytes(t, r, 7)

			h := mustNewHash(t, newHash)
			h.Write(head)
			h.Sum(nil)
			h.Write(tail)

			require.Equal(t, mustSum(t, newHash, concat(head, tail)), h.Sum(nil), "intermediate Sum corrupted the state")
		})
	}
}

func (ht *HashTest) newSumTest(digestSize int, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		t.Run("length", func(t *testing.T) {
			for _, size := range []int{0, 1, 2, 3, 7, 8, 15, 16, 31, 32, 63, 64, 127, 128, 255, 1000} {
				h := mustNewHash(t, newHash)
				h.Write(bytes.Repeat([]byte{byte(size)}, size))
				require.Lenf(t, h.Sum(nil), digestSize, "digest length does not match Size, input=%d", size)
			}
		})

		t.Run("idempotent", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write([]byte("idempotent"))

			want := slices.Clone(h.Sum(nil))
			for round := range randomInputsPerRun {
				require.Equalf(t, want, h.Sum(nil), "Sum is not idempotent, round=%d", round)
			}
		})

		t.Run("append", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write([]byte("append"))
			want := slices.Clone(h.Sum(nil))

			prefix := []byte("PREFIX-0123456789")
			got := h.Sum(prefix)

			require.Len(t, got, len(prefix)+digestSize)
			require.Equal(t, prefix, got[:len(prefix)], "Sum overwrote the prefix")
			require.Equal(t, want, got[len(prefix):], "Sum appended a wrong digest")
		})

		t.Run("append with spare capacity", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write([]byte("spare capacity"))
			want := slices.Clone(h.Sum(nil))

			buf := bytes.Repeat([]byte{0xaa}, 4+digestSize+32)
			copy(buf, "HEAD")

			got := h.Sum(buf[:4])
			require.Len(t, got, 4+digestSize)
			require.Equal(t, []byte("HEAD"), got[:4], "Sum overwrote the existing content")
			require.Equal(t, want, got[4:], "Sum appended a wrong digest")
		})

		t.Run("nil and empty slice", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write([]byte("x"))
			require.Equal(t, h.Sum(nil), h.Sum([]byte{}), "Sum(nil) differs from Sum([]byte{})")
		})

		t.Run("not aliased", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write([]byte("alias"))

			got := h.Sum(nil)
			want := slices.Clone(got)
			for idx := range got {
				got[idx] ^= 0xff
			}
			require.Equal(t, want, h.Sum(nil), "Sum returned a slice aliasing the internal buffer")
		})
	}
}

func (ht *HashTest) newResetTest(r io.Reader, rand *rand.Rand, blockSize int, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		want := mustNewHash(t, newHash).Sum(nil)

		t.Run("after full use", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write(mustRandomBytes(t, r, 3*blockSize+5))
			h.Sum(nil)
			h.Reset()
			require.Equal(t, want, h.Sum(nil), "Reset did not restore the initial state")
		})

		t.Run("after partial write", func(t *testing.T) {
			for _, size := range []int{1, blockSize / 2, blockSize - 1, blockSize, blockSize + 1} {
				if size <= 0 {
					continue
				}
				h := mustNewHash(t, newHash)
				h.Write(mustRandomBytes(t, r, size))
				h.Reset()
				require.Equalf(t, want, h.Sum(nil), "Reset did not restore the initial state, size=%d", size)
			}
		})

		t.Run("idempotent", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write([]byte("data"))
			for range randomInputsPerRun {
				h.Reset()
			}
			require.Equal(t, want, h.Sum(nil), "repeated Reset changed the state")
		})

		t.Run("reuse", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			for round := range randomChunkTrials {
				src := mustRandomBytes(t, r, rand.IntN(4*blockSize+64))

				h.Reset()
				h.Write(src)
				require.Equalf(t, mustSum(t, newHash, src), h.Sum(nil), "reuse after Reset failed, round=%d", round)
			}
		})
	}
}

func (ht *HashTest) newStreamingTest(r io.Reader, rand *rand.Rand, blockSize int, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		src := mustRandomBytes(t, r, 5*blockSize+13)
		want := mustSum(t, newHash, src)

		t.Run("all splits", func(t *testing.T) {
			for split := range len(src) + 1 {
				h := mustNewHash(t, newHash)
				h.Write(src[:split])
				h.Write(src[split:])
				require.Equalf(t, want, h.Sum(nil), "split write mismatch, split=%d", split)
			}
		})

		t.Run("byte by byte", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			for idx := range src {
				h.Write(src[idx : idx+1])
			}
			require.Equal(t, want, h.Sum(nil), "byte-by-byte write mismatch")
		})

		t.Run("random chunks", func(t *testing.T) {
			for round := range randomChunkTrials {
				h := mustNewHash(t, newHash)
				for off := 0; off < len(src); {
					n := rand.IntN(len(src)-off) + 1
					h.Write(src[off : off+n])
					off += n
				}
				require.Equalf(t, want, h.Sum(nil), "random chunk write mismatch, round=%d", round)
			}
		})

		t.Run("large input", func(t *testing.T) {
			if testing.Short() {
				t.Skip("skipping large input in short mode")
			}

			large := mustRandomBytes(t, r, largeInputSize)
			wantLarge := mustSum(t, newHash, large)

			h := mustNewHash(t, newHash)
			for off := 0; off < len(large); off += 7919 {
				h.Write(large[off:min(off+7919, len(large))])
			}
			require.Equal(t, wantLarge, h.Sum(nil), "large split write mismatch")
		})
	}
}

func (ht *HashTest) newBoundaryTest(r io.Reader, blockSize int, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		sizes := make([]int, 0, 64)
		for _, blocks := range []int{0, 1, 2, 3, 4, 8} {
			for _, delta := range []int{-2, -1, 0, 1, 2} {
				sizes = append(sizes, blocks*blockSize+delta)
			}
		}

		sizes = append(sizes, blockSize-17, blockSize-16, blockSize-9, blockSize-8, 2*blockSize-9, 2*blockSize-8)

		slices.Sort(sizes)
		for _, size := range slices.Compact(sizes) {
			if size < 0 {
				continue
			}

			src := mustRandomBytes(t, r, size)
			want := mustSum(t, newHash, src)

			step := blockSize/2 + 1
			h := mustNewHash(t, newHash)
			for off := 0; off < size; off += step {
				h.Write(src[off:min(off+step, size)])
			}
			require.Equalf(t, want, h.Sum(nil), "boundary split write mismatch, size=%d", size)
		}
	}
}

func (ht *HashTest) newSensitivityTest(r io.Reader, rand *rand.Rand, blockSize int, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		t.Run("deterministic", func(t *testing.T) {
			for round := range randomChunkTrials {
				src := mustRandomBytes(t, r, rand.IntN(4*blockSize+64))
				sum1 := mustSum(t, newHash, src)
				sum2 := mustSum(t, newHash, src)
				require.Equalf(t, sum1, sum2, "hash is not deterministic, round=%d", round)
			}
		})

		t.Run("distinct inputs", func(t *testing.T) {
			seen := make(map[string]string)
			for range 2048 {
				src := mustRandomBytes(t, r, rand.IntN(200))

				digest := hex.EncodeToString(mustSum(t, newHash, src))
				input := hex.EncodeToString(src)
				if prev, ok := seen[digest]; ok {
					require.Equalf(t, prev, input, "distinct inputs produced the same digest: %s", digest)
				}
				seen[digest] = input
			}
		})

		t.Run("length", func(t *testing.T) {
			seen := make(map[string]int)
			for size := range 4*blockSize + 9 {
				digest := hex.EncodeToString(mustSum(t, newHash, make([]byte, size)))
				prev, ok := seen[digest]
				require.Falsef(t, ok, "zero inputs of length %d and %d produced the same digest", prev, size)
				seen[digest] = size
			}
		})

		t.Run("single bit", func(t *testing.T) {
			base := mustRandomBytes(t, r, 2*blockSize+3)
			want := mustSum(t, newHash, base)

			for idx := range base {
				for _, bit := range []byte{0x01, 0x80} {
					src := slices.Clone(base)
					src[idx] ^= bit
					require.NotEqualf(t, want, mustSum(t, newHash, src), "flipping bit %#x of byte %d did not change the digest", bit, idx)
				}
			}
		})
	}
}

func (ht *HashTest) newIOTest(newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		t.Run("io.Copy", func(t *testing.T) {
			src := bytes.Repeat([]byte("io.Copy"), 5000)

			h := mustNewHash(t, newHash)
			n, err := io.Copy(h, bytes.NewReader(src))
			require.NoError(t, err)
			require.EqualValues(t, len(src), n)
			require.Equal(t, mustSum(t, newHash, src), h.Sum(nil), "io.Copy mismatch")
		})

		t.Run("io.WriteString", func(t *testing.T) {
			const src = "io.WriteString"

			h := mustNewHash(t, newHash)
			n, err := io.WriteString(h, src)
			require.NoError(t, err)
			require.Equal(t, len(src), n)
			require.Equal(t, mustSum(t, newHash, []byte(src)), h.Sum(nil), "io.WriteString mismatch")
		})

		t.Run("all byte values", func(t *testing.T) {
			src := make([]byte, 256)
			for idx := range src {
				src[idx] = byte(idx)
			}
			want := mustSum(t, newHash, src)

			h := mustNewHash(t, newHash)
			for idx := range src {
				h.Write(src[idx : idx+1])
			}
			require.Equal(t, want, h.Sum(nil), "byte-by-byte write mismatch over all byte values")
		})
	}
}

func (ht *HashTest) newMarshalTest(r io.Reader, rand *rand.Rand, blockSize int, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		h := mustNewHash(t, newHash)
		_, marshaler := h.(encoding.BinaryMarshaler)
		_, appender := h.(encoding.BinaryAppender)
		_, unmarshaler := h.(encoding.BinaryUnmarshaler)

		if !marshaler && !appender && !unmarshaler {
			t.Skip("state marshaling is not implemented")
		}
		require.True(t, unmarshaler, "encoding.BinaryUnmarshaler must be implemented alongside the marshaler")
		require.True(t, marshaler || appender, "encoding.BinaryMarshaler or encoding.BinaryAppender must be implemented")

		if marshaler {
			t.Run("MarshalBinary", ht.newMarshalBinaryTest(r, blockSize, newHash))
		}
		if appender {
			t.Run("AppendBinary", ht.newAppendBinaryTest(r, blockSize, newHash))
		}
		t.Run("UnmarshalBinary", ht.newUnmarshalBinaryTest(r, rand, blockSize, newHash))
	}
}

func (ht *HashTest) newMarshalBinaryTest(r io.Reader, blockSize int, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		head := mustRandomBytes(t, r, 2*blockSize+5)
		tail := mustRandomBytes(t, r, 7)
		want := mustSum(t, newHash, concat(head, tail))

		h := mustNewHash(t, newHash)
		h.Write(head)

		state := mustMarshal(t, h)
		require.NotEmpty(t, state, "MarshalBinary returned an empty state")

		t.Run("does not change the state", func(t *testing.T) {
			h.Write(tail)
			require.Equal(t, want, h.Sum(nil), "MarshalBinary corrupted the hash state")
		})

		t.Run("deterministic", func(t *testing.T) {
			other := mustNewHash(t, newHash)
			other.Write(head)

			first := mustMarshal(t, other)
			second := mustMarshal(t, other)
			require.Equal(t, first, second, "MarshalBinary is not deterministic")
			require.Equal(t, state, first, "state differs between instances for the same input")
		})

		t.Run("not aliased", func(t *testing.T) {
			other := mustNewHash(t, newHash)
			other.Write(head)

			got := mustMarshal(t, other)
			for idx := range got {
				got[idx] ^= 0xff
			}
			require.Equal(t, state, mustMarshal(t, other), "MarshalBinary returned a slice aliasing the internal buffer")
		})

		t.Run("after Reset", func(t *testing.T) {
			other := mustNewHash(t, newHash)
			other.Write(head)
			other.Reset()
			require.Equal(t, mustMarshal(t, mustNewHash(t, newHash)), mustMarshal(t, other), "state after Reset differs from a fresh instance")
		})
	}
}

func (ht *HashTest) newAppendBinaryTest(r io.Reader, blockSize int, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		head := mustRandomBytes(t, r, 2*blockSize+5)
		tail := mustRandomBytes(t, r, 7)
		want := mustSum(t, newHash, concat(head, tail))

		h := mustNewHash(t, newHash)
		h.Write(head)

		appender, ok := h.(encoding.BinaryAppender)
		require.True(t, ok)

		var state []byte
		var err error
		require.NotPanics(t, func() { state, err = appender.AppendBinary(nil) }, "AppendBinary panicked")
		require.NoError(t, err)
		require.NotEmpty(t, state, "AppendBinary returned an empty state")

		t.Run("matches MarshalBinary", func(t *testing.T) {
			if _, ok := h.(encoding.BinaryMarshaler); !ok {
				t.Skip("encoding.BinaryMarshaler is not implemented")
			}
			require.Equal(t, mustMarshal(t, h), state, "AppendBinary differs from MarshalBinary")
		})

		t.Run("append", func(t *testing.T) {
			prefix := []byte("PREFIX-0123456789")

			got, err := appender.AppendBinary(prefix)
			require.NoError(t, err)
			require.Len(t, got, len(prefix)+len(state))
			require.Equal(t, prefix, got[:len(prefix)], "AppendBinary overwrote the prefix")
			require.Equal(t, state, got[len(prefix):], "AppendBinary appended a wrong state")
		})

		t.Run("append with spare capacity", func(t *testing.T) {
			buf := bytes.Repeat([]byte{0xaa}, 4+len(state)+32)
			copy(buf, "HEAD")

			got, err := appender.AppendBinary(buf[:4])
			require.NoError(t, err)
			require.Len(t, got, 4+len(state))
			require.Equal(t, []byte("HEAD"), got[:4], "AppendBinary overwrote the existing content")
			require.Equal(t, state, got[4:], "AppendBinary appended a wrong state")
		})

		t.Run("does not change the state", func(t *testing.T) {
			h.Write(tail)
			require.Equal(t, want, h.Sum(nil), "AppendBinary corrupted the hash state")
		})

		t.Run("round trip", func(t *testing.T) {
			other := mustNewHash(t, newHash)
			mustUnmarshal(t, other, state)
			other.Write(tail)
			require.Equal(t, want, other.Sum(nil), "restoring an AppendBinary state failed")
		})

		t.Run("empty hash", func(t *testing.T) {
			fresh, ok := mustNewHash(t, newHash).(encoding.BinaryAppender)
			require.True(t, ok)

			_, err := fresh.AppendBinary(nil)
			require.NoError(t, err, "AppendBinary failed on a fresh hash")
		})
	}
}

func (ht *HashTest) newUnmarshalBinaryTest(r io.Reader, rand *rand.Rand, blockSize int, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		tail := []byte("--continued--")

		t.Run("write after restore does not panic", func(t *testing.T) {
			src := mustRandomBytes(t, r, blockSize+7)

			h := mustNewHash(t, newHash)
			h.Write(src)
			state := mustMarshal(t, h)

			restored := mustNewHash(t, newHash)
			mustUnmarshal(t, restored, state)

			var (
				got []byte
				n   int
				err error
			)
			require.NotPanics(t, func() {
				n, err = restored.Write(tail)
				got = restored.Sum(nil)
			}, "Write or Sum panicked after UnmarshalBinary")
			require.NoError(t, err)
			require.Equal(t, len(tail), n)
			require.Equal(t, mustSum(t, newHash, concat(src, tail)), got)
		})

		t.Run("offsets", func(t *testing.T) {
			src := mustRandomBytes(t, r, 4*blockSize+16)

			offsets := []int{0, 1, blockSize - 1, blockSize, blockSize + 1, 2*blockSize - 1, 2 * blockSize, 3*blockSize + 5}
			for _, off := range offsets {
				if off < 0 || off > len(src) {
					continue
				}

				h := mustNewHash(t, newHash)
				h.Write(src[:off])
				state := mustMarshal(t, h)

				other := mustNewHash(t, newHash)
				mustUnmarshal(t, other, state)
				require.Equalf(t, mustSum(t, newHash, src[:off]), other.Sum(nil), "digest mismatch right after restoring, off=%d", off)

				other.Write(tail)
				require.Equalf(t, mustSum(t, newHash, concat(src[:off], tail)), other.Sum(nil), "digest mismatch after restoring and writing, off=%d", off)
			}
		})

		t.Run("overwrites dirty state", func(t *testing.T) {
			src := mustRandomBytes(t, r, blockSize+3)

			h := mustNewHash(t, newHash)
			h.Write(src)
			state := mustMarshal(t, h)

			dst := mustNewHash(t, newHash)
			dst.Write(mustRandomBytes(t, r, 3*blockSize+7))
			mustUnmarshal(t, dst, state)

			require.Equal(t, mustSum(t, newHash, src), dst.Sum(nil), "UnmarshalBinary did not fully overwrite the previous state")
		})

		t.Run("does not retain src", func(t *testing.T) {
			src := mustRandomBytes(t, r, blockSize+3)
			want := mustSum(t, newHash, src)

			h := mustNewHash(t, newHash)
			h.Write(src)
			state := mustMarshal(t, h)

			dst := mustNewHash(t, newHash)
			mustUnmarshal(t, dst, state)
			for idx := range state {
				state[idx] ^= 0xff
			}
			require.Equal(t, want, dst.Sum(nil), "UnmarshalBinary retained the caller's mutable state slice")
		})

		t.Run("canonical", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write(mustRandomBytes(t, r, 2*blockSize+11))
			state := mustMarshal(t, h)

			other := mustNewHash(t, newHash)
			mustUnmarshal(t, other, state)
			require.Equal(t, state, mustMarshal(t, other), "re-marshaling a restored state produced different bytes")
		})

		t.Run("reset after restore", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write(mustRandomBytes(t, r, blockSize+3))
			state := mustMarshal(t, h)

			other := mustNewHash(t, newHash)
			mustUnmarshal(t, other, state)
			other.Reset()
			require.Equal(t, mustNewHash(t, newHash).Sum(nil), other.Sum(nil), "Reset after restoring did not restore the initial state")
		})

		t.Run("rejects truncated state", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write(mustRandomBytes(t, r, blockSize+3))
			state := mustMarshal(t, h)

			for size := range state {
				u, ok := mustNewHash(t, newHash).(encoding.BinaryUnmarshaler)
				require.True(t, ok)

				var err error
				require.NotPanicsf(t, func() { err = u.UnmarshalBinary(state[:size]) }, "UnmarshalBinary panicked for truncated state, size=%d", size)
				require.Errorf(t, err, "UnmarshalBinary accepted a truncated state, size=%d/%d", size, len(state))
			}
		})

		t.Run("rejects corrupt state", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write(mustRandomBytes(t, r, blockSize+3))
			state := mustMarshal(t, h)

			corrupt := [][]byte{
				nil,
				{},
				{0x00},
				bytes.Repeat([]byte{0x00}, len(state)),
				bytes.Repeat([]byte{0xff}, len(state)),
				append(slices.Clone(state), 0xff),
			}
			for idx := range min(8, len(state)) {
				bad := slices.Clone(state)
				bad[idx] ^= 0xff
				corrupt = append(corrupt, bad)
			}

			for idx, bad := range corrupt {
				u, ok := mustNewHash(t, newHash).(encoding.BinaryUnmarshaler)
				require.True(t, ok)
				require.NotPanicsf(t, func() { _ = u.UnmarshalBinary(bad) }, "UnmarshalBinary panicked for corrupt state #%d", idx)
			}
		})

		t.Run("accepted mutations remain usable", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			h.Write(mustRandomBytes(t, r, blockSize+7))
			state := mustMarshal(t, h)

			for idx := range state {
				bad := slices.Clone(state)
				bad[idx] ^= 0xff

				candidate := mustNewHash(t, newHash)
				digestSize := candidate.Size()
				u, ok := candidate.(encoding.BinaryUnmarshaler)
				require.True(t, ok)

				var unmarshalErr error
				require.NotPanicsf(t, func() {
					unmarshalErr = u.UnmarshalBinary(bad)
				}, "UnmarshalBinary panicked for mutation at byte %d", idx)
				if unmarshalErr != nil {
					continue
				}

				var (
					digest   []byte
					n        int
					writeErr error
				)
				require.NotPanicsf(t, func() {
					n, writeErr = candidate.Write(tail)
					digest = candidate.Sum(nil)
				}, "accepted mutation at byte %d left an unsafe hash state", idx)
				require.NoError(t, writeErr)
				require.Equal(t, len(tail), n)
				require.Equal(t, digestSize, candidate.Size())
				require.Len(t, digest, digestSize)
			}
		})

		t.Run("usable after failure", func(t *testing.T) {
			h := mustNewHash(t, newHash)
			u, ok := h.(encoding.BinaryUnmarshaler)
			require.True(t, ok)

			if err := u.UnmarshalBinary([]byte{0x00}); err == nil {
				t.Skip("UnmarshalBinary accepts an invalid state")
			}

			src := []byte("recovered")
			h.Reset()
			h.Write(src)
			require.Equal(t, mustSum(t, newHash, src), h.Sum(nil), "hash is unusable after a failed UnmarshalBinary")
		})

		t.Run("round trip", func(t *testing.T) {
			for round := range randomChunkTrials {
				src := mustRandomBytes(t, r, rand.IntN(4*blockSize+64))
				cut := rand.IntN(len(src) + 1)

				h := mustNewHash(t, newHash)
				h.Write(src[:cut])
				state := mustMarshal(t, h)

				other := mustNewHash(t, newHash)
				mustUnmarshal(t, other, state)
				other.Write(src[cut:])

				require.Equalf(t, mustSum(t, newHash, src), other.Sum(nil), "round trip failed, round=%d len=%d cut=%d", round, len(src), cut)
			}
		})
	}
}

func (ht *HashTest) newConcurrentUseTest(r io.Reader, blockSize int, newHash func() hash.Hash) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		input := mustRandomBytes(t, r, 2*blockSize+7)
		want := mustSum(t, newHash, input)

		errs := make(chan error, concurrentWorkers)

		var wg sync.WaitGroup
		for worker := range concurrentWorkers {
			wg.Go(func() {
				h := newHash()
				for round := range concurrentRounds {
					h.Reset()
					h.Write(input)
					if !bytes.Equal(want, h.Sum(nil)) {
						errs <- fmt.Errorf("worker %d round %d: concurrent Sum mismatch", worker, round)
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

func (ht *HashTest) newSumFuncTest(r io.Reader, rand *rand.Rand, digestSize int) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()

		newHash := ht.getNew(t, digestSize)
		sum := ht.getSum(t, digestSize)
		blockSize := ht.HashInfo.BlockSize

		t.Run("empty", func(t *testing.T) {
			for _, src := range [][]byte{nil, {}} {
				got := sum(src)
				require.Len(t, got, digestSize, "digest length does not match the declared size")
				require.Equal(t, mustNewHash(t, newHash).Sum(nil), got, "Sum differs from the streaming implementation for an empty input")
			}
		})

		t.Run("matches New", func(t *testing.T) {
			sizes := make([]int, 0, 64)
			for _, blocks := range []int{0, 1, 2, 3, 4, 8} {
				for _, delta := range []int{-2, -1, 0, 1, 2} {
					sizes = append(sizes, blocks*blockSize+delta)
				}
			}
			sizes = append(sizes, blockSize-17, blockSize-16, blockSize-9, blockSize-8, 2*blockSize-9, 2*blockSize-8)

			slices.Sort(sizes)
			for _, size := range slices.Compact(sizes) {
				if size < 0 {
					continue
				}

				src := mustRandomBytes(t, r, size)
				got := sum(src)
				require.Lenf(t, got, digestSize, "digest length does not match the declared size, size=%d", size)
				require.Equalf(t, mustSum(t, newHash, src), got, "Sum differs from the streaming implementation, size=%d", size)
			}
		})

		t.Run("random", func(t *testing.T) {
			for round := range randomChunkTrials {
				src := mustRandomBytes(t, r, rand.IntN(4*blockSize+64))
				require.Equalf(t, mustSum(t, newHash, src), sum(src), "Sum differs from the streaming implementation, round=%d", round)
			}
		})

		t.Run("deterministic", func(t *testing.T) {
			src := mustRandomBytes(t, r, 2*blockSize+7)
			want := slices.Clone(sum(src))
			for round := range randomInputsPerRun {
				require.Equalf(t, want, sum(src), "Sum is not deterministic, round=%d", round)
			}
		})

		t.Run("does not modify src", func(t *testing.T) {
			for _, size := range []int{1, blockSize - 1, blockSize, blockSize + 1, 8 * blockSize} {
				if size < 0 {
					continue
				}
				src := mustRandomBytes(t, r, size)
				want := slices.Clone(src)

				sum(src)
				require.Equalf(t, want, src, "Sum modified src, size=%d", size)
			}
		})

		t.Run("not aliased", func(t *testing.T) {
			src := mustRandomBytes(t, r, 2*blockSize+7)

			first := sum(src)
			want := slices.Clone(first)
			for idx := range first {
				first[idx] ^= 0xff
			}
			require.Equal(t, want, sum(src), "Sum returned a slice aliasing a shared buffer")
		})

		t.Run("large input", func(t *testing.T) {
			if testing.Short() {
				t.Skip("skipping large input in short mode")
			}

			src := mustRandomBytes(t, r, largeInputSize)
			require.Equal(t, mustSum(t, newHash, src), sum(src), "Sum differs from the streaming implementation for a large input")
		})
	}
}

func mustNewHash(t testing.TB, newHash func() hash.Hash) hash.Hash {
	t.Helper()

	var h hash.Hash
	require.NotPanics(t, func() { h = newHash() }, "constructor panicked")
	require.NotNil(t, h, "constructor returned a nil hash")
	require.Positive(t, h.Size(), "Size must be positive")
	require.Positive(t, h.BlockSize(), "BlockSize must be positive")
	return h
}

func mustSum(t testing.TB, newHash func() hash.Hash, src []byte) []byte {
	t.Helper()

	h := mustNewHash(t, newHash)
	n, err := h.Write(src)
	require.NoError(t, err)
	require.Equal(t, len(src), n)
	return h.Sum(nil)
}

func mustMarshal(t testing.TB, h hash.Hash) []byte {
	t.Helper()

	switch v := h.(type) {
	case encoding.BinaryMarshaler:
		var state []byte
		var err error
		require.NotPanics(t, func() { state, err = v.MarshalBinary() }, "MarshalBinary panicked")
		require.NoError(t, err)
		return state

	case encoding.BinaryAppender:
		var state []byte
		var err error
		require.NotPanics(t, func() { state, err = v.AppendBinary(nil) }, "AppendBinary panicked")
		require.NoError(t, err)
		return state
	}

	require.FailNow(t, "state marshaling is not implemented")
	return nil
}

func mustUnmarshal(t testing.TB, h hash.Hash, state []byte) {
	t.Helper()

	u, ok := h.(encoding.BinaryUnmarshaler)
	require.True(t, ok, "encoding.BinaryUnmarshaler is not implemented")

	var err error
	require.NotPanics(t, func() { err = u.UnmarshalBinary(state) }, "UnmarshalBinary panicked")
	require.NoError(t, err)
}

func mustRandomBytes(t testing.TB, reader io.Reader, size int) []byte {
	t.Helper()

	buf := make([]byte, size)
	_, err := io.ReadFull(reader, buf)
	require.NoError(t, err)
	return buf
}

func concat(parts ...[]byte) []byte {
	var out []byte
	for _, part := range parts {
		out = append(out, part...)
	}
	return out
}
