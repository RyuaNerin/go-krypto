//go:build !noasm && !noasm_krypto && !noasm_krypto_lea && (amd64 || arm64)

package lea

import (
	"unsafe"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/alias"
)

type leaContextAsm struct{ leaContext }

var _ krypto.Block4 = (*leaContextAsm)(nil)

func (ctx *leaContextAsm) Encrypt(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < BlockSize {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:BlockSize], src[:BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[BlockSize-1]
	_ = src[BlockSize-1]

	leaEnc1Go(&ctx.rk, ctx.round, (*[BlockSize]byte)(dst), (*[BlockSize]byte)(src))
}

func (ctx *leaContextAsm) Decrypt(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < BlockSize {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:BlockSize], src[:BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[BlockSize-1]
	_ = src[BlockSize-1]

	leaDec1Go(&ctx.rk, ctx.round, (*[BlockSize]byte)(dst), (*[BlockSize]byte)(src))
}

func (ctx *leaContextAsm) Encrypt4(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < BlockSize*4 {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize*4 {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:4*BlockSize], src[:4*BlockSize]) {
		panic(msgBufferOverlap)
	}

	__lea_encrypt_4block(unsafe.Pointer(&ctx.rk), int64(ctx.round), unsafe.Pointer(&dst[0]), unsafe.Pointer(&src[0]))
}

func (ctx *leaContextAsm) Decrypt4(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < BlockSize*4 {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize*4 {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:4*BlockSize], src[:4*BlockSize]) {
		panic(msgBufferOverlap)
	}

	__lea_decrypt_4block(unsafe.Pointer(&ctx.rk), int64(ctx.round), unsafe.Pointer(&dst[0]), unsafe.Pointer(&src[0]))
}
