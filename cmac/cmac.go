package cmac

import (
	"crypto/cipher"
	"crypto/subtle"
	"hash"

	"github.com/RyuaNerin/go-krypto/internal/kryptoutil"
)

func Equal(mac1, mac2 []byte) bool {
	return subtle.ConstantTimeCompare(mac1, mac2) == 1
}

func New(b cipher.Block) hash.Hash {
	blockSize := b.BlockSize()

	arr := make([]byte, 4*blockSize)
	var (
		k1   = arr[0*blockSize : 1*blockSize]
		k2   = arr[1*blockSize : 2*blockSize]
		ciph = arr[2*blockSize : 3*blockSize]
		m    = arr[3*blockSize : 4*blockSize]
	)

	// k1
	b.Encrypt(k1, k1)
	makeCMACSubkey(k1)

	// k2
	copy(k2, k1)
	makeCMACSubkey(k2)

	return &cmac{
		block: b,
		k1:    k1,
		k2:    k2,
		ciph:  ciph,
		m:     m,
	}
}

func makeCMACSubkey(k []byte) {
	var carry byte
	for i := len(k) - 1; i >= 1; i -= 2 {
		carry2 := k[i] >> 7
		k[i] += k[i] + carry

		carry = k[i-1] >> 7
		k[i-1] += k[i-1] + carry2
	}
	if carry > 0 {
		switch len(k) {
		case 8:
			k[7] ^= 0x1b
		case 16:
			k[15] ^= 0x87
		case 32:
			k[30] ^= 4
			k[31] ^= 0x25
		case 64:
			k[62] ^= 1
			k[63] ^= 0x25
		case 128:
			k[125] ^= 8
			k[126] ^= 0x00
			k[127] ^= 0x43
		default:
			panic(msgUnsupportedCipher)
		}
	}
}

type cmac struct {
	block cipher.Block

	k1, k2 []byte
	ciph   []byte
	m      []byte
	mIdx   int
}

func (c *cmac) Size() int { return c.block.BlockSize() }

func (c *cmac) BlockSize() int { return c.block.BlockSize() }

func (c *cmac) Reset() {
	kryptoutil.MemsetByte(c.ciph, 0)
	c.mIdx = 0
}

func (c *cmac) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}

	blockSize := c.block.BlockSize()

	if c.mIdx > 0 && c.mIdx+len(p) > blockSize {
		copied := copy(c.m[c.mIdx:], p)
		subtle.XORBytes(c.ciph, c.ciph, c.m)
		c.block.Encrypt(c.ciph, c.ciph)

		p = p[copied:]
		c.mIdx = 0
	}

	for len(p) > blockSize { // 마지막 블록 저장을 위해서
		subtle.XORBytes(c.ciph, c.ciph, p)
		c.block.Encrypt(c.ciph, c.ciph)

		p = p[blockSize:]
	}

	if len(p) > 0 {
		copy(c.m[c.mIdx:], p)
		c.mIdx += len(p)
	}

	return len(p), nil
}

func (c *cmac) Sum(b []byte) []byte {
	blockSize := c.block.BlockSize()

	mac := make([]byte, blockSize)

	if c.mIdx == blockSize {
		subtle.XORBytes(mac, c.ciph, c.k1) // CIPH ^ K1
		subtle.XORBytes(mac, mac, c.m)     // CIPH ^ K1 ^ M
	} else {
		subtle.XORBytes(mac, c.ciph, c.k2)      // CIPH ^ K2
		subtle.XORBytes(mac, mac, c.m[:c.mIdx]) // CIPH ^ K2 ^ M
		mac[c.mIdx] ^= 0x80
	}
	c.block.Encrypt(mac, mac)

	return append(b, mac...)
}
