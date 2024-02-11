package internal

import (
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal/kryptoutil"
)

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

func Add(dst []byte, src ...[]byte) {
	n := len(dst)

	var value uint64
	for idx := 0; idx < n; idx++ {
		for _, v := range src {
			if idx < len(v) {
				value += uint64(v[len(v)-idx-1])
			}
		}

		dst[len(dst)-idx-1] = byte(value & 0xFF)
		value = value >> 8
	}
	kryptoutil.MemsetByte(dst[:len(dst)-n], 0)
}

// Int returns a uniform random value in [0, max). It panics if max <= 0.
func Int(dst *big.Int, rand io.Reader, buf []byte, max *big.Int) (bytes []byte, err error) {
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

	bytes = Expand(buf[:0], k)

	for {
		_, err = io.ReadFull(rand, bytes)
		if err != nil {
			return bytes, err
		}

		// Clear bits in the first byte to increase the probability
		// that the candidate is < max.
		bytes[0] &= uint8(int(1<<b) - 1)

		dst.SetBytes(bytes)
		if dst.Cmp(max) < 0 {
			return
		}
	}
}
