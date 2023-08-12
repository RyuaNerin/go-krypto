//go:build amd64

package lea

import (
	"crypto/cipher"
	"testing"
)

func Benchmark_BlockMode_CBC_Decrypt_1Block_Go(b *testing.B) {
	benchAll(b, cbcGo(1, cipher.NewCBCDecrypter))
}
func Benchmark_BlockMode_CBC_Decrypt_1Block_Asm(b *testing.B) {
	benchAll(b, cbcAsm(1, cipher.NewCBCDecrypter))
}

func Benchmark_BlockMode_CBC_Decrypt_4Blocks_Go(b *testing.B) {
	benchAll(b, cbcGo(4, cipher.NewCBCDecrypter))
}
func Benchmark_BlockMode_CBC_Decrypt_4Blocks_Asm(b *testing.B) {
	benchAll(b, cbcAsm(4, cipher.NewCBCDecrypter))
}

func Benchmark_BlockMode_CBC_Decrypt_8Blocks_Go(b *testing.B) {
	benchAll(b, cbcGo(8, cipher.NewCBCDecrypter))
}
func Benchmark_BlockMode_CBC_Decrypt_8Blocks_Asm(b *testing.B) {
	benchAll(b, cbcAsm(8, cipher.NewCBCDecrypter))
}

func cbcGo(blocks int, newBlockMode func(cipher.Block, []byte) cipher.BlockMode) func(b *testing.B, keySize int) {
	return func(b *testing.B, keySize int) {
		var ctx leaContext
		err := ctx.initContext(make([]byte, keySize/8))
		if err != nil {
			b.Error(err)
		}

		bm := newBlockMode(&ctx, make([]byte, BlockSize))

		src := make([]byte, BlockSize*blocks)
		dst := make([]byte, BlockSize*blocks)

		b.ReportAllocs()
		b.ResetTimer()
		b.SetBytes(int64(len(src)))
		for i := 0; i < b.N; i++ {
			bm.CryptBlocks(dst, src)
			copy(src, dst)
		}
	}
}

func cbcAsm(blocks int, newBlockMode func(cipher.Block, []byte) cipher.BlockMode) func(b *testing.B, keySize int) {
	return func(b *testing.B, keySize int) {
		var ctx leaContextAsm
		err := ctx.g.initContext(make([]byte, keySize/8))
		if err != nil {
			b.Error(err)
		}

		bm := newBlockMode(&ctx, make([]byte, BlockSize))

		src := make([]byte, BlockSize*blocks)
		dst := make([]byte, BlockSize*blocks)

		b.ReportAllocs()
		b.ResetTimer()
		b.SetBytes(int64(len(src)))
		for i := 0; i < b.N; i++ {
			bm.CryptBlocks(dst, src)
			copy(src, dst)
		}
	}
}
