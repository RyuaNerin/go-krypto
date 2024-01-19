package kx509

import (
	"crypto/elliptic"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"
	"golang.org/x/crypto/cryptobyte"
	cryptobyte_asn1 "golang.org/x/crypto/cryptobyte/asn1"
)

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/x509.go#L58
type pkixPublicKey struct {
	Algo      pkix.AlgorithmIdentifier
	BitString asn1.BitString
}

// ParsePKIXPublicKey parses a public key in PKIX, ASN.1 DER form. The encoded
// public key is a SubjectPublicKeyInfo structure (see RFC 5280, Section 4.1).
//
// It returns a *eckcdsa.PublicKey, or *kcdsa.PublicKey.
//
// This kind of key is commonly encoded in PEM blocks of type "PUBLIC KEY".
func ParsePKIXPublicKey(derBytes []byte) (pub interface{}, err error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/x509.go#L71
	var pki publicKeyInfo
	if rest, err := asn1.Unmarshal(derBytes, &pki); err != nil {
		return nil, err
	} else if len(rest) != 0 {
		return nil, errors.New("kx509: trailing data after ASN.1 of public-key")
	}
	return parsePublicKey(&pki)
}

// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/x509.go#L84
func marshalPublicKey(pub interface{}) (publicKeyBytes []byte, publicKeyAlgorithm pkix.AlgorithmIdentifier, err error) {
	switch pub := pub.(type) {
	case *eckcdsa.PublicKey:
		oid, ok := oidFromNamedCurve(pub.Curve)
		if !ok {
			return nil, pkix.AlgorithmIdentifier{}, errors.New("kx509: unsupported elliptic curve")
		}
		if !pub.Curve.IsOnCurve(pub.X, pub.Y) {
			return nil, pkix.AlgorithmIdentifier{}, errors.New("kx509: invalid elliptic curve public key")
		}
		publicKeyBytes = elliptic.Marshal(pub.Curve, pub.X, pub.Y)
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
			return nil, pkix.AlgorithmIdentifier{}, errors.New("kx509: invalid public key")
		}

		var paramBytes []byte
		paramBytes, err = asn1.Marshal(kcdsaParameters{
			P: pub.P,
			Q: pub.Q,
			G: pub.G,
			// TODO: Read KCDSA Parameters J, Seed, Count
		})
		if err != nil {
			return nil, pkix.AlgorithmIdentifier{}, errors.New("kx509: invalid paramerter")
		}

		publicKeyAlgorithm.Algorithm = oidPublicKeyKCDSA
		publicKeyAlgorithm.Parameters.FullBytes = paramBytes

	default:
		return nil, pkix.AlgorithmIdentifier{}, fmt.Errorf("kx509: unsupported public key type: %T", pub)
	}

	return publicKeyBytes, publicKeyAlgorithm, nil
}

func addASN1IntBytes(b *cryptobyte.Builder, bytes []byte) {
	for len(bytes) > 0 && bytes[0] == 0 {
		bytes = bytes[1:]
	}
	if len(bytes) == 0 {
		b.SetError(errors.New("invalid integer"))
		return
	}
	b.AddASN1(cryptobyte_asn1.INTEGER, func(c *cryptobyte.Builder) {
		if bytes[0]&0x80 != 0 {
			c.AddUint8(0)
		}
		c.AddBytes(bytes)
	})
}

// MarshalPKIXPublicKey converts a public key to PKIX, ASN.1 DER form.
// The encoded public key is a SubjectPublicKeyInfo structure
// (see RFC 5280, Section 4.1).
//
// The following key types are currently supported: *eckcdsa.PublicKey,
// *kcdsa.PublicKey.
//
// This kind of key is commonly encoded in PEM blocks of type "PUBLIC KEY".
func MarshalPKIXPublicKey(pub interface{}) ([]byte, error) {
	// https://github.com/golang/go/blob/go1.21.6/src/crypto/x509/x509.go#L150
	var publicKeyBytes []byte
	var publicKeyAlgorithm pkix.AlgorithmIdentifier
	var err error

	if publicKeyBytes, publicKeyAlgorithm, err = marshalPublicKey(pub); err != nil {
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
