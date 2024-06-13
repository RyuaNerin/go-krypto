package kipher

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal/alias"
	ikipher "github.com/RyuaNerin/go-krypto/internal/kipher"
)

func NewECBEncryptor(b cipher.Block) cipher.BlockMode {
	kb := ikipher.WrapKipher(b)
	bs := kb.BlockSize()

	return &ecb{
		b:        kb,
		bs1:      1 * bs,
		bs4:      4 * bs,
		bs8:      8 * bs,
		process1: kb.Encrypt,
		process4: kb.Encrypt4,
		process8: kb.Encrypt8,
	}
}

func NewECBDecryptor(b cipher.Block) cipher.BlockMode {
	kb := ikipher.WrapKipher(b)
	bs := kb.BlockSize()

	return &ecb{
		b:        kb,
		bs1:      1 * bs,
		bs4:      4 * bs,
		bs8:      8 * bs,
		process1: kb.Decrypt,
		process4: kb.Decrypt4,
		process8: kb.Decrypt8,
	}
}

type ecb struct {
	b             ikipher.Block
	bs1, bs4, bs8 int
	process1      func(dst, src []byte)
	process4      func(dst, src []byte)
	process8      func(dst, src []byte)
}

func (ecb *ecb) BlockSize() int {
	return ecb.b.BlockSize()
}

func (ecb *ecb) CryptBlocks(dst, src []byte) {
	if len(src)%ecb.bs1 != 0 {
		panic(msgNotFullBlocks)
	}
	if len(dst) < len(src) {
		panic(msgSmallDst)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgBufferOverlap)
	}

	for len(src) >= ecb.bs8 {
		ecb.process8(dst, src)

		dst, src = dst[ecb.bs8:], src[ecb.bs8:]
	}

	for len(src) >= ecb.bs4 {
		ecb.process4(dst, src)

		dst, src = dst[ecb.bs4:], src[ecb.bs4:]
	}

	for len(src) > 0 {
		ecb.process1(dst, src)

		dst, src = dst[ecb.bs1:], src[ecb.bs1:]
	}
}
