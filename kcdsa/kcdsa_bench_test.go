package kcdsa_test

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/randutil"
	"github.com/RyuaNerin/go-krypto/kcdsa"

	"github.com/stretchr/testify/require"
)

func Benchmark_GenerateParameters(b *testing.B) {
	rnd, _ := randutil.NewBenchRand()
	for _, tc := range testConfigList {
		b.Run(tc.Name, func(b *testing.B) {
			var params kcdsa.Parameters

			b.ReportAllocs()
			for b.Loop() {
				_ = kcdsa.GenerateParameters(&params, rnd, tc.Sizes)
			}
		})
	}
}

func Benchmark_GenerateParametersTTAK(b *testing.B) {
	rnd, _ := randutil.NewBenchRand()
	for _, tc := range testConfigList {
		b.Run(tc.Name, func(b *testing.B) {
			var params kcdsa.Parameters

			b.ReportAllocs()
			for b.Loop() {
				_ = kcdsa.GenerateParametersTTAK(&params, rnd, tc.Sizes)
			}
		})
	}
}

func Benchmark_RegenerateParameters(b *testing.B) {
	rnd, _ := randutil.NewBenchRand()
	for _, tc := range testConfigList {
		b.Run(tc.Name, func(b *testing.B) {
			var params kcdsa.Parameters
			err := kcdsa.GenerateParametersTTAK(&params, rnd, tc.Sizes)
			require.NoError(b, err, "GenerateParametersTTAK failed")

			b.ReportAllocs()
			for b.Loop() {
				_ = kcdsa.RegenerateParameters(&params, rnd, tc.Sizes)
			}
		})
	}
}

func Benchmark_GenerateKey(b *testing.B) {
	rnd, _ := randutil.NewBenchRand()
	for _, tc := range testConfigList {
		b.Run(tc.Name, func(b *testing.B) {
			var priv kcdsa.PrivateKey
			err := kcdsa.GenerateParameters(&priv.Parameters, rnd, tc.Sizes)
			require.NoError(b, err, "GenerateParameters failed")

			b.ReportAllocs()
			for b.Loop() {
				_ = kcdsa.GenerateKey(&priv, rnd)
			}
		})
	}
}

func Benchmark_GenerateKeyTTAK(b *testing.B) {
	rnd, _ := randutil.NewBenchRand()
	for _, tc := range testConfigList {
		b.Run(tc.Name, func(b *testing.B) {
			var priv kcdsa.PrivateKey
			err := kcdsa.GenerateParameters(&priv.Parameters, rnd, tc.Sizes)
			require.NoError(b, err, "GenerateParameters failed")

			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = kcdsa.GenerateKeyWithSeed(&priv, rnd, nil, nil, tc.Sizes)
			}
		})
	}
}

func Benchmark_Sign(b *testing.B) {
	data := []byte(`text`)

	rnd, _ := randutil.NewBenchRand()

	for _, tc := range testConfigList {
		b.Run(tc.Name, func(b *testing.B) {
			var priv kcdsa.PrivateKey
			err := kcdsa.GenerateParameters(&priv.Parameters, rnd, tc.Sizes)
			require.NoError(b, err, "GenerateParameters failed")

			err = kcdsa.GenerateKey(&priv, rnd)
			require.NoError(b, err, "GenerateKey failed")

			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = kcdsa.Sign(rnd, &priv, tc.Sizes, data)
			}
		})
	}
}

func Benchmark_Verify(b *testing.B) {
	data := []byte(`text`)

	rnd, _ := randutil.NewBenchRand()

	for _, tc := range testConfigList {
		b.Run(tc.Name, func(b *testing.B) {
			var priv kcdsa.PrivateKey
			err := kcdsa.GenerateParameters(&priv.Parameters, rnd, tc.Sizes)
			require.NoError(b, err, "GenerateParameters failed")

			err = kcdsa.GenerateKey(&priv, rnd)
			require.NoError(b, err, "GenerateKey failed")

			r, s, err := kcdsa.Sign(rnd, &priv, tc.Sizes, data)
			require.NoError(b, err, "Sign failed")

			b.ReportAllocs()
			for b.Loop() {
				_ = kcdsa.Verify(&priv.PublicKey, tc.Sizes, data, r, s)
			}
		})
	}
}
