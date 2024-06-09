package lea

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

var as = []CipherSize{
	{Name: "128", Size: 128},
	{Name: "196", Size: 196},
	{Name: "256", Size: 256},
}

func Test_LEA128_Encrypt(t *testing.T) { BTE(t, BIW(NewCipher), CE, testCases128, false) }
func Test_LEA128_Decrypt(t *testing.T) { BTD(t, BIW(NewCipher), CD, testCases128, false) }

func Test_LEA196_Encrypt(t *testing.T) { BTE(t, BIW(NewCipher), CE, testCases196, false) }
func Test_LEA196_Decrypt(t *testing.T) { BTD(t, BIW(NewCipher), CD, testCases196, false) }

func Test_LEA256_Encrypt(t *testing.T) { BTE(t, BIW(NewCipher), CE, testCases256, false) }
func Test_LEA256_Decrypt(t *testing.T) { BTD(t, BIW(NewCipher), CD, testCases256, false) }

//nolint:unused
func tb(blocks int, funcGo, funcAsm funcBlock) func(t *testing.T, keySize int) {
	return func(t *testing.T, keySize int) {
		BTTC(
			t,
			keySize,
			0,
			blocks*BlockSize,
			0,
			func(key, additional []byte) (interface{}, error) {
				var ctx leaContext
				return &ctx, ctx.initContext(key)
			},
			func(ctx interface{}, dst, src []byte) { funcGo(ctx.(*leaContext), dst, src) },
			func(ctx interface{}, dst, src []byte) { funcAsm(ctx.(*leaContext), dst, src) },
			false,
		)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New(b *testing.B) { BBNA(b, as, 0, BIW(NewCipher), false) }

func Benchmark_Encrypt_1Block(b *testing.B) {
	BBDA(b, as, 0, BlockSize, BIW(NewCipher), bb(leaEnc1Go), false)
}

func Benchmark_Decrypt_1Blocks(b *testing.B) {
	BBDA(b, as, 0, BlockSize, BIW(NewCipher), bb(leaDec1Go), false)
}

func bb(f funcBlock) func(c interface{}, dst, src []byte) {
	return func(c interface{}, dst, src []byte) {
		f(c.(*leaContext), dst, src)
	}
}
