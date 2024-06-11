package internal

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
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
