//go:build (arm64 || amd64) && !purego && (!gccgo || go1.18)
// +build arm64 amd64
// +build !purego
// +build !gccgo go1.18

package lsh256

import (
	"encoding"
	"errors"

	"github.com/RyuaNerin/go-krypto/internal"
)

const (
	magic         = "lsh\x00"
	marshaledSize = len(magic) +
		1 + /* algType */
		1 + /* remainDataByteLen */
		4*8 + /* cvL */
		4*8 + /* cvR */
		BlockSize /* lastBlock */
)

func (ctx *lsh256ContextAsm) MarshalBinary() ([]byte, error) {
	b := make([]byte, 0, marshaledSize)
	b = append(b, magic...)

	b = internal.AppendBigUint8(b, uint8(ctx.algType))

	b = internal.AppendBigUint8(b, uint8(ctx.remainDataByteLen))

	for _, v := range ctx.cvL {
		b = internal.AppendBigUint64(b, v)
	}

	for _, v := range ctx.cvR {
		b = internal.AppendBigUint64(b, v)
	}

	b = append(b, ctx.lastBlock[:ctx.remainDataByteLen]...)
	b = b[:len(b)+len(ctx.lastBlock)-ctx.remainDataByteLen] // already zero

	return b, nil
}

func (ctx *lsh256ContextAsm) UnmarshalBinary(b []byte) error {
	if len(b) < len(magic) || string(b[:len(magic)]) != magic {
		return errors.New(msgInvalidHashStateIdentifier)
	}
	if len(b) != marshaledSize {
		return errors.New(msgInvalidHashStateSize)
	}

	b = b[len(magic):]

	ctx.algType = int(b[0])
	ctx.remainDataByteLen = int(b[1])
	b = b[2:]

	for i := range ctx.cvL {
		b, ctx.cvL[i] = internal.ConsumeUint64(b)
	}

	for i := range ctx.cvR {
		b, ctx.cvR[i] = internal.ConsumeUint64(b)
	}

	copy(ctx.lastBlock[:], b)

	return nil
}

var (
	_ encoding.BinaryMarshaler   = (*lsh256ContextAsm)(nil)
	_ encoding.BinaryUnmarshaler = (*lsh256ContextAsm)(nil)
)
