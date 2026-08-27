package gmac

import (
	"encoding"
	"encoding/binary"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/kipher/kipherimpl"
)

const (
	// Version 1 rejects legacy states whose GHASH representation was ambiguous.
	magic         = "gmac\x01"
	marshaledSize = len(magic) +
		/* implementation */ 1 +
		/* gcm */ 16*(8+8) +
		/* tagMask */ kipherimpl.GCMBlockSize +
		/* tagSize */ 1 +
		/* remainIdx */ 1 +
		/* remains */ kipherimpl.GCMBlockSize +
		/* y */ (8 + 8) +
		/* written */ 8

	implementationOffset = len(magic)
	productTableOffset   = implementationOffset + 1
	tagMaskOffset        = productTableOffset + 16*(8+8)
	tagSizeOffset        = tagMaskOffset + kipherimpl.GCMBlockSize
	remainIdxOffset      = tagSizeOffset + 1
	remainsOffset        = remainIdxOffset + 1
	yOffset              = remainsOffset + kipherimpl.GCMBlockSize
	writtenOffset        = yOffset + 8 + 8
)

func (ctx *ghash) MarshalBinary() ([]byte, error) {
	return ctx.AppendBinary(nil)
}

func (ctx *ghash) AppendBinary(b []byte) ([]byte, error) {
	if err := ctx.z.ReturnIfZeroized(); err != nil {
		return nil, err
	}
	b = internal.IncreaseCapacity(b, marshaledSize)
	b = append(b, magic...)
	b = append(b, kipherimpl.GCMStateID())
	for _, v := range ctx.gcm.ProductTable {
		b = internal.AppendBigUint64(b, v.High)
		b = internal.AppendBigUint64(b, v.Low)
	}
	b = append(b, ctx.tagMask[:]...)
	b = append(b, ctx.tagSize)
	b = append(b, ctx.remainIdx)
	remainsStart := len(b)
	b = append(b, ctx.remains[:]...)
	clear(b[remainsStart+int(ctx.remainIdx):])
	b = internal.AppendBigUint64(b, ctx.y.High)
	b = internal.AppendBigUint64(b, ctx.y.Low)
	b = internal.AppendBigUint64(b, ctx.written)
	return b, nil
}

func (ctx *ghash) UnmarshalBinary(b []byte) error {
	if err := ctx.z.ReturnIfZeroized(); err != nil {
		return err
	}
	if len(b) < len(magic) || string(b[:len(magic)]) != magic {
		return ErrInvalidHashStateIdentifier
	}
	if len(b) != marshaledSize {
		return ErrInvalidHashStateSize
	}
	if b[implementationOffset] != kipherimpl.GCMStateID() {
		return ErrIncompatibleHashState
	}

	tagSize := b[tagSizeOffset]
	remainIdx := b[remainIdxOffset]
	written := binary.BigEndian.Uint64(b[writtenOffset:])
	if tagSize != ctx.tagSize || remainIdx >= kipherimpl.GCMBlockSize || written%kipherimpl.GCMBlockSize != 0 {
		return ErrInvalidHashStateSize
	}
	for _, v := range b[remainsOffset+int(remainIdx) : yOffset] {
		if v != 0 {
			return ErrInvalidHashStateSize
		}
	}

	next := ctx.ghashData
	state := b[productTableOffset:]
	for i := range next.gcm.ProductTable {
		state, next.gcm.ProductTable[i].High = internal.ConsumeBigU64(state)
		state, next.gcm.ProductTable[i].Low = internal.ConsumeBigU64(state)
	}
	copy(next.tagMask[:], b[tagMaskOffset:tagSizeOffset])
	next.tagSize = tagSize
	next.remainIdx = remainIdx
	copy(next.remains[:], b[remainsOffset:yOffset])
	next.y.High = binary.BigEndian.Uint64(b[yOffset:])
	next.y.Low = binary.BigEndian.Uint64(b[yOffset+8:])
	next.written = written

	ctx.ghashData = next
	return nil
}

var (
	_ encoding.BinaryMarshaler   = (*ghash)(nil)
	_ encoding.BinaryUnmarshaler = (*ghash)(nil)
	_ encoding.BinaryAppender    = (*ghash)(nil)
)
