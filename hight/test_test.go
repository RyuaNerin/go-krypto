package hight

import (
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_New(b *testing.B)     { BBNew(b, KeySize*8, 0, BIW(NewCipher)) }
func Benchmark_Encrypt(b *testing.B) { BBDo(b, KeySize*8, 0, BlockSize, BIW(NewCipher), CE) }
func Benchmark_Decrypt(b *testing.B) { BBDo(b, KeySize*8, 0, BlockSize, BIW(NewCipher), CD) }
