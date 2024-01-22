package kcdsattak

import (
	"hash"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/kcdsa"
	kcdsainternal "github.com/RyuaNerin/go-krypto/kcdsa/internal"
)

// Generate the paramters
func GenerateParameters(params *Parameters, rand io.Reader, sizes kcdsa.ParameterSizes) (seed []byte, count int, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, 0, ErrInvalidParameterSizes
	}

	// p. 13
	for {
		seed, err := internal.ReadBits(seed[:0], rand, domain.B)
		if err != nil {
			return nil, 0, err
		}

		J, err := generateJ(seed, domain)
		if err != nil {
			continue
		}

		P, Q, count, err := generatePQ(J, seed, domain)
		if err != nil {
			continue
		}

		_, G, err := GenerateHG(rand, P, J)
		if err != nil {
			continue
		}

		params.P = P
		params.Q = Q
		params.G = G
		return seed, count, nil
	}
}

// Generate PublicKey and PrivateKey
func GenerateKey(priv *PrivateKey, rand io.Reader, userProvidedRandomInput []byte, sizes kcdsa.ParameterSizes) error {
	if priv.P == nil || priv.Q == nil || priv.G == nil {
		return kcdsa.ErrParametersNotSetUp
	}
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return ErrInvalidParameterSizes
	}

	// p.16
	xkey, err := internal.ReadBits(nil, rand, domain.B)
	if err != nil {
		return err
	}

	X, Y, _, _, err := generateXYZ(priv.P, priv.Q, priv.G, userProvidedRandomInput, xkey, domain)
	if err != nil {
		return err
	}

	priv.X = X
	priv.Y = Y
	return nil
}

// Sign data using K generated randomly like in crypto/dsa packages.
func Sign(rand io.Reader, priv *PrivateKey, data []byte, sizes kcdsa.ParameterSizes) (r, s *big.Int, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, nil, kcdsa.ErrInvalidParameterSizes
	}

	machineGeneratedRandomInput := make([]byte, internal.Bytes(domain.B))
	_, err = rand.Read(machineGeneratedRandomInput)
	if err != nil {
		return
	}

	return sign(rand, priv, machineGeneratedRandomInput, data, domain)
}

func sign(rand io.Reader, priv *PrivateKey, machineGeneratedRandomInput []byte, data []byte, domain kcdsainternal.Domain) (r, s *big.Int, err error) {
	if priv.Q.Sign() <= 0 || priv.P.Sign() <= 0 || priv.G.Sign() <= 0 || priv.X.Sign() <= 0 || priv.Q.BitLen()%8 != 0 {
		err = ErrInvalidPublicKey
		return
	}

	J, err := generateK(rand, priv.Q, machineGeneratedRandomInput, domain)
	if err != nil {
		return nil, nil, err
	}

	return kcdsainternal.Sign(priv.P, priv.Q, priv.G, priv.Y, priv.X, J, domain.NewHash(), data)
}

func Verify(pub *PublicKey, h hash.Hash, data []byte, R, S *big.Int) bool {
	return kcdsainternal.Verify(pub.P, pub.Q, pub.G, pub.Y, h, data, R, S)
}
