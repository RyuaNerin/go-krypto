package eckcdsa

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

/**
EC-DSA		Q = d * G
EC-KCDSA	Q = {d ^ {-1}} * G
*/

func FromECDSA(dpk *ecdsa.PrivateKey) *PrivateKey {
	kpk := &PrivateKey{
		D: new(big.Int).Set(dpk.D),
		PublicKey: PublicKey{
			Curve: dpk.Curve,
		},
	}

	dInv := internal.FermatInverse(kpk.D, kpk.Curve.Params().N)
	kpk.X, kpk.Y = kpk.Curve.ScalarBaseMult(dInv.Bytes())

	return kpk
}

func (kpk *PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	dpk := &ecdsa.PrivateKey{
		D: new(big.Int).Set(kpk.D),
		PublicKey: ecdsa.PublicKey{
			Curve: kpk.Curve,
		},
	}

	dpk.X, dpk.Y = dpk.Curve.ScalarBaseMult(dpk.D.Bytes())

	return dpk
}
