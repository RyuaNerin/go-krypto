package aria

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_ARIA128_Encrypt_Go(t *testing.T) { BTE(t, BIW(newCipherGo), CE, testCases128, false) }
func Test_ARIA128_Decrypt_Go(t *testing.T) { BTD(t, BIW(newCipherGo), CD, testCases128, false) }

func Test_ARIA196_Encrypt_Go(t *testing.T) { BTE(t, BIW(newCipherGo), CE, testCases196, false) }
func Test_ARIA196_Decrypt_Go(t *testing.T) { BTD(t, BIW(newCipherGo), CD, testCases196, false) }

func Test_ARIA256_Encrypt_Go(t *testing.T) { BTE(t, BIW(newCipherGo), CE, testCases256, false) }
func Test_ARIA256_Decrypt_Go(t *testing.T) { BTD(t, BIW(newCipherGo), CD, testCases256, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New_Go(b *testing.B) { BBNA(b, as, 0, BIW(newCipherGo), false) }

func Benchmark_Encrypt_Go(b *testing.B) { BBDA(b, as, 0, BlockSize, BIW(newCipherGo), CE, false) }
func Benchmark_Decrypt_Go(b *testing.B) { BBDA(b, as, 0, BlockSize, BIW(newCipherGo), CD, false) }
