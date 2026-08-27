package kx509

import (
	"crypto/elliptic"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	kcdsa1 "github.com/RyuaNerin/go-krypto/internal/kcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/x509.go#L58
type pkixPublicKey struct {
	Raw       asn1.RawContent
	Algo      pkix.AlgorithmIdentifier
	BitString asn1.BitString
}

// ParsePKIXPublicKey parses a public key in PKIX, ASN.1 DER form.
//
// It returns an *eckcdsa.PublicKey, an *kcdsa.PublicKey,
// or the result of crypto/x509.ParsePKIXPublicKey
//
// This kind of key is commonly encoded in PEM blocks of type "PUBLIC KEY".
// Standard id-ecPublicKey keys are parsed by crypto/x509 as ECDSA.
func ParsePKIXPublicKey(derBytes []byte) (pub any, err error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/x509.go#L71-L82
	var pki pkixPublicKey
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
func marshalPublicKey(pub any) (publicKeyBytes []byte, publicKeyAlgorithm pkix.AlgorithmIdentifier, isSupportedType bool, err error) {
	switch pub := pub.(type) {
	case *eckcdsa.PublicKey:
		if pub == nil || pub.Curve == nil {
			return nil, pkix.AlgorithmIdentifier{}, true, ErrUnknownEllipticCurve
		}
		oid, ok := oidFromNamedCurve(pub.Curve)
		if !ok {
			return nil, pkix.AlgorithmIdentifier{}, true, ErrUnknownEllipticCurve
		}
		if pub.X == nil || pub.Y == nil || !pub.Curve.IsOnCurve(pub.X, pub.Y) {
			return nil, pkix.AlgorithmIdentifier{}, true, ErrInvalidPublicKey
		}
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
		publicKeyAlgorithm.Algorithm = oidPublicKeyECKCDSA
		var paramBytes []byte
		paramBytes, err = asn1.Marshal(oid)
		if err != nil {
			return nil, pkix.AlgorithmIdentifier{}, true, ErrUnknownEllipticCurve
		}
		publicKeyAlgorithm.Parameters.FullBytes = paramBytes

	case *kcdsa.PublicKey:
		if pub == nil {
			return nil, pkix.AlgorithmIdentifier{}, true, ErrInvalidPublicKeyY
		}
		var paramBytes []byte
		paramBytes, err = marshalKCDSAParameters(&pub.Parameters)
		if err != nil {
			return nil, pkix.AlgorithmIdentifier{}, true, err
		}
		if !kcdsa1.IsSubgroupElement(pub.Y, pub.P, pub.Q) {
			return nil, pkix.AlgorithmIdentifier{}, true, ErrInvalidPublicKeyY
		}
		publicKeyBytes, err = asn1.Marshal(pub.Y)
		if err != nil {
			return nil, pkix.AlgorithmIdentifier{}, true, ErrInvalidPublicKeyY
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
func MarshalPKIXPublicKey(pub any) ([]byte, error) {
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

	return asn1.Marshal(pkix)
}

func parsePublicKey(keyData *pkixPublicKey) (publicKey any, isSupportedType bool, err error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/parser.go#L217-L220
	oid := keyData.Algo.Algorithm
	params := keyData.Algo.Parameters
	der := keyData.BitString.Bytes

	switch {
	case isECKCDSAPublicKeyAlgorithm(oid):
		if keyData.BitString.BitLength != len(der)*8 {
			return nil, true, ErrInvalidPublicKey
		}
		// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/parser.go#L252-L271
		namedCurveOID := new(asn1.ObjectIdentifier)
		if rest, err := asn1.Unmarshal(params.FullBytes, namedCurveOID); err != nil || len(rest) != 0 {
			return nil, true, ErrInvalidParametersEC
		}
		namedCurve := namedCurveFromOID(*namedCurveOID)
		if namedCurve == nil {
			return nil, true, ErrUnknownEllipticCurve
		}
		x, y := unmarshalECPoint(namedCurve, der)
		if x == nil {
			return nil, true, ErrFailedToUnmarshalEllipticPoint
		}
		pub := &eckcdsa.PublicKey{
			Curve: namedCurve,
			X:     x,
			Y:     y,
		}
		return pub, true, nil

	case isKCDSAPublicKeyAlgorithm(oid):
		if keyData.BitString.BitLength != len(der)*8 {
			return nil, true, ErrInvalidPublicKeyY
		}
		y := new(big.Int)
		if rest, err := asn1.Unmarshal(der, &y); err != nil || len(rest) != 0 {
			return nil, true, ErrInvalidPublicKeyY
		}
		asn1Params, err := parseKCDSAParameters(params.FullBytes)
		if err != nil {
			return nil, true, err
		}
		pub := &kcdsa.PublicKey{
			Y:          y,
			Parameters: asn1Params.parameters(),
		}
		if !kcdsa1.IsSubgroupElement(pub.Y, pub.P, pub.Q) {
			return nil, true, ErrInvalidPublicKeyY
		}

		return pub, true, nil

	default:
		return nil, false, nil
	}
}

func unmarshalECPoint(curve elliptic.Curve, der []byte) (x, y *big.Int) {
	if len(der) == 0 {
		return nil, nil
	}
	if der[0] == 2 || der[0] == 3 {
		return ellipticUnmarshalCompressed(curve, der)
	}
	return ellipticUnmarshal(curve, der)
}
