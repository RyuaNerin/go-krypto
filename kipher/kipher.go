package kipher

import (
	"crypto/cipher"
	"crypto/subtle"
)

type kryptoBlock interface {
	cipher.Block

	Encrypt4(dst, src []byte)
	Decrypt4(dst, src []byte)

	Encrypt8(dst, src []byte)
	Decrypt8(dst, src []byte)
}

func Equal(mac1, mac2 []byte) bool {
	return subtle.ConstantTimeCompare(mac1, mac2) == 1
}
