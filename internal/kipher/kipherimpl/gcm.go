package kipherimpl

import "github.com/RyuaNerin/go-krypto"

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/gcm.go

const (
	GCMBlockSize         = 16
	GCMTagSize           = 16
	GCMStandardNonceSize = 12
)

// counterCrypt crypts in to out using g.cipher in counter mode.
func GCMCounterCrypt(out, in []byte, c krypto.Block8, counter *[GCMBlockSize]byte) {
	var ctr CTR
	defer ctr.Zeroize()
	ctr.Init(c, counter[:], 4)
	ctr.Xor(out, in)
	ctr.CopyCTR(counter[:])
}

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/gcm.go

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

func (g *GCMFieldElement) Zeroize() {
	g.Low = 0
	g.High = 0
}

// GCM represents a Galois Counter Mode with a specific key. See
// https://csrc.nist.gov/groups/ST/toolkit/BCM/documents/proposedmodes/GCM/GCM-revised-spec.pdf
type GCM struct {
	// ProductTable contains the first sixteen powers of the key, H.
	// However, they are in bit reversed order. See NewGCMWithNonceSize.
	ProductTable [16]GCMFieldElement
}

// GCMStateID identifies the layout used by the active GHASH implementation.
// The generic representation is portable, while assembly tables are specific
// to an architecture. Persisted internal tables must carry this identifier.
func GCMStateID() byte { return gcmStateID }

var gcmStateID byte // 0: generic; 1: amd64 assembly; 2: arm64 assembly

var _ krypto.Zeroizable = (*GCM)(nil)

func (g *GCM) Zeroize() {
	clear(g.ProductTable[:])
}

var (
	gcmInit          func(g *GCM, cipher krypto.Block)
	gcmDeriveCounter func(g *GCM, counter *[GCMBlockSize]byte, nonce []byte)
	gcmUpdate        func(g *GCM, y *GCMFieldElement, blocks []byte)
	gcmAuth          func(g *GCM, out, ciphertext, additionalData []byte, tagMask *[GCMTagSize]byte)
	gcmFinish        func(g *GCM, out []byte, y *GCMFieldElement, ciphertextLen, additionalDataLen uint64, tagMask *[GCMTagSize]byte)
)

func (g *GCM) Init(cipher krypto.Block) {
	gcmInit(g, cipher)
}

func (g *GCM) DeriveCounter(counter *[GCMBlockSize]byte, nonce []byte) {
	gcmDeriveCounter(g, counter, nonce)
}

func (g *GCM) Update(y *GCMFieldElement, blocks []byte) {
	if len(blocks) == 0 {
		return
	}

	gcmUpdate(g, y, blocks)
}

func (g *GCM) Auth(out, ciphertext, additionalData []byte, tagMask *[GCMTagSize]byte) {
	gcmAuth(g, out, ciphertext, additionalData, tagMask)
}

func (g *GCM) Finish(out []byte, y *GCMFieldElement, ciphertextLen, additionalDataLen uint64, tagMask *[GCMTagSize]byte) {
	gcmFinish(g, out, y, ciphertextLen, additionalDataLen, tagMask)
}
