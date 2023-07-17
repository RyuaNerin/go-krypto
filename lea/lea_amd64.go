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

func leaEnc8SSE2(round int, rk []uint32, dst, src []byte) {
	leaEnc4SSE2(round, rk, dst[0x00:], src[0x00:])
	leaEnc4SSE2(round, rk, dst[0x40:], src[0x40:])
}
func leaDec8SSE2(round int, rk []uint32, dst, src []byte) {
	leaDec4SSE2(round, rk, dst[0x00:], src[0x00:])
	leaDec4SSE2(round, rk, dst[0x40:], src[0x40:])
}

type leaContextAsm struct {
	g leaContextGo
}

func newCipherAsm(key []byte) (cipher.Block, error) {
	leaCtx := new(leaContextAsm)

	return leaCtx, leaCtx.g.initContext(key)
}

func newCipherAsmECB(key []byte) (cipher.Block, error) {
	leaCtx := new(leaContextAsm)
	leaCtx.g.ecb = true

	return leaCtx, leaCtx.g.initContext(key)
}

func (leaCtx *leaContextAsm) BlockSize() int {
	return BlockSize
}

func (leaCtx *leaContextAsm) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (dst)", len(dst)))
	}

	if !leaCtx.g.ecb {
		leaEnc1(leaCtx.g.round, leaCtx.g.rk, dst, src)
	} else {
		if len(src)%BlockSize != 0 {
			panic("krypto/lea: input not full blocks")
		}

		remainBlock := len(src) / leaCtx.BlockSize()

		for remainBlock >= 8 {
			remainBlock -= 8
			leaEnc8(leaCtx.g.round, leaCtx.g.rk, dst, src)

			dst, src = dst[0x80:], src[0x80:]
		}

		for remainBlock >= 4 {
			remainBlock -= 4
			leaEnc4(leaCtx.g.round, leaCtx.g.rk, dst, src)

			dst, src = dst[0x40:], src[0x40:]
		}

		for remainBlock > 0 {
			remainBlock -= 1
			leaEnc1(leaCtx.g.round, leaCtx.g.rk, dst, src)

			dst, src = dst[0x10:], src[0x10:]
		}
	}
}

func (leaCtx *leaContextAsm) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (dst)", len(dst)))
	}

	if !leaCtx.g.ecb {
		leaDec1(leaCtx.g.round, leaCtx.g.rk, dst, src)
	} else {
		if len(src)%BlockSize != 0 {
			panic("krypto/lea: input not full blocks")
		}

		remainBlock := len(src) / leaCtx.BlockSize()

		for remainBlock >= 8 {
			remainBlock -= 8
			leaDec8(leaCtx.g.round, leaCtx.g.rk, dst, src)

			dst, src = dst[0x80:], src[0x80:]
		}

		for remainBlock >= 4 {
			remainBlock -= 4
			leaDec4(leaCtx.g.round, leaCtx.g.rk, dst, src)

			dst, src = dst[0x40:], src[0x40:]
		}

		for remainBlock > 0 {
			remainBlock -= 1
			leaDec1(leaCtx.g.round, leaCtx.g.rk, dst, src)

			dst, src = dst[0x10:], src[0x10:]
		}
	}
}
