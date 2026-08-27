package internal_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
	"strconv"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/randutil"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	var dst [8]byte
	var xb, yb [8]byte

	_, rnd := randutil.NewTestRand()

	for xs := 1; xs <= 8; xs++ {
		for ys := 8; ys > 0; ys-- {
			for range 1000 {
				x, y := rnd.Uint64(), rnd.Uint64()
				x >>= 64 - xs*8
				y >>= 64 - ys*8

				binary.BigEndian.PutUint64(xb[:], x)
				binary.BigEndian.PutUint64(yb[:], y)

				internal.Add(dst[:], xb[len(xb)-xs:], yb[len(yb)-ys:])
				require.Equal(t, x+y, binary.BigEndian.Uint64(dst[:]))
			}
		}
	}
}

func TestIncCtr(t *testing.T) {
	const maxLen = 16

	test := func(l int) func(t *testing.T) {
		r, _ := randutil.NewTestRand()

		return func(t *testing.T) {
			expect := make([]byte, maxLen)
			answer := make([]byte, maxLen)

			for range 1000 {
				_, err := r.Read(expect)
				require.NoError(t, err)

				copy(answer, expect)
				high := binary.BigEndian.Uint64(expect[0:])
				low := binary.BigEndian.Uint64(expect[8:])

				for range 1000 {
					low++
					if low == 0 {
						high++
					}
					binary.BigEndian.PutUint64(expect[0:], high)
					binary.BigEndian.PutUint64(expect[8:], low)

					internal.IncCtr(answer[maxLen-l:])
					require.Equal(t, expect[maxLen-l:], answer[maxLen-l:])
				}
			}
		}
	}

	for i := 1; i <= maxLen; i++ {
		t.Run(strconv.Itoa(i), test(i))
	}
}

func TestFillCtr(t *testing.T) {
	for size := 1; size <= 16; size++ {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			ctr := make([]byte, size)
			for i := range ctr {
				ctr[i] = byte(i + 1)
			}
			ctr[size-1] = 0xFC

			expectCtr := slices.Clone(ctr)
			expectDst := make([]byte, 8*size)
			for off := 0; off < len(expectDst); off += size {
				internal.IncCtr(expectCtr)
				copy(expectDst[off:], expectCtr)
			}

			dst := make([]byte, 8*size)
			internal.FillCtr(dst, ctr)

			require.Equal(t, expectCtr, ctr)
			require.Equal(t, expectDst, dst)
		})
	}
}

func TestFillCtrRejectsInvalidLengths(t *testing.T) {
	require.Panics(t, func() { internal.FillCtr(make([]byte, 1), nil) })
	require.Panics(t, func() { internal.FillCtr(make([]byte, 5), make([]byte, 4)) })
	require.Panics(t, func() { internal.FillCtr(make([]byte, 10), make([]byte, 3)) })
}

func BenchmarkAdd(b *testing.B) {
	for _, size := range []int{8, 32, 64, 256, 1024} {
		for _, n := range []int{2, 4, 8} {
			b.Run(fmt.Sprintf("size=%d/src=%d", size, n), func(b *testing.B) {
				rnd, _ := randutil.NewBenchRand()

				dst := make([]byte, size)
				src := make([][]byte, n)
				for i := range src {
					src[i] = make([]byte, size)
					rnd.Read(src[i])
				}

				b.SetBytes(int64(size))
				b.ReportAllocs()
				b.ResetTimer()
				for b.Loop() {
					internal.Add(dst, src...)
				}
			})
		}
	}
}

func BenchmarkAddCarry(b *testing.B) {
	const size = 256
	dst := make([]byte, size)
	x := bytes.Repeat([]byte{0xFF}, size)
	y := make([]byte, size)
	y[size-1] = 1

	b.SetBytes(size)
	b.ReportAllocs()
	for b.Loop() {
		internal.Add(dst, x, y)
	}
}

func BenchmarkAddUnaligned(b *testing.B) {
	const size = 256
	buf := make([]byte, size*3+3)
	dst := buf[1 : 1+size]
	x := buf[size+2 : size+2+size]
	y := buf[size*2+3:]

	b.SetBytes(size)
	for b.Loop() {
		internal.Add(dst, x, y)
	}
}

func BenchmarkIncCtr(b *testing.B) {
	bench := func(size int) func(b *testing.B) {
		return func(b *testing.B) {
			rnd, _ := randutil.NewBenchRand()

			ctr := make([]byte, size)
			_, err := rnd.Read(ctr)
			require.NoError(b, err)

			b.SetBytes(int64(size))
			for b.Loop() {
				internal.IncCtr(ctr)
			}
		}
	}

	for i := 1; i <= 16; i++ {
		b.Run(strconv.Itoa(i), bench(i))
	}
}

func BenchmarkFillCtr(b *testing.B) {
	for _, size := range []int{4, 8, 16, 32} {
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			rnd, _ := randutil.NewBenchRand()
			ctr := make([]byte, size)
			dst := make([]byte, 8*size)
			_, err := rnd.Read(ctr)
			require.NoError(b, err)

			b.SetBytes(int64(len(dst)))
			b.ReportAllocs()
			for b.Loop() {
				internal.FillCtr(dst, ctr)
			}
		})
	}
}
