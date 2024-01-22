package kcdsattak

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/kcdsa"
	. "github.com/RyuaNerin/testingutil"
)

func Benchmark_GenerateParameters(b *testing.B) {
	BA(b, as, func(b *testing.B, ps int) {
		var params Parameters
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, _, err := GenerateParameters(&params, rnd, kcdsa.ParameterSizes(ps)); err != nil {
				b.Error(err)
				return
			}
		}
	}, false)
}

func Benchmark_GenerateKey(b *testing.B) {
	BA(b, as, func(b *testing.B, ps int) {
		var priv PrivateKey
		if _, _, err := GenerateParameters(&priv.Parameters, rnd, kcdsa.ParameterSizes(ps)); err != nil {
			b.Error(err)
			return
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if err := GenerateKey(&priv, rnd, UserProvidedRandomInput, kcdsa.ParameterSizes(ps)); err != nil {
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
		if _, _, err := GenerateParameters(&priv.Parameters, rnd, kcdsa.ParameterSizes(sz)); err != nil {
			b.Error(err)
			return
		}
		if err := GenerateKey(&priv, rnd, UserProvidedRandomInput, kcdsa.ParameterSizes(sz)); err != nil {
			b.Error(err)
			return
		}

		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			r, _, err := Sign(rnd, &priv, data, kcdsa.ParameterSizes(sz))
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
		if _, _, err := GenerateParameters(&priv.Parameters, rnd, kcdsa.ParameterSizes(sz)); err != nil {
			b.Error(err)
		}
		if err := GenerateKey(&priv, rnd, UserProvidedRandomInput, kcdsa.ParameterSizes(sz)); err != nil {
			b.Error(err)
		}

		r, s, err := Sign(rnd, &priv, data, kcdsa.ParameterSizes(sz))
		if err != nil {
			b.Error(err)
		}

		h := kcdsa.ParameterSizes(sz).Hash()

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ok := Verify(&priv.PublicKey, h, data, r, s)
			if !ok {
				b.Errorf("%d: Verify failed", i)
				return
			}
		}
	}, false)
}
