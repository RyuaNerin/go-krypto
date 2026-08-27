package eckcdsa_test

import (
	"crypto/elliptic"
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/internal/randutil"

	"github.com/stretchr/testify/require"
)

func Benchmark_GenerateKey(b *testing.B) {
	rnd, _ := randutil.NewBenchRand()

	for _, c := range curves {
		b.Run(c.Params().Name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, err := eckcdsa.GenerateKey(c, rnd)
				require.NoError(b, err)
			}
		})
	}
}

func benchmarkAllSizes(b *testing.B, f func(*testing.B, elliptic.Curve, hash.Hash)) {
	for _, test := range curveAndHashPairs {
		b.Run(test.name, func(b *testing.B) {
			f(b, test.curve, test.hash)
		})
	}
}

func Benchmark_Sign(b *testing.B) {
	data := []byte(`text`)

	benchmarkAllSizes(b, func(b *testing.B, c elliptic.Curve, h hash.Hash) {
		rnd, _ := randutil.NewBenchRand()

		key, err := eckcdsa.GenerateKey(c, rnd)
		require.NoError(b, err)

		b.ReportAllocs()
		for b.Loop() {
			r, _, err := eckcdsa.Sign(rnd, key, h, data)
			require.NoError(b, err)
			data = r.Bytes()
		}
	})
}

func Benchmark_Verify(b *testing.B) {
	data := []byte(`text`)

	benchmarkAllSizes(b, func(b *testing.B, c elliptic.Curve, h hash.Hash) {
		rnd, _ := randutil.NewBenchRand()

		key, err := eckcdsa.GenerateKey(c, rnd)
		require.NoError(b, err)

		r, s, err := eckcdsa.Sign(rnd, key, h, data)
		require.NoError(b, err)

		b.ReportAllocs()
		for b.Loop() {
			ok := eckcdsa.Verify(&key.PublicKey, h, data, r, s)
			require.True(b, ok)
		}
	})
}
