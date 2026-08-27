package hmac_test

import (
	"crypto/sha256"
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/hmac"
	"github.com/RyuaNerin/go-krypto/internal/cavp"
	"github.com/RyuaNerin/go-krypto/internal/testing/thash"
	"github.com/RyuaNerin/go-krypto/lsh256"
	"github.com/RyuaNerin/go-krypto/lsh512"

	"github.com/stretchr/testify/require"
)

var htSHA256 = thash.HashTest{
	HashInfo: thash.HashInfo{
		BlockSize: sha256.BlockSize,
		Sizes: []thash.CipherSize{
			{Name: "TagSize=32", DigestSize: sha256.Size}, // DigestSize = TagSize
		},
	},
	New: map[int]func() hash.Hash{
		sha256.Size: func() hash.Hash {
			h := hmac.New(sha256.New, make([]byte, 16))
			return h
		},
	},
}

const (
	katHMACSHA256    = "08. HMAC/HMAC_SHA-256_KAT"
	katHMACLSH256256 = "08. HMAC/HMAC_LSH-256-256_KAT"
	katHMACLSH512512 = "08. HMAC/HMAC_LSH-512-512_KAT"
)

func Test_HMAC_SHA256_KAT(t *testing.T)    { cavp.TestHMAC(t, katHMACSHA256, sha256.New) }
func Test_HMAC_LSH256256_KAT(t *testing.T) { cavp.TestHMAC(t, katHMACLSH256256, lsh256.New) }

// LSH-512 has a 256-byte block, the largest block size HMAC supports.
func Test_HMAC_LSH512512_KAT(t *testing.T) { cavp.TestHMAC(t, katHMACLSH512512, lsh512.New) }

func Test_HMAC_KeyDependence(t *testing.T) {
	for _, newHash := range []func() hash.Hash{sha256.New, lsh256.New, lsh512.New} {
		a := hmac.New(newHash, []byte("key-A"))
		b := hmac.New(newHash, []byte("key-B"))
		a.Write([]byte("message"))
		b.Write([]byte("message"))
		require.NotEqual(t, a.Sum(nil), b.Sum(nil), "block size %d", a.BlockSize())
	}
}

func Test_HMAC_TagSizeLargerThanHash(t *testing.T) {
	require.Panics(t, func() { hmac.NewWithSize(sha256.New, nil, sha256.Size+1) })
}

func Test_HMAC_SHA256_New(t *testing.T)  { htSHA256.TestNew(t, false) }
func Test_HMAC_SHA256_Hash(t *testing.T) { htSHA256.TestHash(t, false) }

func Fuzz_HMAC_SHA256_New(f *testing.F)  { htSHA256.FuzzNew(f, false) }
func Fuzz_HMAC_SHA256_Hash(f *testing.F) { htSHA256.FuzzHash(f, false) }

func Benchmark_HMAC_SHA256_New(b *testing.B)  { htSHA256.BenchNew(b, false) }
func Benchmark_HMAC_SHA256_Hash(b *testing.B) { htSHA256.BenchHash(b, false) }
