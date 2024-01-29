package eckcdsa

import (
	"crypto/ecdsa"
	"math/big"

	eckcdsainternal "github.com/RyuaNerin/go-krypto/internal/eckcdsa"
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

	kpk.X, kpk.Y = eckcdsainternal.XY(kpk.D, kpk.Curve)
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
