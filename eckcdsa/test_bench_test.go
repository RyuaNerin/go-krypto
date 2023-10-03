package eckcdsa

import (
	"crypto/elliptic"
	"crypto/sha256"
	"hash"
	"testing"
)

func Benchmark_GenerateKey(b *testing.B) {
	tests := []struct {
		name  string
		curve elliptic.Curve
	}{
		{"P224", elliptic.P224()},
		{"P256", elliptic.P256()},
	}
	for _, test := range tests {
		test := test
		b.Run(test.name, func(b *testing.B) {

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := GenerateKey(test.curve, rnd); err != nil {
					b.Error(err)
				}
			}
		})
	}
}

func benchmarkAllSizes(b *testing.B, f func(*testing.B, elliptic.Curve, hash.Hash)) {
	tests := []struct {
		name  string
		curve elliptic.Curve
		h     hash.Hash
	}{
		{"P224_SHA224", elliptic.P224(), sha256.New224()},
		{"P224_SHA256", elliptic.P224(), sha256.New()},
		{"P256_SHA224", elliptic.P256(), sha256.New224()},
		{"P256_SHA256", elliptic.P256(), sha256.New()},
	}
	for _, test := range tests {
		test := test
		b.Run(test.name, func(b *testing.B) {
			f(b, test.curve, test.h)
		})
	}
}

func Benchmark_Sign(b *testing.B) {
	benchmarkAllSizes(b, func(b *testing.B, c elliptic.Curve, h hash.Hash) {
		data := []byte(`text`)

		key, err := GenerateKey(c, rnd)
		if err != nil {
			b.Error(err)
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r, _, err := Sign(rnd, key, h, data)
			if err != nil {
				b.Error(err)
			}
			data = r.Bytes()
		}
	})
}

func Benchmark_Verify(b *testing.B) {
	benchmarkAllSizes(b, func(b *testing.B, c elliptic.Curve, h hash.Hash) {
		data := []byte(`text`)

		key, err := GenerateKey(c, rnd)
		if err != nil {
			b.Error(err)
		}
		r, s, err := Sign(rnd, key, h, data)
		if err != nil {
			b.Error(err)
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			ok := Verify(&key.PublicKey, h, data, r, s)
			if !ok {
				b.Errorf("%d: Verify failed", i)
			}
		}
	})
}
