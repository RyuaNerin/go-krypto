package lea

import (
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

var (
	as = []CipherSize{
		{Name: "128", Size: 128},
		{Name: "196", Size: 196},
		{Name: "256", Size: 256},
	}
)

func Test_Encrypt_1Block_Src(t *testing.T) {
	BTSCA(t, as, 0, BlockSize, BIW(NewCipher), bb(leaEnc1Go), false)
}
func Test_Decrypt_1Block_Src(t *testing.T) {
	BTSCA(t, as, 0, BlockSize, BIW(NewCipher), bb(leaDec1Go), false)
}

func Test_Encrypt_4Block_Src(t *testing.T) {
	BTSCA(t, as, 0, 4*BlockSize, BIW(NewCipher), bb(leaEnc4Go), false)
}
func Test_Decrypt_4Block_Src(t *testing.T) {
	BTSCA(t, as, 0, 4*BlockSize, BIW(NewCipher), bb(leaDec4Go), false)
}

func Test_Encrypt_8Block_Src(t *testing.T) {
	BTSCA(t, as, 0, 8*BlockSize, BIW(NewCipher), bb(leaEnc8Go), false)
}
func Test_Decrypt_8Block_Src(t *testing.T) {
	BTSCA(t, as, 0, 8*BlockSize, BIW(NewCipher), bb(leaDec8Go), false)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New(b *testing.B) { BBNA(b, as, 0, BIW(NewCipher), false) }

func Benchmark_Encrypt_1Block(b *testing.B) {
	BBDA(b, as, 0, 1*BlockSize, BIW(newCipherGo), bb(leaEnc1Go), false)
}
func Benchmark_Encrypt_4Block(b *testing.B) {
	BBDA(b, as, 0, 4*BlockSize, BIW(newCipherGo), bb(leaEnc4Go), false)
}
func Benchmark_Encrypt_8Blocks(b *testing.B) {
	BBDA(b, as, 0, 8*BlockSize, BIW(newCipherGo), bb(leaEnc8Go), false)
}

func Bechmark_Decrypt_1Blocks(b *testing.B) {
	BBDA(b, as, 0, 1*BlockSize, BIW(newCipherGo), bb(leaDec1Go), false)
}
func Benchmark_Decrypt_4Blocks(b *testing.B) {
	BBDA(b, as, 0, 4*BlockSize, BIW(newCipherGo), bb(leaDec4Go), false)
}
func Benchmark_Decrypt_8Blocks(b *testing.B) {
	BBDA(b, as, 0, 8*BlockSize, BIW(newCipherGo), bb(leaDec8Go), false)
}

func bb(f funcBlock) func(c interface{}, dst, src []byte) {
	return func(c interface{}, dst, src []byte) {
		f(c.(*leaContext), dst, src)
	}
}
