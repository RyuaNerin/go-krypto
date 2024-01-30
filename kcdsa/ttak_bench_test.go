package kcdsa

import (
	"crypto/rand"
	"io"
	"testing"

	kcdsainternal "github.com/RyuaNerin/go-krypto/internal/kcdsa"
	. "github.com/RyuaNerin/testingutil"
)

func Benchmark_GenerateParametersTTAK(b *testing.B) {
	BA(b, as, func(b *testing.B, sz int) {
		var params Parameters
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, err := GenerateParametersTTAK(&params, rnd, ParameterSizes(sz)); err != nil {
				b.Error(err)
				return
			}
		}
	}, false)
}

func Benchmark_RegenerateParametersTTAK(b *testing.B) {
	BA(b, as, func(b *testing.B, sz int) {
		var params Parameters
		if _, err := GenerateParametersTTAK(&params, rnd, ParameterSizes(sz)); err != nil {
			b.Error(err)
			return
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := RegenerateParametersTTAK(&params, rnd, ParameterSizes(sz)); err != nil {
				b.Error(err)
				return
			}
		}
	}, false)
}

func Benchmark_GenerateKeyTTAK(b *testing.B) {
	BA(b, as, func(b *testing.B, sz int) {
		var err error

		var priv PrivateKey
		if _, err = GenerateParametersTTAK(&priv.Parameters, rnd, ParameterSizes(sz)); err != nil {
			b.Error(err)
			return
		}

		var xkey, upri []byte

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if xkey, upri, err = GenerateKeyTTAK(&priv, rnd, xkey, upri, ParameterSizes(sz)); err != nil {
				b.Error(err)
				return
			}
		}
	}, false)
}

const testBits = 4096

func Benchmark_ppgf(b *testing.B) {
	BA(b, as, func(b *testing.B, sz int) {
		buf := make([]byte, testBits/8)
		seed := make([]byte, testBits/8)
		if _, err := io.ReadFull(rnd, seed); err != nil {
			b.Error(err)
			return
		}

		d, _ := kcdsainternal.GetDomain(sz)
		h := d.NewHash()

		b.ReportAllocs()
		b.ResetTimer()
		b.SetBytes(testBits)
		for i := 0; i < b.N; i++ {
			buf = ppgf(buf, seed, testBits, h)
			copy(seed, buf)
		}
	}, false)
}

func Benchmark_ppgf_readfull(b *testing.B) {
	buf := make([]byte, testBits/8)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetBytes(testBits)
	for i := 0; i < b.N; i++ {
		if _, err := io.ReadFull(rand.Reader, buf); err != nil {
			b.Error(err)
			return
		}
	}
}
