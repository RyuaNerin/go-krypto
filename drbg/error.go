package drbg

import (
	"errors"
	"fmt"

	"github.com/RyuaNerin/go-krypto/internal/drbg/ctrdrbg"
)

var (
	ErrErrorUninstantiated  = errors.New("krypto/drbg: uninstantiated")
	ErrCTROptionUnsupported = errors.New("krypto/drbg: CTR option is unsupported")
	ErrEntropyLengthInvalid = errors.New("krypto/drbg: invalid entropy length")
)

type (
	InvalidEntropyError struct {
		Min uint64
		Max uint64
	}
	InvalidNonceError struct {
		Min int // Minimum nonce length in bytes for the supplied entropy input.
	}
	InvalidReseedIntervalError          uint64 // Maximum permitted interval.
	InvalidKeySizeError                 int
	InvalidBlockSizeError               int
	InvalidStrengthError                int
	PersonalizationStringIsTooLongError uint64
	TooManyBytesRequestedError          uint64
	AdditionalInputIsTooLongError       uint64
	InvalidCtrLengthError               struct {
		Min int
		Max int
	}
)

// DerivationInputIsTooLongError reports the maximum combined input length,
// in bytes, representable by the CTR derivation function's 32-bit L field.
type DerivationInputIsTooLongError = ctrdrbg.DerivationInputIsTooLongError

var (
	_ error = InvalidEntropyError{}
	_ error = InvalidNonceError{}
	_ error = InvalidReseedIntervalError(0)
	_ error = DerivationInputIsTooLongError(0)
	_ error = InvalidKeySizeError(0)
	_ error = InvalidBlockSizeError(0)
	_ error = InvalidStrengthError(0)
	_ error = PersonalizationStringIsTooLongError(0)
	_ error = TooManyBytesRequestedError(0)
	_ error = AdditionalInputIsTooLongError(0)
	_ error = InvalidCtrLengthError{}
)

func (e InvalidEntropyError) Error() string {
	return fmt.Sprintf("krypto/drbg: invalid entropy size. must be between %d and %d bytes", e.Min, e.Max)
}

func (e InvalidNonceError) Error() string {
	return fmt.Sprintf("krypto/drbg: invalid nonce size. minimum length is %d bytes for this entropy input", e.Min)
}

func (e InvalidReseedIntervalError) Error() string {
	return fmt.Sprintf("krypto/drbg: invalid reseed interval. maximum is %d requests", e)
}

func (e InvalidKeySizeError) Error() string {
	return fmt.Sprintf("krypto/drbg: invalid key size %d", e)
}

func (e InvalidBlockSizeError) Error() string {
	return fmt.Sprintf("krypto/drbg: invalid block size %d bytes. must be 8 or 16 bytes", e)
}

func (e InvalidStrengthError) Error() string {
	return fmt.Sprintf("krypto/drbg: invalid strength. maximum length is %d bytes", e)
}

func (e PersonalizationStringIsTooLongError) Error() string {
	return fmt.Sprintf("krypto/drbg: personalization_string is too long. maximum length is %d bytes", e)
}

func (e TooManyBytesRequestedError) Error() string {
	return fmt.Sprintf("krypto/drbg: too many bytes requested. maximum length is %d bytes", e)
}

func (e AdditionalInputIsTooLongError) Error() string {
	return fmt.Sprintf("krypto/drbg: additionalInput is too long. maximum length is %d bytes", e)
}

func (e InvalidCtrLengthError) Error() string {
	return fmt.Sprintf("krypto/drbg: invalid counter length. must be between %d and %d", e.Min, e.Max)
}
