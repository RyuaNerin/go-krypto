package lea

import "fmt"

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/lea: invalid key size %d", int(k))
}

const (
	msgInvalidBlockSizeSrcFormat = "krypto/lea: invalid block size %d (src)"
	msgInvalidBlockSizeDstFormat = "krypto/lea: invalid block size %d (dst)"
	msgInputNotFullBlocks        = "krypto/lea: input not full blocks"
)
