//go:build amd64 || arm64 || loong64 || mips || mips64 || mips64le || mipsle || ppc64 || ppc64le || riscv64

package seed

import (
	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/alias"
)

var _ krypto.Block4 = (*seedContext)(nil)

func (ctx *seedContext) Encrypt4(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < 4*BlockSize {
		panic(msgNotFullBlock)
	}
	if len(dst) < 4*BlockSize {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:4*BlockSize], src[:4*BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[4*BlockSize-1]
	_ = src[4*BlockSize-1]

	encrypt2Go(&ctx.pdwRoundKey, (*block2)(dst[0*BlockSize:2*BlockSize]), (*block2)(src[0*BlockSize:2*BlockSize]))
	encrypt2Go(&ctx.pdwRoundKey, (*block2)(dst[2*BlockSize:4*BlockSize]), (*block2)(src[2*BlockSize:4*BlockSize]))
}

func (ctx *seedContext) Decrypt4(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < 4*BlockSize {
		panic(msgNotFullBlock)
	}
	if len(dst) < 4*BlockSize {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:4*BlockSize], src[:4*BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[4*BlockSize-1]
	_ = src[4*BlockSize-1]

	decrypt2Go(&ctx.pdwRoundKey, (*block2)(dst[0*BlockSize:2*BlockSize]), (*block2)(src[0*BlockSize:2*BlockSize]))
	decrypt2Go(&ctx.pdwRoundKey, (*block2)(dst[2*BlockSize:4*BlockSize]), (*block2)(src[2*BlockSize:4*BlockSize]))
}
