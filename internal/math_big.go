package internal

import (
	"crypto/subtle"
	"math/big"
)

func BigIntEqual(a, b *big.Int) bool {
	return subtle.ConstantTimeCompare(a.Bytes(), b.Bytes()) == 1
}
