//go:build amd64 && !purego
// +build amd64,!purego

package aria

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_ARIA128_Encrypt_SSSE3(t *testing.T) { BTE(t, BIW(newCipherAsm), CE, testCases128, !hasSSSE3) }
func Test_ARIA128_Decrypt_SSSE3(t *testing.T) { BTD(t, BIW(newCipherAsm), CD, testCases128, !hasSSSE3) }

func Test_ARIA196_Encrypt_SSSE3(t *testing.T) { BTE(t, BIW(newCipherAsm), CE, testCases196, !hasSSSE3) }
func Test_ARIA196_Decrypt_SSSE3(t *testing.T) { BTD(t, BIW(newCipherAsm), CD, testCases196, !hasSSSE3) }

func Test_ARIA256_Encrypt_SSSE3(t *testing.T) { BTE(t, BIW(newCipherAsm), CE, testCases256, !hasSSSE3) }
func Test_ARIA256_Decrypt_SSSE3(t *testing.T) { BTD(t, BIW(newCipherAsm), CD, testCases256, !hasSSSE3) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New_SSSE3(b *testing.B) { BBNA(b, as, 0, BIW(newCipherAsm), !hasSSSE3) }

func Benchmark_Encrypt_SSSE3(b *testing.B) {
	BBDA(b, as, 0, BlockSize, BIW(newCipherAsm), CE, !hasSSSE3)
}
func Benchmark_Decrypt_SSSE3(b *testing.B) {
	BBDA(b, as, 0, BlockSize, BIW(newCipherAsm), CD, !hasSSSE3)
}
