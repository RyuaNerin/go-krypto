package lea

import (
	"crypto/cipher"
	"fmt"
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("kipher/lea: invalid key size %d", int(k))
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

	s := lea_key{}
	s.lea_set_key_generic(key, l)

	return &s, nil
}

type lea_key struct {
	rk    [192]uint32
	round int
}

func (s *lea_key) BlockSize() int {
	return BlockSize
}

func (s *lea_key) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/lea: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/lea: invalid block size %d (dst)", len(dst)))
	}

	s.lea_encrypt(dst, src)
}

func (s *lea_key) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("kipher/lea: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("kipher/lea: invalid block size %d (dst)", len(dst)))
	}

	s.lea_decrypt(dst, src)
}
