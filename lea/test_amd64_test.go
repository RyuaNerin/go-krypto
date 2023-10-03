//go:build amd64 && gc && !purego

package lea

import (
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

func Test_Encrypt_4Blocks_SSE2(t *testing.T) { TA(t, as, tb(4, leaEnc4Go, leaEnc4, true)) }
func Test_Decrypt_4Blocks_SSE2(t *testing.T) { TA(t, as, tb(4, leaDec4Go, leaDec4, true)) }

func Test_Encrypt_8Blocks_SSE2(t *testing.T) { TA(t, as, tb(8, leaEnc8Go, leaEnc8SSE2, true)) }
func Test_Decrypt_8Blocks_SSE2(t *testing.T) { TA(t, as, tb(8, leaDec8Go, leaDec8SSE2, true)) }

func Test_Encrypt_8Blocks_AVX2(t *testing.T) { TA(t, as, tb(8, leaEnc8Go, leaEnc8AVX2, hasAVX2)) }
func Test_Decrypt_8Blocks_AVX2(t *testing.T) { TA(t, as, tb(8, leaDec8Go, leaDec8AVX2, hasAVX2)) }

func tb(blocks int, funcGo, funcAsm funcBlock, do bool) func(t *testing.T, keySize int) {
	return func(t *testing.T, keySize int) {
		if !do {
			t.Skip()
			return
		}

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
		)
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Encrypt_4Blocks_SSE2(b *testing.B) { BA(b, as, bb(b, 4, leaEnc4SSE2, true)) }
func Benchmark_Encrypt_8Blocks_SSE2(b *testing.B) { BA(b, as, bb(b, 8, leaEnc8SSE2, true)) }
func Benchmark_Encrypt_8Blocks_AVX2(b *testing.B) { BA(b, as, bb(b, 8, leaEnc8AVX2, hasAVX2)) }

func Benchmark_Decrypt_4Blocks_SSE2(b *testing.B) { BA(b, as, bb(b, 4, leaDec4SSE2, true)) }
func Benchmark_Decrypt_8Blocks_SSE2(b *testing.B) { BA(b, as, bb(b, 8, leaDec8SSE2, true)) }
func Benchmark_Decrypt_8Blocks_AVX2(b *testing.B) { BA(b, as, bb(b, 8, leaDec8AVX2, hasAVX2)) }
