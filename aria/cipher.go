package aria

import (
	"crypto/cipher"
	"fmt"
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("kipher/aria: invalid key size %d", int(k))
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

	s := aria{}

	s.rounds = encKeySetup(key, s.ek[:], l*8)
	decKeySetup(key, s.dk[:], l*8)

	return &s, nil
}

type aria struct {
	rounds int
	ek     [rkSize]byte
	dk     [rkSize]byte
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

	crypt(src, s.rounds, s.ek[:], dst)
}

func (s *aria) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/aria: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/aria: invalid block size %d (dst)", len(dst)))
	}

	crypt(src, s.rounds, s.dk[:], dst)
}
