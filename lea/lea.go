// Package lea implements LEA encryption, as defined in TTAK.KO-12.0223
package lea

import (
	"crypto/cipher"
	"fmt"
)

type funcNew func(key []byte) (cipher.Block, error)
type funcBlock func(ctx *leaContext, dst, src []byte)

type leaContext struct {
	round uint8
	rk    [192]uint32
	ecb   bool
}

var (
	leaEnc1 funcBlock = leaEnc1Go
	leaEnc4 funcBlock = leaEnc4Go
	leaEnc8 funcBlock = leaEnc8Go

	leaDec1 funcBlock = leaDec1Go
	leaDec4 funcBlock = leaDec4Go
	leaDec8 funcBlock = leaDec8Go

	leaNew    funcNew = newCipherGo
	leaNewECB funcNew = newCipherECBGo
)

const (
	// The LEA block size in bytes.
	BlockSize = 16
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/lea: invalid key size %d", int(k))
}

// NewCipher creates and returns a new cipher.Block.
// The key argument should be the LEA key, either 16, 24, or 32 bytes to select LEA-128, LEA-192, or LEA-256.
func NewCipher(key []byte) (cipher.Block, error) {
	return leaNew(key)
}

// NewCipherECB creates and returns a new cipher.Block by ECB mode.
// This function can be useful in amd64.
// The key argument should be the LEA key, either 16, 24, or 32 bytes to select LEA-128, LEA-192, or LEA-256.
func NewCipherECB(key []byte) (cipher.Block, error) {
	return leaNewECB(key)
}
