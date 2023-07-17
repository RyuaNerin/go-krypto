package main

import (
	. "github.com/mmcloughlin/avo/build"

	"github.com/RyuaNerin/go-krypto/lsh256/lsh256/lsh256sse2"
)

func main() {
	Package("github.com/RyuaNerin/go-krypto/lsh256/lsh256/lsh256avoconst")

	lsh256sse2.LSH256InitSSE2()
	lsh256sse2.LSH256UpdateSSE2()
	lsh256sse2.LSH256FinalSSE2()

	/**
	lsh256ssse3.LSH256InitSSSE3()
	lsh256ssse3.LSH256UpdateSSSE3()
	lsh256ssse3.LSH256FinalSSSE3()
	*/

	Generate()
}
