package eckcdsa

import (
	"crypto/elliptic"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

func GenerateXY(D *big.Int, c elliptic.Curve) (X, Y *big.Int) {
	dInv := internal.FermatInverse(D, c.Params().N)
	defer zeroize.Zeroize(dInv)
	return c.ScalarBaseMult(dInv.Bytes())
}
