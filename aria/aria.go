// Package aria implements ARIA encryption, as defined in KS X 1213-1
package aria

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

const (
	// The ARIA block size in bytes.
	BlockSize = 16
)

// NewCipher creates and returns a new cipher.Block.
//
// The returned value is a [krypto.Block] and therefore implements Zeroize():
// call it to erase the round-key schedule when the cipher is no longer needed.
func NewCipher(key []byte) (cipher.Block, error) {
	l := len(key)
	switch l {
	case 16:
	case 24:
	case 32:
	default:
		return nil, KeySizeError(l)
	}

	return newCipher(key)
}

type (
	roundKey [16 * 17 / 4]uint32
	block1   [BlockSize]byte
	block2   [2 * BlockSize]byte
)

type ariaContext struct {
	z zeroize.ZeroizedFlag

	ek     roundKey
	dk     roundKey
	rounds int
}

var _ krypto.Block = (*ariaContext)(nil)

func (ctx *ariaContext) Zeroize() {
	clear(ctx.ek[:])
	clear(ctx.dk[:])
	ctx.rounds = 0

	ctx.z.Zeroize()
}

func (*ariaContext) BlockSize() int { return BlockSize }
