package kcdsa

import (
	"bytes"
	"crypto/dsa"
	"crypto/rand"
	"testing"
)

const (
	testDsaSize   = dsa.L2048N256
	testKcdsaSize = L2048N256WithSHA256
)

func Benchmark_DSA_GenerateParameters_L2048N256(b *testing.B) {
	var params dsa.Parameters
	for i := 0; i < b.N; i++ {
		if err := dsa.GenerateParameters(&params, rand.Reader, testDsaSize); err != nil {
			b.Error(err)
		}
	}
}
func Benchmark_KCDSA_GenerateParameters_L2048N256(b *testing.B) {
	var params Parameters
	for i := 0; i < b.N; i++ {
		if err := GenerateParameters(&params, rand.Reader, testKcdsaSize); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_DSA_GenerateKey(b *testing.B) {
	var priv dsa.PrivateKey
	if err := dsa.GenerateParameters(&priv.Parameters, rand.Reader, testDsaSize); err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		if err := dsa.GenerateKey(&priv, rand.Reader); err != nil {
			b.Error(err)
		}
	}
}
func Benchmark_KCDSA_GenerateKey(b *testing.B) {
	var priv PrivateKey
	if err := GenerateParameters(&priv.Parameters, rand.Reader, testKcdsaSize); err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		if err := GenerateKey(&priv, rand.Reader); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_DSA_Sign(b *testing.B) {
	data := []byte(`text`)

	var priv dsa.PrivateKey
	if err := dsa.GenerateParameters(&priv.Parameters, rand.Reader, dsa.ParameterSizes(testKcdsaSize)); err != nil {
		b.Error(err)
	}
	if err := dsa.GenerateKey(&priv, rand.Reader); err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		if _, _, err := dsa.Sign(rand.Reader, &priv, data); err != nil {
			b.Error(err)
		}
	}
}
func Benchmark_KCDSA_Sign(b *testing.B) {
	data := []byte(`text`)

	var priv PrivateKey
	if err := GenerateParameters(&priv.Parameters, rand.Reader, ParameterSizes(testKcdsaSize)); err != nil {
		b.Error(err)
	}
	if err := GenerateKey(&priv, rand.Reader); err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		if _, _, err := Sign(rand.Reader, &priv, bytes.NewReader(data)); err != nil {
			b.Error(err)
		}
	}
}

func Benchmark_DSA_Verify(b *testing.B) {
	data := []byte(`text`)

	var priv dsa.PrivateKey
	if err := dsa.GenerateParameters(&priv.Parameters, rand.Reader, dsa.ParameterSizes(testKcdsaSize)); err != nil {
		b.Error(err)
	}
	if err := dsa.GenerateKey(&priv, rand.Reader); err != nil {
		b.Error(err)
	}

	r, s, err := dsa.Sign(rand.Reader, &priv, data)
	if err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		if !dsa.Verify(&priv.PublicKey, data, r, s) {
			b.Errorf("%d: Verify failed", i)
		}
	}
}
func Benchmark_KCDSA_Verify(b *testing.B) {
	data := []byte(`text`)

	var priv PrivateKey
	if err := GenerateParameters(&priv.Parameters, rand.Reader, ParameterSizes(testKcdsaSize)); err != nil {
		b.Error(err)
	}
	if err := GenerateKey(&priv, rand.Reader); err != nil {
		b.Error(err)
	}

	r, s, err := Sign(rand.Reader, &priv, bytes.NewReader(data))
	if err != nil {
		b.Error(err)
	}
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
