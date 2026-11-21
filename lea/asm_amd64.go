//go:build !noasm && !noasm_krypto && !noasm_krypto_lea && amd64

package lea

import (
	"crypto/cipher"
	"unsafe"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal/alias"

	"golang.org/x/sys/cpu"
)

var (
	hasAVX2   = cpu.X86.HasAVX2 && cpu.X86.HasAVX && cpu.X86.HasSSSE3 && cpu.X86.HasSSE2
	hasAVX512 = cpu.X86.HasAVX512F && hasAVX2
)

var newCipherAsm = newCipherAsmSSE2

func init() {
	switch {
	case hasAVX512:
		newCipherAsm = newCipherAsmAVX512
	case hasAVX2:
		newCipherAsm = newCipherAsmAVX2
	}
}

func newCipherAsmSSE2(key []byte) (cipher.Block, error) {
	ctx := new(leaContextAsm)
	ctx.z.SetValid()
	ctx.round = setKeyGo(&ctx.rk, key)

	return ctx, nil
}

type (
	leaContextAsmAVX2   struct{ leaContextAsm }
	leaContextAsmAVX512 struct{ leaContextAsmAVX2 }
)

var (
	_ krypto.Block4  = (*leaContextAsm)(nil)
	_ krypto.Block8  = (*leaContextAsmAVX2)(nil)
	_ krypto.Block16 = (*leaContextAsmAVX512)(nil)
)

func newCipherAsmAVX2(key []byte) (cipher.Block, error) {
	ctx := new(leaContextAsmAVX2)
	ctx.z.SetValid()
	ctx.round = setKeyGo(&ctx.rk, key)

	return ctx, nil
}

func newCipherAsmAVX512(key []byte) (cipher.Block, error) {
	ctx := new(leaContextAsmAVX512)
	ctx.z.SetValid()
	ctx.round = setKeyGo(&ctx.rk, key)

	return ctx, nil
}

func (ctx *leaContextAsmAVX2) Encrypt8(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < BlockSize*8 {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize*8 {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:8*BlockSize], src[:8*BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[BlockSize*8-1]
	_ = src[BlockSize*8-1]

	__lea_encrypt_8block(unsafe.Pointer(&ctx.rk), int64(ctx.round), unsafe.Pointer(&dst[0]), unsafe.Pointer(&src[0]))
}

func (ctx *leaContextAsmAVX2) Decrypt8(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < BlockSize*8 {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize*8 {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:8*BlockSize], src[:8*BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[BlockSize*8-1]
	_ = src[BlockSize*8-1]

	__lea_decrypt_8block(unsafe.Pointer(&ctx.rk), int64(ctx.round), unsafe.Pointer(&dst[0]), unsafe.Pointer(&src[0]))
}

func (ctx *leaContextAsmAVX512) Encrypt16(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < BlockSize*16 {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize*16 {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:16*BlockSize], src[:16*BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[BlockSize*16-1]
	_ = src[BlockSize*16-1]

	__lea_encrypt_16block(unsafe.Pointer(&ctx.rk), int64(ctx.round), unsafe.Pointer(&dst[0]), unsafe.Pointer(&src[0]))
}

func (ctx *leaContextAsmAVX512) Decrypt16(dst, src []byte) {
	ctx.z.ThrowIfZeroized()

	if len(src) < BlockSize*16 {
		panic(msgNotFullBlock)
	}
	if len(dst) < BlockSize*16 {
		panic(msgNotFullBlock)
	}
	if alias.InexactOverlap(dst[:16*BlockSize], src[:16*BlockSize]) {
		panic(msgBufferOverlap)
	}

	_ = dst[BlockSize*16-1]
	_ = src[BlockSize*16-1]

	__lea_decrypt_16block(unsafe.Pointer(&ctx.rk), int64(ctx.round), unsafe.Pointer(&dst[0]), unsafe.Pointer(&src[0]))
}
