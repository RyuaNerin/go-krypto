//go:build amd64 && gc && !purego

package lea

import (
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/testingutil"
)

func Test_BlockMode_CBC_Decrypt(t *testing.T) {
	testingutil.TA(t, as, testCBC(cipher.NewCBCDecrypter))
}

func testCBC(newBlockMode func(cipher.Block, []byte) cipher.BlockMode) func(t *testing.T, keySize int) {
	return func(t *testing.T, keySize int) {
		type ctr struct {
			std, asm cipher.BlockMode
		}

		testingutil.BTTC(
			t,
			keySize,
			BlockSize, // iv
			BlockSize*8,
			BlockSize,
			func(key, iv []byte) (interface{}, error) {
				var ctxStd nonCipherContext
				var ctxLib leaContext

				err := ctxStd.ctx.initContext(key)
				if err != nil {
					return nil, err
				}
				err = ctxLib.initContext(key)
				if err != nil {
					return nil, err
				}

				data := &ctr{
					std: newBlockMode(&ctxStd, iv),
					asm: newBlockMode(&ctxLib, iv),
				}
				return data, nil
			},
			func(data interface{}, dst, src []byte) { data.(*ctr).std.CryptBlocks(dst, src) },
			func(data interface{}, dst, src []byte) { data.(*ctr).asm.CryptBlocks(dst, src) },
		)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_BlockMode_CBC_Decrypt_Std(b *testing.B) {
	testingutil.BA(b, as, benchCBC(
		func(key []byte) (cipher.Block, error) {
			var ctx nonCipherContext
			return &ctx, ctx.ctx.initContext(key)
		},
		cipher.NewCBCDecrypter,
	))
}
func Benchmark_BlockMode_CBC_Decrypt_Asm(b *testing.B) {
	testingutil.BA(b, as, benchCBC(
		func(key []byte) (cipher.Block, error) {
			var ctx leaContext
			return &ctx, ctx.initContext(key)
		},
		cipher.NewCBCDecrypter,
	))
}

func benchCBC(newCipher func(key []byte) (cipher.Block, error), newBlockMode func(cipher.Block, []byte) cipher.BlockMode) func(b *testing.B, keySize int) {
	return func(b *testing.B, keySize int) {
		testingutil.BBDo(
			b,
			keySize,
			BlockSize,
			BlockSize*8,
			func(key, iv []byte) (interface{}, error) {
				ctx, err := newCipher(key)
				if err != nil {
					return nil, err
				}

				cbc := newBlockMode(ctx, iv)
				return cbc, nil
			},
			func(c interface{}, dst, src []byte) {
				c.(cipher.BlockMode).CryptBlocks(dst, src)
			},
		)
	}
}
