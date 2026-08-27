package gmac_test

import (
	"crypto/aes"
	"fmt"
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/gmac"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/testing/thash"
	"github.com/RyuaNerin/go-krypto/lea"

	"github.com/stretchr/testify/require"
)

func newAES128GMAC(tagSize int) func() hash.Hash {
	return func() hash.Hash {
		c, _ := aes.NewCipher(make([]byte, 16))
		h, _ := gmac.NewGMACWithSize(c, make([]byte, 16), tagSize)
		return h
	}
}

func newLEA128GMAC(tagSize int) func() hash.Hash {
	return func() hash.Hash {
		c, _ := lea.NewCipher(make([]byte, 16))
		h, _ := gmac.NewGMACWithSize(c, make([]byte, 16), tagSize)
		return h
	}
}

var (
	htAES256 = thash.HashTest{
		HashInfo: thash.HashInfo{
			BlockSize: aes.BlockSize,
			Sizes: []thash.CipherSize{
				{Name: "TagSize=16", DigestSize: 16}, // DigestSize = TagSize
				{Name: "TagSize=12", DigestSize: 12}, // DigestSize = TagSize
				{Name: "TagSize=8", DigestSize: 8},   // DigestSize = TagSize
				{Name: "TagSize=4", DigestSize: 4},   // DigestSize = TagSize
			},
			SkipUnalignedTest: true,
		},
		New: map[int]func() hash.Hash{
			16: newAES128GMAC(16),
			12: newAES128GMAC(12),
			8:  newAES128GMAC(8),
			4:  newAES128GMAC(4),
		},
	}
	htLEA128 = thash.HashTest{
		HashInfo: thash.HashInfo{
			BlockSize: lea.BlockSize,
			Sizes: []thash.CipherSize{
				{Name: "TagSize=16", DigestSize: 16}, // DigestSize = TagSize
				{Name: "TagSize=12", DigestSize: 12}, // DigestSize = TagSize
				{Name: "TagSize=8", DigestSize: 8},   // DigestSize = TagSize
				{Name: "TagSize=4", DigestSize: 4},   // DigestSize = TagSize
			},
			SkipUnalignedTest: true,
		},
		New: map[int]func() hash.Hash{
			16: newLEA128GMAC(16),
			12: newLEA128GMAC(12),
			8:  newLEA128GMAC(8),
			4:  newLEA128GMAC(4),
		},
	}
)

func Test_GMAC(t *testing.T) {
	var (
		/**
		https://cryptopp.com/wiki/GMAC

		$ ./test.exe
		Key: 00000000000000000000000000000000
		IV:  00000000000000000000000000000000
		Message: Yoda said, do or do not. There is no try.
		GMAC: E7EE2C63B4DC328EED4A86B3FB3490AF
		*/
		key    = internal.HB(`00000000000000000000000000000000`)
		iv     = internal.HB(`00000000000000000000000000000000`)
		msg    = []byte(`Yoda said, do or do not. There is no try.`)
		expect = internal.HB(`E7EE2C63B4DC328EED4A86B3FB3490AF`)
	)

	for tagSize := 16; tagSize >= 4; tagSize -= 2 {
		t.Run(fmt.Sprintf("tagSize=%d", tagSize), func(t *testing.T) {
			expect := expect[:tagSize]

			c, err := aes.NewCipher(key)
			require.NoError(t, err, "failed to create AES cipher")

			h, err := gmac.NewGMACWithSize(c, iv, tagSize)
			require.NoError(t, err, "failed to create GMAC")
			require.NotNil(t, h, "failed to create GMAC")

			for idx := 0; idx < len(msg); {
				n, err := h.Write(msg[idx:])
				require.NoError(t, err, "failed to write GMAC")
				idx += n
			}
			actual := h.Sum(nil)
			require.Equal(t, expect, actual, "GMAC mismatch")
			actual = h.Sum(nil)
			require.Equal(t, expect, actual, "GMAC mismatch after Sum()")

			// reset test
			h.Reset()
			for idx := 0; idx < len(msg); {
				n, err := h.Write(msg[idx:])
				require.NoError(t, err, "failed to write GMAC")
				idx += n
			}
			actual = h.Sum(actual[:0])
			require.Equal(t, expect, actual, "GMAC mismatch after Reset()")
		})
	}
}

func Test_GMAC_AES128_New(t *testing.T)  { htAES256.TestNew(t, false) }
func Test_GMAC_AES128_Hash(t *testing.T) { htAES256.TestHash(t, false) }

func Test_GMAC_LEA128_New(t *testing.T)  { htLEA128.TestNew(t, false) }
func Test_GMAC_LEA128_Hash(t *testing.T) { htLEA128.TestHash(t, false) }

func Fuzz_GMAC_AES128_New(f *testing.F)  { htAES256.FuzzNew(f, false) }
func Fuzz_GMAC_AES128_Hash(f *testing.F) { htAES256.FuzzHash(f, false) }

func Fuzz_GMAC_LEA128_New(f *testing.F)  { htLEA128.FuzzNew(f, false) }
func Fuzz_GMAC_LEA128_Hash(f *testing.F) { htLEA128.FuzzHash(f, false) }

func Benchmark_GMAC_AES128_New(b *testing.B)  { htAES256.BenchNew(b, false) }
func Benchmark_GMAC_AES128_Hash(b *testing.B) { htAES256.BenchHash(b, false) }

func Benchmark_GMAC_LEA128_New(b *testing.B)  { htLEA128.BenchNew(b, false) }
func Benchmark_GMAC_LEA128_Hash(b *testing.B) { htLEA128.BenchHash(b, false) }
