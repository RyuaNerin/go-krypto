package drbg

import (
	"encoding/binary"
	"io"
	"math/bits"
)

type entropy struct {
	rand     io.Reader
	min, max int

	buf    []byte
	remain int
}

func newEntropy(rand io.Reader, min, max int) *entropy {
	e := &entropy{
		rand: rand,
		min:  min,
		max:  max,
		buf:  make([]byte, max),
	}

	return e
}

func (e *entropy) Get() (p []byte, err error) {
	length, err := uintN(e.rand, e.min, e.max)
	if err != nil {
		return nil, err
	}

	_, err = io.ReadFull(e.rand, e.buf[:length])
	if err != nil {
		return nil, err
	}

	return e.buf[:length], nil
}

// [min, max]
func uintN(rand io.Reader, min, max int) (int, error) {
	if min == max {
		return min, nil
	}

	var buf [8]byte

	rangeSize := uint64(max - min + 1)
	bitSize := uint(bits.Len64(rangeSize))

	var randomValue uint64
	for {
		if _, err := io.ReadFull(rand, buf[:]); err != nil {
			return 0, err
		}

		randomValue = binary.BigEndian.Uint64(buf[:])
		randomValue >>= 64 - bitSize

		if randomValue < rangeSize {
			return min + int(randomValue), nil
		}
	}
}
