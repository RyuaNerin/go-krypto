package gmac

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/gcm.go

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"hash"

	"github.com/RyuaNerin/go-krypto/internal"
	igcm "github.com/RyuaNerin/go-krypto/internal/gcm"
	"github.com/RyuaNerin/go-krypto/internal/kryptoutil"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

var defaultIV [igcm.GCMStandardNonceSize]byte

// new MAC using GMAC by only passing additional data(aad data).
func NewGMAC(b cipher.Block, iv []byte) (hash.Hash, error) {
	kb := internal.WrapBlock(b)

	if kb.BlockSize() != igcm.GCMBlockSize {
		return nil, errors.New(msgRequired128BitBlockCipher)
	}

	if len(iv) == 0 {
		iv = defaultIV[:]
	}

	g := &ghash{}
	igcm.Init(&g.gcm, kb)

	var counter [igcm.GCMBlockSize]byte
	g.gcm.DeriveCounter(&counter, iv)
	g.gcm.Cipher.Encrypt(g.tagMask[:], counter[:])

	return g, nil
}

type ghash struct {
	gcm igcm.GCM

	tagMask [igcm.GCMBlockSize]byte

	y         igcm.GCMFieldElement
	remains   [igcm.GCMBlockSize]byte
	remainIdx int
	written   int
}

func (g ghash) Size() int {
	return igcm.GCMBlockSize
}

func (g ghash) BlockSize() int {
	return igcm.GCMBlockSize
}

func (g *ghash) Reset() {
	g.y = igcm.GCMFieldElement{}
	g.remainIdx = 0
	g.written = 0
}

func (g *ghash) Write(b []byte) (n int, err error) {
	if g.remainIdx > 0 {
		n = copy(g.remains[g.remainIdx:], b)
		g.written += n
		g.remainIdx += n

		if g.remainIdx < igcm.GCMBlockSize {
			return n, nil
		}
		b = b[n:]

		g.gcm.UpdateBlocks(&g.y, g.remains[:])
		g.remainIdx = 0
	}

	fullBlocks := (len(b) / igcm.GCMBlockSize) * igcm.GCMBlockSize
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

	var block [igcm.GCMBlockSize]byte

	if g.remainIdx > 0 {
		n := copy(block[:], g.remains[:g.remainIdx])
		g.gcm.UpdateBlocks(&yy, block[:])

		kryptoutil.MemsetByte(block[:n], 0)
	}

	yy.Low ^= uint64(written) * 8
	g.gcm.Mul(&yy)

	ret, out := internal.SliceForAppend(b, len(b)+igcm.GCMBlockSize)
	binary.BigEndian.PutUint64(out, yy.Low)
	binary.BigEndian.PutUint64(out[8:], yy.High)

	subtle.XORBytes(out, out, g.tagMask[:])

	return ret
}
