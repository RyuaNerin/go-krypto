package hight

import "fmt"

// KeySizeError is returned when an invalid HIGHT key size is supplied.
type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/hight: invalid key size %d", int(k))
}

const (
	msgNotFullBlock  = "krypto/hight: not full block"
	msgBufferOverlap = "krypto/hight: invalid buffer overlap"
)
