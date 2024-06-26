package kx509

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	kcdsainternal "github.com/RyuaNerin/go-krypto/internal/kcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L21-L26
type pkcs8 struct {
	Version    int
	Algo       pkix.AlgorithmIdentifier
	PrivateKey []byte
	// optional attributes omitted.
}

// ParsePKCS8PrivateKey parses an unencrypted private key in PKCS #8, ASN.1 DER form.
//
// It returns an *eckcdsa.PrivateKey, an *kcdsa.PrivateKey,
// or the result of crypto/x509.ParsePKCS8PrivateKey
//
// This kind of key is commonly encoded in PEM blocks of type "PRIVATE KEY".
func ParsePKCS8PrivateKey(der []byte) (key interface{}, err error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L35-L45
	var privKey pkcs8
	if _, err := asn1.Unmarshal(der, &privKey); err != nil {
		if _, err := asn1.Unmarshal(der, &eckcPrivateKey{}); err == nil {
			return nil, errors.New(msgUseParseECKCPrivateKey)
		}
		return nil, err
	}

	switch {
	case privKey.Algo.Algorithm.Equal(oidPublicKeyECKCDSA):
		fallthrough
	case privKey.Algo.Algorithm.Equal(oidPublicKeyECKCDSAAlteGOV):
		// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L54-L64
		bytes := privKey.Algo.Parameters.FullBytes
		namedCurveOID := new(asn1.ObjectIdentifier)
		if _, err := asn1.Unmarshal(bytes, namedCurveOID); err != nil {
			namedCurveOID = nil
		}
		key, err = parseECKCPrivateKey(namedCurveOID, privKey.PrivateKey)
		if err != nil {
			return nil, err
		}
		return key, nil

	case privKey.Algo.Algorithm.Equal(oidPublicKeyKCDSA):
		fallthrough
	case privKey.Algo.Algorithm.Equal(oidPublicKeyKCDSAAlteGOV):
		// Parse X
		X := new(big.Int)
		_, err = asn1.Unmarshal(privKey.PrivateKey, &X)
		if err != nil {
			return nil, errors.New(msgInvalidPrivateKeyX)
		}

		// Parse parameters
		bytes := privKey.Algo.Parameters.FullBytes
		var params kcdsaParameters
		if _, err = asn1.Unmarshal(bytes, &params); err != nil {
			return nil, errors.New(msgInvalidParametersDSA)
		}

		priv := &kcdsa.PrivateKey{
			X: X,
			PublicKey: kcdsa.PublicKey{
				Parameters: kcdsa.Parameters{
					P: params.P,
					Q: params.Q,
					G: params.G,
					// TODO: Read KCDSA Parameters J, Seed, Count
					GenParameters: kcdsa.GenerationParameters{
						J:     params.J,
						Seed:  params.Seed,
						Count: params.Count,
					},
				},
			},
		}

		priv.Y = new(big.Int)
		kcdsainternal.GenerateY(priv.Y, priv.P, priv.Q, priv.G, priv.X)

		return priv, nil

	default:
		return x509.ParsePKCS8PrivateKey(der)
	}
}

func marshalPKCS8PrivateKeyECKCDSA(privKey *pkcs8, k *eckcdsa.PrivateKey) error {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L112-L129
	oid, ok := oidFromNamedCurve(k.Curve)
	if !ok {
		return errors.New(msgUnknownEllipticCurve)
	}
	oidBytes, err := asn1.Marshal(oid)
	if err != nil {
		return fmt.Errorf(msgUnknownEllipticCurveOIDFormat, oid.String())
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

func marshalPKCS8PrivateKeyKCDSA(privKey *pkcs8, k *kcdsa.PrivateKey) error {
	params := kcdsaParameters{
		P: k.P,
		Q: k.Q,
		G: k.G,
	}
	if k.GenParameters.IsValid() {
		params.J = k.GenParameters.J
		params.Seed = k.GenParameters.Seed
		params.Count = k.GenParameters.Count
	}
	paramBytes, err := asn1.Marshal(params)
	if err != nil {
		return errors.New(msgInvalidParametersDSA)
	}

	privKey.Algo = pkix.AlgorithmIdentifier{
		Algorithm: oidPublicKeyKCDSA,
		Parameters: asn1.RawValue{
			FullBytes: paramBytes,
		},
	}
	privKey.PrivateKey, err = asn1.Marshal(k.X)
	if err != nil {
		return errors.New(msgInvalidPublicKey)
	}

	return nil
}

// MarshalPKCS8PrivateKey converts a private key to PKCS #8, ASN.1 DER form.
//
// supported key types : *eckcdsa.PublicKey and *kcdsa.PublicKey.
// for unsupported types, returns the result of crypto/x509.MarshalPKCS8PrivateKey
//
// This kind of key is commonly encoded in PEM blocks of type "PRIVATE KEY".
func MarshalPKCS8PrivateKey(key interface{}) ([]byte, error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L101-L102
	var privKey pkcs8

	var err error
	switch k := key.(type) {
	case eckcdsa.PrivateKey:
		err = marshalPKCS8PrivateKeyECKCDSA(&privKey, &k)
	case *eckcdsa.PrivateKey: //nolint:typecheck
		err = marshalPKCS8PrivateKeyECKCDSA(&privKey, k)

	case kcdsa.PrivateKey:
		err = marshalPKCS8PrivateKeyKCDSA(&privKey, &k)
	case *kcdsa.PrivateKey: //nolint:typecheck
		err = marshalPKCS8PrivateKeyKCDSA(&privKey, k)

	default:
		return x509.MarshalPKCS8PrivateKey(key)
	}

	if err != nil {
		return nil, err
	}
	return asn1.Marshal(privKey)
}
