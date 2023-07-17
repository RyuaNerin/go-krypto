package kcdsa

import (
	"bytes"
	"crypto/dsa"
	"math/rand"
	"testing"
)

const (
	testDsaSize   = dsa.L2048N256
	testKcdsaSize = L2048N256SHA256
)

func Benchmark_DSA_GenerateParameters_L2048N256(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))

	var params dsa.Parameters
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dsa.GenerateParameters(&params, rnd, testDsaSize); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_KCDSA_GenerateParameters_GO_L2048N256(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))

	var params Parameters
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := GenerateParameters(&params, rnd, testKcdsaSize); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_KCDSA_GenerateParameters_KISA_L2048N256(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))

	var params Parameters
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := GenerateParametersKISA(&params, rnd, testKcdsaSize); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_DSA_GenerateKey(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))

	var priv dsa.PrivateKey
	if err := dsa.GenerateParameters(&priv.Parameters, rnd, testDsaSize); err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dsa.GenerateKey(&priv, rnd); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_KCDSA_GO_GenerateKey(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))

	var priv PrivateKey
	if err := GenerateParameters(&priv.Parameters, rnd, testKcdsaSize); err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := GenerateKey(&priv, rnd); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_KCDSA_KISA_GenerateKey(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))

	var priv PrivateKey
	if _, _, err := GenerateParametersKISA(&priv.Parameters, rnd, testKcdsaSize); err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := GenerateKeyKISA(&priv, rnd, UserProvidedRandomInput); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_DSA_Sign(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))
	data := []byte(`text`)

	var priv dsa.PrivateKey
	if err := dsa.GenerateParameters(&priv.Parameters, rnd, dsa.ParameterSizes(testKcdsaSize)); err != nil {
		b.Error(err)
	}
	if err := dsa.GenerateKey(&priv, rnd); err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := dsa.Sign(rnd, &priv, data); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_KCDSA_Sign(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))
	data := []byte(`text`)

	var priv PrivateKey
	if err := GenerateParameters(&priv.Parameters, rnd, ParameterSizes(testKcdsaSize)); err != nil {
		b.Error(err)
	}
	if err := GenerateKey(&priv, rnd); err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, _, err := Sign(rnd, &priv, bytes.NewReader(data)); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_DSA_Verify(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))
	data := []byte(`text`)

	var priv dsa.PrivateKey
	if err := dsa.GenerateParameters(&priv.Parameters, rnd, dsa.ParameterSizes(testKcdsaSize)); err != nil {
		b.Error(err)
	}
	if err := dsa.GenerateKey(&priv, rnd); err != nil {
		b.Error(err)
	}

	r, s, err := dsa.Sign(rnd, &priv, data)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if !dsa.Verify(&priv.PublicKey, data, r, s) {
			b.Errorf("%d: Verify failed", i)
		}
	}
}

func Benchmark_KCDSA_Verify(b *testing.B) {
	rnd := rand.New(rand.NewSource(0))
	data := []byte(`text`)

	var priv PrivateKey
	if err := GenerateParameters(&priv.Parameters, rnd, ParameterSizes(testKcdsaSize)); err != nil {
		b.Error(err)
	}
	if err := GenerateKey(&priv, rnd); err != nil {
		b.Error(err)
	}

	r, s, err := Sign(rnd, &priv, bytes.NewReader(data))
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ok, err := Verify(&priv.PublicKey, bytes.NewReader(data), r, s)
		if err != nil {
			b.Error(err)
		}
		if !ok {
			b.Errorf("%d: Verify failed", i)
		}
	}
}
