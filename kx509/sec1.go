package kx509

import (
	"crypto/elliptic"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/internal"
)

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/sec1.go#L27-L32
type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

// ParseECPrivateKey parses an EC private key in SEC 1, ASN.1 DER form.
//
// This kind of key is commonly encoded in PEM blocks of type "EC PRIVATE KEY".
func ParseECPrivateKey(der []byte) (*eckcdsa.PrivateKey, error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/sec1.go#L37-L39
	return parseECPrivateKey(nil, der)
}

// MarshalECPrivateKey converts an EC private key to SEC 1, ASN.1 DER form.
//
// This kind of key is commonly encoded in PEM blocks of type "EC PRIVATE KEY".
// For a more flexible key format which is not EC specific, use
// MarshalPKCS8PrivateKey.
func MarshalECPrivateKey(key *eckcdsa.PrivateKey) ([]byte, error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/sec1.go#L46-L53
	oid, ok := oidFromNamedCurve(key.Curve)
	if !ok {
		return nil, errors.New("kx509: unknown elliptic curve")
	}

	return marshalECPrivateKeyWithOID(key, oid)
}

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/sec1.go#L55-L68
func marshalECPrivateKeyWithOID(key *eckcdsa.PrivateKey, oid asn1.ObjectIdentifier) ([]byte, error) {
	if !key.Curve.IsOnCurve(key.X, key.Y) {
		return nil, errors.New("kx509: invalid elliptic key public key")
	}
	privateKey := make([]byte, (key.Curve.Params().N.BitLen()+7)/8)
	return asn1.Marshal(ecPrivateKey{
		Version:       1,
		PrivateKey:    key.D.FillBytes(privateKey),
		NamedCurveOID: oid,
		PublicKey:     asn1.BitString{Bytes: elliptic.Marshal(key.Curve, key.X, key.Y)},
	})
}

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/sec1.go#L84-L136
func parseECPrivateKey(namedCurveOID *asn1.ObjectIdentifier, der []byte) (key *eckcdsa.PrivateKey, err error) {
	var privKey ecPrivateKey
	if _, err := asn1.Unmarshal(der, &privKey); err != nil {
		if _, err := asn1.Unmarshal(der, &pkcs8{}); err == nil {
			return nil, errors.New("kx509: failed to parse private key (use ParsePKCS8PrivateKey instead for this key format)")
		}
		return nil, errors.New("kx509: failed to parse EC private key: " + err.Error())
	}
	if privKey.Version != ecPrivKeyVersion {
		return nil, fmt.Errorf("kx509: unknown EC private key version %d", privKey.Version)
	}

	var curve elliptic.Curve
	if namedCurveOID != nil {
		curve = namedCurveFromOID(*namedCurveOID)
	} else {
		curve = namedCurveFromOID(privKey.NamedCurveOID)
	}
	if curve == nil {
		return nil, errors.New("kx509: unknown elliptic curve")
	}

	k := new(big.Int).SetBytes(privKey.PrivateKey)
	curveOrder := curve.Params().N
	if k.Cmp(curveOrder) >= 0 {
		return nil, errors.New("kx509: invalid elliptic curve private key value")
	}

	priv := new(eckcdsa.PrivateKey)
	priv.Curve = curve
	priv.D = k

	privateKey := make([]byte, (curveOrder.BitLen()+7)/8)

	// Some private keys have leading zero padding. This is invalid
	// according to [SEC1], but this code will ignore it.
	for len(privKey.PrivateKey) > len(privateKey) {
		if privKey.PrivateKey[0] != 0 {
			return nil, errors.New("kx509: invalid private key length")
		}
		privKey.PrivateKey = privKey.PrivateKey[1:]
	}

	// Some private keys remove all leading zeros, this is also invalid
	// according to [SEC1] but since OpenSSL used to do this, we ignore
	// this too.
	copy(privateKey[len(privateKey)-len(privKey.PrivateKey):], privKey.PrivateKey)

	d := new(big.Int).SetBytes(privateKey)
	dInv := internal.FermatInverse(d, curve.Params().N)
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(dInv.Bytes())

	return priv, nil
}
