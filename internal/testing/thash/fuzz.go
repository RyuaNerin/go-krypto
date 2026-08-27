package thash

import (
	"bytes"
	"encoding"
	"hash"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	fuzzMinBlockSize = 8
	fuzzMaxBlockSize = 256 // LSH-512 는 256 바이트 블록을 사용한다

	fuzzMinDigestSize = 2
	fuzzMaxDigestSize = 64

	fuzzBufSize  = 1024
	fuzzGuardLen = 32
)

// newValid builds a hash and checks the properties every iteration relies on.
func (ht *HashTest) newValid(t *testing.T, digestSize int) (hash.Hash, func() hash.Hash, int) {
	t.Helper()

	newHash := ht.getNew(t, digestSize)

	h := newHash()
	require.NotNil(t, h)

	bs, ds := h.BlockSize(), h.Size()
	require.GreaterOrEqual(t, bs, fuzzMinBlockSize, "BlockSize is too small")
	require.LessOrEqual(t, bs, fuzzMaxBlockSize, "BlockSize is too large")
	require.GreaterOrEqual(t, ds, fuzzMinDigestSize, "Size is too small")
	require.LessOrEqual(t, ds, fuzzMaxDigestSize, "Size is too large")
	require.Equal(t, digestSize, ds, "Size() does not match the HashTest.New key")

	if ht.HashInfo.BlockSize != 0 {
		require.Equal(t, ht.HashInfo.BlockSize, bs, "BlockSize does not match HashInfo.BlockSize")
	}
	return h, newHash, bs
}

// sumOf hashes src with a fresh instance.
func sumOf(newHash func() hash.Hash, src []byte) []byte {
	h := newHash()
	h.Write(src)
	return h.Sum(nil)
}

// guardedAppend returns an empty slice with spare capacity carved out of a
// larger buffer with zero guards on BOTH sides, plus a func asserting that
// neither guard was touched. Sum must write exactly Size() bytes into it.
func guardedAppend(t *testing.T, n int) ([]byte, func(what string)) {
	raw := make([]byte, fuzzGuardLen+n+fuzzGuardLen)
	zero := make([]byte, fuzzGuardLen)
	return raw[fuzzGuardLen:fuzzGuardLen:len(raw)], func(what string) {
		t.Helper()
		require.Equalf(t, zero, raw[:fuzzGuardLen], "%s wrote before dst", what)
		require.Equalf(t, zero, raw[fuzzGuardLen+n:], "%s wrote past dst", what)
	}
}

// marshalerOf returns a func producing the serialized state, or nil when the
// hash does not support state marshaling.
func marshalerOf(h hash.Hash) func() ([]byte, error) {
	switch v := h.(type) {
	case encoding.BinaryMarshaler:
		return v.MarshalBinary
	case encoding.BinaryAppender:
		return func() ([]byte, error) { return v.AppendBinary(nil) }
	}
	return nil
}

// anyNew returns any constructor, for probing optional interfaces before the
// fuzz body runs.
func (ht *HashTest) anyNew(f *testing.F) func() hash.Hash {
	f.Helper()

	for _, fn := range ht.New {
		return fn
	}

	f.Fatal("no constructor found")
	return nil
}

func (ht *HashTest) addHashSeeds(f *testing.F) {
	f.Add([]byte(nil), 1)
	f.Add([]byte("x"), 1)
	f.Add([]byte("The quick brown fox jumps over the lazy dog"), 7)
	f.Add(make([]byte, fuzzBufSize), 3)
	f.Add(bytes.Repeat([]byte{0xff}, fuzzBufSize), fuzzBufSize)
	f.Add(bytes.Repeat([]byte{0xa5}, fuzzMaxBlockSize), fuzzMaxBlockSize)
	f.Add(bytes.Repeat([]byte{0x5a}, fuzzMaxBlockSize+1), fuzzMaxBlockSize-1)
}

// FuzzNew exercises the constructor contract and input ownership.
func (ht *HashTest) FuzzNew(f *testing.F, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}

	f.Add([]byte(nil))
	f.Add([]byte("short"))
	f.Add(make([]byte, fuzzBufSize))
	f.Add(bytes.Repeat([]byte{0xff}, fuzzMaxBlockSize+1))

	f.Fuzz(func(t *testing.T, src []byte) {
		for _, size := range ht.Sizes {
			h, newHash, bs := ht.newValid(t, size.DigestSize)

			// Every instance must agree on the geometry.
			other, _, obs := ht.newValid(t, size.DigestSize)
			require.Equal(t, bs, obs, "BlockSize is not stable across instances")
			require.Equal(t, other.Sum(nil), newHash().Sum(nil), "the empty digest is not deterministic")

			orig := slices.Clone(src)
			want := sumOf(newHash, orig)

			n, err := h.Write(src)
			require.NoError(t, err, "Write returned an error")
			require.Equal(t, len(src), n, "Write returned a wrong count")
			require.Equal(t, orig, src, "Write mutated the caller's src")

			require.Equal(t, bs, h.BlockSize(), "BlockSize changed after Write")
			require.Equal(t, size.DigestSize, h.Size(), "Size changed after Write")

			// The internal state must not alias the caller's slice.
			for i := range src {
				src[i] ^= 0xff
			}
			got := h.Sum(nil)
			require.Len(t, got, size.DigestSize, "digest length does not match Size")
			require.Equal(t, want, got, "the hash retains a reference to the caller's src")

			copy(src, orig)
		}
	})
}

// FuzzHash covers streaming and the append semantics of Sum.
func (ht *HashTest) FuzzHash(f *testing.F, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}

	ht.addHashSeeds(f)

	f.Fuzz(func(t *testing.T, src []byte, chunk int) {
		chunk = min(max(chunk, 1), fuzzBufSize)

		for _, size := range ht.Sizes {
			h, newHash, _ := ht.newValid(t, size.DigestSize)
			want := slices.Clone(src)

			// One-shot reference.
			h.Write(src)
			digest := h.Sum(nil)
			require.Len(t, digest, size.DigestSize, "digest length does not match Size")
			require.Equal(t, want, src, "Write mutated src")

			// Sum must be idempotent and must not advance the state.
			require.Equal(t, digest, h.Sum(nil), "Sum is not idempotent")
			require.Equal(t, digest, h.Sum([]byte{}), "Sum(nil) differs from Sum([]byte{})")

			// Determinism.
			require.Equal(t, digest, sumOf(newHash, want), "the hash is not deterministic")

			// The returned slice must not alias the internal buffer.
			alias := h.Sum(nil)
			for i := range alias {
				alias[i] ^= 0xff
			}
			require.Equal(t, digest, h.Sum(nil), "Sum returned a slice aliasing the internal buffer")

			// Append semantics, and exactly Size() bytes may be written.
			g, check := guardedAppend(t, size.DigestSize)
			g = h.Sum(g)
			require.Equal(t, digest, g, "Sum wrote a wrong digest into a spare-capacity dst")
			check("Sum")

			// Chunked writes must match the one-shot result.
			chunked := newHash()
			for off := 0; off < len(want); off += chunk {
				chunked.Write(want[off:min(off+chunk, len(want))])
			}
			require.Equalf(t, digest, chunked.Sum(nil), "chunked write differs, chunk=%d", chunk)

			// Interleaved Sum calls must not disturb the running state.
			mixed := newHash()
			for off := 0; off < len(want); off += chunk {
				mixed.Sum(nil)
				mixed.Write(want[off:min(off+chunk, len(want))])
				mixed.Sum(nil)
			}
			require.Equalf(t, digest, mixed.Sum(nil), "interleaved Sum corrupted the state, chunk=%d", chunk)

			// Empty writes are no-ops in any position.
			empty := newHash()
			empty.Write(nil)
			empty.Write(want)
			empty.Write([]byte{})
			require.Equal(t, digest, empty.Sum(nil), "an empty Write changed the state")

			// Reset must return the instance to the initial state, and it must
			// be reusable afterwards.
			h.Reset()
			require.Equal(t, newHash().Sum(nil), h.Sum(nil), "Reset did not restore the initial state")
			h.Write(want)
			require.Equal(t, digest, h.Sum(nil), "reuse after Reset differs")
		}
	})
}

// FuzzSum compares the one-shot Sum functions against the streaming ones.
func (ht *HashTest) FuzzSum(f *testing.F, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}
	if len(ht.Sum) == 0 {
		f.Skip("no sum function")
	}

	f.Add([]byte(nil))
	f.Add([]byte("x"))
	f.Add([]byte("The quick brown fox jumps over the lazy dog"))
	f.Add(make([]byte, fuzzMaxBlockSize))
	f.Add(bytes.Repeat([]byte{0xff}, fuzzBufSize))

	f.Fuzz(func(t *testing.T, src []byte) {
		for _, size := range ht.Sizes {
			if _, ok := ht.Sum[size.DigestSize]; !ok {
				continue
			}

			_, newHash, _ := ht.newValid(t, size.DigestSize)
			sum := ht.getSum(t, size.DigestSize)

			orig := slices.Clone(src)
			want := sumOf(newHash, orig)

			got := sum(src)
			require.Len(t, got, size.DigestSize, "digest length does not match the declared size")
			require.Equal(t, want, got, "Sum differs from the streaming implementation")
			require.Equal(t, orig, src, "Sum mutated the caller's src")

			// Deterministic and not aliasing a shared buffer.
			again := sum(src)
			require.Equal(t, want, again, "Sum is not deterministic")
			for i := range again {
				again[i] ^= 0xff
			}
			require.Equal(t, want, sum(src), "Sum returned a slice aliasing a shared buffer")
		}
	})
}

// FuzzMarshal covers snapshot/restore round trips at arbitrary cut points.
func (ht *HashTest) FuzzMarshal(f *testing.F, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}

	// Without state marshaling the body below is a silent no-op, which would
	// report as a pass while verifying nothing.
	probe := ht.anyNew(f)()
	require.NotNil(f, probe)
	if marshalerOf(probe) == nil {
		f.Skip("hash implements neither encoding.BinaryMarshaler nor encoding.BinaryAppender")
	}
	_, ok := probe.(encoding.BinaryUnmarshaler)
	require.True(f, ok, "encoding.BinaryUnmarshaler must be implemented alongside the marshaler")

	ht.addHashSeeds(f)

	f.Fuzz(func(t *testing.T, src []byte, cut int) {
		cut = min(max(cut, 0), len(src))

		for _, size := range ht.Sizes {
			h, newHash, _ := ht.newValid(t, size.DigestSize)
			want := sumOf(newHash, src)

			marshal := marshalerOf(h)
			require.NotNil(t, marshal, "state marshaling varies between instances of the same hash")

			h.Write(src[:cut])

			state, err := marshal()
			require.NoError(t, err, "marshaling the state failed")
			require.NotEmpty(t, state, "marshaling returned an empty state")

			// Marshaling must not disturb the running state.
			h.Write(src[cut:])
			require.Equalf(t, want, h.Sum(nil), "marshaling corrupted the state, cut=%d", cut)

			// Marshaling is deterministic and does not alias the internal buffer.
			other := newHash()
			other.Write(src[:cut])
			again, err := marshalerOf(other)()
			require.NoError(t, err)
			require.Equal(t, state, again, "the state differs between instances for the same input")

			for i := range again {
				again[i] ^= 0xff
			}
			third, err := marshalerOf(other)()
			require.NoError(t, err)
			require.Equal(t, state, third, "marshaling returned a slice aliasing the internal buffer")

			// Restore, then finish the input.
			restored := newHash()
			u, ok := restored.(encoding.BinaryUnmarshaler)
			require.True(t, ok)
			require.NoErrorf(t, u.UnmarshalBinary(state), "restoring a valid state failed, cut=%d", cut)
			require.Equalf(t, sumOf(newHash, src[:cut]), restored.Sum(nil), "digest mismatch right after restoring, cut=%d", cut)

			var (
				got      []byte
				n        int
				writeErr error
			)
			require.NotPanicsf(t, func() {
				n, writeErr = restored.Write(src[cut:])
				got = restored.Sum(nil)
			}, "Write or Sum panicked after restoring, cut=%d", cut)
			require.NoError(t, writeErr)
			require.Equal(t, len(src)-cut, n)
			require.Equalf(t, want, got, "round trip failed, cut=%d", cut)

			// Re-marshaling a restored state must be canonical, and restoring
			// must not retain the caller's slice.
			canonical := newHash()
			require.NoError(t, canonical.(encoding.BinaryUnmarshaler).UnmarshalBinary(state))
			remarshaled, err := marshalerOf(canonical)()
			require.NoError(t, err)
			require.Equal(t, state, remarshaled, "re-marshaling a restored state produced different bytes")

			for i := range state {
				state[i] ^= 0xff
			}
			require.Equal(t, sumOf(newHash, src[:cut]), canonical.Sum(nil), "the hash retains a reference to the caller's state")

			canonical.Reset()
			require.Equal(t, newHash().Sum(nil), canonical.Sum(nil), "Reset after restoring did not restore the initial state")
		}
	})
}

// FuzzUnmarshal feeds arbitrary bytes to UnmarshalBinary: it may reject
// them, but it must never panic and must never leave a broken instance behind.
func (ht *HashTest) FuzzUnmarshal(f *testing.F, skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}

	probeNew := ht.anyNew(f)
	probe := probeNew()
	require.NotNil(f, probe)
	if marshalerOf(probe) == nil {
		f.Skip("hash implements neither encoding.BinaryMarshaler nor encoding.BinaryAppender")
	}
	if _, ok := probe.(encoding.BinaryUnmarshaler); !ok {
		f.Skip("hash does not implement encoding.BinaryUnmarshaler")
	}

	// Seed with real states so the fuzzer has something to mutate.
	for _, size := range []int{0, 1, fuzzMaxBlockSize - 1, fuzzMaxBlockSize, fuzzMaxBlockSize + 1} {
		h := probeNew()
		h.Write(bytes.Repeat([]byte{0x5a}, size))
		if state, err := marshalerOf(h)(); err == nil {
			f.Add(state)
			f.Add(state[:len(state)/2])
			f.Add(append(slices.Clone(state), 0xff))
		}
	}
	f.Add([]byte(nil))
	f.Add([]byte{0x00})
	f.Add(bytes.Repeat([]byte{0xff}, fuzzBufSize))

	f.Fuzz(func(t *testing.T, state []byte) {
		if len(state) > fuzzBufSize {
			state = state[:fuzzBufSize]
		}

		for _, size := range ht.Sizes {
			h, newHash, _ := ht.newValid(t, size.DigestSize)

			u, ok := h.(encoding.BinaryUnmarshaler)
			require.True(t, ok)

			orig := slices.Clone(state)

			var err error
			require.NotPanicsf(t, func() { err = u.UnmarshalBinary(state) }, "UnmarshalBinary panicked, len=%d", len(state))
			require.Equal(t, orig, state, "UnmarshalBinary mutated the caller's state")

			if err != nil {
				// A rejected state must leave a usable instance behind.
				h.Reset()
				h.Write(orig)
				require.Equal(t, sumOf(newHash, orig), h.Sum(nil), "the hash is unusable after a rejected state")
				continue
			}

			// An accepted state must behave like any other live instance.
			var digest []byte
			require.NotPanics(t, func() {
				digest = h.Sum(nil)
			}, "Sum panicked after accepting a state")
			require.Len(t, digest, size.DigestSize, "digest length does not match Size after restoring")
			require.Equal(t, digest, h.Sum(nil), "Sum is not idempotent after restoring")

			remarshaled, mErr := marshalerOf(h)()
			require.NoError(t, mErr, "marshaling a restored state failed")
			require.Equal(t, orig, remarshaled, "an accepted state is not canonical")

			var (
				afterWrite []byte
				n          int
				writeErr   error
			)
			require.NotPanics(t, func() {
				n, writeErr = h.Write(orig)
				afterWrite = h.Sum(nil)
			}, "Write or Sum panicked after accepting a state")
			require.NoError(t, writeErr)
			require.Equal(t, len(orig), n)
			require.NotEmpty(t, afterWrite)

			h.Reset()
			require.Equal(t, newHash().Sum(nil), h.Sum(nil), "Reset after restoring did not restore the initial state")
		}
	})
}
