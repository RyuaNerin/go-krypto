package eckcdsa

import (
	"crypto/elliptic"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

func XY(D *big.Int, c elliptic.Curve) (X, Y *big.Int) {
	dInv := internal.FermatInverse(D, c.Params().N)
	return c.ScalarBaseMult(dInv.Bytes())
}
