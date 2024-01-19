package kx509

import (
	"crypto/elliptic"
	"encoding/asn1"
	"errors"
	"math/big"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"

	"golang.org/x/crypto/cryptobyte"
	cryptobyte_asn1 "golang.org/x/crypto/cryptobyte/asn1"
)

func parsePublicKey(keyData *publicKeyInfo) (interface{}, error) {
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
			return nil, errors.New("kx509: invalid ECDSA parameters")
		}
		namedCurve := namedCurveFromOID(*namedCurveOID)
		if namedCurve == nil {
			return nil, errors.New("kx509: unsupported elliptic curve")
		}
		x, y := elliptic.Unmarshal(namedCurve, der)
		if x == nil {
			return nil, errors.New("kx509: failed to unmarshal elliptic curve point")
		}
		pub := &eckcdsa.PublicKey{
			Curve: namedCurve,
			X:     x,
			Y:     y,
		}
		return pub, nil

	case oid.Equal(oidPublicKeyKCDSA):
		y := new(big.Int)
		if !der.ReadASN1Integer(y) {
			return nil, errors.New("kx509: invalid KCDSA public key")
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
			return nil, errors.New("kx509: invalid KCDSA parameters")
		}
		if pub.Y.Sign() <= 0 || pub.Parameters.P.Sign() <= 0 ||
			pub.Parameters.Q.Sign() <= 0 || pub.Parameters.G.Sign() <= 0 {
			return nil, errors.New("kx509: zero or negative KCDSA parameter")
		}
		return pub, nil

	default:
		return nil, errors.New("kx509: unknown public key algorithm")
	}
}
