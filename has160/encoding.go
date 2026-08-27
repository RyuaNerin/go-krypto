package has160

import (
	"encoding"
	"encoding/binary"

	"github.com/RyuaNerin/go-krypto/internal"
)

const (
	magic         = "has\x00"
	marshaledSize = len(magic) + /*H*/ 5*4 + /*block*/ BlockSize + /*boff*/ 1 + /*length*/ 8
)

var (
	_ encoding.BinaryMarshaler   = (*has160Context)(nil)
	_ encoding.BinaryUnmarshaler = (*has160Context)(nil)
)

func (ctx *has160Context) MarshalBinary() ([]byte, error) {
	return ctx.AppendBinary(nil)
}

func (ctx *has160Context) AppendBinary(b []byte) ([]byte, error) {
	b = internal.IncreaseCapacity(b, marshaledSize)
	b = append(b, magic...)
	for _, h := range ctx.h {
		b = internal.AppendBigUint32(b, h)
	}
	n := len(b)
	b = append(b, ctx.block[:]...)
	clear(b[n+ctx.boff:]) // boff 이후는 명시적으로 0
	b = append(b, byte(ctx.boff))
	b = internal.AppendBigUint64(b, uint64(ctx.length))
	return b, nil
}

func (ctx *has160Context) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic) || string(b[:len(magic)]) != magic {
		return ErrInvalidHashStateIdentifier
	}
	if len(b) != marshaledSize {
		return ErrInvalidHashStateSize
	}
	_ = b[:marshaledSize]

	const boffAt = len(magic) + 5*4 + BlockSize
	boff := int(b[boffAt])
	length := binary.BigEndian.Uint64(b[boffAt+1:])
	if boff >= BlockSize || length%BlockSize != uint64(boff) {
		return ErrInvalidHashStateSize
	}

	const blockStart = len(magic) + 5*4
	for _, v := range b[blockStart+boff : boffAt] {
		if v != 0 {
			return ErrInvalidHashStateSize
		}
	}

	b = b[len(magic):]
	for i := range ctx.h {
		b, ctx.h[i] = internal.ConsumeBigU32(b)
	}
	copy(ctx.block[:], b)
	ctx.boff = boff
	ctx.length = length
	return nil
}
