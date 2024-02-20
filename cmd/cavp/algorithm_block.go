package main

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/kipher"
)

func blockECB(c cipher.Block, iv []byte) funcProcessBlock {
	return func(dst, src []byte, encrypt bool) {
		if encrypt {
			for i := 0; i < len(dst); i += c.BlockSize() {
				c.Encrypt(dst[i:i+c.BlockSize()], src[i:i+c.BlockSize()])
			}
		} else {
			for i := 0; i < len(dst); i += c.BlockSize() {
				c.Decrypt(dst[i:i+c.BlockSize()], src[i:i+c.BlockSize()])
			}
		}
	}
}

func blockCBC(c cipher.Block, iv []byte) funcProcessBlock {
	cbcEnc := cipher.NewCBCEncrypter(c, iv)
	cbcDec := cipher.NewCBCDecrypter(c, iv)

	return func(dst, src []byte, encrypt bool) {
		if encrypt {
			cbcEnc.CryptBlocks(dst, src)
		} else {
			cbcDec.CryptBlocks(dst, src)
		}
	}
}

func blockCFB(blockBits int) func(c cipher.Block, iv []byte) funcProcessBlock {
	blockBits /= 8

	return func(c cipher.Block, iv []byte) funcProcessBlock {
		cbcEnc := kipher.NewCFBEncrypter(c, iv, blockBits)
		cbcDec := kipher.NewCFBDecrypter(c, iv, blockBits)

		return func(dst, src []byte, encrypt bool) {
			if encrypt {
				cbcEnc.XORKeyStream(dst, src)
			} else {
				cbcDec.XORKeyStream(dst, src)
			}
		}
	}
}

func blockOFB(c cipher.Block, iv []byte) funcProcessBlock {
	ofb := cipher.NewOFB(c, iv)

	return func(dst, src []byte, encrypt bool) {
		ofb.XORKeyStream(dst, src)
	}
}

func blockCTR(c cipher.Block, iv []byte) funcProcessBlock {
	ctr := cipher.NewCTR(c, iv)

	return func(dst, src []byte, encrypt bool) {
		ctr.XORKeyStream(dst, src)
	}
}
