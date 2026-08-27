package gmac

import "errors"

const (
	msgInvalidInternalState = "krypto/gmac: invalid internal state"
)

var (
	ErrRequired128BitBlockCipher  = errors.New("krypto/gmac: requires 128-bit block cipher")
	ErrInvalidIV                  = errors.New("krypto/gmac: IV must be nonempty and unique for each message")
	ErrInvalidTagSizeGMAC         = errors.New("krypto/gmac: tagSize must be between 1 and 16")
	ErrInvalidHashStateIdentifier = errors.New("krypto/gmac: invalid hash state identifier")
	ErrInvalidHashStateSize       = errors.New("krypto/gmac: invalid hash state size")
	ErrIncompatibleHashState      = errors.New("krypto/gmac: incompatible hash state implementation")
)
