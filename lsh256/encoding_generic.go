//go:build !(arm64 || amd64) || purego || (gccgo && !go1.18)
// +build !arm64,!amd64 purego gccgo,!go1.18

package lsh256

import (
	"encoding"
	"encoding/binary"
	"errors"

	"github.com/RyuaNerin/go-krypto/internal"
)

const (
	magic         = "lsh\x01"
	marshaledSize = len(magic) +
		/*cv*/ 16*4 +
		/*tcv*/ 16*4 +
		/*msg*/ 16*(numStep+1)*4 +
		/*block*/ BlockSize +
		/*boff*/ 1 +
		/*outlenbytes*/ 2
)

func (ctx *lsh256ContextGo) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, marshaledSize)
	b = append(b, magic...)
	for _, v := range ctx.cv {
		b = internal.AppendBigUint32(b, v)
	}
	for _, v := range ctx.tcv {
		b = internal.AppendBigUint32(b, v)
	}
	for _, v := range ctx.msg {
		b = internal.AppendBigUint32(b, v)
	}
	b = append(b, ctx.block[:ctx.boff]...)
	b = b[:len(b)+len(ctx.block)-ctx.boff] // already zero
	b = append(b, byte(ctx.boff))
	b = internal.AppendBigUint16(b, uint16(ctx.outlenbytes))
	return b, nil
}

func (ctx *lsh256ContextGo) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic) || string(b[:len(magic)]) != magic {
		return errors.New(msgInvalidHashStateIdentifier)
	}
	if len(b) != marshaledSize {
		return errors.New(msgInvalidHashStateSize)
	}

	b = b[len(magic):]
	for i := range ctx.cv {
		b, ctx.cv[i] = internal.ConsumeBigU32(b)
	}
	for i := range ctx.tcv {
		b, ctx.tcv[i] = internal.ConsumeBigU32(b)
	}
	for i := range ctx.msg {
		b, ctx.msg[i] = internal.ConsumeBigU32(b)
	}
	b = b[copy(ctx.block[:], b[:BlockSize]):]
	ctx.boff = int(b[0])
	ctx.outlenbytes = int(binary.BigEndian.Uint16(b[1:]))
	return nil
}

var (
	_ encoding.BinaryMarshaler   = (*lsh256ContextGo)(nil)
	_ encoding.BinaryUnmarshaler = (*lsh256ContextGo)(nil)
)
