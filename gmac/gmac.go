// Package gmac implements a Galois/Counter Mode (GCM) based MAC, as defined in KS X ISO/IEC 9797-3, NIST SP 800-38D.
package gmac

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/gcm.go

import (
	"crypto/cipher"
	"crypto/subtle"
	"hash"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/kipher/kipherimpl"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

const (
	// the block size of GCM, in bytes.
	BlockSize = 16
	// the minimum tag size of GCM, in bytes.
	MinimumTagSize = 1
	// the standard tag size of GCM, in bytes.
	StandardTagSize = 16
	// the standard nonce size of GCM, in bytes.
	StandardNonceSize = 12
)

// Equal reports whether mac1 and mac2 are identical, comparing them in constant time.
func Equal(mac1, mac2 []byte) bool {
	return subtle.ConstantTimeCompare(mac1, mac2) == 1
}

// NewGMAC returns the given 128-bit block cipher in Galois/Counter Mode (GCM)
// as a GMAC, using the given IV.
//
// The returned value is a [krypto.Hash] and therefore implements Zeroize(): call it
// to erase the key-derived GCM state and accumulator when the GMAC is no longer
// needed. NewGMACWithSize is the same.
func NewGMAC(b cipher.Block, iv []byte) (hash.Hash, error) {
	return newGMAC(b, iv, StandardTagSize)
}

// NewGMACWithSize returns a new hash.Hash computing the GMAC checksum, with the
// given tag size.
func NewGMACWithSize(b cipher.Block, iv []byte, tagSize int) (hash.Hash, error) {
	return newGMAC(b, iv, tagSize)
}

func newGMAC(b cipher.Block, iv []byte, tagSize int) (*ghash, error) {
	if b.BlockSize() != BlockSize {
		return nil, ErrRequired128BitBlockCipher
	}

	if len(iv) == 0 {
		return nil, ErrInvalidIV
	}
	if tagSize < MinimumTagSize || tagSize > BlockSize || tagSize%2 != 0 {
		return nil, ErrInvalidTagSizeGMAC
	}

	kb := internal.ToBlock8(b)

	g := &ghash{
		ghashData: ghashData{
			tagSize: byte(tagSize),
		},
	}
	g.gcm.Init(kb)

	var counter [BlockSize]byte
	g.gcm.DeriveCounter(&counter, iv)
	kb.Encrypt(g.tagMask[:], counter[:])

	return g, nil
}

type ghash struct {
	z zeroize.ZeroizedFlag
	ghashData
}
type ghashData struct {
	gcm kipherimpl.GCM

	tagMask [BlockSize]byte
	tagSize byte

	remainIdx byte
	remains   [BlockSize]byte

	y       kipherimpl.GCMFieldElement
	written uint64
}

var _ krypto.Hash = (*ghash)(nil)

func (g *ghashData) Zeroize() {
	g.gcm.Zeroize()

	clear(g.tagMask[:])
	g.tagSize = 0

	g.remainIdx = 0
	clear(g.remains[:])

	zeroize.Zeroize(&g.y)
	g.written = 0
}

func (g *ghash) Zeroize() {
	g.z.SetInvalid()
	g.ghashData.Zeroize()
}

func (g *ghash) Size() int {
	return int(g.tagSize)
}

func (g *ghash) BlockSize() int {
	return BlockSize
}

func (g *ghash) Reset() {
	g.z.ThrowIfZeroized()

	g.y = kipherimpl.GCMFieldElement{}
	g.remainIdx = 0
	g.written = 0
}

func (g *ghash) Write(b []byte) (n int, err error) {
	if err = g.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	remainIdx := int(g.remainIdx)
	if remainIdx > BlockSize {
		panic(msgInvalidInternalState)
	}
	if remainIdx > 0 {
		n = copy(g.remains[remainIdx:], b)
		g.remainIdx += byte(n)

		if g.remainIdx < BlockSize {
			return n, nil
		}
		b = b[n:]

		g.gcm.Update(&g.y, g.remains[:])
		g.remainIdx = 0
		g.written += BlockSize
	}

	fullBlocks := len(b) &^ (BlockSize - 1)
	if fullBlocks > 0 {
		g.gcm.Update(&g.y, b[:fullBlocks])
		n += fullBlocks
		g.written += uint64(fullBlocks)
		b = b[fullBlocks:]
	}

	if len(b) > 0 {
		g.remainIdx = byte(copy(g.remains[:], b))
		n += int(g.remainIdx)
	}

	return n, nil
}

func (g *ghash) Sum(b []byte) []byte {
	g.z.ThrowIfZeroized()

	gg := g.ghashData
	defer gg.Zeroize()
	remainIdx, tagSize := int(gg.remainIdx), int(gg.tagSize)
	if remainIdx > BlockSize || tagSize < MinimumTagSize || tagSize > BlockSize {
		panic(msgInvalidInternalState)
	}

	written := gg.written + uint64(remainIdx)

	if remainIdx > 0 {
		gg.gcm.Update(&gg.y, gg.remains[:remainIdx])
	}

	var full [BlockSize]byte
	gg.gcm.Finish(full[:], &gg.y, 0, written, &gg.tagMask)

	return append(b, full[:tagSize]...)
}
