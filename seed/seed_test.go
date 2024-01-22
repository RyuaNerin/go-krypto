package seed

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

var (
	as = []CipherSize{
		{Name: "128", Size: 128},
	}
)

func Test_SEED_Encrypt(t *testing.T) { BTE(t, BIW(NewCipher), CE, testCases128, false) }
func Test_SEED_Decrypt(t *testing.T) { BTD(t, BIW(NewCipher), CD, testCases128, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New(b *testing.B) { BBNA(b, as, 0, BIW(NewCipher), false) }

func Benchmark_Encrypt(b *testing.B) { BBDA(b, as, 0, BlockSize, BIW(NewCipher), CE, false) }
func Benchmark_Decrypt(b *testing.B) { BBDA(b, as, 0, BlockSize, BIW(NewCipher), CD, false) }
