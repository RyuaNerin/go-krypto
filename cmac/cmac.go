// Package cmac implements CMAC(Cipher Message Authentication Code) as defined in NIST SP 800-38B, KS X ISO/IEC 9797-1.
package cmac

import (
	"crypto/cipher"
	"crypto/subtle"
	"hash"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

const (
	// the standard tag size of CMAC, in bytes.
	StandardTagSize = 16
)

const maxBlockSize = 16

// Equal reports whether mac1 and mac2 are identical, comparing them in constant time.
func Equal(mac1, mac2 []byte) bool {
	return subtle.ConstantTimeCompare(mac1, mac2) == 1
}

// [New] returns a new hash.Hash computing the CMAC checksum.
//
// The returned value is a [krypto.Hash] and therefore implements Zeroize(): call it
// to erase the derived subkeys (k1, k2) and accumulator when the MAC is no longer
// needed. After Zeroize() further use returns krypto.ErrStateCleared. NewWithSize is
// the same.
func New(b cipher.Block) hash.Hash {
	return newCMAC(b, b.BlockSize())
}

func NewWithSize(b cipher.Block, tagSize int) krypto.Hash {
	return newCMAC(b, tagSize)
}

func newCMAC(b cipher.Block, tagSize int) *cmac {
	blockSize := b.BlockSize()
	switch blockSize {
	case 8:
	case 16:
	default:
		panic(msgUnsupportedBlockCipher)
	}
	if tagSize <= 0 || tagSize > blockSize {
		panic(msgInvalidTagSize)
	}

	c := &cmac{
		block:   b,
		tagSize: tagSize,
	}

	// k1
	b.Encrypt(c.k1[:blockSize], c.k1[:blockSize])
	makeCMACSubkey(c.k1[:blockSize])

	// k2
	c.k2 = c.k1
	makeCMACSubkey(c.k2[:blockSize])

	c.z.SetValid()
	return c
}

func makeCMACSubkey(k []byte) {
	var carry byte
	for i := len(k) - 1; i >= 1; i -= 2 {
		carry2 := k[i] >> 7
		k[i] += k[i] + carry

		carry = k[i-1] >> 7
		k[i-1] += k[i-1] + carry2
	}
	mask := -carry
	switch len(k) {
	case 8:
		k[7] ^= 0x1b & mask
	case 16:
		k[15] ^= 0x87 & mask
	}
}

type cmac struct {
	z zeroize.ZeroizedFlag

	block cipher.Block

	k1, k2  [maxBlockSize]byte
	ciph    [maxBlockSize]byte
	m       [maxBlockSize]byte
	mIdx    int
	tagSize int
}

var _ krypto.Hash = (*cmac)(nil)

func (ctx *cmac) Zeroize() {
	ctx.z.Zeroize()

	zeroize.Zeroize(ctx.block)
	ctx.block = nil

	clear(ctx.k1[:])
	clear(ctx.k2[:])
	clear(ctx.ciph[:])
	clear(ctx.m[:])
	ctx.mIdx = 0
	ctx.tagSize = 0
}

func (c *cmac) Size() int { return c.tagSize }

func (c *cmac) BlockSize() int { return c.block.BlockSize() }

func (c *cmac) Reset() {
	clear(c.ciph[:])
	c.mIdx = 0
}

func (c *cmac) Write(p []byte) (n int, err error) {
	if err = c.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	n = len(p)
	if n == 0 {
		return 0, nil
	}

	blockSize := c.block.BlockSize()
	switch blockSize {
	case 8:
	case 16:
	default:
		panic(msgInvalidInternalState)
	}

	ciph, m := c.ciph[:blockSize], c.m[:blockSize]

	mIdx := c.mIdx
	if mIdx > blockSize {
		panic(msgInvalidInternalState)
	}
	if mIdx > 0 {
		if mIdx+len(p) <= blockSize {
			c.mIdx += copy(m[mIdx:blockSize], p)
			return n, nil
		}
		p = p[copy(m[mIdx:blockSize], p):]
		c.mIdx = 0

		subtle.XORBytes(ciph, ciph, m[:blockSize])
		c.block.Encrypt(ciph, ciph)
	}

	for len(p) > blockSize {
		subtle.XORBytes(ciph, ciph, p[:blockSize])
		c.block.Encrypt(ciph, ciph)
		p = p[blockSize:]
	}

	c.mIdx = copy(m, p)
	return n, nil
}

func (c *cmac) Sum(b []byte) []byte {
	c.z.ThrowIfZeroized()

	blockSize, mIdx, tagSize := c.block.BlockSize(), c.mIdx, c.tagSize
	switch blockSize {
	case 8:
	case 16:
	default:
		panic(msgInvalidInternalState)
	}
	if mIdx > blockSize || tagSize < 1 || tagSize > blockSize {
		panic(msgInvalidInternalState)
	}

	k1, k2, ciph, m := c.k1[:blockSize], c.k2[:blockSize], c.ciph[:blockSize], c.m[:blockSize]

	var mac [maxBlockSize]byte
	defer clear(mac[:])

	if mIdx == blockSize {
		subtle.XORBytes(mac[:], ciph, k1)  // CIPH ^ K1
		subtle.XORBytes(mac[:], mac[:], m) // CIPH ^ K1 ^ M
	} else {
		subtle.XORBytes(mac[:], ciph[:], k2)      // CIPH ^ K2
		subtle.XORBytes(mac[:], mac[:], m[:mIdx]) // CIPH ^ K2 ^ M
		mac[mIdx] ^= 0x80
	}
	c.block.Encrypt(mac[:blockSize], mac[:blockSize])

	return append(b, mac[:tagSize]...)
}
