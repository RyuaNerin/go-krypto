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

type blockWrap struct {
	cipher.Block
}

var _ kryptoBlock = (*blockWrap)(nil)

func (b blockWrap) Encrypt4(dst, src []byte) {
	bs := b.BlockSize()
	b.Encrypt(dst[0*bs:1*bs], src[0*bs:1*bs])
	b.Encrypt(dst[1*bs:2*bs], src[1*bs:2*bs])
	b.Encrypt(dst[2*bs:3*bs], src[2*bs:3*bs])
	b.Encrypt(dst[3*bs:4*bs], src[3*bs:4*bs])
}

func (b blockWrap) Decrypt4(dst, src []byte) {
	bs := b.BlockSize()
	b.Decrypt(dst[0*bs:1*bs], src[0*bs:1*bs])
	b.Decrypt(dst[1*bs:2*bs], src[1*bs:2*bs])
	b.Decrypt(dst[2*bs:3*bs], src[2*bs:3*bs])
	b.Decrypt(dst[3*bs:4*bs], src[3*bs:4*bs])
}

func (b blockWrap) Encrypt8(dst, src []byte) {
	bs := b.BlockSize()
	b.Encrypt(dst[0*bs:1*bs], src[0*bs:1*bs])
	b.Encrypt(dst[1*bs:2*bs], src[1*bs:2*bs])
	b.Encrypt(dst[2*bs:3*bs], src[2*bs:3*bs])
	b.Encrypt(dst[3*bs:4*bs], src[3*bs:4*bs])
	b.Encrypt(dst[4*bs:5*bs], src[4*bs:5*bs])
	b.Encrypt(dst[5*bs:6*bs], src[5*bs:6*bs])
	b.Encrypt(dst[6*bs:7*bs], src[6*bs:7*bs])
	b.Encrypt(dst[7*bs:8*bs], src[7*bs:8*bs])
}

func (b blockWrap) Decrypt8(dst, src []byte) {
	bs := b.BlockSize()
	b.Decrypt(dst[0*bs:1*bs], src[0*bs:1*bs])
	b.Decrypt(dst[1*bs:2*bs], src[1*bs:2*bs])
	b.Decrypt(dst[2*bs:3*bs], src[2*bs:3*bs])
	b.Decrypt(dst[3*bs:4*bs], src[3*bs:4*bs])
	b.Decrypt(dst[4*bs:5*bs], src[4*bs:5*bs])
	b.Decrypt(dst[5*bs:6*bs], src[5*bs:6*bs])
	b.Decrypt(dst[6*bs:7*bs], src[6*bs:7*bs])
	b.Decrypt(dst[7*bs:8*bs], src[7*bs:8*bs])
}
