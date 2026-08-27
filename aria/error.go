package aria

import "fmt"

// KeySizeError is returned when an invalid ARIA key size is supplied.
type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/aria: invalid key size %d", int(k))
}

const (
	msgNotFullBlock  = "krypto/aria: not full block"
	msgBufferOverlap = "krypto/aria: invalid buffer overlap"
)
