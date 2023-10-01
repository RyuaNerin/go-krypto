//go:build amd64 && gc && !purego

package lea

import (
	"crypto/cipher"
	"testing"
)

func Benchmark_BlockMode_CBC_Decrypt_Std(b *testing.B) {
	benchAll(
		b,
		benchCBC(
			func(key []byte) (cipher.Block, error) {
				var ctx nonCipherContext
				return &ctx, ctx.ctx.initContext(key)
			},
			cipher.NewCBCDecrypter,
		),
	)
}
func Benchmark_BlockMode_CBC_Decrypt_Asm(b *testing.B) {
	benchAll(
		b,
		benchCBC(
			func(key []byte) (cipher.Block, error) {
				var ctx leaContext
				return &ctx, ctx.initContext(key)
			},
			cipher.NewCBCDecrypter,
		),
	)
}

func benchCBC(newCipher func(key []byte) (cipher.Block, error), newBlockMode func(cipher.Block, []byte) cipher.BlockMode) func(b *testing.B, keySize int) {
	return func(b *testing.B, keySize int) {
		key := make([]byte, keySize/8)
		iv := make([]byte, BlockSize)
		src := make([]byte, BlockSize*8)
		dst := make([]byte, BlockSize*8)
		rnd.Read(key)
		rnd.Read(iv)
		rnd.Read(src)

		ctx, err := newCipher(key)
		if err != nil {
			b.Error(err)
		}

		bm := newBlockMode(ctx, iv)

		b.ReportAllocs()
		b.SetBytes(int64(len(src)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bm.CryptBlocks(dst, src)
			copy(src, dst)
		}
	}
}
