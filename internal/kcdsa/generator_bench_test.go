package kcdsa

import (
	"io"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/randutil"

	"github.com/stretchr/testify/require"
)

func BenchmarkGenerateParametersTTAK(b *testing.B) {
	rnd, _ := randutil.NewBenchRand()

	domain, ok := GetDomain(A2048B224SHA224)
	require.True(b, ok, "KCDSA domain not found")

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, err := GenerateParametersTTAK(rnd, domain)
		require.NoError(b, err, "GenerateParametersTTAK failed")
	}
}

func BenchmarkGenerateJ(b *testing.B) {
	domain, ok := GetDomain(A2048B224SHA224)
	require.True(b, ok, "KCDSA domain not found")

	rnd, _ := randutil.NewBenchRand()

	j := new(big.Int)

	seed := make([]byte, internal.BitsToBytes(domain.B))
	h := domain.NewHash()

	var buf []byte

	b.ReportAllocs()
	for b.Loop() {
		_, err := io.ReadFull(rnd, seed)
		require.NoError(b, err, "failed to read seed")
		buf, _ = GenerateJ(j, buf, seed, h, domain)
	}
}

func BenchmarkGeneratePQ(b *testing.B) {
	domain, ok := GetDomain(A2048B224SHA224)
	require.True(b, ok, "KCDSA domain not found")

	rnd, _ := randutil.NewBenchRand()

	parameters, err := GenerateParametersTTAK(rnd, domain)
	require.NoError(b, err, "GenerateParametersTTAK failed")

	p, q := new(big.Int), new(big.Int)

	var buf []byte
	b.ReportAllocs()
	for b.Loop() {
		buf, _, ok = GeneratePQ(p, q, buf, parameters.J, parameters.Seed, domain.NewHash(), domain)
		require.True(b, ok, "GeneratePQ failed")
	}
}
