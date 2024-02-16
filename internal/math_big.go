package internal

import (
	"crypto/subtle"
	"io"
	"math/big"
)

const NumMRTests = 64

// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/crypto/dsa/dsa.go;l=188-192
// fermatInverse calculates the inverse of k in GF(P) using Fermat's method.
// This has better constant-time properties than Euclid's method (implemented
// in math/big.Int.ModInverse) although math/big itself isn't strictly
// constant-time so it's not perfect.
func FermatInverse(a, N *big.Int) *big.Int {
	nMinus2 := new(big.Int).Sub(N, Two)
	return new(big.Int).Exp(a, nMinus2, N)
}

func BigIntEqual(a, b *big.Int) bool {
	return subtle.ConstantTimeCompare(a.Bytes(), b.Bytes()) == 1
}

// resize dst, ReadFull, cut from right
func ReadBits(dst []byte, rand io.Reader, bits int) ([]byte, error) {
	bytes := Bytes(bits)

	dst = ResizeBuffer(dst, bytes)

	_, err := io.ReadFull(rand, dst)
	if err != nil {
		return dst, err
	}

	bytes = bits & 0x07
	if bytes != 0 {
		dst[0] &= byte((1 << bytes) - 1)
	}

	return dst, nil
}

// resize dst, ReadFull, cut from right
func ReadBytes(dst []byte, rand io.Reader, bytes int) ([]byte, error) {
	dst = ResizeBuffer(dst, bytes)

	_, err := io.ReadFull(rand, dst)
	if err != nil {
		return dst, err
	}

	return dst, nil
}
