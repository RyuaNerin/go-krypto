package hight

import "fmt"

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/hight: invalid key size %d", int(k))
}

const (
	msgInvalidBlockSizeSrcFormat = "krypto/hight: invalid block size %d (src)"
	msgInvalidBlockSizeDstFormat = "krypto/hight: invalid block size %d (dst)"
)
