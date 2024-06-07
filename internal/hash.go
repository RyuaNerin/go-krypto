package internal

import (
	"hash"
)

type Hash interface {
	hash.Hash
	// encoding.BinaryMarshaler
	MarshalBinary() (data []byte, err error)
	// encoding.BinaryUnmarshaler
	UnmarshalBinary(data []byte) error
}
