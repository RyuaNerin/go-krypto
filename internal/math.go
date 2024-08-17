package internal

import (
	"encoding/binary"
	"io"
	"math/big"
	"math/bits"

	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

const NumMRTests = 64

var (
	One   = big.NewInt(1)
	Two   = big.NewInt(2)
	Three = big.NewInt(3)
)

func BitsToBytes(bits int) int {
	// 32bit: 0xFFFF_FFF8 / 64bit: 0xFFFF FFFF FFFF FFF8
	const MAX_UINT_MINUS_7 = ^uint(0) - 7

	return int(((uint(bits) + 7) & MAX_UINT_MINUS_7) / 8)
}

// same with `int(math.Ceil(float64(a) / float64(b)))`
func CeilDiv(a, b int) int {
	if b == 0 {
		return 0
	}
	return (a + b - 1) / b
}

// [0, max)
func ReadInt(r io.Reader, max int) (int, error) { //nolint:predeclared
	var buf [8]byte

	bitSize := uint(bits.Len64(uint64(max)))

	var randomValue uint64
	for {
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return 0, err
		}
		buf[0] &= 0x7F

		randomValue = binary.BigEndian.Uint64(buf[:])
		randomValue >>= 64 - bitSize

		if randomValue < uint64(max) {
			return int(randomValue), nil
		}
	}
}

// ReadBigInt returns a uniform random value in [0, max). It panics if max <= 0.
func ReadBigInt(dst *big.Int, rand io.Reader, buf []byte, max *big.Int) (bufNew []byte, err error) { //nolint:predeclared
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
	k := BitsToBytes(bitLen)
	// b is the number of bits in the most significant byte of max-1.
	b := uint(bitLen % 8)
	if b == 0 {
		b = 8
	}

	bufNew = Grow(buf, k)

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

// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/crypto/dsa/dsa.go;l=188-192
// fermatInverse calculates the inverse of k in GF(P) using Fermat's method.
// This has better constant-time properties than Euclid's method (implemented
// in math/big.Int.ModInverse) although math/big itself isn't strictly
// constant-time so it's not perfect.
func FermatInverse(k, P *big.Int) *big.Int {
	tmp := new(big.Int).Sub(P, Two)
	return tmp.Exp(k, tmp, P)
}

func BigEqual(a, b *big.Int) bool {
	return subtle.ConstantTimeCompare(a.Bytes(), b.Bytes()) == 1
}

// log(n) / log(2)
func BigLog2(n *big.Int) int {
	return n.BitLen() - 1
}

// log(n) / log(2 ** e)
func BigLog2e(n *big.Int, e int) int {
	return (n.BitLen() - 1) / e
}

// = ceil(log(n) / log(2 ** e))
func BigCeilLog2(n *big.Int, e int) int {
	x := BigLog2e(n, e)
	var c big.Int
	c.Lsh(One, uint(x)*uint(e))

	// if n > c, return x+1
	// else return x
	if n.Cmp(&c) > 0 {
		return x + 1
	}
	return x
}
