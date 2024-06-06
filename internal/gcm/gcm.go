package gcm

import (
	"encoding/binary"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/gcm.go

const (
	GCMBlockSize         = 16
	GCMTagSize           = 16
	GCMMinimumTagSize    = 12 // NIST SP 800-38D recommends tags with 12 or more bytes.
	GCMStandardNonceSize = 12
)

// GCMFieldElement represents a value in GF(2¹²⁸). In order to reflect the GCM
// standard and make binary.BigEndian suitable for marshaling these values, the
// bits are stored in big endian order. For example:
//
//	the coefficient of x⁰ can be obtained by v.low >> 63.
//	the coefficient of x⁶³ can be obtained by v.low & 1.
//	the coefficient of x⁶⁴ can be obtained by v.high >> 63.
//	the coefficient of x¹²⁷ can be obtained by v.high & 1.
type GCMFieldElement struct {
	Low, High uint64
}

// GCM represents a Galois Counter Mode with a specific key. See
// https://csrc.nist.gov/groups/ST/toolkit/BCM/documents/proposedmodes/GCM/GCM-revised-spec.pdf
type GCM struct {
	Cipher internal.Block
	// ProductTable contains the first sixteen powers of the key, H.
	// However, they are in bit reversed order. See NewGCMWithNonceSize.
	ProductTable [16]GCMFieldElement
}

func Init(g *GCM, cipher internal.Block) {
	var key [GCMBlockSize]byte
	cipher.Encrypt(key[:], key[:])

	g.Cipher = cipher

	// We precompute 16 multiples of |key|. However, when we do lookups
	// into this table we'll be using bits from a field element and
	// therefore the bits will be in the reverse order. So normally one
	// would expect, say, 4*key to be in index 4 of the table but due to
	// this bit ordering it will actually be in index 0010 (base 2) = 2.
	x := GCMFieldElement{
		binary.BigEndian.Uint64(key[:8]),
		binary.BigEndian.Uint64(key[8:]),
	}
	g.ProductTable[reverseBits(1)] = x

	for i := 2; i < 16; i += 2 {
		g.ProductTable[reverseBits(i)] = gcmDouble(&g.ProductTable[reverseBits(i/2)])
		g.ProductTable[reverseBits(i+1)] = gcmAdd(&g.ProductTable[reverseBits(i)], &x)
	}
}

// reverseBits reverses the order of the bits of 4-bit number in i.
func reverseBits(i int) int {
	i = ((i << 2) & 0xc) | ((i >> 2) & 0x3)
	i = ((i << 1) & 0xa) | ((i >> 1) & 0x5)
	return i
}

// gcmAdd adds two elements of GF(2¹²⁸) and returns the sum.
func gcmAdd(x, y *GCMFieldElement) GCMFieldElement {
	// Addition in a characteristic 2 field is just XOR.
	return GCMFieldElement{x.Low ^ y.Low, x.High ^ y.High}
}

// gcmDouble returns the result of doubling an element of GF(2¹²⁸).
func gcmDouble(x *GCMFieldElement) (double GCMFieldElement) {
	msbSet := x.High&1 == 1

	// Because of the bit-ordering, doubling is actually a right shift.
	double.High = x.High >> 1
	double.High |= x.Low << 63
	double.Low = x.Low >> 1

	// If the most-significant bit was set before shifting then it,
	// conceptually, becomes a term of x^128. This is greater than the
	// irreducible polynomial so the result has to be reduced. The
	// irreducible polynomial is 1+x+x^2+x^7+x^128. We can subtract that to
	// eliminate the term at x^128 which also means subtracting the other
	// four terms. In characteristic 2 fields, subtraction == addition ==
	// XOR.
	if msbSet {
		double.Low ^= 0xe100000000000000
	}

	return
}

var gcmReductionTable = []uint16{
	0x0000, 0x1c20, 0x3840, 0x2460, 0x7080, 0x6ca0, 0x48c0, 0x54e0,
	0xe100, 0xfd20, 0xd940, 0xc560, 0x9180, 0x8da0, 0xa9c0, 0xb5e0,
}

// Mul sets y to y*H, where H is the GCM key, fixed during NewGCMWithNonceSize.
func (g *GCM) Mul(y *GCMFieldElement) {
	var z GCMFieldElement

	for i := 0; i < 2; i++ {
		word := y.High
		if i == 1 {
			word = y.Low
		}

		// Multiplication works by multiplying z by 16 and adding in
		// one of the precomputed multiples of H.
		for j := 0; j < 64; j += 4 {
			msw := z.High & 0xf
			z.High >>= 4
			z.High |= z.Low << 60
			z.Low >>= 4
			z.Low ^= uint64(gcmReductionTable[msw]) << 48

			// the values in |table| are ordered for
			// little-endian bit positions. See the comment
			// in NewGCMWithNonceSize.
			t := &g.ProductTable[word&0xf]

			z.Low ^= t.Low
			z.High ^= t.High
			word >>= 4
		}
	}

	*y = z
}

// UpdateBlocks extends y with more polynomial terms from blocks, based on
// Horner's rule. There must be a multiple of gcmBlockSize bytes in blocks.
func (g *GCM) UpdateBlocks(y *GCMFieldElement, blocks []byte) {
	for len(blocks) > 0 {
		y.Low ^= binary.BigEndian.Uint64(blocks)
		y.High ^= binary.BigEndian.Uint64(blocks[8:])
		g.Mul(y)
		blocks = blocks[GCMBlockSize:]
	}
}

// update extends y with more polynomial terms from data. If data is not a
// multiple of gcmBlockSize bytes long then the remainder is zero padded.
func (g *GCM) Update(y *GCMFieldElement, data []byte) {
	fullBlocks := (len(data) >> 4) << 4
	g.UpdateBlocks(y, data[:fullBlocks])

	if len(data) != fullBlocks {
		var partialBlock [GCMBlockSize]byte
		copy(partialBlock[:], data[fullBlocks:])
		g.UpdateBlocks(y, partialBlock[:])
	}
}

// GCMInc32 treats the final four bytes of counterBlock as a big-endian value
// and increments it.
func GCMInc32(counterBlock *[GCMBlockSize]byte) {
	ctr := counterBlock[len(counterBlock)-4:]
	binary.BigEndian.PutUint32(ctr, binary.BigEndian.Uint32(ctr)+1)
}

func gcmSub32(counterBlock *[GCMBlockSize]byte, value uint32) {
	ctr := counterBlock[len(counterBlock)-4:]
	binary.BigEndian.PutUint32(ctr, binary.BigEndian.Uint32(ctr)+value)
}

// counterCrypt crypts in to out using g.cipher in counter mode.
func (g *GCM) CounterCrypt(out, in []byte, counter *[GCMBlockSize]byte) {
	const (
		bs8 = 8 * GCMBlockSize
		bs4 = 4 * GCMBlockSize
		bs3 = 3 * GCMBlockSize
		bs1 = 1 * GCMBlockSize
	)
	var maskBuf [bs8]byte

	for len(in) >= bs8 {
		for i := 0; i < 8; i++ {
			copy(maskBuf[i*bs1:], counter[:])
			GCMInc32(counter)
		}
		g.Cipher.Encrypt8(maskBuf[:], maskBuf[:])

		subtle.XORBytes(out, in, maskBuf[:])
		out = out[bs8:]
		in = in[bs8:]
	}

	mask := maskBuf[:0]
	for len(in) >= bs4 {
		if len(mask) == 0 {
			for i := 0; i < 8; i++ {
				copy(maskBuf[i*bs1:], counter[:])
				GCMInc32(counter)
			}
			g.Cipher.Encrypt8(maskBuf[:], maskBuf[:])
			mask = maskBuf[:]
		}

		subtle.XORBytes(out, in, mask[:bs4])
		out = out[bs4:]
		in = in[bs4:]
		mask = mask[bs4:]
	}

	for len(in) >= bs1 {
		if len(mask) == 0 {
			for i := 0; i < 3; i++ {
				copy(maskBuf[i*bs1:], counter[:])
				GCMInc32(counter)
			}
			g.Cipher.Encrypt4(maskBuf[:bs4], maskBuf[:bs4])
			mask = maskBuf[:bs3]
		}

		subtle.XORBytes(out, in, mask[:bs1])
		out = out[bs1:]
		in = in[bs1:]
		mask = mask[bs1:]
	}

	if len(in) > 0 {
		if len(mask) == 0 {
			mask = maskBuf[:bs1]
			g.Cipher.Encrypt(mask, counter[:])
		}
		subtle.XORBytes(out, in, mask)
		mask = nil
	}

	if len(mask) > 0 {
		gcmSub32(counter, uint32(len(mask)/bs1))
	}
}

// deriveCounter computes the initial GCM counter state from the given nonce.
// See NIST SP 800-38D, section 7.1. This assumes that counter is filled with
// zeros on entry.
func (g *GCM) DeriveCounter(counter *[GCMBlockSize]byte, nonce []byte) {
	// GCM has two modes of operation with respect to the initial counter
	// state: a "fast path" for 96-bit (12-byte) nonces, and a "slow path"
	// for nonces of other lengths. For a 96-bit nonce, the nonce, along
	// with a four-byte big-endian counter starting at one, is used
	// directly as the starting counter. For other nonce sizes, the counter
	// is computed by passing it through the GHASH function.
	if len(nonce) == GCMStandardNonceSize {
		copy(counter[:], nonce)
		counter[GCMBlockSize-1] = 1
	} else {
		var y GCMFieldElement
		g.Update(&y, nonce)
		y.High ^= uint64(len(nonce)) * 8
		g.Mul(&y)
		binary.BigEndian.PutUint64(counter[:8], y.Low)
		binary.BigEndian.PutUint64(counter[8:], y.High)
	}
}
