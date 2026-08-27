package krypto_test

import (
	"testing"

	"github.com/RyuaNerin/go-krypto"
	_ "github.com/RyuaNerin/go-krypto/has160"
	_ "github.com/RyuaNerin/go-krypto/lsh256"
	_ "github.com/RyuaNerin/go-krypto/lsh512"

	"github.com/stretchr/testify/require"
)

func TestLSHHashRegistration(t *testing.T) {
	tests := []struct {
		hash      krypto.KryptoHash
		size      int
		blockSize int
	}{
		{hash: krypto.HAS160, size: 20, blockSize: 64},
		{hash: krypto.LSH256_224, size: 28, blockSize: 128},
		{hash: krypto.LSH256, size: 32, blockSize: 128},
		{hash: krypto.LSH512_224, size: 28, blockSize: 256},
		{hash: krypto.LSH512_256, size: 32, blockSize: 256},
		{hash: krypto.LSH512_384, size: 48, blockSize: 256},
		{hash: krypto.LSH512, size: 64, blockSize: 256},
	}

	for _, tt := range tests {
		require.Equal(t, tt.size, tt.hash.Size())
		require.Equal(t, tt.blockSize, tt.hash.BlockSize())

		h := tt.hash.New()
		require.Equal(t, tt.size, h.Size())
		require.Equal(t, tt.blockSize, h.BlockSize())
	}
}
