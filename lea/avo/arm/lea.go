package main

import (
	. "github.com/mmcloughlin/avo/build"
)

func main() {
	Package("kryptosimd/lea/avo")
	ConstraintExpr("arm,gc,!purego")

	leaEnc4NEON()
	leaDec4NEON()

	Generate()
	print("done")
}
