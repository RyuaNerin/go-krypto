package internal

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"strconv"
	"testing"
)

func TestAdd(t *testing.T) {
	var dst [8]byte
	var xb, yb [8]byte

	for xs := 1; xs <= 8; xs++ {
		for ys := 8; ys > 0; ys-- {
			for i := 0; i < 1000; i++ {
				x, y := uint64(rand.Int63()), uint64(rand.Int63())
				x >>= 64 - xs*8
				y >>= 64 - ys*8

				binary.BigEndian.PutUint64(xb[:], x)
				binary.BigEndian.PutUint64(yb[:], y)

				Add(dst[:], xb[len(xb)-xs:], yb[len(yb)-ys:])
				if binary.BigEndian.Uint64(dst[:]) != x+y {
					t.Fail()
					return
				}
			}
		}
	}
}

func TestIncCtr(t *testing.T) {
	const maxLen = 16

	test := func(l int) func(t *testing.T) {
		return func(t *testing.T) {
			expect := make([]byte, maxLen)
			answer := make([]byte, maxLen)

			for i := 0; i < 1000; i++ {
				crand.Read(expect)
				copy(answer, expect)
				high := binary.BigEndian.Uint64(expect[0:])
				low := binary.BigEndian.Uint64(expect[8:])

				for i := 0; i < 1000; i++ {
					low++
					if low == 0 {
						high++
					}
					binary.BigEndian.PutUint64(expect[0:], high)
					binary.BigEndian.PutUint64(expect[8:], low)

					IncCtr(answer[maxLen-l:])
					if !bytes.Equal(expect[maxLen-l:], answer[maxLen-l:]) {
						t.Errorf("test failed\nvalue: %16x %16x\nexpect: %x\nanswer: %x", high, low, expect[maxLen-l:], answer[maxLen-l:])
						return
					}
				}
			}
		}
	}

	for i := 1; i <= maxLen; i++ {
		t.Run(strconv.Itoa(i), test(i))
	}
}

func BenchmarkAdd(b *testing.B) {
	var dst [8]byte
	var xb, yb [8]byte

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		crand.Read(xb[:])
		crand.Read(yb[:])

		Add(dst[:], xb[:], yb[:])
	}
}

func BenchmarkIncCtr(b *testing.B) {
	bench := func(size int) func(b *testing.B) {
		return func(b *testing.B) {
			ctr := make([]byte, size)
			crand.Read(ctr)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				IncCtr(ctr)
			}
		}
	}

	for i := 1; i <= 10; i++ {
		b.Run(strconv.Itoa(i), bench(i))
	}
}
