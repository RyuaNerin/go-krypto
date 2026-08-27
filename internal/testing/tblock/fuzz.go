package tblock

import (
	"bytes"
	"crypto/cipher"
	"slices"
	"testing"

	"github.com/RyuaNerin/go-krypto"

	"github.com/stretchr/testify/require"
)

const (
	fuzzMinBlockSize = 8
	fuzzMaxBlockSize = 32

	fuzzMaxWide  = 8                              // Encrypt8
	fuzzBufSize  = fuzzMaxBlockSize * fuzzMaxWide // 256
	fuzzGuardLen = fuzzMaxBlockSize
)

type wideFn struct {
	n   int
	enc func(dst, src []byte)
	dec func(dst, src []byte)
}

func wideFns(b cipher.Block) []wideFn {
	var out []wideFn
	if w, ok := b.(krypto.Block4); ok {
		out = append(out, wideFn{4, w.Encrypt4, w.Decrypt4})
	}
	if w, ok := b.(krypto.Block8); ok {
		out = append(out, wideFn{8, w.Encrypt8, w.Decrypt8})
	}
	return out
}

func wideName(n int, dec bool) string {
	op := "Encrypt"
	if dec {
		op = "Decrypt"
	}
	if n == 1 {
		return op
	}
	return op + string(rune('0'+n))
}

func refN(b cipher.Block, n int, src []byte, one func(dst, src []byte)) []byte {
	bs := b.BlockSize()
	dst := make([]byte, bs*n)
	for i := range n {
		one(dst[i*bs:(i+1)*bs], src[i*bs:(i+1)*bs])
	}
	return dst
}

func fitTo(buf []byte, need int) []byte {
	if len(buf) == 0 || need <= 0 {
		return nil
	}
	if len(buf) < need {
		buf = bytes.Repeat(buf, need/len(buf)+1)
	}
	return slices.Clone(buf[:need])
}

// validKeySize reports whether n is a documented key length for this cipher.
func (ci *CipherInfo) validKeySize(n int) bool {
	return slices.ContainsFunc(ci.Sizes, func(s CipherSize) bool { return s.KeySize == n })
}

// normKey turns arbitrary fuzz input into a key of a documented length so that
// every iteration reaches the cipher logic instead of being skipped.
func (ci *CipherInfo) normKey(key []byte) []byte {
	if len(key) == 0 || len(ci.Sizes) == 0 {
		return nil
	}
	return fitTo(key, ci.Sizes[int(key[0])%len(ci.Sizes)].KeySize)
}

// newValid builds a cipher from a normalized key, or skips the iteration.
func (ci *CipherInfo) newValid(t *testing.T, newCipher func([]byte) (cipher.Block, error), key []byte) (cipher.Block, int) {
	t.Helper()
	key = ci.normKey(key)
	if key == nil {
		t.Skip()
	}
	b, err := newCipher(key)
	require.NoErrorf(t, err, "rejected a valid %d-byte key", len(key))
	require.NotNil(t, b)

	bs := b.BlockSize()
	require.GreaterOrEqual(t, bs, fuzzMinBlockSize)
	require.LessOrEqual(t, bs, fuzzMaxBlockSize)
	return b, bs
}

// guarded returns an n-byte working buffer surrounded by zero guards on BOTH
// sides, plus a func asserting neither guard was touched.
func guarded(t *testing.T, n int) ([]byte, func(what string)) {
	raw := make([]byte, fuzzGuardLen+n+fuzzGuardLen)
	zero := make([]byte, fuzzGuardLen)
	return raw[fuzzGuardLen : fuzzGuardLen+n], func(what string) {
		t.Helper()
		require.Equalf(t, zero, raw[:fuzzGuardLen], "%s wrote before dst", what)
		require.Equalf(t, zero, raw[fuzzGuardLen+n:], "%s wrote past dst", what)
	}
}

func (ci *CipherInfo) addBlockSeeds(f *testing.F) {
	for _, ks := range ci.Sizes {
		f.Add(bytes.Repeat([]byte{0x00}, ks.KeySize), make([]byte, fuzzBufSize))
		f.Add(bytes.Repeat([]byte{0xff}, ks.KeySize), bytes.Repeat([]byte{0xa5}, fuzzBufSize))
		f.Add(bytes.Repeat([]byte{0x5a}, ks.KeySize), []byte("The quick brown fox jumps over the lazy dog"))
	}
	f.Add([]byte{0x01}, []byte("x"))
}

// FuzzNew exercises the constructor contract with arbitrary key lengths.
// Only the exported NewCipher should be passed here; internal per-implementation
// constructors assume a pre-validated key.
func (ci *CipherInfo) FuzzNew(f *testing.F, newCipher func(key []byte) (cipher.Block, error), skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}

	for _, ks := range ci.Sizes {
		f.Add(bytes.Repeat([]byte{0x00}, ks.KeySize))
		f.Add(bytes.Repeat([]byte{0xff}, ks.KeySize-1))
		f.Add(bytes.Repeat([]byte{0xff}, ks.KeySize+1))
	}
	f.Add([]byte(nil))
	f.Add([]byte("short"))
	f.Add(make([]byte, fuzzBufSize))

	f.Fuzz(func(t *testing.T, key []byte) {
		orig := slices.Clone(key)

		b, err := newCipher(key)
		if err != nil {
			require.Falsef(t, ci.validKeySize(len(key)), "rejected a valid %d-byte key", len(key))
			require.Nil(t, b, "returned a non-nil Block alongside an error")
			require.Equal(t, orig, key, "mutated the key before failing")
			return
		}
		require.Truef(t, ci.validKeySize(len(key)), "accepted an invalid %d-byte key", len(key))
		require.NotNil(t, b)
		require.Equal(t, orig, key, "mutated the caller's key")

		bs := b.BlockSize()
		require.GreaterOrEqual(t, bs, fuzzMinBlockSize)
		require.LessOrEqual(t, bs, fuzzMaxBlockSize)

		// The key schedule must not alias the caller's slice.
		pt := bytes.Repeat([]byte{0x42}, bs)
		want := make([]byte, bs)
		b.Encrypt(want, pt)

		for i := range key {
			key[i] ^= 0xff
		}
		got := make([]byte, bs)
		b.Encrypt(got, pt)
		require.Equal(t, want, got, "the Block retains a reference to the caller's key")
	})
}

func (ci *CipherInfo) FuzzBlock(f *testing.F, newCipher func(key []byte) (cipher.Block, error), skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}

	ci.addBlockSeeds(f)

	f.Fuzz(func(t *testing.T, key, buf []byte) {
		b, bs := ci.newValid(t, newCipher, key)

		src := fitTo(buf, bs)
		if src == nil {
			t.Skip()
		}
		want := slices.Clone(src)

		// Round trip.
		ct := make([]byte, bs)
		b.Encrypt(ct, src)
		require.Equal(t, want, src, "Encrypt mutated src")
		require.NotEqual(t, want, ct, "Encrypt is the identity function")

		pt := make([]byte, bs)
		b.Decrypt(pt, ct)
		require.Equal(t, want, pt)

		// Determinism.
		ct2 := make([]byte, bs)
		b.Encrypt(ct2, want)
		require.Equal(t, ct, ct2)

		// dst aliasing src.
		ip := slices.Clone(want)
		b.Encrypt(ip, ip)
		require.Equal(t, ct, ip, "in-place Encrypt differs")
		b.Decrypt(ip, ip)
		require.Equal(t, want, ip, "in-place Decrypt differs")

		// src longer than one block: only the first block may be consumed.
		if long := fitTo(buf, bs+fuzzGuardLen); long != nil {
			copy(long, want)
			ctl := make([]byte, bs)
			b.Encrypt(ctl, long)
			require.Equal(t, ct, ctl, "Encrypt read past one block of src")
		}

		// Nothing may be written outside one block, in either direction.
		g, check := guarded(t, bs)
		b.Encrypt(g, want)
		require.Equal(t, ct, g)
		check("Encrypt")

		g, check = guarded(t, bs)
		b.Decrypt(g, ct)
		require.Equal(t, want, g)
		check("Decrypt")
	})
}

func (ci *CipherInfo) FuzzWide(f *testing.F, newCipher func(key []byte) (cipher.Block, error), skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}

	// Without a wide implementation the body below is a silent no-op, which
	// would report as a pass while verifying nothing.
	require.NotEmpty(f, ci.Sizes)
	probe, err := newCipher(make([]byte, ci.Sizes[0].KeySize))
	require.NoError(f, err)
	if len(wideFns(probe)) == 0 {
		f.Skip("cipher implements neither krypto.Block4 nor krypto.Block8")
	}

	ci.addBlockSeeds(f)

	f.Fuzz(func(t *testing.T, key, buf []byte) {
		b, bs := ci.newValid(t, newCipher, key)

		fns := wideFns(b)
		require.NotEmpty(t, fns, "wide interfaces vary between instances of the same cipher")

		for _, w := range fns {
			need := bs * w.n
			src := fitTo(buf, need)
			if src == nil {
				continue
			}
			want := slices.Clone(src)
			wantCT := refN(b, w.n, want, b.Encrypt)

			// Batch result must match the single-block reference.
			ct := make([]byte, need)
			w.enc(ct, src)
			require.Equalf(t, wantCT, ct, "%s differs from %dx Encrypt", wideName(w.n, false), w.n)
			require.Equalf(t, want, src, "%s mutated src", wideName(w.n, false))

			pt := make([]byte, need)
			w.dec(pt, ct)
			require.Equalf(t, want, pt, "%s round trip failed", wideName(w.n, true))

			// dst aliasing src.
			ip := slices.Clone(want)
			w.enc(ip, ip)
			require.Equalf(t, wantCT, ip, "in-place %s differs", wideName(w.n, false))
			w.dec(ip, ip)
			require.Equalf(t, want, ip, "in-place %s differs", wideName(w.n, true))

			// src longer than the batch: only the batch may be consumed.
			if long := fitTo(buf, need+fuzzGuardLen); long != nil {
				copy(long, want)
				ctl := make([]byte, need)
				w.enc(ctl, long)
				require.Equalf(t, wantCT, ctl, "%s read past the batch", wideName(w.n, false))
			}

			// Nothing may be written outside the batch, in either direction.
			g, check := guarded(t, need)
			w.enc(g, want)
			require.Equal(t, wantCT, g)
			check(wideName(w.n, false))

			g, check = guarded(t, need)
			w.dec(g, wantCT)
			require.Equal(t, want, g)
			check(wideName(w.n, true))
		}
	})
}

func (ci *CipherInfo) FuzzBlockShortBuffer(f *testing.F, newCipher func(key []byte) (cipher.Block, error), skip bool) {
	f.Helper()
	if skip {
		f.Skip("skipping fuzz test due to missing CPU feature")
	}

	for _, ks := range ci.Sizes {
		f.Add(bytes.Repeat([]byte{0x01}, ks.KeySize), 3, 5)
		f.Add(bytes.Repeat([]byte{0xff}, ks.KeySize), 0, 0)
		f.Add(bytes.Repeat([]byte{0x5a}, ks.KeySize), fuzzBufSize, fuzzBufSize)
		f.Add(bytes.Repeat([]byte{0x5a}, ks.KeySize), -1, fuzzBufSize)
	}

	f.Fuzz(func(t *testing.T, key []byte, dstLen, srcLen int) {
		clamp := func(v int) int { return min(max(v, 0), fuzzBufSize) }
		dstLen, srcLen = clamp(dstLen), clamp(srcLen)

		b, bs := ci.newValid(t, newCipher, key)
		dst, src := make([]byte, dstLen), make([]byte, srcLen)

		check := func(n int, enc, dec func(dst, src []byte)) {
			need := bs * n
			for _, c := range []struct {
				name string
				fn   func(dst, src []byte)
			}{
				{wideName(n, false), enc},
				{wideName(n, true), dec},
			} {
				if dstLen >= need && srcLen >= need {
					// Large enough: a panic here is just as much a bug.
					require.NotPanicsf(t, func() { c.fn(dst, src) },
						"%s panicked with dst=%d src=%d need=%d", c.name, dstLen, srcLen, need)
					continue
				}
				require.Panicsf(t, func() { c.fn(dst, src) },
					"%s did not panic with dst=%d src=%d need=%d", c.name, dstLen, srcLen, need)
			}
		}

		check(1, b.Encrypt, b.Decrypt)
		for _, w := range wideFns(b) {
			check(w.n, w.enc, w.dec)
		}
	})
}
