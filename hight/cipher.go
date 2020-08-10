package hight

import (
	"crypto/cipher"
	"fmt"
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("kipher/hight: invalid key size %d", int(k))
}

func NewCipher(key []byte) (cipher.Block, error) {
	l := len(key)
	if l != 16 {
		return nil, KeySizeError(l)
	}

	s := hight{}

	hight_KeySched(key, 16, s.pdwRoundKey[:])

	return &s, nil
}

type hight struct {
	pdwRoundKey [pdwRoundKeySize]byte
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

	copy(dst, src[:BlockSize])
	hight_Encrypt(s.pdwRoundKey[:], dst)
}

func (s *hight) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/hight: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/hight: invalid block size %d (dst)", len(dst)))
	}

	copy(dst, src[:BlockSize])
	hight_Decrypt(s.pdwRoundKey[:], dst)
}
