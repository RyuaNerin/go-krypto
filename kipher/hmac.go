package kipher

import (
	"crypto/hmac"
	"hash"
)

func NewHMAC(h func() hash.Hash, key []byte) hash.Hash {
	return hmac.New(h, key)
}
