package has160

import "errors"

const (
	msgInvalidInternalState = "krypto/has160: invalid internal state"
)

var (
	ErrInvalidHashStateSize       = errors.New("krypto/has160: invalid hash state size")
	ErrInvalidHashStateIdentifier = errors.New("krypto/has160: invalid hash state identifier")
)
