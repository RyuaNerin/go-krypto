package main

import (
	"kryptosimd/lsh512/avo/lsh512common"
	"kryptosimd/lsh512/avo/lsh512sse2"

	. "github.com/mmcloughlin/avo/build"
)

func main() {
	Package("kryptosimd/lsh512/avo/lsh512avoconst")

	lsh512common.LSH512Init("SSE2", lsh512sse2.Lsh512_sse2_init)
	lsh512common.LSH512Update("SSE2", lsh512sse2.Lsh512_sse2_update)
	lsh512common.LSH512Final("SSE2", lsh512sse2.Lsh512_sse2_final)

	Generate()
	print("done")
}
