package cavp

import (
	"crypto/cipher"
	"hash"
	"slices"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/drbg"
)

type drbgVector struct {
	EntropyInput          []byte
	Nonce                 []byte
	PersonalizationString []byte
	EntropyInputReseed    []byte
	AdditionalInputReseed []byte
	EntropyInputPR        [2][]byte
	AdditionalInput       [2][]byte
	ReturnedBitsLen       int
	PredictionResistance  bool
}

type drbgEntropy struct {
	values [][]byte
	next   int
}

func newDRBGEntropy(values ...[]byte) *drbgEntropy { return &drbgEntropy{values: values} }

func (e *drbgEntropy) Get() ([]byte, error) {
	v := slices.Clone(e.values[e.next])
	e.next++
	return v, nil
}

func generate(v drbgVector, newDRBG func(drbg.Entropy) (krypto.DRBG, error)) ([]byte, error) {
	entropy := newDRBGEntropy(v.EntropyInputPR[0], v.EntropyInputPR[1])
	g, err := newDRBG(entropy)
	if err != nil {
		return nil, err
	}
	out := make([]byte, v.ReturnedBitsLen/8)
	if v.PredictionResistance {
		if _, err = g.Generate(out, v.AdditionalInput[0]); err != nil {
			return nil, err
		}
		_, err = g.Generate(out, v.AdditionalInput[1])
		return out, err
	}
	if err = g.ReseedWithEntropy(v.AdditionalInputReseed, v.EntropyInputReseed); err != nil {
		return nil, err
	}
	if _, err = g.Generate(out, v.AdditionalInput[0]); err != nil {
		return nil, err
	}
	_, err = g.Generate(out, v.AdditionalInput[1])
	return out, err
}

func drbgHash(v drbgVector, newHash func() hash.Hash) ([]byte, error) {
	return generate(v, func(e drbg.Entropy) (krypto.DRBG, error) {
		return drbg.NewHashDRBG(
			newHash(),
			drbg.WithInstantiateEntropy(v.EntropyInput),
			drbg.WithNonce(v.Nonce),
			drbg.WithPersonalizationString(v.PersonalizationString),
			drbg.WithPredictionResistance(v.PredictionResistance),
			drbg.WithCustomEntropy(e),
		)
	})
}

func drbgHMAC(v drbgVector, newHash func() hash.Hash) ([]byte, error) {
	return generate(v, func(e drbg.Entropy) (krypto.DRBG, error) {
		return drbg.NewHMACDRBG(
			newHash,
			drbg.WithInstantiateEntropy(v.EntropyInput),
			drbg.WithNonce(v.Nonce),
			drbg.WithPersonalizationString(v.PersonalizationString),
			drbg.WithPredictionResistance(v.PredictionResistance),
			drbg.WithCustomEntropy(e),
		)
	})
}

func drbgCTR(v drbgVector, newCipher func([]byte) (cipher.Block, error), keySize int, useDF bool) ([]byte, error) {
	return generate(v, func(e drbg.Entropy) (krypto.DRBG, error) {
		return drbg.NewCTRDRBG(
			newCipher, keySize,
			drbg.WithInstantiateEntropy(v.EntropyInput),
			drbg.WithNonce(v.Nonce),
			drbg.WithPersonalizationString(v.PersonalizationString),
			drbg.WithPredictionResistance(v.PredictionResistance),
			drbg.WithDerivationFunction(useDF),
			drbg.WithCustomEntropy(e),
		)
	})
}
