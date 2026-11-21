//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh256 && (arm64 || amd64)

package lsh256

import (
	"encoding"
	"hash"
	"unsafe"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/memory"
)

type simdSet interface {
	update(ctx, data unsafe.Pointer, databytelen int64)
	final(ctx *lsh256ContextData) (sum [Size]byte)
}

var defaultSimd simdSet

func newContextAsm(size int) hash.Hash {
	return newContextWith(size, defaultSimd)
}

func sumAsm(size int, data []byte) [Size]byte {
	return sumWith(size, data, defaultSimd)
}

func newContextWith(size int, simd simdSet) hash.Hash {
	ctx := &lsh256ContextAsm{
		lsh256Context: lsh256Context{
			lsh256ContextData: lsh256ContextData{
				algType: uint64(size),
			},
		},
		backend: simd,
	}
	ctx.z.SetValid()
	ctx.Reset()
	return ctx
}

func sumWith(size int, data []byte, simd simdSet) (sum [Size]byte) {
	ctx := lsh256ContextAsm{
		lsh256Context: lsh256Context{
			lsh256ContextData: lsh256ContextData{
				algType: uint64(size),
			},
		},
		backend: simd,
	}
	defer ctx.Zeroize()

	ctx.z.SetValid()
	ctx.Reset()
	ctx.Write(data)
	return ctx.backend.final(&ctx.lsh256ContextData)
}

type lsh256ContextAsm struct {
	lsh256Context

	backend simdSet
}

var (
	_ krypto.Hash                = (*lsh256ContextAsm)(nil)
	_ encoding.BinaryMarshaler   = (*lsh256ContextAsm)(nil)
	_ encoding.BinaryUnmarshaler = (*lsh256ContextAsm)(nil)
	_ encoding.BinaryAppender    = (*lsh256ContextAsm)(nil)
)

func (ctx *lsh256ContextAsm) Zeroize() {
	ctx.backend = nil
	ctx.lsh256Context.Zeroize()
}

func (ctx *lsh256ContextAsm) Size() int {
	return int(ctx.algType)
}

func (ctx *lsh256ContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *lsh256ContextAsm) Write(data []byte) (n int, err error) {
	if err = ctx.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, nil
	}
	ctx.backend.update(unsafe.Pointer(ctx), memory.Pointer(data), int64(len(data)))

	return len(data), nil
}

func (ctx *lsh256ContextAsm) Sum(p []byte) []byte {
	ctx.z.ThrowIfZeroized()

	ctx2 := ctx.lsh256ContextData
	defer ctx2.Zeroize()

	sum := ctx.backend.final(&ctx2)
	size := int(ctx.algType)
	if size < 0 || uint64(size) != ctx.algType || size > len(sum) {
		panic(msgInvalidInternalState)
	}
	return append(p, sum[:size]...)
}
