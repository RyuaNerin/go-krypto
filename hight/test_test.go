package hight

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_Encrypt_Src(t *testing.T) { BTSC(t, KeySize*8, 0, BlockSize, BIW(NewCipher), CE, false) }
func Test_Decrypt_Src(t *testing.T) { BTSC(t, KeySize*8, 0, BlockSize, BIW(NewCipher), CD, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New(b *testing.B)     { BBN(b, KeySize*8, 0, BIW(NewCipher), false) }
func Benchmark_Encrypt(b *testing.B) { BBD(b, KeySize*8, 0, BlockSize, BIW(NewCipher), CE, false) }
func Benchmark_Decrypt(b *testing.B) { BBD(b, KeySize*8, 0, BlockSize, BIW(NewCipher), CD, false) }
