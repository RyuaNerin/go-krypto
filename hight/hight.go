package hight

import (
	"crypto/cipher"
	"fmt"
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("kipher/hight: invalid key size %d", int(k))
}

type hight struct {
	pdwRoundKey [pdwRoundKeySize]byte
}

func NewCipher(key []byte) (cipher.Block, error) {
	l := len(key)
	if l != 16 {
		return nil, KeySizeError(l)
	}

	block := new(hight)

	var i, j int32

	for i = 0; i < 4; i++ {
		block.pdwRoundKey[i] = key[i+12]
		block.pdwRoundKey[i+4] = key[i]
	}

	for i = 0; i < 8; i++ {
		for j = 0; j < 8; j++ {
			block.pdwRoundKey[8+16*i+j] = key[(j-i)&7] + delta[16*i+j]
		}
		for j = 0; j < 8; j++ {
			block.pdwRoundKey[8+16*i+j+8] = key[((j-i)&7)+8] + delta[16*i+j+8]
		}
	}

	return block, nil
}

func (s *hight) BlockSize() int {
	return BlockSize
}

func (s *hight) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/hight: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/hight: invalid block size %d (dst)", len(dst)))
	}

	XX := [...]byte{
		src[0] + s.pdwRoundKey[0],
		src[1],
		src[2] ^ s.pdwRoundKey[1],
		src[3],
		src[4] + s.pdwRoundKey[2],
		src[5],
		src[6] ^ s.pdwRoundKey[3],
		src[7],
	}

	var XX000 byte
	for i := 2; i <= 33; i++ {
		XX000 = XX[0]
		XX[0] = XX[7] ^ (hight_F0[XX[6]] + s.pdwRoundKey[4*i+3])
		XX[7] = XX[6]
		XX[6] = XX[5] + (hight_F1[XX[4]] ^ s.pdwRoundKey[4*i+2])
		XX[5] = XX[4]
		XX[4] = XX[3] ^ (hight_F0[XX[2]] + s.pdwRoundKey[4*i+1])
		XX[3] = XX[2]
		XX[2] = XX[1] + (hight_F1[XX000] ^ s.pdwRoundKey[4*i+0])
		XX[1] = XX000
	}

	// Final Round
	dst[0] = XX[1] + s.pdwRoundKey[4]
	dst[1] = XX[2]
	dst[2] = XX[3] ^ s.pdwRoundKey[5]
	dst[3] = XX[4]
	dst[4] = XX[5] + s.pdwRoundKey[6]
	dst[5] = XX[6]
	dst[6] = XX[7] ^ s.pdwRoundKey[7]
	dst[7] = XX[0]
}

func (s *hight) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/hight: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/hight: invalid block size %d (dst)", len(dst)))
	}

	XX := [...]byte{
		src[7],
		src[0] - s.pdwRoundKey[4],
		src[1],
		src[2] ^ s.pdwRoundKey[5],
		src[3],
		src[4] - s.pdwRoundKey[6],
		src[5],
		src[6] ^ s.pdwRoundKey[7],
	}

	var XX_1_ byte
	for i := 33; i >= 2; i-- {
		XX_1_ = XX[1]
		XX[1] = XX[2] - (hight_F1[XX[1]] ^ s.pdwRoundKey[4*i+0])
		XX[2] = XX[3]
		XX[3] = XX[4] ^ (hight_F0[XX[3]] + s.pdwRoundKey[4*i+1])
		XX[4] = XX[5]
		XX[5] = XX[6] - (hight_F1[XX[5]] ^ s.pdwRoundKey[4*i+2])
		XX[6] = XX[7]
		XX[7] = XX[0] ^ (hight_F0[XX[7]] + s.pdwRoundKey[4*i+3])
		XX[0] = XX_1_
	}

	// Final Round
	dst[0] = XX[0] - s.pdwRoundKey[0]
	dst[1] = XX[1]
	dst[2] = XX[2] ^ s.pdwRoundKey[1]
	dst[3] = XX[3]
	dst[4] = XX[4] - s.pdwRoundKey[2]
	dst[5] = XX[5]
	dst[6] = XX[6] ^ s.pdwRoundKey[3]
	dst[7] = XX[7]
}
