package internal

import (
	"io"
	"math/big"
)

func Bytes(bits int) int {
	// 32bit: 0xFFFF_FFF8 / 64bit: 0xFFFF FFFF FFFF FFF8
	const MaxUint_m7 = ^uint(0) - 7

	return int(((uint(bits) + 7) & MaxUint_m7) / 8)
}

// same with `int(math.Ceil(float64(a) / float64(b)))`
func CeilDiv(a, b int) int {
	if b == 0 {
		return 0
	}
	return (a + b - 1) / b
}

func IncCtr(b []byte) {
	for i := len(b) - 1; i >= 0; i-- {
		c := b[i]
		c++
		b[i] = c
		if c > 0 {
			return
		}
	}
}

// Int returns a uniform random value in [0, max). It panics if max <= 0.
func Int(dst *big.Int, rand io.Reader, buf []byte, max *big.Int) (bufNew []byte, err error) {
	if max.Sign() <= 0 {
		panic("crypto/rand: argument to Int is <= 0")
	}
	dst.Sub(max, dst.SetUint64(1))
	// bitLen is the maximum bit length needed to encode a value < max.
	bitLen := dst.BitLen()
	if bitLen == 0 {
		// the only valid result is 0
		return
	}
	// k is the maximum byte length needed to encode a value < max.
	k := (bitLen + 7) / 8
	// b is the number of bits in the most significant byte of max-1.
	b := uint(bitLen % 8)
	if b == 0 {
		b = 8
	}

	bufNew = ResizeBuffer(buf, k)

	mask := uint8(int(1<<b) - 1)

	for {
		_, err = io.ReadFull(rand, bufNew)
		if err != nil {
			return bufNew, err
		}

		// Clear bits in the first byte to increase the probability
		// that the candidate is < max.
		bufNew[0] &= mask

		dst.SetBytes(bufNew)
		if dst.Cmp(max) < 0 {
			return
		}
	}
}
