package kx509

import (
	"crypto/elliptic"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"math/big"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/internal/golang.org/x/crypto/cryptobyte"
	cryptobyte_asn1 "github.com/RyuaNerin/go-krypto/internal/golang.org/x/crypto/cryptobyte/asn1"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/x509.go#L58
type pkixPublicKey struct {
	Algo      pkix.AlgorithmIdentifier
	BitString asn1.BitString
}

// ParsePKIXPublicKey parses a public key in PKIX, ASN.1 DER form.
//
// It returns an *eckcdsa.PublicKey, an *kcdsa.PublicKey,
// or the result of crypto/x509.ParsePKIXPublicKey
//
// This kind of key is commonly encoded in PEM blocks of type "PUBLIC KEY".
func ParsePKIXPublicKey(derBytes []byte) (pub interface{}, err error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/x509.go#L71-L82
	var pki publicKeyInfo
	var isSupportedType bool
	rest, err := asn1.Unmarshal(derBytes, &pki)
	if err == nil && len(rest) == 0 {
		pub, isSupportedType, err = parsePublicKey(&pki)
		if isSupportedType {
			return
		}
	}
	return x509.ParsePKIXPublicKey(derBytes)
}

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/x509.go#L84
func marshalPublicKey(pub interface{}) (publicKeyBytes []byte, publicKeyAlgorithm pkix.AlgorithmIdentifier, isSupportedType bool, err error) {
	switch pub := pub.(type) {
	case *eckcdsa.PublicKey: //nolint:typecheck
		oid, ok := oidFromNamedCurve(pub.Curve)
		if !ok {
			return nil, pkix.AlgorithmIdentifier{}, true, errors.New(msgUnknownEllipticCurve)
		}
		if !pub.Curve.IsOnCurve(pub.X, pub.Y) { //nolint:staticcheck
			return nil, pkix.AlgorithmIdentifier{}, true, errors.New(msgInvalidPublicKey)
		}
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y) //nolint:staticcheck
		publicKeyAlgorithm.Algorithm = oidPublicKeyECKCDSA
		var paramBytes []byte
		paramBytes, err = asn1.Marshal(oid)
		if err != nil {
			return
		}
		publicKeyAlgorithm.Parameters.FullBytes = paramBytes

	case *kcdsa.PublicKey:
		publicKeyBytes, err = asn1.Marshal(pub.Y)
		if err != nil {
			return nil, pkix.AlgorithmIdentifier{}, true, errors.New(msgInvalidPrivateKeyY)
		}

		params := kcdsaParameters{
			P: pub.P,
			Q: pub.Q,
			G: pub.G,
		}
		if pub.GenParameters.IsValid() {
			params.J = pub.GenParameters.J
			params.Seed = pub.GenParameters.Seed
			params.Count = pub.GenParameters.Count
		}

		var paramBytes []byte
		paramBytes, err = asn1.Marshal(params)
		if err != nil {
			return nil, pkix.AlgorithmIdentifier{}, true, errors.New(msgInvalidParametersDSA)
		}

		publicKeyAlgorithm.Algorithm = oidPublicKeyKCDSA
		publicKeyAlgorithm.Parameters.FullBytes = paramBytes

	default:
		return nil, pkix.AlgorithmIdentifier{}, false, nil
	}

	return publicKeyBytes, publicKeyAlgorithm, true, nil
}

// MarshalPKIXPublicKey converts a public key to PKIX, ASN.1 DER form.
//
// supported key types : *eckcdsa.PublicKey and *kcdsa.PublicKey.
// for unsupported types, returns the result of crypto/x509.MarshalPKIXPublicKey
//
// This kind of key is commonly encoded in PEM blocks of type "PUBLIC KEY".
func MarshalPKIXPublicKey(pub interface{}) ([]byte, error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/x509.go#L150
	var publicKeyBytes []byte
	var publicKeyAlgorithm pkix.AlgorithmIdentifier
	var isSupportedType bool
	var err error

	publicKeyBytes, publicKeyAlgorithm, isSupportedType, err = marshalPublicKey(pub)
	if !isSupportedType {
		return x509.MarshalPKIXPublicKey(pub)
	}
	if err != nil {
		return nil, err
	}

	pkix := pkixPublicKey{
		Algo: publicKeyAlgorithm,
		BitString: asn1.BitString{
			Bytes:     publicKeyBytes,
			BitLength: 8 * len(publicKeyBytes),
		},
	}

	ret, _ := asn1.Marshal(pkix)
	return ret, nil
}

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/x509.go#L201-L205
type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

func parsePublicKey(keyData *publicKeyInfo) (privateKey interface{}, isSupportedType bool, err error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/parser.go#L217-L220
	oid := keyData.Algorithm.Algorithm
	params := keyData.Algorithm.Parameters
	der := cryptobyte.String(keyData.PublicKey.RightAlign())

	switch {
	case oid.Equal(oidPublicKeyECKCDSA):
		// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/parser.go#L252-L271
		paramsDer := cryptobyte.String(params.FullBytes)
		namedCurveOID := new(asn1.ObjectIdentifier)
		if !paramsDer.ReadASN1ObjectIdentifier(namedCurveOID) {
			return nil, true, errors.New(msgInvalidParametersEC)
		}
		namedCurve := namedCurveFromOID(*namedCurveOID)
		if namedCurve == nil {
			return nil, true, errors.New(msgUnknownEllipticCurve)
		}
		x, y := elliptic.Unmarshal(namedCurve, der) //nolint:staticcheck
		if x == nil {
			return nil, true, errors.New(msgFailedToUnmarshalEllipticPoint)
		}
		pub := &eckcdsa.PublicKey{
			Curve: namedCurve,
			X:     x,
			Y:     y,
		}
		return pub, true, nil

	case oid.Equal(oidPublicKeyKCDSA):
		y := new(big.Int)
		if !der.ReadASN1Integer(y) {
			return nil, true, errors.New(msgInvalidPublicKeyY)
		}
		pub := &kcdsa.PublicKey{
			Y: y,
			Parameters: kcdsa.Parameters{
				P: new(big.Int),
				Q: new(big.Int),
				G: new(big.Int),
			},
		}
		paramsDer := cryptobyte.String(params.FullBytes)
		// TODO: Read KCDSA Parameters J, Seed, Count
		if !paramsDer.ReadASN1(&paramsDer, cryptobyte_asn1.SEQUENCE) ||
			!paramsDer.ReadASN1Integer(pub.Parameters.P) ||
			!paramsDer.ReadASN1Integer(pub.Parameters.Q) ||
			!paramsDer.ReadASN1Integer(pub.Parameters.G) {
			return nil, true, errors.New(msgInvalidParametersDSA)
		}
		if pub.Y.Sign() <= 0 || pub.Parameters.P.Sign() <= 0 ||
			pub.Parameters.Q.Sign() <= 0 || pub.Parameters.G.Sign() <= 0 {
			return nil, true, errors.New(msgZeroOrNegativeParameterDSA)
		}

		// TODO: Read KCDSA Parameters J, Seed, Count
		J := new(big.Int)
		seed := make([]byte, 0, 32)
		var count int
		if paramsDer.ReadASN1Integer(J) &&
			paramsDer.ReadASN1Bytes(&seed, cryptobyte_asn1.OCTET_STRING) &&
			paramsDer.ReadASN1Integer(&count) {
			pub.Parameters.GenParameters = kcdsa.GenerationParameters{
				J:     J,
				Seed:  seed,
				Count: count,
			}
		}

		return pub, true, nil

	default:
		return nil, false, nil
	}
}
