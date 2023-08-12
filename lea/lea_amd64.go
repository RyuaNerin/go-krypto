//go:build amd64

package lea

import (
	"crypto/cipher"
	"fmt"

	"golang.org/x/sys/cpu"
)

var (
	hasAVX2 = cpu.X86.HasAVX2
)

func init() {
	leaEnc4 = leaEnc4SSE2
	leaDec4 = leaDec4SSE2

	leaEnc8 = leaEnc8SSE2
	leaDec8 = leaDec8SSE2

	if hasAVX2 {
		leaEnc8 = leaEnc8AVX2
		leaDec8 = leaDec8AVX2
	}

	leaNew = newCipherAsm
	leaNewECB = newCipherAsmECB
}

func leaEnc8SSE2(ctx *leaContext, dst, src []byte) {
	leaEnc4SSE2(ctx, dst[0x00:], src[0x00:])
	leaEnc4SSE2(ctx, dst[0x40:], src[0x40:])
}
func leaDec8SSE2(ctx *leaContext, dst, src []byte) {
	leaDec4SSE2(ctx, dst[0x00:], src[0x00:])
	leaDec4SSE2(ctx, dst[0x40:], src[0x40:])
}

type leaContextAsm struct {
	g leaContext
}

func newCipherAsm(key []byte) (cipher.Block, error) {
	ctx := new(leaContextAsm)

	if err := ctx.g.initContext(key); err != nil {
		return nil, err
	}
	return ctx, nil
}

func newCipherAsmECB(key []byte) (cipher.Block, error) {
	ctx := new(leaContextAsm)
	ctx.g.ecb = true

	if err := ctx.g.initContext(key); err != nil {
		return nil, err
	}
	return ctx, nil
}

func (ctx *leaContextAsm) BlockSize() int {
	return BlockSize
}

func (ctx *leaContextAsm) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (dst)", len(dst)))
	}

	if !ctx.g.ecb {
		leaEnc1(&ctx.g, dst, src)
	} else {
		if len(src)%BlockSize != 0 {
			panic("krypto/lea: input not full blocks")
		}

		remainBlock := len(src) / ctx.BlockSize()

		for remainBlock >= 8 {
			remainBlock -= 8
			leaEnc8(&ctx.g, dst, src)

			dst, src = dst[0x80:], src[0x80:]
		}

		for remainBlock >= 4 {
			remainBlock -= 4
			leaEnc4(&ctx.g, dst, src)

			dst, src = dst[0x40:], src[0x40:]
		}

		for remainBlock > 0 {
			remainBlock -= 1
			leaEnc1(&ctx.g, dst, src)

			dst, src = dst[0x10:], src[0x10:]
		}
	}
}

func (ctx *leaContextAsm) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (dst)", len(dst)))
	}

	if !ctx.g.ecb {
		leaDec1(&ctx.g, dst, src)
	} else {
		if len(src)%BlockSize != 0 {
			panic("krypto/lea: input not full blocks")
		}

		remainBlock := len(src) / ctx.BlockSize()

		for remainBlock >= 8 {
			remainBlock -= 8
			leaDec8(&ctx.g, dst, src)

			dst, src = dst[0x80:], src[0x80:]
		}

		for remainBlock >= 4 {
			remainBlock -= 4
			leaDec4(&ctx.g, dst, src)

			dst, src = dst[0x40:], src[0x40:]
		}

		for remainBlock > 0 {
			remainBlock -= 1
			leaDec1(&ctx.g, dst, src)

			dst, src = dst[0x10:], src[0x10:]
		}
	}
}
