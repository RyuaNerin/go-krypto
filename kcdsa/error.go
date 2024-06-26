package kcdsa

const (
	msgInvalidPublicKey            = "krypto/kcdsa: invalid public key"
	msgInvalidGenerationParameters = "krypto/kcdsa: invalid generation parameters"
	msgInvalidParameterSizes       = "krypto/kcdsa: invalid ParameterSizes"
	msgErrorParametersNotSetUp     = "krypto/kcdsa: parameters not set up before generating key"
	msgErrorShortXKEY              = "krypto/kcdsa: XKEY is too small."
	msgInvalidInteger              = "krypto/kcdsa: invalid integer"
	msgInvalidASN1                 = "krypto/kcdsa: invalid ASN.1"
	msgInvalidSignerOpts           = "kcdsa: opts must be *kcdsa.SignerOpts"
)
