package kcdsa

import (
	"crypto/dsa"
	"errors"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

/**
DSA		Y = { G^x } mod P
KCDSA	Y = G^{X^{-1} mod Q} mod P
*/

func FromDSA(dpk *dsa.PrivateKey) (*PrivateKey, error) {
	kpk := &PrivateKey{
		X: new(big.Int).Set(dpk.X),
		PublicKey: PublicKey{
			Y: new(big.Int),
			Parameters: Parameters{
				P: new(big.Int).Set(dpk.P),
				Q: new(big.Int).Set(dpk.Q),
				G: new(big.Int).Set(dpk.G),
			},
		},
	}

	xInv := internal.FermatInverse(kpk.X, kpk.Q)
	if xInv == nil {
		return nil, errors.New("kcdsa: Unsupported X")
	}
	kpk.PublicKey.Y.Exp(kpk.G, xInv, kpk.P)

	return kpk, nil
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
