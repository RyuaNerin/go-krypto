package eckcdsa

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

/**
EC-DSA		Q = d * G
EC-KCDSA	Q = {d ^ {-1}} * G
*/

func FromECDSA(dpk *ecdsa.PrivateKey) (*PrivateKey, error) {
	kpk := &PrivateKey{
		D: new(big.Int).Set(dpk.D),
		PublicKey: PublicKey{
			Curve: dpk.Curve,
		},
	}

	dInv := internal.FermatInverse(kpk.D, kpk.Curve.Params().N)
	if dInv == nil {
		return nil, errors.New("eckcdsa: Unsupported D")
	}
	kpk.X, kpk.Y = kpk.Curve.ScalarBaseMult(dInv.Bytes())

	return kpk, nil
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
