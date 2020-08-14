package seed

import (
	"crypto/cipher"
	"fmt"
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/seed128: invalid key size %d", int(k))
}

func NewCipher(key []byte) (cipher.Block, error) {
	l := len(key)
	switch l {
	case 16:
		return new128(key), nil
	case 32:
		return new256(key), nil
	}

	return nil, KeySizeError(l)
}
