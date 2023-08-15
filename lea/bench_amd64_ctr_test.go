//go:build amd64 && gc && !purego

package lea

import (
	"crypto/cipher"
	"testing"
)

func Benchmark_BlockMode_CTR_Go(b *testing.B)  { benchAll(b, ctrGo) }
func Benchmark_BlockMode_CTR_Asm(b *testing.B) { benchAll(b, ctrAsm) }

func ctrGo(b *testing.B, keySize int) {
	var ctx leaContext
	err := ctx.initContext(make([]byte, keySize/8))
	if err != nil {
		b.Error(err)
	}

	ctr := cipher.NewCTR(&ctx, make([]byte, BlockSize))

	src := make([]byte, BlockSize)
	dst := make([]byte, BlockSize)

	rnd.Read(src)

	b.ReportAllocs()
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
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

	rnd.Read(src)

	b.ReportAllocs()
	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctr.XORKeyStream(dst, src)
		copy(src, dst)
	}
}
