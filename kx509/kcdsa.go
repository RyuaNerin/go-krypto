package kx509

import (
	"encoding/asn1"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal/zeroize"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

// kcdsaParameters reflects the KCDSAParameters structure used in the
// AlgorithmIdentifier of both PrivateKeyInfo and SubjectPublicKeyInfo.
//
//	KCDSAParameters ::= SEQUENCE {
//	    p       INTEGER,
//	    q       INTEGER,
//	    g       INTEGER,
//	    J       INTEGER OPTIONAL,
//	    Seed    OCTET STRING OPTIONAL,
//	    Count   INTEGER OPTIONAL
//	}
type kcdsaParameters struct {
	P     *big.Int
	Q     *big.Int
	G     *big.Int
	J     *big.Int `asn1:"optional"`
	Seed  []byte   `asn1:"optional"`
	Count int      `asn1:"optional"`
}

func (params *kcdsaParameters) Zeroize() {
	zeroize.Zeroize(params.P)
	zeroize.Zeroize(params.Q)
	zeroize.Zeroize(params.G)
	zeroize.Zeroize(params.J)
	zeroize.Zeroize(params.Seed)
	params.Count = 0
}

func newKCDSAParameters(params *kcdsa.Parameters) kcdsaParameters {
	if params == nil {
		return kcdsaParameters{}
	}
	return kcdsaParameters{
		P:     params.P,
		Q:     params.Q,
		G:     params.G,
		J:     params.GenParameters.J,
		Seed:  params.GenParameters.Seed,
		Count: params.GenParameters.Count,
	}
}

func (params *kcdsaParameters) parameters() kcdsa.Parameters {
	return kcdsa.Parameters{
		P: params.P,
		Q: params.Q,
		G: params.G,
		GenParameters: kcdsa.GenerationParameters{
			J:     params.J,
			Seed:  params.Seed,
			Count: params.Count,
		},
	}
}

func parseKCDSAParameters(der []byte) (kcdsaParameters, error) {
	var params kcdsaParameters
	rest, err := asn1.Unmarshal(der, &params)
	if err != nil || len(rest) != 0 || params.P == nil || params.Q == nil || params.G == nil {
		return kcdsaParameters{}, ErrInvalidParametersDSA
	}
	if err := validateKCDSAParameters(params); err != nil {
		return kcdsaParameters{}, err
	}
	return params, nil
}

func marshalKCDSAParameters(params *kcdsa.Parameters) ([]byte, error) {
	asn1Params := newKCDSAParameters(params)
	if err := validateKCDSAParameters(asn1Params); err != nil {
		return nil, err
	}
	return asn1.Marshal(asn1Params)
}

func validateKCDSAParameters(params kcdsaParameters) error {
	// Restricting P and Q to the supported sizes also bounds the cost of the
	// modular exponentiations done on parsed keys.
	domain := params.parameters()
	if err := domain.Validate(); err != nil {
		return ErrInvalidParametersDSA
	}
	if (params.J != nil && params.J.Sign() <= 0) || params.Count < 0 || params.Count > 1<<24 {
		return ErrInvalidParametersDSA
	}
	return nil
}
