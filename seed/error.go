package seed

import "fmt"

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/seed: invalid key size %d", int(k))
}

const (
	msgInvalidBlockSizeSrcFormat = "krypto/seed: invalid block size %d (src)"
	msgInvalidBlockSizeDstFormat = "krypto/seed: invalid block size %d (dst)"
)
