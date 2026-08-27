package ctrdrbg

import "fmt"

// DerivationInputIsTooLongError reports the maximum combined input size in bytes.
type DerivationInputIsTooLongError uint64

func (e DerivationInputIsTooLongError) Error() string {
	return fmt.Sprintf("krypto/drbg: derivation input is too long. maximum combined length is %d bytes", e)
}
