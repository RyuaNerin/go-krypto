package internal

import (
	"encoding"
	"hash"
)

type Hash interface {
	hash.Hash
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	encoding.BinaryAppender
}
