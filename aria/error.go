package aria

import "fmt"

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/aria: invalid key size %d", int(k))
}

const (
	msgFormatInvalidBlockSizeSrc = "krypto/aria: invalid block size %d (src)"
	msgFormatInvalidBlockSizeDst = "krypto/aria: invalid block size %d (dst)"
)
