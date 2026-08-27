// Package tblockmode provides reusable contract tests, fuzz targets, and
// benchmarks for implementations of cipher.BlockMode, cipher.Stream, and
// cipher.AEAD.
package tblockmode

import "crypto/cipher"

type (
	// NewBlockMode constructs a fresh block mode using iv. Implementations must
	// copy iv if they need it after the constructor returns.
	NewBlockMode func(iv []byte) cipher.BlockMode

	// NewStream constructs a fresh stream using iv. Implementations must copy iv
	// if they need it after the constructor returns.
	NewStream func(iv []byte) cipher.Stream

	// NewAEAD constructs a fresh, fully configured AEAD.
	NewAEAD func() (cipher.AEAD, error)
)

type (
	BlockModeInfo struct {
		BlockSize int
		IVSize    int
	}

	StreamInfo struct {
		IVSize int
	}

	AEADInfo struct {
		NonceSize int
		Overhead  int
	}
)
