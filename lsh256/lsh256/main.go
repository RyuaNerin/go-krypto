package main

import (
	"github.com/RyuaNerin/go-krypto/lsh256/lsh256/lsh256avx2"
	"github.com/RyuaNerin/go-krypto/lsh256/lsh256/lsh256sse2"

	. "github.com/mmcloughlin/avo/build"
)

func main() {
	Package("github.com/RyuaNerin/go-krypto/lsh256/lsh256/lsh256avoconst")

	lsh256sse2.LSH256InitSSE2()
	lsh256sse2.LSH256UpdateSSE2()
	lsh256sse2.LSH256FinalSSE2()

	lsh256avx2.LSH256InitAVX2()
	lsh256avx2.LSH256UpdateAVX2()
	lsh256avx2.LSH256FinalAVX2()

	Generate()
	print("done")
}
