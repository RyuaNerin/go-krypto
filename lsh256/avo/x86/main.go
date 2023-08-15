package main

import (
	"kryptosimd/lsh256/avo/x86/lsh256avx2"
	. "kryptosimd/lsh256/avo/x86/lsh256common"
	"kryptosimd/lsh256/avo/x86/lsh256sse2"
	"kryptosimd/lsh256/avo/x86/lsh256ssse3"

	. "github.com/mmcloughlin/avo/build"
)

func main() {
	Package("kryptosimd/lsh256/avo/x86/lsh256avoconst")
	ConstraintExpr("amd64,gc,!purego")

	LSH256Init("SSE2", lsh256sse2.Lsh256_sse2_init)
	LSH256Update("SSE2", lsh256sse2.Lsh256_sse2_update)
	LSH256Final("SSE2", lsh256sse2.Lsh256_sse2_final)

	//LSH256Init("SSSE3", lsh256ssse3.Lsh256_ssse3_init) // same with SSE2
	LSH256Update("SSSE3", lsh256ssse3.Lsh256_ssse3_update)
	LSH256Final("SSSE3", lsh256ssse3.Lsh256_ssse3_final)

	LSH256Init("AVX2", lsh256avx2.Lsh256_avx2_init)
	LSH256Update("AVX2", lsh256avx2.Lsh256_avx2_update)
	LSH256Final("AVX2", lsh256avx2.Lsh256_avx2_final)

	Generate()
	print("done")
}
