// Pakcage gmac implements a Galois/Counter Mode (GCM) based MAC, as defined in KS X ISO/IEC 9797-3, NIST SP 800-38D.
package gmac

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/gcm.go

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"hash"

	"github.com/RyuaNerin/go-krypto/internal"
	ikipher "github.com/RyuaNerin/go-krypto/internal/kipher"
	"github.com/RyuaNerin/go-krypto/internal/memory"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

var defaultIV [ikipher.GCMStandardNonceSize]byte

// new MAC using GMAC by only passing additional data(aad data).
func NewGMAC(b cipher.Block, iv []byte) (hash.Hash, error) {
	kb := ikipher.WrapKipher(b)

	if kb.BlockSize() != ikipher.GCMBlockSize {
		return nil, errors.New(msgRequired128BitBlockCipher)
	}

	if len(iv) == 0 {
		iv = defaultIV[:]
	}

	g := &ghash{}
	ikipher.Init(&g.gcm, kb)

	var counter [ikipher.GCMBlockSize]byte
	g.gcm.DeriveCounter(&counter, iv)
	g.gcm.Cipher.Encrypt(g.tagMask[:], counter[:])

	return g, nil
}

type ghash struct {
	gcm ikipher.GCM

	tagMask [ikipher.GCMBlockSize]byte

	y         ikipher.GCMFieldElement
	remains   [ikipher.GCMBlockSize]byte
	remainIdx int
	written   int
}

func (g ghash) Size() int {
	return ikipher.GCMBlockSize
}

func (g ghash) BlockSize() int {
	return ikipher.GCMBlockSize
}

func (g *ghash) Reset() {
	g.y = ikipher.GCMFieldElement{}
	g.remainIdx = 0
	g.written = 0
}

func (g *ghash) Write(b []byte) (n int, err error) {
	if g.remainIdx > 0 {
		n = copy(g.remains[g.remainIdx:], b)
		g.written += n
		g.remainIdx += n

		if g.remainIdx < ikipher.GCMBlockSize {
			return n, nil
		}
		b = b[n:]

		g.gcm.UpdateBlocks(&g.y, g.remains[:])
		g.remainIdx = 0
	}

	fullBlocks := (len(b) / ikipher.GCMBlockSize) * ikipher.GCMBlockSize
	g.gcm.UpdateBlocks(&g.y, b[:fullBlocks])
	n += fullBlocks
	g.written += fullBlocks
	b = b[fullBlocks:]

	if len(b) > 0 {
		g.remainIdx = copy(g.remains[:], b)
		n += g.remainIdx
	}

	return
}

func (g *ghash) Sum(b []byte) []byte {
	yy := g.y

	written := g.written + g.remainIdx

	var block [ikipher.GCMBlockSize]byte

	if g.remainIdx > 0 {
		n := copy(block[:], g.remains[:g.remainIdx])
		g.gcm.UpdateBlocks(&yy, block[:])

		memory.Memclr(block[:n])
	}

	yy.Low ^= uint64(written) * 8
	g.gcm.Mul(&yy)

	ret, out := internal.SliceForAppend(b, len(b)+ikipher.GCMBlockSize)
	binary.BigEndian.PutUint64(out, yy.Low)
	binary.BigEndian.PutUint64(out[8:], yy.High)

	subtle.XORBytes(out, out, g.tagMask[:])

	return ret
}
