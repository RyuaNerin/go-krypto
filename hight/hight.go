package hight

import (
	"crypto/cipher"
	"fmt"
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/hight: invalid key size %d", int(k))
}

type hight struct {
	pdwRoundKey [136]byte
}

// NewCipher creates and returns a new cipher.Block.
func NewCipher(key []byte) (cipher.Block, error) {
	l := len(key)
	if l != 16 {
		return nil, KeySizeError(l)
	}

	block := new(hight)

	for i := 0; i < 4; i++ {
		block.pdwRoundKey[i+0] = key[i+12]
		block.pdwRoundKey[i+4] = key[i+0]
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			block.pdwRoundKey[8+16*i+j] = key[(j-i)&7] + delta[16*i+j]
		}
		for j := 0; j < 8; j++ {
			block.pdwRoundKey[8+16*i+j+8] = key[((j-i)&7)+8] + delta[16*i+j+8]
		}
	}

	return block, nil
}

func (s *hight) BlockSize() int {
	return BlockSize
}

func (s *hight) encryptStep(XX []byte, k, i0, i1, i2, i3, i4, i5, i6, i7 int) {
	XX[i0] = (XX[i0] ^ (hight_F0[XX[i1]] + s.pdwRoundKey[4*k+3]))
	XX[i2] = (XX[i2] + (hight_F1[XX[i3]] ^ s.pdwRoundKey[4*k+2]))
	XX[i4] = (XX[i4] ^ (hight_F0[XX[i5]] + s.pdwRoundKey[4*k+1]))
	XX[i6] = (XX[i6] + (hight_F1[XX[i7]] ^ s.pdwRoundKey[4*k+0]))
}

func (s *hight) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/hight: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/hight: invalid block size %d (dst)", len(dst)))
	}

	XX := []byte{
		src[0] + s.pdwRoundKey[0],
		src[1],
		src[2] ^ s.pdwRoundKey[1],
		src[3],
		src[4] + s.pdwRoundKey[2],
		src[5],
		src[6] ^ s.pdwRoundKey[3],
		src[7],
	}

	s.encryptStep(XX, 2, 7, 6, 5, 4, 3, 2, 1, 0)
	s.encryptStep(XX, 3, 6, 5, 4, 3, 2, 1, 0, 7)
	s.encryptStep(XX, 4, 5, 4, 3, 2, 1, 0, 7, 6)
	s.encryptStep(XX, 5, 4, 3, 2, 1, 0, 7, 6, 5)
	s.encryptStep(XX, 6, 3, 2, 1, 0, 7, 6, 5, 4)
	s.encryptStep(XX, 7, 2, 1, 0, 7, 6, 5, 4, 3)
	s.encryptStep(XX, 8, 1, 0, 7, 6, 5, 4, 3, 2)
	s.encryptStep(XX, 9, 0, 7, 6, 5, 4, 3, 2, 1)
	s.encryptStep(XX, 10, 7, 6, 5, 4, 3, 2, 1, 0)
	s.encryptStep(XX, 11, 6, 5, 4, 3, 2, 1, 0, 7)
	s.encryptStep(XX, 12, 5, 4, 3, 2, 1, 0, 7, 6)
	s.encryptStep(XX, 13, 4, 3, 2, 1, 0, 7, 6, 5)
	s.encryptStep(XX, 14, 3, 2, 1, 0, 7, 6, 5, 4)
	s.encryptStep(XX, 15, 2, 1, 0, 7, 6, 5, 4, 3)
	s.encryptStep(XX, 16, 1, 0, 7, 6, 5, 4, 3, 2)
	s.encryptStep(XX, 17, 0, 7, 6, 5, 4, 3, 2, 1)
	s.encryptStep(XX, 18, 7, 6, 5, 4, 3, 2, 1, 0)
	s.encryptStep(XX, 19, 6, 5, 4, 3, 2, 1, 0, 7)
	s.encryptStep(XX, 20, 5, 4, 3, 2, 1, 0, 7, 6)
	s.encryptStep(XX, 21, 4, 3, 2, 1, 0, 7, 6, 5)
	s.encryptStep(XX, 22, 3, 2, 1, 0, 7, 6, 5, 4)
	s.encryptStep(XX, 23, 2, 1, 0, 7, 6, 5, 4, 3)
	s.encryptStep(XX, 24, 1, 0, 7, 6, 5, 4, 3, 2)
	s.encryptStep(XX, 25, 0, 7, 6, 5, 4, 3, 2, 1)
	s.encryptStep(XX, 26, 7, 6, 5, 4, 3, 2, 1, 0)
	s.encryptStep(XX, 27, 6, 5, 4, 3, 2, 1, 0, 7)
	s.encryptStep(XX, 28, 5, 4, 3, 2, 1, 0, 7, 6)
	s.encryptStep(XX, 29, 4, 3, 2, 1, 0, 7, 6, 5)
	s.encryptStep(XX, 30, 3, 2, 1, 0, 7, 6, 5, 4)
	s.encryptStep(XX, 31, 2, 1, 0, 7, 6, 5, 4, 3)
	s.encryptStep(XX, 32, 1, 0, 7, 6, 5, 4, 3, 2)
	s.encryptStep(XX, 33, 0, 7, 6, 5, 4, 3, 2, 1)

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

func (s *hight) decryptStep(XX [8]byte, k, i0, i1, i2, i3, i4, i5, i6, i7 int) {
	XX[i1] = (XX[i1] - (hight_F1[XX[i2]] ^ s.pdwRoundKey[4*k+2]))
	XX[i3] = (XX[i3] ^ (hight_F0[XX[i4]] + s.pdwRoundKey[4*k+1]))
	XX[i5] = (XX[i5] - (hight_F1[XX[i6]] ^ s.pdwRoundKey[4*k+0]))
	XX[i7] = (XX[i7] ^ (hight_F0[XX[i0]] + s.pdwRoundKey[4*k+3]))
}

func (s *hight) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/hight: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/hight: invalid block size %d (dst)", len(dst)))
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

	s.decryptStep(XX, 33, 7, 6, 5, 4, 3, 2, 1, 0)
	s.decryptStep(XX, 32, 0, 7, 6, 5, 4, 3, 2, 1)
	s.decryptStep(XX, 31, 1, 0, 7, 6, 5, 4, 3, 2)
	s.decryptStep(XX, 30, 2, 1, 0, 7, 6, 5, 4, 3)
	s.decryptStep(XX, 29, 3, 2, 1, 0, 7, 6, 5, 4)
	s.decryptStep(XX, 28, 4, 3, 2, 1, 0, 7, 6, 5)
	s.decryptStep(XX, 27, 5, 4, 3, 2, 1, 0, 7, 6)
	s.decryptStep(XX, 26, 6, 5, 4, 3, 2, 1, 0, 7)
	s.decryptStep(XX, 25, 7, 6, 5, 4, 3, 2, 1, 0)
	s.decryptStep(XX, 24, 0, 7, 6, 5, 4, 3, 2, 1)
	s.decryptStep(XX, 23, 1, 0, 7, 6, 5, 4, 3, 2)
	s.decryptStep(XX, 22, 2, 1, 0, 7, 6, 5, 4, 3)
	s.decryptStep(XX, 21, 3, 2, 1, 0, 7, 6, 5, 4)
	s.decryptStep(XX, 20, 4, 3, 2, 1, 0, 7, 6, 5)
	s.decryptStep(XX, 19, 5, 4, 3, 2, 1, 0, 7, 6)
	s.decryptStep(XX, 18, 6, 5, 4, 3, 2, 1, 0, 7)
	s.decryptStep(XX, 17, 7, 6, 5, 4, 3, 2, 1, 0)
	s.decryptStep(XX, 16, 0, 7, 6, 5, 4, 3, 2, 1)
	s.decryptStep(XX, 15, 1, 0, 7, 6, 5, 4, 3, 2)
	s.decryptStep(XX, 14, 2, 1, 0, 7, 6, 5, 4, 3)
	s.decryptStep(XX, 13, 3, 2, 1, 0, 7, 6, 5, 4)
	s.decryptStep(XX, 12, 4, 3, 2, 1, 0, 7, 6, 5)
	s.decryptStep(XX, 11, 5, 4, 3, 2, 1, 0, 7, 6)
	s.decryptStep(XX, 10, 6, 5, 4, 3, 2, 1, 0, 7)
	s.decryptStep(XX, 9, 7, 6, 5, 4, 3, 2, 1, 0)
	s.decryptStep(XX, 8, 0, 7, 6, 5, 4, 3, 2, 1)
	s.decryptStep(XX, 7, 1, 0, 7, 6, 5, 4, 3, 2)
	s.decryptStep(XX, 6, 2, 1, 0, 7, 6, 5, 4, 3)
	s.decryptStep(XX, 5, 3, 2, 1, 0, 7, 6, 5, 4)
	s.decryptStep(XX, 4, 4, 3, 2, 1, 0, 7, 6, 5)
	s.decryptStep(XX, 3, 5, 4, 3, 2, 1, 0, 7, 6)
	s.decryptStep(XX, 2, 6, 5, 4, 3, 2, 1, 0, 7)

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
