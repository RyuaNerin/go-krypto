package pbkdf2_test

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"runtime"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
	"github.com/RyuaNerin/go-krypto/lsh512"
	"github.com/RyuaNerin/go-krypto/pbkdf2"

	"github.com/stretchr/testify/require"
	golangpbkdf2 "golang.org/x/crypto/pbkdf2"
)

const (
	katSHA2256   = "20. PBKDF/PBKDF_HMAC_SHA2-256_KAT"
	katLSH512512 = "20. PBKDF/PBKDF_HMAC_LSH-512-512_KAT"
)

func Test_PBKDF_HMAC_SHA2256_KAT(t *testing.T)   { cavp.TestPbkdfKAT(t, katSHA2256, sha256.New) }
func Test_PBKDF_HMAC_LSH512512_KAT(t *testing.T) { cavp.TestPbkdfKAT(t, katLSH512512, lsh512.New) }

func Test_PBKDF_Parallel(t *testing.T) {
	if runtime.NumCPU() <= 1 {
		t.Skip("skipping parallel test on single CPU")
	}

	rnd, rand := randutil.NewTestRand()

	const keySize = 512
	const maxLen = 512

	var password [maxLen]byte
	var seed [maxLen]byte

	dstSingle := make([]byte, keySize)
	dstParallel := make([]byte, keySize)

	for range 100 {
		passwordLen := rand.IntN(maxLen-1) + 1
		seedLen := rand.IntN(maxLen-1) + 1

		_, err := rnd.Read(password[:passwordLen])
		require.NoError(t, err)
		_, err = rnd.Read(seed[:seedLen])
		require.NoError(t, err)

		dstSingle = pbkdf2.Generate(dstSingle, password[:passwordLen], seed[:seedLen], 1000, keySize, sha256.New)
		dstParallel = pbkdf2.GenerateParallel(dstParallel, password[:passwordLen], seed[:seedLen], 1000, keySize, sha256.New, 0)

		require.Equal(t, dstSingle, dstParallel, "parallel PBKDF2 output does not match single-threaded output")
	}
}

func Fuzz_PBKDF_HMAC_SHA2256(f *testing.F) {
	f.Add([]byte("password"), []byte("salt"), uint8(0), uint16(31), uint8(0))
	f.Add([]byte{}, []byte{}, uint8(1), uint16(63), uint8(3))
	f.Add([]byte{0x00, 0xff}, []byte{0xff, 0x00}, uint8(31), uint16(255), uint8(7))

	f.Fuzz(func(t *testing.T, password, salt []byte, iterationSeed uint8, keyLenSeed uint16, threadSeed uint8) {
		const maxInputLen = 1024
		password = password[:min(len(password), maxInputLen)]
		salt = salt[:min(len(salt), maxInputLen)]

		iteration := int(iterationSeed%32) + 1
		keyLen := int(keyLenSeed%512) + 1
		nThreads := int(threadSeed % 8)

		want := golangpbkdf2.Key(password, salt, iteration, keyLen, sha256.New)
		got := pbkdf2.Generate(nil, password, salt, iteration, keyLen, sha256.New)
		gotParallel := pbkdf2.GenerateParallel(nil, password, salt, iteration, keyLen, sha256.New, nThreads)

		require.Equal(t, want, got)
		require.Equal(t, want, gotParallel)
	})
}

func Benchmark_PBKDF_HMAC_SHA2256(b *testing.B) {
	password := []byte("password")
	salt := bytes.Repeat([]byte{0xA5}, 16)

	b.Run("SingleThreaded", func(b *testing.B) {
		for _, iter := range []int{1_000, 10_000, 100_000} {
			for _, klen := range []int{32, 64, 100} {
				b.Run(fmt.Sprintf("iter=%d/klen=%d", iter, klen), func(b *testing.B) {
					blocks := (klen + sha256.Size - 1) / sha256.Size
					b.SetBytes(int64(iter * blocks * sha256.Size))
					b.ReportAllocs()
					for b.Loop() {
						k := pbkdf2.Generate(nil, password, salt, iter, klen, sha256.New)
						_ = k
					}
				})
			}
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		for _, iter := range []int{1_000, 10_000, 100_000} {
			for _, klen := range []int{32, 64, 100} {
				b.Run(fmt.Sprintf("iter=%d/klen=%d", iter, klen), func(b *testing.B) {
					blocks := (klen + sha256.Size - 1) / sha256.Size
					b.SetBytes(int64(iter * blocks * sha256.Size))
					b.ReportAllocs()
					for b.Loop() {
						k := pbkdf2.GenerateParallel(nil, password, salt, iter, klen, sha256.New, 0)
						_ = k
					}
				})
			}
		}
	})

	b.Run("stdlib", func(b *testing.B) {
		for _, iter := range []int{1_000, 10_000, 100_000} {
			for _, klen := range []int{32, 64, 100} {
				b.Run(fmt.Sprintf("iter=%d/klen=%d", iter, klen), func(b *testing.B) {
					blocks := (klen + sha256.Size - 1) / sha256.Size
					b.SetBytes(int64(iter * blocks * sha256.Size))
					b.ReportAllocs()
					for b.Loop() {
						k := golangpbkdf2.Key(password, salt, iter, klen, sha256.New)
						_ = k
					}
				})
			}
		}
	})
}
