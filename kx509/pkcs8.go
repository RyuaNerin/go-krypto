package kx509

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	kcdsa1 "github.com/RyuaNerin/go-krypto/internal/kcdsa"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L21-L26
type pkcds8PrivateKey struct {
	Version    int
	Algo       pkix.AlgorithmIdentifier
	PrivateKey []byte
}

// ParsePKCS8PrivateKey parses an unencrypted private key in PKCS #8, ASN.1 DER form.
//
// It returns an *eckcdsa.PrivateKey, a *kcdsa.PrivateKey,
// or the result of crypto/x509.ParsePKCS8PrivateKey.
//
// This kind of key is commonly encoded in PEM blocks of type "PRIVATE KEY".
// Standard id-ecPublicKey keys are parsed by crypto/x509 as ECDSA.
//
// The returned *kcdsa.PrivateKey / *eckcdsa.PrivateKey implements Zeroize():
// call it to erase the private key when it is no longer needed.
func ParsePKCS8PrivateKey(der []byte) (key any, err error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L35-L45
	var privKey pkcds8PrivateKey
	rest, unmarshalErr := asn1.Unmarshal(der, &privKey)
	if unmarshalErr != nil {
		return x509.ParsePKCS8PrivateKey(der)
	}
	if !isECKCDSAPublicKeyAlgorithm(privKey.Algo.Algorithm) &&
		!isKCDSAPublicKeyAlgorithm(privKey.Algo.Algorithm) {
		return x509.ParsePKCS8PrivateKey(der)
	}
	if len(rest) != 0 {
		return nil, ErrInvalidPKCS8PrivateKey
	}
	if privKey.Version != 0 {
		return nil, UnknownPKCS8PrivateKeyVersionError(privKey.Version)
	}

	switch {
	case isECKCDSAPublicKeyAlgorithm(privKey.Algo.Algorithm):
		// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L54-L64
		bytes := privKey.Algo.Parameters.FullBytes
		namedCurveOID := new(asn1.ObjectIdentifier)
		if rest, oidErr := asn1.Unmarshal(bytes, namedCurveOID); oidErr != nil || len(rest) != 0 {
			namedCurveOID = nil
		}
		key, err = parseECKCPrivateKey(namedCurveOID, privKey.PrivateKey)
		if err != nil {
			return nil, err
		}
		return key, nil

	case isKCDSAPublicKeyAlgorithm(privKey.Algo.Algorithm):
		// Parse X
		X := new(big.Int)
		rest, unmarshalErr := asn1.Unmarshal(privKey.PrivateKey, &X)
		if unmarshalErr != nil || len(rest) != 0 {
			zeroize.Zeroize(X)
			return nil, ErrInvalidPrivateKeyX
		}

		// Parse parameters
		params, parseErr := parseKCDSAParameters(privKey.Algo.Parameters.FullBytes)
		if parseErr != nil {
			zeroize.Zeroize(X)
			return nil, parseErr
		}
		if X.Sign() <= 0 || X.Cmp(params.Q) >= 0 {
			zeroize.Zeroize(X)
			return nil, ErrInvalidPrivateKeyX
		}

		priv := &kcdsa.PrivateKey{
			X: X,
			PublicKey: kcdsa.PublicKey{
				Parameters: params.parameters(),
			},
		}

		priv.Y = new(big.Int)
		kcdsa1.GenerateY(priv.Y, priv.P, priv.Q, priv.G, priv.X)

		return priv, nil

	default: // fallback
		return x509.ParsePKCS8PrivateKey(der)
	}
}

func marshalPKCS8PrivateKeyECKCDSA(privKey *pkcds8PrivateKey, k *eckcdsa.PrivateKey) error {
	if k == nil || k.Curve == nil {
		return ErrUnknownEllipticCurve
	}
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L112-L129
	oid, ok := oidFromNamedCurve(k.Curve)
	if !ok {
		return ErrUnknownEllipticCurve
	}
	oidBytes, err := asn1.Marshal(oid)
	if err != nil {
		return UnknownEllipticCurveOIDFormatError(oid.String())
	}
	privKey.Algo = pkix.AlgorithmIdentifier{
		Algorithm: oidPublicKeyECKCDSA,
		Parameters: asn1.RawValue{
			FullBytes: oidBytes,
		},
	}
	if privKey.PrivateKey, err = marshalECKCPrivateKeyWithOID(k, nil); err != nil {
		return err
	}

	return nil
}

func marshalPKCS8PrivateKeyKCDSA(privKey *pkcds8PrivateKey, k *kcdsa.PrivateKey) error {
	if k == nil || k.X == nil || k.Q == nil || k.X.Sign() <= 0 || k.X.Cmp(k.Q) >= 0 {
		return ErrInvalidPrivateKeyX
	}
	paramBytes, err := marshalKCDSAParameters(&k.Parameters)
	if err != nil {
		return err
	}

	privKey.Algo = pkix.AlgorithmIdentifier{
		Algorithm: oidPublicKeyKCDSA,
		Parameters: asn1.RawValue{
			FullBytes: paramBytes,
		},
	}
	privKey.PrivateKey, err = asn1.Marshal(k.X)
	if err != nil {
		return ErrInvalidPublicKey
	}

	return nil
}

// MarshalPKCS8PrivateKey converts a private key to PKCS #8, ASN.1 DER form.
//
// supported key types : *eckcdsa.PrivateKey and *kcdsa.PrivateKey.
// for unsupported types, returns the result of crypto/x509.MarshalPKCS8PrivateKey
//
// This kind of key is commonly encoded in PEM blocks of type "PRIVATE KEY".
func MarshalPKCS8PrivateKey(key any) ([]byte, error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L101-L102
	var privKey pkcds8PrivateKey

	var err error
	switch k := key.(type) {
	case eckcdsa.PrivateKey:
		err = marshalPKCS8PrivateKeyECKCDSA(&privKey, &k)
	case *eckcdsa.PrivateKey:
		err = marshalPKCS8PrivateKeyECKCDSA(&privKey, k)

	case kcdsa.PrivateKey:
		err = marshalPKCS8PrivateKeyKCDSA(&privKey, &k)
	case *kcdsa.PrivateKey:
		err = marshalPKCS8PrivateKeyKCDSA(&privKey, k)

	default: // fallback
		return x509.MarshalPKCS8PrivateKey(key)
	}

	if err != nil {
		return nil, err
	}
	return asn1.Marshal(privKey)
}
