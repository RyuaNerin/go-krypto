//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh512 && (arm64 || amd64)

package lsh512

import (
	"hash"
	"unsafe"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/memory"
)

type simdSet interface {
	update(ctx, data unsafe.Pointer, databytelen int64)
	final(ctx *lsh512ContextData) (sum [Size]byte)
}

var defaultSimd simdSet

func newContextAsm(size int) hash.Hash {
	return newContextWith(size, defaultSimd)
}

func sumAsm(size int, data []byte) [Size]byte {
	return sumWith(size, data, defaultSimd)
}

func newContextWith(size int, simd simdSet) hash.Hash {
	ctx := &lsh512ContextAsm{
		lsh512Context: lsh512Context{
			lsh512ContextData: lsh512ContextData{
				algType: uint64(size),
			},
		},
		backend: simd,
	}
	ctx.Reset()
	return ctx
}

func sumWith(size int, data []byte, simd simdSet) [Size]byte {
	ctx := lsh512ContextAsm{
		lsh512Context: lsh512Context{
			lsh512ContextData: lsh512ContextData{
				algType: uint64(size),
			},
		},
		backend: simd,
	}
	defer ctx.Zeroize()

	ctx.Reset()
	ctx.Write(data)
	return ctx.backend.final(&ctx.lsh512ContextData)
}

type lsh512ContextAsm struct {
	lsh512Context

	backend simdSet
}

var _ krypto.Hash = (*lsh512ContextAsm)(nil)

func (ctx *lsh512ContextAsm) Zeroize() {
	ctx.backend = nil
	ctx.lsh512Context.Zeroize()
}

func (ctx *lsh512ContextAsm) Size() int {
	return int(ctx.algType)
}

func (ctx *lsh512ContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *lsh512ContextAsm) Write(data []byte) (n int, err error) {
	if err = ctx.z.ReturnIfZeroized(); err != nil {
		return 0, err
	}

	if len(data) == 0 {
		return 0, nil
	}
	ctx.backend.update(unsafe.Pointer(ctx), memory.Pointer(data), int64(len(data)))

	return len(data), nil
}

func (ctx *lsh512ContextAsm) Sum(p []byte) []byte {
	ctx.z.ThrowIfZeroized()

	ctx2 := ctx.lsh512ContextData
	defer ctx2.Zeroize()

	sum := ctx.backend.final(&ctx2)
	size := int(ctx.algType)
	if size < 0 || uint64(size) != ctx.algType || size > len(sum) {
		panic(msgInvalidInternalState)
	}
	return append(p, sum[:size]...)
}
