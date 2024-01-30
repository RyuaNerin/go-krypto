package kipher

import "crypto/cipher"

type Stream interface {
	cipher.Stream
	IV() []byte
}

type BlockMode interface {
	cipher.BlockMode
	IV() []byte
}

type kryptoBlock interface {
	cipher.Block

	Encrypt4(dst, src []byte)
	Decrypt4(dst, src []byte)

	Encrypt8(dst, src []byte)
	Decrypt8(dst, src []byte)
}
