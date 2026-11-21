// Package seed implements SEED encryption, as defined in TTAS.KO-12.0004/R1
package seed

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

const (
	// The SEED block size in bytes.
	BlockSize = 16
	// The SEED key size in bytes.
	KeySize = 16
)

type (
	roundKey    [32]uint32
	seedContext struct {
		z           zeroize.ZeroizedFlag
		pdwRoundKey roundKey
	}

	block1 [1 * BlockSize]byte
	block2 [2 * BlockSize]byte
)

var _ krypto.Block = (*seedContext)(nil)

// NewCipher creates and returns a new SEED cipher.Block. The key argument must be the
// 16-byte SEED-128 key.
//
// The returned value is a [krypto.Block] and therefore implements Zeroize():
// call it to erase the round-key schedule when the cipher is no longer needed.
func NewCipher(key []byte) (cipher.Block, error) {
	if l := len(key); l != 16 {
		return nil, KeySizeError(l)
	}

	c := new(seedContext)
	c.z.SetValid()
	keySchedKeyGo(&c.pdwRoundKey, (*[16]byte)(key))
	return c, nil
}

func (s *seedContext) Zeroize() {
	clear(s.pdwRoundKey[:])
	s.z.Zeroize()
}

func (s *seedContext) BlockSize() int {
	return BlockSize
}

func (s *seedContext) Encrypt(dst, src []byte) {
	s.z.ThrowIfZeroized()
	if len(src) < BlockSize {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:BlockSize], src[:BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[BlockSize-1]
	_ = src[BlockSize-1]

	encrypt1Go(&s.pdwRoundKey, (*block1)(dst), (*block1)(src))
}

func (s *seedContext) Decrypt(dst, src []byte) {
	s.z.ThrowIfZeroized()
	if len(src) < BlockSize {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:BlockSize], src[:BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[BlockSize-1]
	_ = src[BlockSize-1]

	decrypt1Go(&s.pdwRoundKey, (*block1)(dst), (*block1)(src))
}
