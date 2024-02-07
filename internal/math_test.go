package internal

import (
	"encoding/binary"
	"math/rand"
	"testing"
)

func TestAdd(t *testing.T) {
	var xb, yb, zb [8]byte

	t.Run("8+8", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			x, y := uint64(rand.Int63()), uint64(rand.Int63())

			binary.BigEndian.PutUint64(xb[:], x)
			binary.BigEndian.PutUint64(yb[:], y)

			Add(zb[:], xb[:], yb[:])
			if binary.BigEndian.Uint64(zb[:]) != x+y {
				t.Fail()
				return
			}
		}
	})

	t.Run("4+8", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			x, y := uint64(rand.Int31()), uint64(rand.Int63())

			binary.BigEndian.PutUint64(xb[:], x)
			binary.BigEndian.PutUint64(yb[:], y)

			Add(zb[:], xb[4:], yb[:])
			if binary.BigEndian.Uint64(zb[:]) != x+y {
				t.Fail()
				return
			}
		}
	})

	t.Run("8+4", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			x, y := uint64(rand.Int63()), uint64(rand.Int31())

			binary.BigEndian.PutUint64(xb[:], x)
			binary.BigEndian.PutUint64(yb[:], y)

			Add(zb[:], xb[:], yb[4:])
			if binary.BigEndian.Uint64(zb[:]) != x+y {
				t.Fail()
				return
			}
		}
	})
}
