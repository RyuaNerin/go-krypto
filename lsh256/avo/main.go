package main

import (
	"kryptosimd/lsh256/avo/lsh256avx2"
	"kryptosimd/lsh256/avo/lsh256sse2"
	"kryptosimd/lsh256/avo/lsh256ssse3"

	. "github.com/mmcloughlin/avo/build"
)

func main() {
	Package("kryptosimd/lsh256/avo/lsh256avoconst")

	lsh256sse2.LSH256InitSSE2()
	lsh256sse2.LSH256UpdateSSE2()
	lsh256sse2.LSH256FinalSSE2()

	//lsh256ssse3.LSH256InitSSSE3() // same with LSH256InitSSE2
	lsh256ssse3.LSH256UpdateSSSE3()
	lsh256ssse3.LSH256FinalSSSE3()

	lsh256avx2.LSH256InitAVX2()
	lsh256avx2.LSH256UpdateAVX2()
	lsh256avx2.LSH256FinalAVX2()

	Generate()
	print("done")
}
