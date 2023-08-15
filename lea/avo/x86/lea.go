package main

import (
	. "github.com/mmcloughlin/avo/build"
)

func main() {
	Package("kryptosimd/lea/avo")
	ConstraintExpr("amd64,gc,!purego")

	leaEnc4SSE2()
	leaDec4SSE2()

	leaEnc8AVX2()
	leaDec8AVX2()

	Generate()
	print("done")
}
