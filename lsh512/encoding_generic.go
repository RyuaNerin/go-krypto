//go:build !(arm64 || amd64 || amd64p32) || purego
// +build !arm64,!amd64,!amd64p32 purego

package lsh512

import (
	"encoding"
	"encoding/binary"
	"errors"

	"github.com/RyuaNerin/go-krypto/internal"
)

const (
	magic         = "lsh\x01go"
	marshaledSize = len(magic) +
		/*cv*/ 16*8 +
		/*tcv*/ 16*8 +
		/*msg*/ 16*(numStep+1)*8 +
		/*block*/ BlockSize +
		/*boff*/ 2 +
		/*outlenbytes*/ 2
)

func (ctx *lsh512ContextGo) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, marshaledSize)
	b = append(b, magic...)
	for _, v := range ctx.cv {
		b = binary.BigEndian.AppendUint64(b, v)
	}
	for _, v := range ctx.tcv {
		b = binary.BigEndian.AppendUint64(b, v)
	}
	for _, v := range ctx.msg {
		b = binary.BigEndian.AppendUint64(b, v)
	}
	b = append(b, ctx.block[:ctx.boff]...)
	b = b[:len(b)+len(ctx.block)-ctx.boff] // already zero
	b = binary.BigEndian.AppendUint16(b, uint16(ctx.boff))
	b = binary.BigEndian.AppendUint16(b, uint16(ctx.outlenbytes))
	return b, nil
}

func (ctx *lsh512ContextGo) UnmarshalBinary(b []byte) error {
	if string(b[:len(magic)]) != magic {
		return errors.New("krypto/hash160: invalid hash state identifier")
	}
	if len(b) != marshaledSize {
		return errors.New("krypto/hash160: invalid hash state size")
	}

	b = b[len(magic):]
	for i := range ctx.cv {
		b, ctx.cv[i] = internal.ConsumeUint64(b)
	}
	for i := range ctx.tcv {
		b, ctx.tcv[i] = internal.ConsumeUint64(b)
	}
	for i := range ctx.msg {
		b, ctx.msg[i] = internal.ConsumeUint64(b)
	}
	b = b[copy(ctx.block[:], b[:BlockSize]):]
	ctx.boff = int(binary.BigEndian.Uint16(b))
	ctx.outlenbytes = int(binary.BigEndian.Uint16(b[2:]))
	return nil
}

var (
	_ encoding.BinaryMarshaler   = (*lsh512ContextGo)(nil)
	_ encoding.BinaryUnmarshaler = (*lsh512ContextGo)(nil)
)
