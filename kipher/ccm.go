package kipher

import (
	"crypto/cipher"
	"errors"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/memory"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

const ccmBlockSize = 16

type ccm struct {
	cipher    cipher.Block
	nonceSize int
	tagSize   int
}

func NewCCM(b cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	if nonceSize <= 0 {
		return nil, errors.New(msgInvalidNonceZero)
	}
	if nonceSize < 7 || nonceSize > 13 {
		return nil, errors.New(msgInvalidNonceSize)
	}
	if tagSize < 4 || 16 < tagSize || tagSize%2 != 0 {
		return nil, errors.New(msgInvalidTagSizeCCM)
	}

	if b.BlockSize() != ccmBlockSize {
		return nil, errors.New(msgRequire128Bits)
	}

	return &ccm{
		cipher:    b,
		nonceSize: nonceSize,
		tagSize:   tagSize,
	}, nil
}

func (g *ccm) NonceSize() int {
	return g.nonceSize
}

func (g *ccm) Overhead() int {
	return g.tagSize
}

func (g *ccm) tag(ptLen int, N, A []byte) (tag, ctr, S0 [ccmBlockSize]byte) {
	Nlen := g.nonceSize
	Alen := len(A)

	/* Formatting of B0 */
	if len(A) > 0 {
		tag[0] = 0x40
	}
	tag[0] |= byte(g.tagSize-2) << 2
	tag[0] |= byte(14 - Nlen)
	copy(tag[1:14], N)
	tag[14] = byte(ptLen >> 8)
	tag[15] = byte(ptLen)
	if Nlen < 13 {
		tag[13] = byte(ptLen >> 16)
		if Nlen < 12 {
			tag[12] = byte(ptLen >> 24)
		}
	}
	g.cipher.Encrypt(tag[:], tag[:])

	/* Formatting of the Associated Data */
	var i int
	var x [6]byte
	if Alen < 0xff00 {
		x[0] = byte(Alen >> 8)
		x[1] = byte(Alen)

		subtle.XORBytes(tag[:], tag[:], x[:2])

		i = 2
	} else {
		x[0] = 0xff
		x[1] = 0xfe
		x[2] = byte(Alen >> 24)
		x[3] = byte(Alen >> 16)
		x[4] = byte(Alen >> 8)
		x[5] = byte(Alen)

		subtle.XORBytes(tag[:], tag[:], x[:])
		i = 6
	}

	for idx := 0; idx < Alen; i = 0 {
		idx += subtle.XORBytes(tag[i:], tag[i:], A[idx:])
		g.cipher.Encrypt(tag[:], tag[:])
	}

	/* Formatting of the counter blocks */
	ctr[0] = byte(14 - Nlen)
	copy(ctr[1:], N)

	/* Calculate S0 */
	g.cipher.Encrypt(S0[:], ctr[:])
	internal.IncCtr(ctr[:])

	return
}

func (g *ccm) Seal(dst, N, pt, A []byte) []byte {
	if len(N) != g.nonceSize {
		panic(msgInvalidNonce)
	}
	if uint64(len(pt)) > ((1<<32)-2)*uint64(g.cipher.BlockSize()) {
		panic(msgDataTooLarge)
	}

	ret, out := internal.SliceForAppend(dst, len(pt)+g.tagSize)
	if alias.InexactOverlap(out, pt) {
		panic(msgBufferOverlap)
	}
	ct, T := out[:len(pt)], out[len(pt):]

	tag, ctr, S0 := g.tag(len(pt), N, A)

	/* Calculate Yr */
	for idx := 0; idx < len(pt); {
		idx += subtle.XORBytes(tag[:], tag[:], pt[idx:])
		g.cipher.Encrypt(tag[:], tag[:])
	}
	subtle.XORBytes(T, tag[:g.tagSize], S0[:])

	ctrMode := NewCTR(g.cipher, ctr[:])
	ctrMode.XORKeyStream(ct, pt)

	return ret
}

func (g *ccm) Open(dst, N, ctT, A []byte) ([]byte, error) {
	if len(N) != g.nonceSize {
		panic(msgInvalidNonce)
	}
	if uint64(len(ctT)) > ((1<<32)-2)*uint64(g.cipher.BlockSize()) {
		panic(msgDataTooLarge)
	}

	ct, T := ctT[:len(ctT)-g.tagSize], ctT[len(ctT)-g.tagSize:]

	ret, pt := internal.SliceForAppend(dst, len(ct))
	if alias.InexactOverlap(pt, ct) {
		panic(msgBufferOverlap)
	}

	tag, ctr, S0 := g.tag(len(ct), N, A)

	ctrMode := NewCTR(g.cipher, ctr[:])
	ctrMode.XORKeyStream(pt, ct)

	/* Calculate Yr */
	for idx := 0; idx < len(ct); {
		idx += subtle.XORBytes(tag[:], tag[:], pt[idx:])
		g.cipher.Encrypt(tag[:], tag[:])
	}

	subtle.XORBytes(tag[:], tag[:g.tagSize], S0[:])

	if subtle.ConstantTimeCompare(T[:g.tagSize], tag[:g.tagSize]) != 1 {
		memory.Memclr(pt)
		return nil, errors.New(msgOpenFailed)
	}

	return ret, nil
}
