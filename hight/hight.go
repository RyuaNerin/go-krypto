// Package hight implements HIGHT encryption, as defined in TTAS.KO-12.0040/R1
package hight

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

const (
	// The HIGHT block size in bytes.
	BlockSize = 8
	// The HIGHT key size in bytes.
	KeySize = 16
)

type (
	roundKeys    [136]byte
	hightContext struct {
		z           zeroize.ZeroizedFlag
		pdwRoundKey roundKeys
	}

	block1 [1 * BlockSize]byte
	block2 [2 * BlockSize]byte
)

var _ krypto.Block = (*hightContext)(nil)

// NewCipher creates and returns a new cipher.Block.
//
// The returned value is a [krypto.Block] and therefore implements Zeroize():
// call it to erase the round-key schedule when the cipher is no longer needed.
func NewCipher(key []byte) (cipher.Block, error) {
	l := len(key)
	if l != KeySize {
		return nil, KeySizeError(l)
	}

	block := new(hightContext)
	block.z.SetValid()

	for i := range 4 {
		block.pdwRoundKey[i+0] = key[i+12]
		block.pdwRoundKey[i+4] = key[i+0]
	}

	for i := range 8 {
		for k := range 8 {
			block.pdwRoundKey[8+16*i+k+0] = key[((k-i)&7)+0] + delta[16*i+k+0]
		}
		for k := range 8 {
			block.pdwRoundKey[8+16*i+k+8] = key[((k-i)&7)+8] + delta[16*i+k+8]
		}
	}

	return block, nil
}

func (ctx *hightContext) Zeroize() {
	clear(ctx.pdwRoundKey[:])
	ctx.z.Zeroize()
}

func (ctx *hightContext) BlockSize() int {
	return BlockSize
}

func (ctx *hightContext) Encrypt(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

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

	encryptGo(&ctx.pdwRoundKey, (*block1)(dst), (*block1)(src))
}

func (ctx *hightContext) Decrypt(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

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

	decryptGo(&ctx.pdwRoundKey, (*block1)(dst), (*block1)(src))
}
