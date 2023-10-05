package seed

import (
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

var (
	as = []CipherSize{
		{Name: "128", Size: 128},
		{Name: "256", Size: 256},
	}
)

func Test_Encrypt_Src(t *testing.T) { BTSCA(t, as, 0, BlockSize, BIW(NewCipher), CE, false) }
func Test_Decrypt_Src(t *testing.T) { BTSCA(t, as, 0, BlockSize, BIW(NewCipher), CD, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New(b *testing.B) { BBNA(b, as, 0, BIW(NewCipher), false) }

func Benchmark_Encrypt(b *testing.B) { BBDA(b, as, 0, BlockSize, BIW(NewCipher), CE, false) }
func Benchmark_Decrypt(b *testing.B) { BBDA(b, as, 0, BlockSize, BIW(NewCipher), CD, false) }
