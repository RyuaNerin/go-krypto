package kcdsa

import (
	"crypto/sha256"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/randutil"
)

func BenchmarkPPGF(b *testing.B) {
	rnd, _ := randutil.NewBenchRand()

	h := sha256.New()

	dstBits := h.BlockSize()*h.BlockSize() + 3

	src1 := make([]byte, 2048/8)
	src2 := make([]byte, 2048/8)
	dst := make([]byte, dstBits)

	rnd.Read(src1)
	ppgf := newPPGF(h, src1)

	b.ReportAllocs()
	b.SetBytes(int64(dstBits))
	for b.Loop() {
		src1, src2 = src2, src1
		dst = ppgf.Generate(dst[:0], dstBits, src2)
	}
}
