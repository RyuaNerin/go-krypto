package internal

import (
	"crypto/subtle"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal/zeroize"
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

// ReadBigInt returns a uniform random value in [0, max). It panics if max <= 0.
func ReadBigInt(dst *big.Int, rand io.Reader, buf []byte, max *big.Int) (bufNew []byte, err error) { //nolint:predeclared
	if max.Sign() <= 0 {
		panic(msgInvalidIntArgument)
	}
	dst.Sub(max, dst.SetUint64(1))
	// bitLen is the maximum bit length needed to encode a value < max.
	bitLen := dst.BitLen()
	if bitLen == 0 {
		// the only valid result is 0
		return buf, nil
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
			return bufNew, nil
		}
	}
}

// https://github.com/golang/go/blob/go1.27.0/src/crypto/dsa/dsa.go#L197-L201
// fermatInverse calculates the inverse of k in GF(P) using Fermat's method.
// This has better constant-time properties than Euclid's method (implemented
// in math/big.Int.ModInverse) although math/big itself isn't strictly
// constant-time so it's not perfect.
func FermatInverse(k, P *big.Int) *big.Int {
	pMinus2 := new(big.Int).Sub(P, Two)
	return new(big.Int).Exp(k, pMinus2, P)
}

func BigEqual(a, b *big.Int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Sign() == b.Sign() && subtle.ConstantTimeCompare(a.Bytes(), b.Bytes()) == 1
}

// = ceil(log(n) / log(2 ** e))
func DigitLen(n *big.Int) int {
	const e = 8

	var b big.Int
	defer zeroize.Zeroize(&b)

	b.Sub(n, One)
	return (b.BitLen() + e - 1) / e
}
