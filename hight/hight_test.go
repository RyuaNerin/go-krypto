package hight

import (
	"testing"

	. "github.com/RyuaNerin/testingutil"
)

func Test_HIGHT_Encrypt(t *testing.T) { BTE(t, BIW(NewCipher), CE, testCases, false) }
func Test_HIGHT_Decrypt(t *testing.T) { BTD(t, BIW(NewCipher), CD, testCases, false) }

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New(b *testing.B)     { BBN(b, KeySize*8, 0, BIW(NewCipher), false) }
func Benchmark_Encrypt(b *testing.B) { BBD(b, KeySize*8, 0, BlockSize, BIW(NewCipher), CE, false) }
func Benchmark_Decrypt(b *testing.B) { BBD(b, KeySize*8, 0, BlockSize, BIW(NewCipher), CD, false) }
