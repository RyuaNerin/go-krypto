package aria

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

var (
	as = []CipherSize{
		{Name: "128", Size: 128},
		{Name: "196", Size: 196},
		{Name: "256", Size: 256},
	}
)

func Test_ARIA128_Encrypt(t *testing.T) { BTE(t, BIW(newCipherGo), CE, testCases128, false) }
func Test_ARIA128_Decrypt(t *testing.T) { BTD(t, BIW(newCipherGo), CD, testCases128, false) }

func Test_ARIA196_Encrypt(t *testing.T) { BTE(t, BIW(newCipherGo), CE, testCases196, false) }
func Test_ARIA196_Decrypt(t *testing.T) { BTD(t, BIW(newCipherGo), CD, testCases196, false) }

func Test_ARIA256_Encrypt(t *testing.T) { BTE(t, BIW(newCipherGo), CE, testCases256, false) }
func Test_ARIA256_Decrypt(t *testing.T) { BTD(t, BIW(newCipherGo), CD, testCases256, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New(b *testing.B) { BBNA(b, as, 0, BIW(newCipherGo), false) }

func Benchmark_Encrypt(b *testing.B) { BBDA(b, as, 0, BlockSize, BIW(newCipherGo), CE, false) }
func Benchmark_Decrypt(b *testing.B) { BBDA(b, as, 0, BlockSize, BIW(newCipherGo), CD, false) }
