package kipher

import "crypto/cipher"

type kryptoBlock interface {
	cipher.Block

	Encrypt4(dst, src []byte)
	Decrypt4(dst, src []byte)

	Encrypt8(dst, src []byte)
	Decrypt8(dst, src []byte)
}
