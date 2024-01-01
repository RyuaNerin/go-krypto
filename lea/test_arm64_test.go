//go:build arm64 && gc && !purego

package lea

import (
	"testing"
)

func Test_Encrypt_4Blocks_NEON(t *testing.T) { testAll(t, tb(4, leaEnc4Go, leaEnc4NEON, false)) }
func Test_Decrypt_4Blocks_NEON(t *testing.T) { testAll(t, tb(4, leaDec4Go, leaDec4NEON, false)) }

func Benchmark_Encrypt_4Blocks_NEON(b *testing.B) { benchAll(b, bb(4, leaEnc4NEON, false)) }
func Benchmark_Decrypt_4Blocks_NEON(b *testing.B) { benchAll(b, bb(4, leaDec4NEON, false)) }
