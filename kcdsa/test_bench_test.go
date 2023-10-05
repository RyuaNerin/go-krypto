package kcdsa

import (
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

func Benchmark_GenerateParameters_GO(b *testing.B) {
	BA(b, as, func(b *testing.B, ps int) {
		var params Parameters
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := GenerateParameters(&params, rnd, ParameterSizes(ps)); err != nil {
				b.Error(err)
				return
			}
		}
	}, false)
}

func Benchmark_GenerateParameters_TTAK(b *testing.B) {
	BA(b, as, func(b *testing.B, ps int) {
		var params Parameters
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, _, err := GenerateParametersTTAK(&params, rnd, ParameterSizes(ps)); err != nil {
				b.Error(err)
				return
			}
		}
	}, false)
}

func Benchmark_GenerateKey_Go(b *testing.B) {
	BA(b, as, func(b *testing.B, ps int) {
		var priv PrivateKey
		if err := GenerateParameters(&priv.Parameters, rnd, ParameterSizes(ps)); err != nil {
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

func Benchmark_GenerateKey_TTAK(b *testing.B) {
	BA(b, as, func(b *testing.B, ps int) {
		var priv PrivateKey
		if _, _, err := GenerateParametersTTAK(&priv.Parameters, rnd, ParameterSizes(ps)); err != nil {
			b.Error(err)
			return
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := GenerateKeyTTAK(&priv, rnd, UserProvidedRandomInput); err != nil {
				b.Error(err)
				return
			}
		}
	}, false)
}

func Benchmark_Sign(b *testing.B) {
	BA(b, as, func(b *testing.B, ps int) {
		data := []byte(`text`)

		var priv PrivateKey
		if err := GenerateParameters(&priv.Parameters, rnd, ParameterSizes(ps)); err != nil {
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
			r, _, err := Sign(rnd, &priv, data)
			if err != nil {
				b.Error(err)
			}
			data = r.Bytes()
		}
	}, false)
}

func Benchmark_Verify(b *testing.B) {
	BA(b, as, func(b *testing.B, ps int) {
		data := []byte(`text`)

		var priv PrivateKey
		if err := GenerateParameters(&priv.Parameters, rnd, ParameterSizes(ps)); err != nil {
			b.Error(err)
		}
		if err := GenerateKey(&priv, rnd); err != nil {
			b.Error(err)
		}

		r, s, err := Sign(rnd, &priv, data)
		if err != nil {
			b.Error(err)
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ok := Verify(&priv.PublicKey, data, r, s)
			if !ok {
				b.Errorf("%d: Verify failed", i)
			}
		}
	}, false)
}
