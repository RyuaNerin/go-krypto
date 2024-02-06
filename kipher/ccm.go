package kipher

import (
	"crypto/cipher"
	"errors"

	"github.com/RyuaNerin/go-krypto/internal/alias"
)

const ccmBlockSize = 16

type ccm struct {
	cipher    cipher.Block
	nonceSize int
	tagSize   int
}

func NewCCM(b cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	if nonceSize <= 0 {
		return nil, errors.New("krypto/kiphe: the nonce can't have zero length, or the security of the key will be immediately compromised")
	}
	if nonceSize < 7 || nonceSize > 13 {
		return nil, errors.New("krypto/kiphe: invalid nonce size")
	}
	if tagSize < 4 || 16 < tagSize || tagSize%2 != 0 {
		return nil, errors.New("krypto/kiphe: tagSize must be 4, 6, 8, 10, 12, 14 or 16")
	}

	if b.BlockSize() != ccmBlockSize {
		return nil, errors.New("krypto/kipher: NewCCM requires 128-bit block cipher")
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
	for i := 0; i < Nlen; i++ {
		tag[i+1] = N[i] /* Set N */
	}
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
	if Alen < 0xff00 {
		tag[0] ^= byte(Alen >> 8)
		tag[1] ^= byte(Alen)

		i = 2
	} else {
		tag[0] ^= 0xff
		tag[1] ^= 0xfe
		tag[2] ^= byte(Alen >> 24)
		tag[3] ^= byte(Alen >> 16)
		tag[4] ^= byte(Alen >> 8)
		tag[5] ^= byte(Alen)
		i = 6
	}

	for j := 0; j < Alen; i = 0 {
		for i < 16 && j < Alen {
			tag[i] ^= A[j]

			i++
			j++
		}

		g.cipher.Encrypt(tag[:], tag[:])
	}

	/* Formatting of the counter blocks */
	ctr[0] = byte(14 - Nlen)

	for i = 0; i < Nlen; i++ {
		ctr[i+1] = N[i]
	}

	/* Calculate S0 */
	g.cipher.Encrypt(S0[:], ctr[:])
	{
		//ctr64_inc
		for i := 7; i >= 0; i-- {
			c := ctr[8+i]
			c++
			ctr[8+i] = c
			if c > 0 {
				break
			}
		}
	}

	return
}

func (g *ccm) Seal(dst, N, pt, A []byte) []byte {
	if len(N) != g.nonceSize {
		panic("krypto/kipher: incorrect nonce length given to CCM")
	}
	if uint64(len(pt)) > ((1<<32)-2)*uint64(g.cipher.BlockSize()) {
		panic("krypto/kipher: message too large for CCM")
	}

	ret, out := sliceForAppend(dst, len(pt)+g.tagSize)
	if alias.InexactOverlap(out, pt) {
		panic("krypto/kipher: invalid buffer overlap")
	}
	ct := out[:len(pt)]
	T := out[len(out)-g.tagSize:]

	pt_len := len(pt)
	Tlen := g.tagSize
	//Nlen := g.nonceSize
	//Alen := len(A)

	tag, ctr, S0 := g.tag(pt_len, N, A)

	/* Calculate Yr */
	for i := 0; i < pt_len; {
		j := pt_len - i
		if j > 0x10 {
			j = 0x10
		}

		for h := 0; h < j; {
			tag[h] ^= pt[i]
			h++
			i++
		}
		g.cipher.Encrypt(tag[:], tag[:])
	}
	for i := 0; i < Tlen; i++ {
		T[i] = tag[i] ^ S0[i]
	}

	ctrMode := NewCTR(g.cipher, ctr[:])
	ctrMode.XORKeyStream(ct, pt)

	return ret
}

func (g *ccm) Open(dst, N, ctT, A []byte) ([]byte, error) {
	if len(N) != g.nonceSize {
		panic("krypto/kipher: incorrect nonce length given to CCM")
	}
	if uint64(len(ctT)) > ((1<<32)-2)*uint64(g.cipher.BlockSize()) {
		panic("krypto/kipher: message too large for CCM")
	}

	ret, pt := sliceForAppend(dst, len(ctT)-g.tagSize)
	if alias.InexactOverlap(pt, ctT) {
		panic("krypto/kipher: invalid buffer overlap")
	}

	ct := ctT[:len(ctT)-g.tagSize]
	ct_len := len(ct)
	T := ctT[len(ctT)-g.tagSize:]
	Tlen := g.tagSize

	tag, ctr, S0 := g.tag(ct_len, N, A)

	ctrMode := NewCTR(g.cipher, ctr[:])
	ctrMode.XORKeyStream(pt, ct)

	/* Calculate Yr */
	i := 0
	for i < ct_len {
		j := ct_len - i
		if j > 0x10 {
			j = 0x10
		}

		for h := 0; h < j; {
			tag[h] ^= pt[i]
			h++
			i++
		}
		g.cipher.Encrypt(tag[:], tag[:])
	}

	for i := 0; i < Tlen; i++ {
		tag[i] = tag[i] ^ S0[i]
	}

	for i := 0; i < Tlen; i++ {
		if T[i] != tag[i] {
			return nil, errOpen
		}
	}

	return ret, nil
}
