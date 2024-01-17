package internal

import (
	"io"
	"math/big"
)

// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/crypto/dsa/dsa.go;l=188-192
// fermatInverse calculates the inverse of k in GF(P) using Fermat's method.
// This has better constant-time properties than Euclid's method (implemented
// in math/big.Int.ModInverse) although math/big itself isn't strictly
// constant-time so it's not perfect.
func FermatInverse(a, N *big.Int) *big.Int {
	two := big.NewInt(2)
	nMinus2 := new(big.Int).Sub(N, two)
	return new(big.Int).Exp(a, nMinus2, N)
}

const NumMRTests = 64

func Bytes(bits int) int {
	return int((uint(bits) + 7) & ^uint(0xFFFFFFF8) / 8)
}

func ReadBits(dst []byte, rand io.Reader, bits int) ([]byte, error) {
	i := Bytes(bits)
	if cap(dst) < i {
		dst = make([]byte, i)
	} else {
		for len(dst) < i {
			dst = append(dst, 0)
		}
	}

	_, err := rand.Read(dst)
	if err != nil {
		return nil, err
	}

	i = bits & 0x07
	if i != 0 {
		dst[0] &= byte((1 << i) - 1)
	}

	return dst, nil
}
