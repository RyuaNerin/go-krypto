package internal

import (
	"io"
	"math/big"
)

func FermatInverse(a, N *big.Int) *big.Int {
	// https://cs.opensource.google/go/go/+/refs/tags/go1.18:src/crypto/dsa/dsa.go;l=188-192
	two := big.NewInt(2)
	nMinus2 := new(big.Int).Sub(N, two)
	return new(big.Int).Exp(a, nMinus2, N)
}

const NumMRTests = 64

func Bytes(bits int) int {
	return ((bits + 7) & 0xFFFFFFF8) / 8
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
