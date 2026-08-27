package eckcdsa

import "errors"

var (
	ErrParametersNotSetUp  = errors.New("krypto/eckcdsa: parameters not set up before generating key")
	ErrInvalidPrivateKey   = errors.New("krypto/eckcdsa: invalid private key")
	ErrInvalidSignerOpts   = errors.New("krypto/eckcdsa: opts must select an available hash function")
	ErrInvalidInteger      = errors.New("krypto/eckcdsa: invalid integer")
	ErrInvalidASN1         = errors.New("krypto/eckcdsa: invalid ASN.1")
	ErrInvalidWeakKeyRange = errors.New("krypto/eckcdsa: invalid weak key exclusion range")
)
