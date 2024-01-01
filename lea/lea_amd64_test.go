//go:build amd64 && !purego

package lea

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

var (
	leaEnc4SSE2 = __lea_encrypt_4block
	leaDec4SSE2 = __lea_decrypt_4block

	leaEnc8AVX2 = __lea_encrypt_8block
	leaDec8AVX2 = __lea_decrypt_8block
)

func Test_Encrypt_4Blocks_SSE2(t *testing.T) { TA(t, as, tb(4, leaEnc4Go, leaEnc4SSE2), false) }
func Test_Decrypt_4Blocks_SSE2(t *testing.T) { TA(t, as, tb(4, leaDec4Go, leaDec4SSE2), false) }

func Test_Encrypt_8Blocks_AVX2(t *testing.T) { TA(t, as, tb(8, leaEnc8Go, leaEnc8AVX2), !hasAVX2) }
func Test_Decrypt_8Blocks_AVX2(t *testing.T) { TA(t, as, tb(8, leaDec8Go, leaDec8AVX2), !hasAVX2) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Encrypt_4Blocks_SSE2(b *testing.B) {
	BBDA(b, as, 0, 4*BlockSize, BIW(NewCipher), bb(leaEnc4SSE2), false)
}
func Benchmark_Decrypt_4Blocks_SSE2(b *testing.B) {
	BBDA(b, as, 0, 4*BlockSize, BIW(NewCipher), bb(leaDec4SSE2), false)
}

func Benchmark_Encrypt_8Blocks_AVX2(b *testing.B) {
	BBDA(b, as, 0, 8*BlockSize, BIW(NewCipher), bb(leaEnc8AVX2), !hasAVX2)
}
func Benchmark_Decrypt_8Blocks_AVX2(b *testing.B) {
	BBDA(b, as, 0, 8*BlockSize, BIW(NewCipher), bb(leaDec8AVX2), !hasAVX2)
}
