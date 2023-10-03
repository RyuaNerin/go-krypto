//go:build amd64 && gc && !purego

package lea

import (
	"crypto/cipher"
	"testing"

	"github.com/RyuaNerin/go-krypto/testingutil"
)

func Test_BlockMode_CTR(t *testing.T) {
	testingutil.TA(t, as, testCTR)
}

func testCTR(t *testing.T, keySize int) {
	type ctr struct {
		std, asm cipher.Stream
	}

	testingutil.BTTC(
		t,
		keySize,
		BlockSize, // iv
		BlockSize*8,
		1,
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
				std: cipher.NewCTR(&ctxStd, iv),
				asm: cipher.NewCTR(&ctxLib, iv),
			}
			return data, nil
		},
		func(data interface{}, dst, src []byte) { data.(*ctr).std.XORKeyStream(dst, src) },
		func(data interface{}, dst, src []byte) { data.(*ctr).asm.XORKeyStream(dst, src) },
	)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_BlockMode_CTR_Std(b *testing.B) {
	testingutil.BA(b, as, benchCTR(
		func(key []byte) (cipher.Block, error) {
			var ctx nonCipherContext
			return &ctx, ctx.ctx.initContext(key)
		},
	))
}
func Benchmark_BlockMode_CTR_Krypto(b *testing.B) {
	testingutil.BA(b, as, benchCTR(
		func(key []byte) (cipher.Block, error) {
			var ctx leaContext
			return &ctx, ctx.initContext(key)
		},
	))
}

func benchCTR(
	newCipher func(key []byte) (cipher.Block, error),
) func(b *testing.B, keySize int) {
	return func(b *testing.B, keySize int) {
		testingutil.BBDo(
			b,
			keySize,
			BlockSize,
			BlockSize*8,
			func(key, additional []byte) (interface{}, error) {
				ctx, err := newCipher(key)
				if err != nil {
					return nil, err
				}

				ctr := cipher.NewCTR(ctx, make([]byte, BlockSize))
				return ctr, nil
			},
			func(c interface{}, dst, src []byte) {
				c.(cipher.Stream).XORKeyStream(dst, src)
			},
		)
	}
}
