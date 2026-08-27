package kx509

import (
	"crypto/elliptic"
	"encoding/asn1"

	"github.com/RyuaNerin/elliptic2"
	"github.com/RyuaNerin/elliptic2/oids"
)

const ecPrivKeyVersion = 1

var (
	oidPublicKeyECKCDSA        = asn1.ObjectIdentifier{1, 0, 14888, 3, 0, 5}            // {iso(1) standard(0) digital-signature-with-appendix(14888) part3(3) algorithm(0) ec-kcdsa(5)}
	oidPublicKeyECKCDSAAlteGOV = asn1.ObjectIdentifier{1, 2, 410, 200004, 1, 100, 2, 1} // EC-KCDSA public key, eGOV-C01.0008

	oidPublicKeyKCDSA         = asn1.ObjectIdentifier{1, 0, 14888, 3, 0, 2}     // {iso(1) standard(0) digital-signature-with-appendix(14888) part3(3) algorithm(0) kcdsa(2)}
	oidPublicKeyKCDSAAlteGOV  = asn1.ObjectIdentifier{1, 2, 410, 200004, 1, 1}  // kcdsa, eGOV-C01.0008
	oidPublicKeyKCDSA1AlteGOV = asn1.ObjectIdentifier{1, 2, 410, 200004, 1, 21} // kcdsa1

	oidNamedCurveP224 = asn1.ObjectIdentifier{1, 3, 132, 0, 33}          // P-224 / secp224r1 / wap-wsg-idm-ecid-wtls12 / ansip224r1
	oidNamedCurveP256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7} // P-256 / secp256r1 / prime256v1
	oidNamedCurveP384 = asn1.ObjectIdentifier{1, 3, 132, 0, 34}          // P-384 / secp384r1 / ansip384r1
	oidNamedCurveP521 = asn1.ObjectIdentifier{1, 3, 132, 0, 35}          // P-521 / secp521r1 / secp521r1 / ansip521r1
)

func isKCDSAPublicKeyAlgorithm(oid asn1.ObjectIdentifier) bool {
	switch {
	case oid.Equal(oidPublicKeyKCDSA):
	case oid.Equal(oidPublicKeyKCDSAAlteGOV):
	case oid.Equal(oidPublicKeyKCDSA1AlteGOV):
	default:
		return false
	}
	return true
}

func isECKCDSAPublicKeyAlgorithm(oid asn1.ObjectIdentifier) bool {
	switch {
	case oid.Equal(oidPublicKeyECKCDSA):
	case oid.Equal(oidPublicKeyECKCDSAAlteGOV):
	default:
		return false
	}
	return true
}

// https://github.com/golang/go/blob/master/src/crypto/x509/x509.go#L527-L539
func namedCurveFromOID(oid asn1.ObjectIdentifier) elliptic.Curve {
	switch {
	case oid.Equal(oidNamedCurveP224):
		return elliptic.P224()
	case oid.Equal(oidNamedCurveP256):
		return elliptic.P256()
	case oid.Equal(oidNamedCurveP384):
		return elliptic.P384()
	case oid.Equal(oidNamedCurveP521):
		return elliptic.P521()
	}

	return oids.GetCurve(oid)
}

// https://github.com/golang/go/blob/master/src/crypto/x509/x509.go#L541-L554
func oidFromNamedCurve(curve elliptic.Curve) (asn1.ObjectIdentifier, bool) {
	if c, ok := curve.(elliptic2.CurveExtended); ok {
		if oid := c.OID(); len(oid) > 0 {
			return oid, true
		}
	}

	switch curve {
	case elliptic.P224():
		return oidNamedCurveP224, true
	case elliptic.P256():
		return oidNamedCurveP256, true
	case elliptic.P384():
		return oidNamedCurveP384, true
	case elliptic.P521():
		return oidNamedCurveP521, true
	}

	return nil, false
}
