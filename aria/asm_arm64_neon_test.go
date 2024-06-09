//go:build arm64 && !purego && (!gccgo || go1.18)
// +build arm64
// +build !purego
// +build !gccgo go1.18

package aria

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_ARIA128_Encrypt_NEON(t *testing.T) { BTE(t, BIW(newCipherAsm), CE, testCases128, !hasNEON) }
func Test_ARIA128_Decrypt_NEON(t *testing.T) { BTD(t, BIW(newCipherAsm), CD, testCases128, !hasNEON) }

func Test_ARIA196_Encrypt_NEON(t *testing.T) { BTE(t, BIW(newCipherAsm), CE, testCases196, !hasNEON) }
func Test_ARIA196_Decrypt_NEON(t *testing.T) { BTD(t, BIW(newCipherAsm), CD, testCases196, !hasNEON) }

func Test_ARIA256_Encrypt_NEON(t *testing.T) { BTE(t, BIW(newCipherAsm), CE, testCases256, !hasNEON) }
func Test_ARIA256_Decrypt_NEON(t *testing.T) { BTD(t, BIW(newCipherAsm), CD, testCases256, !hasNEON) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New_NEON(b *testing.B) { BBNA(b, as, 0, BIW(newCipherAsm), !hasNEON) }

func Benchmark_Encrypt_NEON(b *testing.B) { BBDA(b, as, 0, BlockSize, BIW(newCipherAsm), CE, !hasNEON) }
func Benchmark_Decrypt_NEON(b *testing.B) { BBDA(b, as, 0, BlockSize, BIW(newCipherAsm), CD, !hasNEON) }
