//go:build amd64

package lea

import (
	"crypto/cipher"
	"testing"
)

func Benchmark_LEA128_CBC_Encrypt_1Block_Go(b *testing.B)  { cbcGo(b, 128, 1, cbcEnc) }
func Benchmark_LEA128_CBC_Encrypt_1Block_Asm(b *testing.B) { cbcAsm(b, 256, 1, cbcEnc) }
func Benchmark_LEA196_CBC_Encrypt_1Block_Go(b *testing.B)  { cbcGo(b, 196, 1, cbcEnc) }
func Benchmark_LEA196_CBC_Encrypt_1Block_Asm(b *testing.B) { cbcAsm(b, 256, 1, cbcEnc) }
func Benchmark_LEA256_CBC_Encrypt_1Block_Go(b *testing.B)  { cbcGo(b, 256, 1, cbcEnc) }
func Benchmark_LEA256_CBC_Encrypt_1Block_Asm(b *testing.B) { cbcAsm(b, 256, 1, cbcEnc) }

func Benchmark_LEA128_CBC_Encrypt_4Blocks_Go(b *testing.B)  { cbcGo(b, 128, 4, cbcEnc) }
func Benchmark_LEA128_CBC_Encrypt_4Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 4, cbcEnc) }
func Benchmark_LEA196_CBC_Encrypt_4Blocks_Go(b *testing.B)  { cbcGo(b, 196, 4, cbcEnc) }
func Benchmark_LEA196_CBC_Encrypt_4Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 4, cbcEnc) }
func Benchmark_LEA256_CBC_Encrypt_4Blocks_Go(b *testing.B)  { cbcGo(b, 256, 4, cbcEnc) }
func Benchmark_LEA256_CBC_Encrypt_4Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 4, cbcEnc) }

func Benchmark_LEA128_CBC_Encrypt_8Blocks_Go(b *testing.B)  { cbcGo(b, 128, 8, cbcEnc) }
func Benchmark_LEA128_CBC_Encrypt_8Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 8, cbcEnc) }
func Benchmark_LEA196_CBC_Encrypt_8Blocks_Go(b *testing.B)  { cbcGo(b, 196, 8, cbcEnc) }
func Benchmark_LEA196_CBC_Encrypt_8Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 8, cbcEnc) }
func Benchmark_LEA256_CBC_Encrypt_8Blocks_Go(b *testing.B)  { cbcGo(b, 256, 8, cbcEnc) }
func Benchmark_LEA256_CBC_Encrypt_8Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 8, cbcEnc) }

func Benchmark_LEA128_CBC_Decrypt_1Block_Go(b *testing.B)  { cbcGo(b, 128, 1, cbcDec) }
func Benchmark_LEA128_CBC_Decrypt_1Block_Asm(b *testing.B) { cbcAsm(b, 256, 1, cbcDec) }
func Benchmark_LEA196_CBC_Decrypt_1Block_Go(b *testing.B)  { cbcGo(b, 196, 1, cbcDec) }
func Benchmark_LEA196_CBC_Decrypt_1Block_Asm(b *testing.B) { cbcAsm(b, 256, 1, cbcDec) }
func Benchmark_LEA256_CBC_Decrypt_1Block_Go(b *testing.B)  { cbcGo(b, 256, 1, cbcDec) }
func Benchmark_LEA256_CBC_Decrypt_1Block_Asm(b *testing.B) { cbcAsm(b, 256, 1, cbcDec) }

func Benchmark_LEA128_CBC_Decrypt_4Blocks_Go(b *testing.B)  { cbcGo(b, 128, 4, cbcDec) }
func Benchmark_LEA128_CBC_Decrypt_4Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 4, cbcDec) }
func Benchmark_LEA196_CBC_Decrypt_4Blocks_Go(b *testing.B)  { cbcGo(b, 196, 4, cbcDec) }
func Benchmark_LEA196_CBC_Decrypt_4Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 4, cbcDec) }
func Benchmark_LEA256_CBC_Decrypt_4Blocks_Go(b *testing.B)  { cbcGo(b, 256, 4, cbcDec) }
func Benchmark_LEA256_CBC_Decrypt_4Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 4, cbcDec) }

func Benchmark_LEA128_CBC_Decrypt_8Blocks_Go(b *testing.B)  { cbcGo(b, 128, 8, cbcDec) }
func Benchmark_LEA128_CBC_Decrypt_8Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 8, cbcDec) }
func Benchmark_LEA196_CBC_Decrypt_8Blocks_Go(b *testing.B)  { cbcGo(b, 196, 8, cbcDec) }
func Benchmark_LEA196_CBC_Decrypt_8Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 8, cbcDec) }
func Benchmark_LEA256_CBC_Decrypt_8Blocks_Go(b *testing.B)  { cbcGo(b, 256, 8, cbcDec) }
func Benchmark_LEA256_CBC_Decrypt_8Blocks_Asm(b *testing.B) { cbcAsm(b, 256, 8, cbcDec) }

func cbcGo(b *testing.B, keySize int, blocks int, newBlockMode func(cipher.Block, []byte) cipher.BlockMode) {
	var ctx leaContextGo
	err := ctx.initContext(make([]byte, keySize/8))
	if err != nil {
		b.Error(err)
	}

	bm := newBlockMode(&ctx, make([]byte, BlockSize))

	src := make([]byte, BlockSize*blocks)
	dst := make([]byte, BlockSize*blocks)
	b.ResetTimer()
	b.SetBytes(int64(len(src)))
	for i := 0; i < b.N; i++ {
		bm.CryptBlocks(dst, src)
		copy(src, dst)
	}
}

func cbcAsm(b *testing.B, keySize int, blocks int, newBlockMode func(cipher.Block, []byte) cipher.BlockMode) {
	var ctx leaContextAsm
	err := ctx.g.initContext(make([]byte, keySize/8))
	if err != nil {
		b.Error(err)
	}

	bm := newBlockMode(&ctx, make([]byte, BlockSize))

	src := make([]byte, BlockSize*blocks)
	dst := make([]byte, BlockSize*blocks)
	b.ResetTimer()
	b.SetBytes(int64(len(src)))
	for i := 0; i < b.N; i++ {
		bm.CryptBlocks(dst, src)
		copy(src, dst)
	}
}
