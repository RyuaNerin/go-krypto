//go:build amd64

package lea

import (
	"bytes"
	"testing"
)

const (
	testBlocks = 16 * 1024
)

func Test_Encrypt_4Blocks_SSE2(t *testing.T) { testAll(t, tb(4, leaEnc4Go, leaEnc4, true)) }
func Test_Decrypt_4Blocks_SSE2(t *testing.T) { testAll(t, tb(4, leaDec4Go, leaDec4, true)) }

func Test_Encrypt_8Blocks_SSE2(t *testing.T) { testAll(t, tb(8, leaEnc8Go, leaEnc8SSE2, true)) }
func Test_Decrypt_8Blocks_SSE2(t *testing.T) { testAll(t, tb(8, leaDec8Go, leaDec8SSE2, true)) }

func Test_Encrypt_8Blocks_AVX2(t *testing.T) { testAll(t, tb(8, leaEnc8Go, leaEnc8AVX2, hasAVX2)) }
func Test_Decrypt_8Blocks_AVX2(t *testing.T) { testAll(t, tb(8, leaDec8Go, leaDec8AVX2, hasAVX2)) }

func tb(blocks int, funcGo, funcAsm funcBlock, do bool) func(t *testing.T, keySize int) {
	return func(t *testing.T, keySize int) {
		if !do {
			t.Skip()
			return
		}

		k := make([]byte, keySize/8)

		srcGo := make([]byte, BlockSize*blocks)
		dstGo := make([]byte, BlockSize*blocks)

		srcAsm := make([]byte, BlockSize*blocks)
		dstAsm := make([]byte, BlockSize*blocks)

		var ctx leaContextAsm
		err := ctx.g.initContext(k)
		if err != nil {
			t.Error(err)
		}

		for i := 0; i < testBlocks/blocks; i++ {
			funcGo(ctx.g.round, ctx.g.rk, dstGo, srcGo)
			funcAsm(ctx.g.round, ctx.g.rk, dstAsm, srcAsm)

			if !bytes.Equal(dstGo, dstAsm) {
				t.Error("did not match")
			}

			copy(srcGo, dstGo)
			copy(srcAsm, dstAsm)
		}
	}
}
