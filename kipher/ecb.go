package kipher

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/alias"
)

func NewECBEncryptor(b cipher.Block) cipher.BlockMode {
	kb := internal.WrapBlock(b)

	return &ecb{
		b:        kb,
		process1: kb.Encrypt,
		process4: kb.Encrypt4,
		process8: kb.Encrypt8,
	}
}

func NewECBDecryptor(b cipher.Block) cipher.BlockMode {
	kb := internal.WrapBlock(b)

	return &ecb{
		b:        kb,
		process1: kb.Decrypt,
		process4: kb.Decrypt4,
		process8: kb.Decrypt8,
	}
}

type ecb struct {
	b        internal.Block
	process1 func(dst, src []byte)
	process4 func(dst, src []byte)
	process8 func(dst, src []byte)
}

func (ecb *ecb) BlockSize() int {
	return ecb.b.BlockSize()
}

func (ecb *ecb) CryptBlocks(dst, src []byte) {
	blockSize := ecb.BlockSize()

	if len(src)%blockSize != 0 {
		panic("krypto/kipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("krypto/kipher: output smaller than input")
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic("krypto/kipher: invalid buffer overlap")
	}

	var (
		blockSize4 = 4 * blockSize
		blockSize8 = 8 * blockSize
	)

	for len(src) >= blockSize8 {
		ecb.process8(dst, src)

		dst, src = dst[blockSize8:], src[blockSize8:]
	}

	for len(src) >= blockSize4 {
		ecb.process4(dst, src)

		dst, src = dst[blockSize4:], src[blockSize4:]
	}

	for len(src) > 0 {
		ecb.process1(dst, src)

		dst, src = dst[blockSize:], src[blockSize:]
	}
}
