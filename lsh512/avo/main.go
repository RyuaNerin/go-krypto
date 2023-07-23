package main

import (
	"kryptosimd/lsh512/avo/lsh512sse2"

	. "github.com/mmcloughlin/avo/build"
)

func main() {
	Package("kryptosimd/lsh512/avo/lsh512avoconst")

	lsh512sse2.LSH512InitSSE2()
	lsh512sse2.LSH512UpdateSSE2()
	lsh512sse2.LSH512FinalSSE2()

	Generate()
	print("done")
}
