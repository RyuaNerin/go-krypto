package kx509

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"math/big"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/internal"
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
		if _, err := asn1.Unmarshal(der, &ecPrivateKey{}); err == nil {
			return nil, errors.New("kx509: failed to parse private key (use ParseECPrivateKey instead for this key format)")
		}
		return nil, err
	}

	switch {
	case privKey.Algo.Algorithm.Equal(oidPublicKeyECKCDSA):
		// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L54-L64
		bytes := privKey.Algo.Parameters.FullBytes
		namedCurveOID := new(asn1.ObjectIdentifier)
		if _, err := asn1.Unmarshal(bytes, namedCurveOID); err != nil {
			namedCurveOID = nil
		}
		key, err = parseECPrivateKey(namedCurveOID, privKey.PrivateKey)
		if err != nil {
			return nil, errors.New("kx509: failed to parse EC private key embedded in PKCS#8: " + err.Error())
		}
		return key, nil

	case privKey.Algo.Algorithm.Equal(oidPublicKeyKCDSA):
		// Parse X
		X := new(big.Int)
		_, err = asn1.Unmarshal(privKey.PrivateKey, &X)
		if err != nil {
			return nil, errors.New("kx509: invalid x")
		}

		// Parse parameters
		bytes := privKey.Algo.Parameters.FullBytes
		var params kcdsaParameters
		if _, err = asn1.Unmarshal(bytes, &params); err != nil {
			return nil, errors.New("kx509: invalid paramerter")
		}

		priv := &kcdsa.PrivateKey{
			X: X,
			PublicKey: kcdsa.PublicKey{
				Parameters: kcdsa.Parameters{
					P: params.P,
					Q: params.Q,
					G: params.G,
					// TODO: Read KCDSA Parameters J, Seed, Count
					TTAKParams: kcdsa.TTAKParameters{
						J:     params.J,
						Seed:  params.Seed,
						Count: params.Count,
					},
				},
			},
		}

		xInv := internal.FermatInverse(priv.X, priv.Q)
		priv.Y = new(big.Int).Exp(priv.G, xInv, priv.P)

		return priv, nil

	default:
		return x509.ParsePKCS8PrivateKey(der)
	}
}

func marshalPKCS8PrivateKeyECKCDSA(privKey *pkcs8, k *eckcdsa.PrivateKey) error {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/pkcs8.go#L112-L129
	oid, ok := oidFromNamedCurve(k.Curve)
	if !ok {
		return errors.New("kx509: unknown curve while marshaling to PKCS#8")
	}
	oidBytes, err := asn1.Marshal(oid)
	if err != nil {
		return errors.New("kx509: failed to marshal curve OID: " + err.Error())
	}
	privKey.Algo = pkix.AlgorithmIdentifier{
		Algorithm: oidPublicKeyECKCDSA,
		Parameters: asn1.RawValue{
			FullBytes: oidBytes,
		},
	}
	if privKey.PrivateKey, err = marshalECPrivateKeyWithOID(k, nil); err != nil {
		return errors.New("kx509: failed to marshal EC private key while building PKCS#8: " + err.Error())
	}

	return nil
}
func marshalPKCS8PrivateKeyKCDSA(privKey *pkcs8, k *kcdsa.PrivateKey) error {
	paramBytes, err := asn1.Marshal(kcdsaParameters{
		P: k.P,
		Q: k.Q,
		G: k.G,
		// TODO: Read KCDSA Parameters J, Seed, Count
		J:     k.TTAKParams.J,
		Seed:  k.TTAKParams.Seed,
		Count: k.TTAKParams.Count,
	})
	if err != nil {
		return errors.New("kx509: invalid paramerter")
	}

	privKey.Algo = pkix.AlgorithmIdentifier{
		Algorithm: oidPublicKeyKCDSA,
		Parameters: asn1.RawValue{
			FullBytes: paramBytes,
		},
	}
	privKey.PrivateKey, err = asn1.Marshal(k.X)
	if err != nil {
		return errors.New("kx509: invalid public key")
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
	case *eckcdsa.PrivateKey:
		err = marshalPKCS8PrivateKeyECKCDSA(&privKey, k)

	case kcdsa.PrivateKey:
		err = marshalPKCS8PrivateKeyKCDSA(&privKey, &k)
	case *kcdsa.PrivateKey:
		err = marshalPKCS8PrivateKeyKCDSA(&privKey, k)

	default:
		return x509.MarshalPKCS8PrivateKey(key)
	}

	if err != nil {
		return nil, err
	}
	return asn1.Marshal(privKey)
}
