package kipher

import (
	"crypto/cipher"
	"fmt"
)

func NewECBEncryptor(b cipher.Block) cipher.BlockMode {
	kb, ok := b.(kryptoBlock)
	if !ok {
		kb = &blockWrap{b}
	}

	return &ecbEnc{
		b:        kb,
		process1: kb.Encrypt,
		process4: kb.Encrypt4,
		process8: kb.Encrypt8,
	}
}
func NewECBDecryptor(b cipher.Block) cipher.BlockMode {
	kb, ok := b.(kryptoBlock)
	if !ok {
		kb = &blockWrap{b}
	}

	return &ecbEnc{
		b:        kb,
		process1: kb.Decrypt,
		process4: kb.Decrypt4,
		process8: kb.Decrypt8,
	}
}

type ecbEnc struct {
	b        kryptoBlock
	process1 func(dst, src []byte)
	process4 func(dst, src []byte)
	process8 func(dst, src []byte)
}

func (ecb *ecbEnc) BlockSize() int {
	return ecb.b.BlockSize()
}
func (ecb *ecbEnc) CryptBlocks(dst, src []byte) {
	BlockSize := ecb.BlockSize()

	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/kipher: invalid block size %d (src)", len(src)))
	}
	if len(dst) < len(src) {
		panic(fmt.Sprintf("krypto/kipher: invalid block size %d (dst)", len(dst)))
	}

	if len(src)%BlockSize != 0 {
		panic("krypto/kipher: input not full blocks")
	}

	remainBlock := len(src) / BlockSize

	for remainBlock >= 8 {
		remainBlock -= 8
		ecb.process8(dst, src)

		dst, src = dst[0x80:], src[0x80:]
	}

	for remainBlock >= 4 {
		remainBlock -= 4
		ecb.process4(dst, src)

		dst, src = dst[0x40:], src[0x40:]
	}

	for remainBlock > 0 {
		remainBlock -= 1
		ecb.process1(dst, src)

		dst, src = dst[0x10:], src[0x10:]
	}
}
