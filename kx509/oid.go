package kx509

import (
	"crypto/elliptic"
	"encoding/asn1"

	"github.com/RyuaNerin/elliptic2/nist"
)

const ecPrivKeyVersion = 1

var (
	oidPublicKeyKCDSA   = asn1.ObjectIdentifier{1, 2, 410, 200004, 1, 1}         // eGOV-C01.0008
	oidPublicKeyECKCDSA = asn1.ObjectIdentifier{1, 2, 410, 200004, 1, 100, 2, 1} // eGOV-C01.0008

	oidPublicKeyKCDSAstd   = asn1.ObjectIdentifier{1, 0, 14888, 3, 0, 2} // {iso(1) standard(0) digital-signature-with-appendix(14888) part3(3) algorithm(0) kcdsa(2)}
	oidPublicKeyECKCDSAstd = asn1.ObjectIdentifier{1, 0, 14888, 3, 0, 5} // {iso(1) standard(0) digital-signature-with-appendix(14888) part3(3) algorithm(0) ec-kcdsa(5)}

	oidNamedCurveP224 = asn1.ObjectIdentifier{1, 3, 132, 0, 33}          // P-224 / secp224r1 / wap-wsg-idm-ecid-wtls12 / ansip224r1
	oidNamedCurveP256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7} // P-256 / secp256r1 / prime256v1
	oidNamedCurveP384 = asn1.ObjectIdentifier{1, 3, 132, 0, 34}          // P-384 / secp384r1 / ansip384r1
	oidNamedCurveP521 = asn1.ObjectIdentifier{1, 3, 132, 0, 35}          // P-521 / secp521r1 / secp521r1 / ansip521r1

	// github.com/RyuaNerin/elliptic2
	oidNamedCurveB233 = asn1.ObjectIdentifier{1, 3, 132, 0, 27} // B-233
	oidNamedCurveB283 = asn1.ObjectIdentifier{1, 3, 132, 0, 17} // B-283
	oidNamedCurveB409 = asn1.ObjectIdentifier{1, 3, 132, 0, 37} // B-409
	oidNamedCurveB571 = asn1.ObjectIdentifier{1, 3, 132, 0, 39} // B-571

	oidNamedCurveK233 = asn1.ObjectIdentifier{1, 3, 132, 0, 26} // K-233
	oidNamedCurveK283 = asn1.ObjectIdentifier{1, 3, 132, 0, 16} // K-283
	oidNamedCurveK409 = asn1.ObjectIdentifier{1, 3, 132, 0, 36} // K-409
	oidNamedCurveK571 = asn1.ObjectIdentifier{1, 3, 132, 0, 38} // K-571
)

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

	// github.com/RyuaNerin/elliptic2
	case oid.Equal(oidNamedCurveB233):
		return nist.B233()
	case oid.Equal(oidNamedCurveB283):
		return nist.B283()
	case oid.Equal(oidNamedCurveB409):
		return nist.B409()
	case oid.Equal(oidNamedCurveB571):
		return nist.B571()

	case oid.Equal(oidNamedCurveK233):
		return nist.K233()
	case oid.Equal(oidNamedCurveK283):
		return nist.K283()
	case oid.Equal(oidNamedCurveK409):
		return nist.K409()
	case oid.Equal(oidNamedCurveK571):
		return nist.K571()
	}
	return nil
}

// https://github.com/golang/go/blob/master/src/crypto/x509/x509.go#L541-L554
func oidFromNamedCurve(curve elliptic.Curve) (asn1.ObjectIdentifier, bool) {
	switch curve {
	case elliptic.P224():
		return oidNamedCurveP224, true
	case elliptic.P256():
		return oidNamedCurveP256, true
	case elliptic.P384():
		return oidNamedCurveP384, true
	case elliptic.P521():
		return oidNamedCurveP521, true

	// github.com/RyuaNerin/elliptic2
	case nist.B233():
		return oidNamedCurveB233, true
	case nist.B283():
		return oidNamedCurveB283, true
	case nist.B409():
		return oidNamedCurveB409, true
	case nist.B571():
		return oidNamedCurveB571, true

	case nist.K233():
		return oidNamedCurveK233, true
	case nist.K283():
		return oidNamedCurveK283, true
	case nist.K409():
		return oidNamedCurveK409, true
	case nist.K571():
		return oidNamedCurveK571, true

	}

	return nil, false
}
