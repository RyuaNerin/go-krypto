package kipher

import (
	"crypto/subtle"
	"encoding/binary"
	"math"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/kipher/kipherimpl"
)

const CCMBlockSize = 16

func NewCCM(b krypto.Block8, nonceSize, tagSize int) *CCM {
	g := &CCM{
		cipher:    b,
		nonceSize: nonceSize,
		tagSize:   tagSize,
	}
	return g
}

type CCM struct {
	cipher    krypto.Block8
	nonceSize int
	tagSize   int
}

var _ krypto.AEAD = (*CCM)(nil)

func (g *CCM) Zeroize() {
	g.cipher = nil
	g.nonceSize = 0
	g.tagSize = 0
}

func (g *CCM) NonceSize() int {
	return g.nonceSize
}

func (g *CCM) Overhead() int {
	return g.tagSize
}

func (g *CCM) maxLength() uint64 {
	q := 15 - g.nonceSize
	if q >= 8 {
		return math.MaxUint64
	}
	return 1<<(8*q) - 1
}

func (g *CCM) tag(ptLen int, nonce, additional []byte) (tag, ctr, S0 [CCMBlockSize]byte) {
	nonceLen := g.nonceSize
	additionalLen := len(additional)

	/* Formatting of B0 */
	if len(additional) > 0 {
		tag[0] = 0x40
	}
	tag[0] |= byte(g.tagSize-2) << 2
	tag[0] |= byte(14 - nonceLen)
	copy(tag[1:1+nonceLen], nonce)

	for i, l := CCMBlockSize-1, uint64(ptLen); i > nonceLen; i, l = i-1, l>>8 {
		tag[i] = byte(l)
	}
	g.cipher.Encrypt(tag[:], tag[:])

	/* Formatting of the Associated Data */
	var i int
	var x [10]byte
	defer clear(x[:])

	switch {
	case additionalLen < 0xff00:
		binary.BigEndian.PutUint16(x[:], uint16(additionalLen))
		i = 2
	case uint64(additionalLen) <= math.MaxUint32:
		x[0], x[1] = 0xff, 0xfe
		binary.BigEndian.PutUint32(x[2:], uint32(additionalLen))
		i = 6
	default:
		x[0], x[1] = 0xff, 0xff
		binary.BigEndian.PutUint64(x[2:], uint64(additionalLen))
		i = 10
	}
	subtle.XORBytes(tag[:], tag[:], x[:i])

	for idx := 0; idx < additionalLen; i = 0 {
		idx += subtle.XORBytes(tag[i:], tag[i:], additional[idx:])
		g.cipher.Encrypt(tag[:], tag[:])
	}

	/* Formatting of the counter blocks */
	ctr[0] = byte(14 - nonceLen)
	copy(ctr[1:], nonce)

	/* Calculate S0 */
	g.cipher.Encrypt(S0[:], ctr[:])
	internal.IncCtr(ctr[:])

	return tag, ctr, S0
}

func (g *CCM) Seal(dst, nonce, plainText, additional []byte) []byte {
	if len(nonce) != g.nonceSize {
		panic(msgInvalidNonce)
	}
	if uint64(len(plainText)) > g.maxLength() {
		panic(msgDataTooLarge)
	}

	ret, out := internal.SliceForAppend(dst, len(plainText)+g.tagSize)
	if alias.InexactOverlap(out, plainText) {
		panic(msgBufferOverlap)
	}
	ct, T := out[:len(plainText)], out[len(plainText):]

	tag, ctr, S0 := g.tag(len(plainText), nonce, additional)
	defer clear(tag[:])
	defer clear(ctr[:])
	defer clear(S0[:])

	/* Calculate Yr */
	for idx := 0; idx < len(plainText); {
		idx += subtle.XORBytes(tag[:], tag[:], plainText[idx:])
		g.cipher.Encrypt(tag[:], tag[:])
	}
	subtle.XORBytes(T, tag[:g.tagSize], S0[:])

	var ctrMode kipherimpl.CTR
	defer ctrMode.Zeroize()
	ctrMode.Init(g.cipher, ctr[:], 0)
	ctrMode.Xor(ct, plainText)

	return ret
}

func (g *CCM) Open(dst, nonce, cipherText, additional []byte) ([]byte, error) {
	if len(nonce) != g.nonceSize {
		panic(msgInvalidNonce)
	}
	if len(cipherText) < g.tagSize {
		return nil, ErrOpenFailed
	}
	if uint64(len(cipherText)-g.tagSize) > g.maxLength() {
		return nil, ErrOpenFailed
	}

	ct, T := cipherText[:len(cipherText)-g.tagSize], cipherText[len(cipherText)-g.tagSize:]

	ret, pt := internal.SliceForAppend(dst, len(ct))
	if alias.InexactOverlap(pt, ct) {
		panic(msgBufferOverlap)
	}
	if alias.AnyOverlap(pt, additional) {
		panic(msgBufferOverlapAdditional)
	}

	tag, ctr, S0 := g.tag(len(ct), nonce, additional)
	defer clear(tag[:])
	defer clear(ctr[:])
	defer clear(S0[:])

	var ctrMode kipherimpl.CTR
	defer ctrMode.Zeroize()
	ctrMode.Init(g.cipher, ctr[:], 0)
	ctrMode.Xor(pt, ct)

	/* Calculate Yr */
	for idx := 0; idx < len(ct); {
		idx += subtle.XORBytes(tag[:], tag[:], pt[idx:])
		g.cipher.Encrypt(tag[:], tag[:])
	}

	subtle.XORBytes(tag[:], tag[:g.tagSize], S0[:])

	if subtle.ConstantTimeCompare(T[:g.tagSize], tag[:g.tagSize]) != 1 {
		clear(pt)
		return nil, ErrOpenFailed
	}

	return ret, nil
}
