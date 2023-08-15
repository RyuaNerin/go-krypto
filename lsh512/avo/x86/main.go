package main

import (
	"kryptosimd/lsh512/avo/x86/lsh512avx2"
	"kryptosimd/lsh512/avo/x86/lsh512common"
	"kryptosimd/lsh512/avo/x86/lsh512sse2"
	"kryptosimd/lsh512/avo/x86/lsh512ssse3"

	. "github.com/mmcloughlin/avo/build"
)

func main() {
	Package("kryptosimd/lsh512/avo/x86/lsh512avoconst")
	ConstraintExpr("amd64,gc,!purego")

	lsh512common.LSH512Init("SSE2", lsh512sse2.Lsh512_sse2_init)
	lsh512common.LSH512Update("SSE2", lsh512sse2.Lsh512_sse2_update)
	lsh512common.LSH512Final("SSE2", lsh512sse2.Lsh512_sse2_final)
	//lsh512sse2.UnitTest()

	//lsh512common.LSH512Init("SSSE3", lsh512ssse3.Lsh512_ssse3_init) // same with sse2
	lsh512common.LSH512Update("SSSE3", lsh512ssse3.Lsh512_ssse3_update)
	lsh512common.LSH512Final("SSSE3", lsh512ssse3.Lsh512_ssse3_final)
	//lsh512ssse3.UnitTest()

	lsh512common.LSH512Init("AVX2", lsh512avx2.Lsh512_avx2_init)
	lsh512common.LSH512Update("AVX2", lsh512avx2.Lsh512_avx2_update)
	lsh512common.LSH512Final("AVX2", lsh512avx2.Lsh512_avx2_final)
	//lsh512avx2.UnitTest()

	Generate()
	print("done")
}
