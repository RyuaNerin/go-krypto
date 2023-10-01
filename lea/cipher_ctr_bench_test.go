//go:build amd64 && gc && !purego

package lea

import (
	"crypto/cipher"
	"testing"
)

func Benchmark_BlockMode_CTR_Std(b *testing.B) {
	benchAll(
		b,
		benchCTR(
			func(key []byte) (cipher.Block, error) {
				var ctx nonCipherContext
				return &ctx, ctx.ctx.initContext(key)
			},
		),
	)
}
func Benchmark_BlockMode_CTR_Krypto(b *testing.B) {
	benchAll(
		b,
		benchCTR(
			func(key []byte) (cipher.Block, error) {
				var ctx leaContext
				return &ctx, ctx.initContext(key)
			},
		),
	)
}

func benchCTR(
	newCipher func(key []byte) (cipher.Block, error),
) func(b *testing.B, keySize int) {
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

		ctr := cipher.NewCTR(ctx, make([]byte, BlockSize))

		b.ReportAllocs()
		b.SetBytes(int64(len(src)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ctr.XORKeyStream(dst, src)
			copy(src, dst)
		}
	}
}
