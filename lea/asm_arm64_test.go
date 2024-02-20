//go:build arm64 && !purego
// +build arm64,!purego

package lea

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

var (
	leaEnc4NEON = toAsmFunc(__lea_encrypt_4block)
	leaDec4NEON = toAsmFunc(__lea_decrypt_4block)
)

func Test_Encrypt_4Blocks_NEON(t *testing.T) { TA(t, as, tb(4, leaEnc4Go, leaEnc4NEON), false) }
func Test_Decrypt_4Blocks_NEON(t *testing.T) { TA(t, as, tb(4, leaDec4Go, leaDec4NEON), false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_Encrypt_4Blocks_NEON(b *testing.B) {
	BBDA(b, as, 0, 4*BlockSize, BIW(NewCipher), bb(leaEnc4NEON), false)
}
func Benchmark_Decrypt_4Blocks_NEON(b *testing.B) {
	BBDA(b, as, 0, 4*BlockSize, BIW(NewCipher), bb(leaDec4NEON), false)
}
