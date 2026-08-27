package lsh512

import "errors"

const (
	msgInvalidSize          = "krypto/lsh512: invalid size"
	msgInvalidInternalState = "krypto/lsh512: invalid internal state"
)

var (
	ErrInvalidHashStateIdentifier = errors.New("krypto/lsh512: invalid hash state identifier")
	ErrInvalidHashStateSize       = errors.New("krypto/lsh512: invalid hash state size")
)
