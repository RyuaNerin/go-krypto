package seed128

import (
	"crypto/cipher"
	"fmt"
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("kipher/seed128: invalid key size %d", int(k))
}

func NewCipher(key []byte) (cipher.Block, error) {
	l := len(key)
	if l != 16 {
		return nil, KeySizeError(l)
	}

	s := seed{}
	seed_KeySchedKey(s.pdwRoundKey[:], key)

	return &s, nil
}

type seed struct {
	pdwRoundKey [pdwRoundKeySize]uint32
}

func (s *seed) BlockSize() int {
	return BlockSize
}

func (s *seed) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/seed128: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/seed128: invalid block size %d (dst)", len(dst)))
	}

	copy(dst, src[:BlockSize])
	seed_Encrypt(dst, s.pdwRoundKey[:])
}

func (s *seed) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/seed128: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/seed128: invalid block size %d (dst)", len(dst)))
	}

	copy(dst, src[:BlockSize])
	seed_Decrypt(dst, s.pdwRoundKey[:])
}
