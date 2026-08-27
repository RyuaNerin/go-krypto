package lsh512

import (
	"encoding"

	"github.com/RyuaNerin/go-krypto/internal"
)

const (
	magic         = "lsh\x1f"
	marshaledSize = len(magic) +
		1 + /* algType */
		1 + /* remainDataByteLen */
		8*8 + /* cvL */
		8*8 + /* cvR */
		BlockSize /* lastBlock */
)

func (ctx *lsh512Context) MarshalBinary() ([]byte, error) {
	return ctx.AppendBinary(nil)
}

func (ctx *lsh512Context) AppendBinary(b []byte) ([]byte, error) {
	b = internal.IncreaseCapacity(b, marshaledSize)
	b = append(b, magic...)

	b = internal.AppendBigUint8(b, uint8(ctx.algType))

	b = internal.AppendBigUint8(b, uint8(ctx.remainDataByteLen))

	for _, v := range ctx.cvL {
		b = internal.AppendBigUint64(b, v)
	}

	for _, v := range ctx.cvR {
		b = internal.AppendBigUint64(b, v)
	}

	blockStart := len(b)
	b = append(b, ctx.lastBlock[:]...)
	clear(b[blockStart+int(ctx.remainDataByteLen):])

	return b, nil
}

func (ctx *lsh512Context) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic) || string(b[:len(magic)]) != magic {
		return ErrInvalidHashStateIdentifier
	}
	if len(b) != marshaledSize {
		return ErrInvalidHashStateSize
	}

	algType := int(b[len(magic)])
	remainDataByteLen := int(b[len(magic)+1])
	if remainDataByteLen >= BlockSize || !isValidHashSize(algType) || algType != int(ctx.algType) {
		return ErrInvalidHashStateSize
	}

	blockStart := len(b) - BlockSize
	for _, v := range b[blockStart+remainDataByteLen:] {
		if v != 0 {
			return ErrInvalidHashStateSize
		}
	}

	b = b[len(magic):]
	ctx.algType = uint64(algType)
	ctx.remainDataByteLen = uint64(remainDataByteLen)
	b = b[2:]

	for i := range ctx.cvL {
		b, ctx.cvL[i] = internal.ConsumeBigU64(b)
	}

	for i := range ctx.cvR {
		b, ctx.cvR[i] = internal.ConsumeBigU64(b)
	}

	copy(ctx.lastBlock[:], b)

	return nil
}

func isValidHashSize(size int) bool {
	return size == Size224 || size == Size256 || size == Size384 || size == Size
}

var (
	_ encoding.BinaryMarshaler   = (*lsh512Context)(nil)
	_ encoding.BinaryUnmarshaler = (*lsh512Context)(nil)
	_ encoding.BinaryAppender    = (*lsh512Context)(nil)
)
