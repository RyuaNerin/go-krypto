package pbkdf2_test

import (
	"crypto/sha256"
	"hash"
	"strconv"
	"sync"
	"testing"

	"github.com/RyuaNerin/go-krypto/pbkdf2"

	"github.com/stretchr/testify/require"
)

type bufferTrackingHash struct {
	hash.Hash
	record func([]byte)
}

func (h *bufferTrackingHash) Sum(dst []byte) []byte {
	result := h.Hash.Sum(dst)
	if cap(dst) == 2*h.Size() {
		h.record(dst[:cap(dst)])
	}
	return result
}

func TestZeroizeDerivationBuffers(t *testing.T) {
	for _, threads := range []int{1, 4} {
		t.Run(strconv.Itoa(threads), func(t *testing.T) {
			var mu sync.Mutex
			var buffers [][]byte
			newHash := func() hash.Hash {
				return &bufferTrackingHash{Hash: sha256.New(), record: func(b []byte) { mu.Lock(); buffers = append(buffers, b); mu.Unlock() }}
			}
			got := pbkdf2.GenerateParallel(nil, []byte("password"), []byte("salt"), 3, 128, newHash, threads)
			want := pbkdf2.Generate(nil, []byte("password"), []byte("salt"), 3, 128, sha256.New)

			require.Equal(t, got, want, "parallel and sequential derivation outputs differ")

			for _, buf := range buffers {
				require.Equal(t, buf, make([]byte, len(buf)), "U/T derivation buffer was not fully zeroized")
			}
		})
	}
}
