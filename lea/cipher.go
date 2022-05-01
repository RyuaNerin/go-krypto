package lea

import (
	"crypto/cipher"
	"fmt"
)

var (
	useASM = false
)

type KeySizeError int

func (k KeySizeError) Error() string {
	return fmt.Sprintf("krypto/lea: invalid key size %d", int(k))
}

// NewCipher creates and returns a new cipher.Block.
// The key argument should be the LEA key, either 16, 24, or 32 bytes to select LEA-128, LEA-192, or LEA-256.
func NewCipher(key []byte) (cipher.Block, error) {
	leaCtx := new(leaContext)
	err := initContext(leaCtx, key)
	return leaCtx, err
}

// NewCipherECB creates and returns a new cipher.Block by ECB mode.
// This function can be useful in amd64.
// The key argument should be the LEA key, either 16, 24, or 32 bytes to select LEA-128, LEA-192, or LEA-256.
func NewCipherECB(key []byte) (cipher.Block, error) {
	leaCtx := new(leaContext)
	leaCtx.ecb = true
	err := initContext(leaCtx, key)
	return leaCtx, err
}

type leaContext struct {
	round int
	rk    []uint32
	ecb   bool
}

func initContext(leaCtx *leaContext, key []byte) error {
	var rkSize int

	l := len(key)
	switch l {
	case 16:
		rkSize = 144
	case 24:
		rkSize = 168
	case 32:
		rkSize = 192
	default:
		return KeySizeError(l)
	}

	if len(leaCtx.rk) < rkSize {
		leaCtx.rk = make([]uint32, rkSize)
	}
	leaCtx.round = leaSetKeyGo(leaCtx.rk, key)

	return nil
}

func (leaCtx *leaContext) BlockSize() int {
	return BlockSize
}

func (leaCtx *leaContext) Encrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (dst)", len(dst)))
	}

	if !leaCtx.ecb {
		leaEnc1Go(leaCtx.round, leaCtx.rk, dst, src)
	} else {
		if len(src)%BlockSize != 0 {
			panic("krypto/lea: input not full blocks")
		}

		remainBlock := len(src) / leaCtx.BlockSize()

		for remainBlock >= 8 {
			remainBlock -= 8
			leaEnc8(leaCtx.round, leaCtx.rk, dst, src)

			dst, src = dst[0x80:], src[0x80:]
		}

		for remainBlock >= 4 {
			remainBlock -= 4
			leaEnc4(leaCtx.round, leaCtx.rk, dst, src)

			dst, src = dst[0x40:], src[0x40:]
		}

		for remainBlock > 0 {
			remainBlock -= 1
			leaEnc1(leaCtx.round, leaCtx.rk, dst, src)

			dst, src = dst[0x10:], src[0x10:]
		}
	}
}

func (leaCtx *leaContext) Decrypt(dst, src []byte) {
	if len(src) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (src)", len(src)))
	}
	if len(dst) < BlockSize {
		panic(fmt.Sprintf("krypto/lea: invalid block size %d (dst)", len(dst)))
	}

	if !leaCtx.ecb {
		leaDec1Go(leaCtx.round, leaCtx.rk, dst, src)
	} else {
		if len(src)%BlockSize != 0 {
			panic("krypto/lea: input not full blocks")
		}

		remainBlock := len(src) / leaCtx.BlockSize()

		for remainBlock >= 8 {
			remainBlock -= 8
			leaDec8(leaCtx.round, leaCtx.rk, dst, src)

			dst, src = dst[0x80:], src[0x80:]
		}

		for remainBlock >= 4 {
			remainBlock -= 4
			leaDec4(leaCtx.round, leaCtx.rk, dst, src)

			dst, src = dst[0x40:], src[0x40:]
		}

		for remainBlock > 0 {
			remainBlock -= 1
			leaDec1(leaCtx.round, leaCtx.rk, dst, src)

			dst, src = dst[0x10:], src[0x10:]
		}
	}
}
