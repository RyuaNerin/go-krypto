package kx509

import (
	"crypto/elliptic"
	"encoding/asn1"
	"math/big"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	eckcdsa1 "github.com/RyuaNerin/go-krypto/internal/eckcdsa"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

type eckcdsaPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

// ParseECKCPrivateKey parses an EC private key in SEC 1, ASN.1 DER form.
//
// This kind of key is commonly encoded in PEM blocks of type "EC PRIVATE KEY".
//
// warning: this is non-normative
//
// The returned *eckcdsa.PrivateKey implements Zeroize(): call it to erase the
// private key when it is no longer needed.
func ParseECKCPrivateKey(der []byte) (*eckcdsa.PrivateKey, error) {
	return parseECKCPrivateKey(nil, der)
}

// MarshalECKCPrivateKey converts an EC private key to SEC 1, ASN.1 DER form.
//
// This kind of key is commonly encoded in PEM blocks of type "EC PRIVATE KEY".
// For a more flexible key format which is not EC specific, use
// MarshalPKCS8PrivateKey.
//
// warning: this is non-normative
func MarshalECKCPrivateKey(key *eckcdsa.PrivateKey) ([]byte, error) {
	if key == nil || key.Curve == nil {
		return nil, ErrUnknownEllipticCurve
	}
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/sec1.go#L46-L53
	oid, ok := oidFromNamedCurve(key.Curve)
	if !ok {
		return nil, ErrUnknownEllipticCurve
	}

	return marshalECKCPrivateKeyWithOID(key, oid)
}

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/sec1.go#L55-L68
func marshalECKCPrivateKeyWithOID(key *eckcdsa.PrivateKey, oid asn1.ObjectIdentifier) ([]byte, error) {
	if key == nil || key.Curve == nil || key.D == nil {
		return nil, ErrInvalidPrivateKeyValue
	}
	curveOrder := key.Curve.Params().N
	if key.D.Sign() <= 0 || key.D.Cmp(curveOrder) >= 0 {
		return nil, ErrInvalidPrivateKeyValue
	}
	if key.X == nil || key.Y == nil || !key.Curve.IsOnCurve(key.X, key.Y) {
		return nil, ErrInvalidPublicKey
	}
	privateKey := make([]byte, (curveOrder.BitLen()+7)/8)
	return asn1.Marshal(eckcdsaPrivateKey{
		Version:       ecPrivKeyVersion,
		PrivateKey:    key.D.FillBytes(privateKey),
		NamedCurveOID: oid,
		PublicKey:     asn1.BitString{Bytes: elliptic.Marshal(key.Curve, key.X, key.Y)},
	})
}

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/sec1.go#L84-L136
func parseECKCPrivateKey(namedCurveOID *asn1.ObjectIdentifier, der []byte) (key *eckcdsa.PrivateKey, err error) {
	var privKey eckcdsaPrivateKey
	// privKey.PrivateKey is only set by Unmarshal, so defer a closure.
	defer func() { clear(privKey.PrivateKey) }()
	rest, err := asn1.Unmarshal(der, &privKey)
	if err != nil || len(rest) != 0 {
		if rest, pkcs8Err := asn1.Unmarshal(der, &pkcds8PrivateKey{}); pkcs8Err == nil && len(rest) == 0 {
			return nil, ErrUseParsePKCS8PrivateKey
		}
		if err != nil {
			return nil, err
		}
		return nil, ErrInvalidECPrivateKey
	}
	if privKey.Version != ecPrivKeyVersion {
		return nil, UnknownECPrivateKeyVersionError(privKey.Version)
	}

	var curve elliptic.Curve
	if namedCurveOID != nil {
		curve = namedCurveFromOID(*namedCurveOID)
	} else {
		curve = namedCurveFromOID(privKey.NamedCurveOID)
	}
	if curve == nil {
		return nil, ErrUnknownEllipticCurve
	}

	k := new(big.Int).SetBytes(privKey.PrivateKey)
	curveOrder := curve.Params().N
	if k.Sign() <= 0 || k.Cmp(curveOrder) >= 0 {
		return nil, ErrInvalidPrivateKeyValue
	}

	priv := new(eckcdsa.PrivateKey)
	priv.Curve = curve
	priv.D = k

	privateKey := make([]byte, (curveOrder.BitLen()+7)/8)
	defer zeroize.Zeroize(privateKey)

	// Some private keys have leading zero padding. This is invalid
	// according to [SEC1], but this code will ignore it.
	for len(privKey.PrivateKey) > len(privateKey) {
		if privKey.PrivateKey[0] != 0 {
			priv.Zeroize()
			return nil, ErrInvalidPrivateKeyLength
		}
		privKey.PrivateKey = privKey.PrivateKey[1:]
	}

	// Some private keys remove all leading zeros, this is also invalid
	// according to [SEC1] but since OpenSSL used to do this, we ignore
	// this too.
	copy(privateKey[len(privateKey)-len(privKey.PrivateKey):], privKey.PrivateKey)

	d := new(big.Int).SetBytes(privateKey)
	defer zeroize.Zeroize(d)
	priv.PublicKey.X, priv.PublicKey.Y = eckcdsa1.GenerateXY(d, curve)

	return priv, nil
}
