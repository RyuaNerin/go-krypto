package kcdsa

import "errors"

const (
	msgInvalidParameterSizes = "krypto/kcdsa: invalid ParameterSizes"
	msgInvalidSignerOpts     = "krypto/kcdsa: opts must be *kcdsa.SignerOpts"
)

var (
	ErrInvalidParameters           = errors.New("krypto/kcdsa: invalid domain parameters")
	ErrInvalidPublicKey            = errors.New("krypto/kcdsa: invalid public key")
	ErrInvalidPrivateKey           = errors.New("krypto/kcdsa: invalid private key")
	ErrInvalidGenerationParameters = errors.New("krypto/kcdsa: invalid generation parameters")
	ErrInvalidParameterSizes       = errors.New("krypto/kcdsa: invalid ParameterSizes")
	ErrErrorParametersNotSetUp     = errors.New("krypto/kcdsa: parameters not set up before generating key")
	ErrErrorShortXKEY              = errors.New("krypto/kcdsa: XKEY is too small")
	ErrInvalidInteger              = errors.New("krypto/kcdsa: invalid integer")
	ErrInvalidASN1                 = errors.New("krypto/kcdsa: invalid ASN.1")
	ErrInvalidSignerOpts           = errors.New(msgInvalidSignerOpts)
)
