package kcdsa

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Benchmark_GenerateParameters(b *testing.B) {
	BA(b, as, func(b *testing.B, sz int) {
		var params Parameters
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := GenerateParameters(&params, rnd, ParameterSizes(sz)); err != nil {
				b.Error(err)
				return
			}
		}
	}, false)
}

func Benchmark_GenerateParametersTTAK(b *testing.B) {
	BA(b, as, func(b *testing.B, sz int) {
		var params Parameters
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := GenerateParametersTTAK(&params, rnd, ParameterSizes(sz)); err != nil {
				b.Error(err)
				return
			}
		}
	}, false)
}

func Benchmark_RegenerateParameters(b *testing.B) {
	BA(b, as, func(b *testing.B, sz int) {
		var params Parameters
		if err := GenerateParametersTTAK(&params, rnd, ParameterSizes(sz)); err != nil {
			b.Error(err)
			return
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := RegenerateParameters(&params, rnd, ParameterSizes(sz)); err != nil {
				b.Error(err)
				return
			}
		}
	}, false)
}

func Benchmark_GenerateKey(b *testing.B) {
	BA(b, as, func(b *testing.B, sz int) {
		var priv PrivateKey
		if err := GenerateParameters(&priv.Parameters, rnd, ParameterSizes(sz)); err != nil {
			b.Error(err)
			return
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := GenerateKey(&priv, rnd); err != nil {
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
		if err = GenerateParameters(&priv.Parameters, rnd, ParameterSizes(sz)); err != nil {
			b.Error(err)
			return
		}

		var xkey, upri []byte

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if xkey, upri, err = GenerateKeyWithSeed(&priv, rnd, xkey, upri, ParameterSizes(sz)); err != nil {
				b.Error(err)
				return
			}
		}
	}, false)
}

func Benchmark_Sign(b *testing.B) {
	BA(b, as, func(b *testing.B, sz int) {
		data := []byte(`text`)

		var priv PrivateKey
		if err := GenerateParameters(&priv.Parameters, rnd, ParameterSizes(sz)); err != nil {
			b.Error(err)
			return
		}
		if err := GenerateKey(&priv, rnd); err != nil {
			b.Error(err)
			return
		}

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			r, _, err := Sign(rnd, &priv, ParameterSizes(sz), data)
			if err != nil {
				b.Error(err)
				return
			}
			data = r.Bytes()
		}
	}, false)
}

func Benchmark_Verify(b *testing.B) {
	BA(b, as, func(b *testing.B, sz int) {
		data := []byte(`text`)

		var priv PrivateKey
		if err := GenerateParameters(&priv.Parameters, rnd, ParameterSizes(sz)); err != nil {
			b.Error(err)
		}
		if err := GenerateKey(&priv, rnd); err != nil {
			b.Error(err)
		}

		r, s, err := Sign(rnd, &priv, ParameterSizes(sz), data)
		if err != nil {
			b.Error(err)
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ok := Verify(&priv.PublicKey, ParameterSizes(sz), data, r, s)
			if !ok {
				b.Errorf("%d: Verify failed", i)
				return
			}
		}
	}, false)
}
