//go:build amd64 && gc && !purego

package lea

import (
	"testing"

	. "github.com/RyuaNerin/go-krypto/internal/testingutil"
)

func Test_Encrypt_4Block_Src_SSE2(t *testing.T) {
	BTSCA(t, as, 0, 4*BlockSize, BIW(NewCipher), bb(leaEnc4SSE2), false)
}
func Test_Decrypt_4Block_Src_SSE2(t *testing.T) {
	BTSCA(t, as, 0, 4*BlockSize, BIW(NewCipher), bb(leaDec4SSE2), false)
}

func Test_Encrypt_8Block_Src_SSE2(t *testing.T) {
	BTSCA(t, as, 0, 8*BlockSize, BIW(NewCipher), bb(leaEnc8SSE2), false)
}
func Test_Decrypt_8Block_Src_SSE2(t *testing.T) {
	BTSCA(t, as, 0, 8*BlockSize, BIW(NewCipher), bb(leaDec8SSE2), false)
}

func Test_Encrypt_8Block_Src_AVX2(t *testing.T) {
	BTSCA(t, as, 0, 8*BlockSize, BIW(NewCipher), bb(leaEnc8AVX2), !hasAVX2)
}
func Test_Decrypt_8Block_Src_AVX2(t *testing.T) {
	BTSCA(t, as, 0, 8*BlockSize, BIW(NewCipher), bb(leaDec8AVX2), !hasAVX2)
}

func Test_Encrypt_4Blocks_SSE2(t *testing.T) { TA(t, as, tb(4, leaEnc4Go, leaEnc4SSE2), false) }
func Test_Decrypt_4Blocks_SSE2(t *testing.T) { TA(t, as, tb(4, leaDec4Go, leaDec4SSE2), false) }

func Test_Encrypt_8Blocks_SSE2(t *testing.T) { TA(t, as, tb(8, leaEnc8Go, leaEnc8SSE2), false) }
func Test_Decrypt_8Blocks_SSE2(t *testing.T) { TA(t, as, tb(8, leaDec8Go, leaDec8SSE2), false) }

func Test_Encrypt_8Blocks_AVX2(t *testing.T) { TA(t, as, tb(8, leaEnc8Go, leaEnc8AVX2), !hasAVX2) }
func Test_Decrypt_8Blocks_AVX2(t *testing.T) { TA(t, as, tb(8, leaDec8Go, leaDec8AVX2), !hasAVX2) }

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

func Benchmark_Encrypt_4Blocks_SSE2(b *testing.B) {
	BBDA(b, as, 0, 4*BlockSize, BIW(newCipherGo), bb(leaEnc4SSE2), false)
}
func Benchmark_Encrypt_8Blocks_SSE2(b *testing.B) {
	BBDA(b, as, 0, 8*BlockSize, BIW(newCipherGo), bb(leaEnc8SSE2), false)
}
func Benchmark_Encrypt_8Blocks_AVX2(b *testing.B) {
	BBDA(b, as, 0, 8*BlockSize, BIW(newCipherGo), bb(leaEnc8AVX2), !hasAVX2)
}

func Benchmark_Decrypt_4Blocks_SSE2(b *testing.B) {
	BBDA(b, as, 0, 4*BlockSize, BIW(newCipherGo), bb(leaDec4SSE2), false)
}
func Benchmark_Decrypt_8Blocks_SSE2(b *testing.B) {
	BBDA(b, as, 0, 8*BlockSize, BIW(newCipherGo), bb(leaDec8SSE2), false)
}
func Benchmark_Decrypt_8Blocks_AVX2(b *testing.B) {
	BBDA(b, as, 0, 8*BlockSize, BIW(newCipherGo), bb(leaDec8AVX2), !hasAVX2)
}
