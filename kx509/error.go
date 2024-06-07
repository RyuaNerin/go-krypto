package kx509

const (
	msgUseParseECKCPrivateKey           = "krypto/kx509: failed to parse private key (use ParseECKCPrivateKey instead for this key format)"
	msgUseParsePKCS8PrivateKey          = "krypto/kx509: failed to parse private key (use ParsePKCS8PrivateKey instead for this key format)"
	msgUnknownEllipticCurve             = "krypto/kx509: unknown elliptic curve"
	msgFailedToUnmarshalEllipticPoint   = "krypto/kx509: failed to unmarshal elliptic curve point"
	msgUnknownEllipticCurveOIDFormat    = "krypto/kx509: unknown elliptic curve (OID: %s)"
	msgInvalidPublicKey                 = "krypto/kx509: invalid public key"
	msgInvalidPublicKeyY                = "krypto/kx509: invalid public key Y"
	msgZeroOrNegativeParameterDSA       = "krypto/kx509: zero or negative KCDSA parameter"
	msgInvalidPrivateKeyValue           = "krypto/kx509: invalid private key value"
	msgInvalidPrivateKeyLength          = "krypto/kx509: invalid private key length"
	msgInvalidPrivateKeyX               = "krypto/kx509: invalid private key X"
	msgInvalidPrivateKeyY               = "krypto/kx509: invalid private key Y"
	msgInvalidParametersDSA             = "krypto/kx509: invalid parameters"
	msgInvalidParametersEC              = "krypto/kx509: invalid EC-KCDSA parameters"
	msgUnknownECPrivateKeyVersionFormat = "krypto/kx509: unknown EC private key version %d"
)
