package cmac_test

import (
	"hash"
	"testing"

	"github.com/RyuaNerin/go-krypto/cmac"
	"github.com/RyuaNerin/go-krypto/hight"
	"github.com/RyuaNerin/go-krypto/internal/cavp"
	"github.com/RyuaNerin/go-krypto/internal/testing/thash"
	"github.com/RyuaNerin/go-krypto/lea"

	"github.com/stretchr/testify/require"
)

var (
	htLEA = thash.HashTest{
		HashInfo: thash.HashInfo{
			BlockSize: lea.BlockSize,
			Sizes: []thash.CipherSize{
				{Name: "TagSize=16", DigestSize: lea.BlockSize}, // DigestSize = TagSize
			},
			SkipUnalignedTest: true,
		},
		New: map[int]func() hash.Hash{
			lea.BlockSize: func() hash.Hash {
				c, _ := lea.NewCipher(make([]byte, 16))
				h := cmac.New(c)
				return h
			},
		},
	}
	htHIGHT = thash.HashTest{
		HashInfo: thash.HashInfo{
			BlockSize: hight.BlockSize,
			Sizes: []thash.CipherSize{
				{Name: "TagSize=8", DigestSize: hight.BlockSize}, // DigestSize = TagSize
			},
			SkipUnalignedTest: true,
		},
		New: map[int]func() hash.Hash{
			hight.BlockSize: func() hash.Hash {
				c, _ := hight.NewCipher(make([]byte, 16))
				h := cmac.New(c)
				return h
			},
		},
	}
)

const (
	genLEA128 = "09. CMAC and GMAC/CMAC/CMAC_LEA-128_Gen"
	verLEA128 = "09. CMAC and GMAC/CMAC/CMAC_LEA-128_Ver"
	genHIGHT  = "09. CMAC and GMAC/CMAC/CMAC_HIGHT_Gen"
	verHIGHT  = "09. CMAC and GMAC/CMAC/CMAC_HIGHT_Ver"
)

func Test_CMAC_LEA128_Gen(t *testing.T) { cavp.TestCMACGen(t, genLEA128, lea.NewCipher) }
func Test_CMAC_LEA128_Ver(t *testing.T) { cavp.TestCMACVer(t, verLEA128, lea.NewCipher) }

func Test_CMAC_HIGHT_Gen(t *testing.T) { cavp.TestCMACGen(t, genHIGHT, hight.NewCipher) }
func Test_CMAC_HIGHT_Ver(t *testing.T) { cavp.TestCMACVer(t, verHIGHT, hight.NewCipher) }

func Test_CMAC_CustomTagSize(t *testing.T) {
	block, err := lea.NewCipher(make([]byte, 16))
	require.NoError(t, err)

	for _, size := range []int{1, 4, 8, 12, 16} {
		var h hash.Hash
		require.NotPanics(t, func() { h = cmac.NewWithSize(block, size) })
		require.NoErrorf(t, err, "size %d", size)

		require.Equal(t, size, h.Size(), "size %d: unexpected Size()", size)
		require.Len(t, h.Sum(nil), size, "size %d: unexpected Sum() length", size)
	}
}

func Test_CMAC_LEA128_New(t *testing.T)  { htLEA.TestNew(t, false) }
func Test_CMAC_LEA128_Hash(t *testing.T) { htLEA.TestHash(t, false) }

func Test_CMAC_HIGHT_New(t *testing.T)  { htHIGHT.TestNew(t, false) }
func Test_CMAC_HIGHT_Hash(t *testing.T) { htHIGHT.TestHash(t, false) }

func Fuzz_CMAC_LEA128_New(f *testing.F)  { htLEA.FuzzNew(f, false) }
func Fuzz_CMAC_LEA128_Hash(f *testing.F) { htLEA.FuzzHash(f, false) }

func Fuzz_CMAC_HIGHT_New(f *testing.F)  { htHIGHT.FuzzNew(f, false) }
func Fuzz_CMAC_HIGHT_Hash(f *testing.F) { htHIGHT.FuzzHash(f, false) }

func Benchmark_CMAC_LEA128_New(b *testing.B)  { htLEA.BenchNew(b, false) }
func Benchmark_CMAC_LEA128_Hash(b *testing.B) { htLEA.BenchHash(b, false) }

func Benchmark_CMAC_HIGHT_New(b *testing.B)  { htHIGHT.BenchNew(b, false) }
func Benchmark_CMAC_HIGHT_Hash(b *testing.B) { htHIGHT.BenchHash(b, false) }
