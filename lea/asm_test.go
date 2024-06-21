//go:build (arm64 || amd64) && !purego && (!gccgo || go1.18)
// +build arm64 amd64
// +build !purego
// +build !gccgo go1.18

package lea

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

//nolint:unused
func leaEnc4Go(ctx *leaContext, dst, src []byte) {
	leaEnc1Go(ctx, dst[BlockSize*0:], src[BlockSize*0:])
	leaEnc1Go(ctx, dst[BlockSize*1:], src[BlockSize*1:])
	leaEnc1Go(ctx, dst[BlockSize*2:], src[BlockSize*2:])
	leaEnc1Go(ctx, dst[BlockSize*3:], src[BlockSize*3:])
}

//nolint:unused
func leaDec4Go(ctx *leaContext, dst, src []byte) {
	leaDec1Go(ctx, dst[BlockSize*0:], src[BlockSize*0:])
	leaDec1Go(ctx, dst[BlockSize*1:], src[BlockSize*1:])
	leaDec1Go(ctx, dst[BlockSize*2:], src[BlockSize*2:])
	leaDec1Go(ctx, dst[BlockSize*3:], src[BlockSize*3:])
}

//nolint:unused
func leaEnc8Go(ctx *leaContext, dst, src []byte) {
	leaEnc1Go(ctx, dst[BlockSize*0:], src[BlockSize*0:])
	leaEnc1Go(ctx, dst[BlockSize*1:], src[BlockSize*1:])
	leaEnc1Go(ctx, dst[BlockSize*2:], src[BlockSize*2:])
	leaEnc1Go(ctx, dst[BlockSize*3:], src[BlockSize*3:])
	leaEnc1Go(ctx, dst[BlockSize*4:], src[BlockSize*4:])
	leaEnc1Go(ctx, dst[BlockSize*5:], src[BlockSize*5:])
	leaEnc1Go(ctx, dst[BlockSize*6:], src[BlockSize*6:])
	leaEnc1Go(ctx, dst[BlockSize*7:], src[BlockSize*7:])
}

//nolint:unused
func leaDec8Go(ctx *leaContext, dst, src []byte) {
	leaDec1Go(ctx, dst[BlockSize*0:], src[BlockSize*0:])
	leaDec1Go(ctx, dst[BlockSize*1:], src[BlockSize*1:])
	leaDec1Go(ctx, dst[BlockSize*2:], src[BlockSize*2:])
	leaDec1Go(ctx, dst[BlockSize*3:], src[BlockSize*3:])
	leaDec1Go(ctx, dst[BlockSize*4:], src[BlockSize*4:])
	leaDec1Go(ctx, dst[BlockSize*5:], src[BlockSize*5:])
	leaDec1Go(ctx, dst[BlockSize*6:], src[BlockSize*6:])
	leaDec1Go(ctx, dst[BlockSize*7:], src[BlockSize*7:])
}

func bb(f funcBlock) func(c interface{}, dst, src []byte) {
	return func(c interface{}, dst, src []byte) {
		f(&c.(*leaContextAsm).leaContext, dst, src)
	}
}

func tb(blocks int, funcGo, funcAsm funcBlock) func(t *testing.T, keySize int) {
	return func(t *testing.T, keySize int) {
		BTTC(
			t,
			keySize,
			0,
			blocks*BlockSize,
			0,
			func(key, additional []byte) (interface{}, error) {
				return newCipherGo(key)
			},
			func(ctx interface{}, dst, src []byte) { funcGo(ctx.(*leaContext), dst, src) },
			func(ctx interface{}, dst, src []byte) { funcAsm(ctx.(*leaContext), dst, src) },
			false,
		)
	}
}
