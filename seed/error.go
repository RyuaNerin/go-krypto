package seed

import "fmt"

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/seed: invalid key size %d", int(k))
}

const (
	msgNotFullBlock  = "krypto/seed: not full block"
	msgBufferOverlap = "krypto/seed: invalid buffer overlap"
)
