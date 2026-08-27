package kipher

import (
	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/alias"
)

func NewECBEncryptor(b krypto.Block8) *ecbEncrypt {
	x := &ecbEncrypt{
		ecb: ecb{
			cipher:    b,
			blockSize: b.BlockSize(),
		},
	}
	return x
}

func NewECBDecryptor(b krypto.Block8) *ecbDecrypt {
	x := &ecbDecrypt{
		ecb: ecb{
			cipher:    b,
			blockSize: b.BlockSize(),
		},
	}
	return x
}

type (
	ecb struct {
		cipher    krypto.Block8
		blockSize int
	}
	ecbEncrypt struct{ ecb }
	ecbDecrypt struct{ ecb }
)

var (
	_ krypto.BlockMode = (*ecbEncrypt)(nil)
	_ krypto.BlockMode = (*ecbDecrypt)(nil)
)

func (ecb *ecb) Zeroize() {
	//zeroize.Zeroize(ecb.b)
	ecb.cipher = nil
	ecb.blockSize = 0
}

func (ecb *ecb) BlockSize() int {
	return ecb.cipher.BlockSize()
}

func (ecb *ecbEncrypt) CryptBlocks(dst, src []byte) { //nolint:dupl
	bs1 := ecb.blockSize
	if bs1 <= 0 || bs1 > int(^uint(0)>>1)/8 {
		panic(msgInvalidBlockSize)
	}
	bs4, bs8 := 4*bs1, 8*bs1
	if len(src)%bs1 != 0 {
		panic(msgNotFullBlocks)
	}
	if len(dst) < len(src) {
		panic(msgSmallDst)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgBufferOverlap)
	}
	dst = dst[:len(src)]

	for len(src) >= bs8 {
		ecb.cipher.Encrypt8(dst, src)

		dst, src = dst[bs8:], src[bs8:]
	}

	for len(src) >= bs4 {
		ecb.cipher.Encrypt4(dst, src)

		dst, src = dst[bs4:], src[bs4:]
	}

	for len(src) >= bs1 {
		ecb.cipher.Encrypt(dst, src)

		dst, src = dst[bs1:], src[bs1:]
	}
}

func (ecb *ecbDecrypt) CryptBlocks(dst, src []byte) { //nolint:dupl
	bs1 := ecb.blockSize
	if bs1 <= 0 || bs1 > int(^uint(0)>>1)/8 {
		panic(msgInvalidBlockSize)
	}
	bs4, bs8 := 4*bs1, 8*bs1
	if len(src)%bs1 != 0 {
		panic(msgNotFullBlocks)
	}
	if len(dst) < len(src) {
		panic(msgSmallDst)
	}
	if alias.InexactOverlap(dst[:len(src)], src) {
		panic(msgBufferOverlap)
	}
	dst = dst[:len(src)]

	for len(src) >= bs8 {
		ecb.cipher.Decrypt8(dst, src)

		dst, src = dst[bs8:], src[bs8:]
	}

	for len(src) >= bs4 {
		ecb.cipher.Decrypt4(dst, src)

		dst, src = dst[bs4:], src[bs4:]
	}

	for len(src) >= bs1 {
		ecb.cipher.Decrypt(dst, src)

		dst, src = dst[bs1:], src[bs1:]
	}
}
