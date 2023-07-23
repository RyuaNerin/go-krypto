// Package lsh512 implements the LSH-512, LSH-384, LSH-512-256, LSH-512-224 hash algorithms as defined in TTAK.KO-12.0276
package lsh512

import (
	"errors"
)

var ErrInvalidDataBitLen = errors.New("krypto/lsh512: bit level update is not allowed")

const (
	// The size of a LSH-512 checksum in bytes.
	Size = 64
	// The size of a LSH-384 checksum in bytes.
	Size384 = 48
	// The size of a LSH-512-256 checksum in bytes.
	Size256 = 32
	// The size of a LSH-512-224 checksum in bytes.
	Size224 = 28

	// The blocksize of LSH-512, LSH-384, LSH-512-256 and LSH-512-224 in bytes.
	BlockSize = 256
)
