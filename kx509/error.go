package kx509

import (
	"errors"
	"fmt"
)

var (
	ErrUseParseECKCPrivateKey         = errors.New("krypto/kx509: failed to parse private key (use ParseECKCDSAPrivateKey instead for this key format)")
	ErrUseParsePKCS8PrivateKey        = errors.New("krypto/kx509: failed to parse private key (use ParsePKCS8PrivateKey instead for this key format)")
	ErrInvalidPKCS8PrivateKey         = errors.New("krypto/kx509: invalid PKCS#8 private key")
	ErrInvalidECPrivateKey            = errors.New("krypto/kx509: invalid EC private key")
	ErrUnknownEllipticCurve           = errors.New("krypto/kx509: unknown elliptic curve")
	ErrFailedToUnmarshalEllipticPoint = errors.New("krypto/kx509: failed to unmarshal elliptic curve point")
	ErrInvalidPublicKey               = errors.New("krypto/kx509: invalid public key")
	ErrInvalidPublicKeyY              = errors.New("krypto/kx509: invalid public key Y")
	ErrInvalidPrivateKeyValue         = errors.New("krypto/kx509: invalid private key value")
	ErrInvalidPrivateKeyLength        = errors.New("krypto/kx509: invalid private key length")
	ErrInvalidPrivateKeyX             = errors.New("krypto/kx509: invalid private key X")
	ErrInvalidParametersDSA           = errors.New("krypto/kx509: invalid parameters")
	ErrInvalidParametersEC            = errors.New("krypto/kx509: invalid EC-KCDSA parameters")
)

type (
	UnknownECPrivateKeyVersionError    int
	UnknownPKCS8PrivateKeyVersionError int
	UnknownEllipticCurveOIDFormatError string
)

var (
	_ error = UnknownECPrivateKeyVersionError(0)
	_ error = UnknownPKCS8PrivateKeyVersionError(0)
	_ error = UnknownEllipticCurveOIDFormatError("")
)

func (e UnknownECPrivateKeyVersionError) Error() string {
	return fmt.Sprintf("krypto/kx509: unknown PKCS#8 EC private key version %d", int(e))
}

func (e UnknownPKCS8PrivateKeyVersionError) Error() string {
	return fmt.Sprintf("krypto/kx509: unknown PKCS#8 private key version %d", int(e))
}

func (e UnknownEllipticCurveOIDFormatError) Error() string {
	return fmt.Sprintf("krypto/kx509: unknown elliptic curve (OID: %s)", string(e))
}
