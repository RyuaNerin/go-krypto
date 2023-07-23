//go:build amd64

package lea

import (
	"crypto/cipher"
	"testing"
)

func Benchmark_LEA128_CTR_Go(b *testing.B)  { ctrGo(b, 128) }
func Benchmark_LEA128_CTR_Asm(b *testing.B) { ctrAsm(b, 128) }

func Benchmark_LEA196_CTR_Go(b *testing.B)  { ctrGo(b, 196) }
func Benchmark_LEA196_CTR_Asm(b *testing.B) { ctrAsm(b, 196) }

func Benchmark_LEA256_CTR_Go(b *testing.B)  { ctrGo(b, 256) }
func Benchmark_LEA256_CTR_Asm(b *testing.B) { ctrAsm(b, 256) }

func ctrGo(b *testing.B, keySize int) {
	var ctx leaContextGo
	err := ctx.initContext(make([]byte, keySize/8))
	if err != nil {
		b.Error(err)
	}

	ctr := cipher.NewCTR(&ctx, make([]byte, BlockSize))

	src := make([]byte, BlockSize)
	dst := make([]byte, BlockSize)
	b.ResetTimer()
	b.SetBytes(int64(len(src)))
	for i := 0; i < b.N; i++ {
		ctr.XORKeyStream(dst, src)
		copy(src, dst)
	}
}

func ctrAsm(b *testing.B, keySize int) {
	var ctx leaContextAsm
	err := ctx.g.initContext(make([]byte, keySize/8))
	if err != nil {
		b.Error(err)
	}

	ctr := cipher.NewCTR(&ctx, make([]byte, BlockSize))

	src := make([]byte, BlockSize)
	dst := make([]byte, BlockSize)
	b.ResetTimer()
	b.SetBytes(int64(len(src)))
	for i := 0; i < b.N; i++ {
		ctr.XORKeyStream(dst, src)
		copy(src, dst)
	}
}
