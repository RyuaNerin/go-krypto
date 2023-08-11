// Package lsh256 implements the LSH-256 hash algorithms as defined in TTAK.KO-12.0276
package lsh256

import (
	"errors"
)

var ErrInvalidDataBitLen = errors.New("krypto/lsh256: bit level update is not allowed")

const (
	// The size of a LSH-256 checksum in bytes.
	Size = 32
	// The size of a LSH-224 checksum in bytes.
	Size224 = 28

	// The blocksize of LSH-256 and LSH-224 in bytes.
	BlockSize = 128
)
