package has160

import (
	"encoding"
	"encoding/binary"
	"errors"

	"github.com/RyuaNerin/go-krypto/internal"
)

const (
	magic         = "has\x01"
	marshaledSize = len(magic) + /*H*/ 5*4 + /*block*/ BlockSize + /*boff*/ 1 + /*length*/ 8
)

func (ctx *has160Context) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, marshaledSize)
	b = append(b, magic...)
	b = binary.BigEndian.AppendUint32(b, ctx.H[0])
	b = binary.BigEndian.AppendUint32(b, ctx.H[1])
	b = binary.BigEndian.AppendUint32(b, ctx.H[2])
	b = binary.BigEndian.AppendUint32(b, ctx.H[3])
	b = binary.BigEndian.AppendUint32(b, ctx.H[4])
	b = append(b, ctx.block[:ctx.boff]...)
	b = b[:len(b)+len(ctx.block)-ctx.boff] // already zero
	b = append(b, byte(ctx.boff))
	b = binary.BigEndian.AppendUint64(b, uint64(ctx.length))
	return b, nil
}

func (ctx *has160Context) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic) || string(b[:len(magic)]) != magic {
		return errors.New("krypto/hash160: invalid hash state identifier")
	}
	if len(b) != marshaledSize {
		return errors.New("krypto/hash160: invalid hash state size")
	}

	b = b[len(magic):]
	b, ctx.H[0] = internal.ConsumeUint32(b)
	b, ctx.H[1] = internal.ConsumeUint32(b)
	b, ctx.H[2] = internal.ConsumeUint32(b)
	b, ctx.H[3] = internal.ConsumeUint32(b)
	b, ctx.H[4] = internal.ConsumeUint32(b)
	b = b[copy(ctx.block[:], b[:BlockSize]):]
	ctx.boff = int(b[0])
	ctx.length = int(binary.BigEndian.Uint64(b[1:]))
	return nil
}

var (
	_ encoding.BinaryMarshaler   = (*has160Context)(nil)
	_ encoding.BinaryUnmarshaler = (*has160Context)(nil)
)
