package aria

import (
	"crypto/cipher"
	"fmt"
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("kipher/aria: invalid key size %d", int(k))
}

type aria struct {
	rounds int
	ek     [rkSize]byte
	dk     [rkSize]byte
}

func NewCipher(key []byte) (cipher.Block, error) {
	l := len(key)
	switch l {
	case 16:
	case 24:
	case 32:
	default:
		return nil, KeySizeError(l)
	}

	block := new(aria)
	block.rounds = (l*8 + 256) / 32

	////////////////////////////////////////

	var t, w1, w2, w3 [16]byte

	q := (l*8 - 128) / 64
	for i := 0; i < 16; i++ {
		t[i] = s[i%4][krk[q][i]^key[i]]
	}

	dl(t[:], w1[:])
	if block.rounds == 14 {
		for i := 0; i < 8; i++ {
			w1[i] ^= key[16+i]
		}
	} else if block.rounds == 16 {
		for i := 0; i < 16; i++ {
			w1[i] ^= key[16+i]
		}
	}

	if q == 2 {
		q = 0
	} else {
		q = q + 1
	}

	for i := 0; i < 16; i++ {
		t[i] = s[(2+i)%4][krk[q][i]^w1[i]]
	}
	dl(t[:], w2[:])
	for i := 0; i < 16; i++ {
		w2[i] ^= key[i]
	}

	if q == 2 {
		q = 0
	} else {
		q = (q + 1)
	}
	for i := 0; i < 16; i++ {
		t[i] = s[i%4][krk[q][i]^w2[i]]
	}
	dl(t[:], w3[:])
	for i := 0; i < 16; i++ {
		w3[i] ^= w1[i]
	}

	for i := 0; i < 16*(block.rounds+1); i++ {
		block.ek[i] = 0
	}

	rotXOR(key, 0, block.ek[:], 0)
	rotXOR(w1[:], 19, block.ek[:], 0)
	rotXOR(w1[:], 0, block.ek[:], 16)
	rotXOR(w2[:], 19, block.ek[:], 16)
	rotXOR(w2[:], 0, block.ek[:], 32)
	rotXOR(w3[:], 19, block.ek[:], 32)
	rotXOR(w3[:], 0, block.ek[:], 48)
	rotXOR(key, 19, block.ek[:], 48)
	rotXOR(key, 0, block.ek[:], 64)
	rotXOR(w1[:], 31, block.ek[:], 64)
	rotXOR(w1[:], 0, block.ek[:], 80)
	rotXOR(w2[:], 31, block.ek[:], 80)
	rotXOR(w2[:], 0, block.ek[:], 96)
	rotXOR(w3[:], 31, block.ek[:], 96)
	rotXOR(w3[:], 0, block.ek[:], 112)
	rotXOR(key, 31, block.ek[:], 112)
	rotXOR(key, 0, block.ek[:], 128)
	rotXOR(w1[:], 67, block.ek[:], 128)
	rotXOR(w1[:], 0, block.ek[:], 144)
	rotXOR(w2[:], 67, block.ek[:], 144)
	rotXOR(w2[:], 0, block.ek[:], 160)
	rotXOR(w3[:], 67, block.ek[:], 160)
	rotXOR(w3[:], 0, block.ek[:], 176)
	rotXOR(key, 67, block.ek[:], 176)
	rotXOR(key, 0, block.ek[:], 192)
	rotXOR(w1[:], 97, block.ek[:], 192)
	if block.rounds > 12 {
		rotXOR(w1[:], 0, block.ek[:], 208)
		rotXOR(w2[:], 97, block.ek[:], 208)
		rotXOR(w2[:], 0, block.ek[:], 224)
		rotXOR(w3[:], 97, block.ek[:], 224)
	}
	if block.rounds > 14 {
		rotXOR(w3[:], 0, block.ek[:], 240)
		rotXOR(key, 97, block.ek[:], 240)
		rotXOR(key, 0, block.ek[:], 256)
		rotXOR(w1[:], 109, block.ek[:], 256)
	}

	////////////////////////////////////////

	copy(block.dk[:], block.ek[:])

	for j := 0; j < 16; j++ {
		t[j] = block.dk[j]
		block.dk[j] = block.dk[16*block.rounds+j]
		block.dk[16*block.rounds+j] = t[j]
	}
	for i := 1; i <= block.rounds/2; i++ {
		dl(block.dk[i*16:], t[:])
		dl(block.dk[(block.rounds-i)*16:], block.dk[i*16:])
		for j := 0; j < 16; j++ {
			block.dk[(block.rounds-i)*16+j] = t[j]
		}
	}

	return block, nil
}

func (s *aria) BlockSize() int {
	return BlockSize
}

func (s *aria) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/aria: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/aria: invalid block size %d (dst)", len(dst)))
	}

	s.crypt(dst, src, true)
}

func (s *aria) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/aria: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/aria: invalid block size %d (dst)", len(dst)))
	}

	s.crypt(dst, src, false)
}

func (cb *aria) crypt(dst, src []byte, encryption bool) {
	var e []byte
	if encryption {
		e = cb.ek[:]
	} else {
		e = cb.dk[:]
	}

	var i, j int
	var t [16]byte

	copy(dst, src[:BlockSize])

	ei := 0
	for i = 0; i < cb.rounds/2; i++ {
		for j = 0; j < 16; j++ {
			t[j] = s[j%4][e[ei+j]^dst[j]]
		}
		dl(t[:], dst)
		ei += 16
		for j = 0; j < 16; j++ {
			t[j] = s[(2+j)%4][e[ei+j]^dst[j]]
		}
		dl(t[:], dst)
		ei += 16
	}
	dl(dst, t[:])
	for j = 0; j < 16; j++ {
		dst[j] = e[ei+j] ^ t[j]
	}
}

func dl(i, o []byte) {
	var T byte

	T = i[3] ^ i[4] ^ i[9] ^ i[14]
	o[0] = i[6] ^ i[8] ^ i[13] ^ T
	o[5] = i[1] ^ i[10] ^ i[15] ^ T
	o[11] = i[2] ^ i[7] ^ i[12] ^ T
	o[14] = i[0] ^ i[5] ^ i[11] ^ T
	T = i[2] ^ i[5] ^ i[8] ^ i[15]
	o[1] = i[7] ^ i[9] ^ i[12] ^ T
	o[4] = i[0] ^ i[11] ^ i[14] ^ T
	o[10] = i[3] ^ i[6] ^ i[13] ^ T
	o[15] = i[1] ^ i[4] ^ i[10] ^ T
	T = i[1] ^ i[6] ^ i[11] ^ i[12]
	o[2] = i[4] ^ i[10] ^ i[15] ^ T
	o[7] = i[3] ^ i[8] ^ i[13] ^ T
	o[9] = i[0] ^ i[5] ^ i[14] ^ T
	o[12] = i[2] ^ i[7] ^ i[9] ^ T
	T = i[0] ^ i[7] ^ i[10] ^ i[13]
	o[3] = i[5] ^ i[11] ^ i[14] ^ T
	o[6] = i[2] ^ i[9] ^ i[12] ^ T
	o[8] = i[1] ^ i[4] ^ i[15] ^ T
	o[13] = i[3] ^ i[6] ^ i[8] ^ T
}

// Right-rotate 128 bit source string s by n bits and XOR it to target string t
func rotXOR(s []byte, n int, t []byte, ti int) {
	q := n / 8
	n %= 8
	for i := 0; i < 16; i++ {
		t[ti+(q+i)%16] ^= (s[i] >> n)
		if n != 0 {
			t[ti+(q+i+1)%16] ^= (s[i] << (8 - n))
		}
	}
}
