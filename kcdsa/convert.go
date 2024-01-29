package kcdsa

import (
	"crypto/dsa"
	"math/big"

	kcdsainternal "github.com/RyuaNerin/go-krypto/internal/kcdsa"
)

/**
DSA		Y = { G^x } mod P
KCDSA	Y = G^{X^{-1} mod Q} mod P
*/

func FromDSA(dpk *dsa.PrivateKey) *PrivateKey {
	kpk := &PrivateKey{
		X: new(big.Int).Set(dpk.X),
		PublicKey: PublicKey{
			Parameters: Parameters{
				P: new(big.Int).Set(dpk.P),
				Q: new(big.Int).Set(dpk.Q),
				G: new(big.Int).Set(dpk.G),
			},
		},
	}

	kpk.PublicKey.Y = kcdsainternal.Y(kpk.P, kpk.Q, kpk.G, kpk.X)

	return kpk
}

func (kpk *PrivateKey) ToDSA() *dsa.PrivateKey {
	dpk := &dsa.PrivateKey{
		X: new(big.Int).Set(kpk.X),
		PublicKey: dsa.PublicKey{
			Y: new(big.Int),
			Parameters: dsa.Parameters{
				P: new(big.Int).Set(kpk.P),
				Q: new(big.Int).Set(kpk.Q),
				G: new(big.Int).Set(kpk.G),
			},
		},
	}

	dpk.Y.Exp(dpk.G, dpk.X, dpk.P)
	return dpk
}
