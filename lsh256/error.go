package lsh256

import "errors"

const (
	msgInvalidSize          = "krypto/lsh256: invalid size"
	msgInvalidInternalState = "krypto/lsh256: invalid internal state"
)

var (
	ErrInvalidHashStateIdentifier = errors.New("krypto/lsh256: invalid hash state identifier")
	ErrInvalidHashStateSize       = errors.New("krypto/lsh256: invalid hash state size")
)
