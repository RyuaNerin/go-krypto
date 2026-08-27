package lea

import "fmt"

// KeySizeError is returned when an invalid LEA key size is supplied.
type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/lea: invalid key size %d", int(k))
}

const (
	msgNotFullBlock  = "krypto/lea: not full block"
	msgBufferOverlap = "krypto/lea: invalid buffer overlap"
)
