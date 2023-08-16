package main

import (
	. "github.com/mmcloughlin/avo/build"
)

func main() {
	Package("kryptosimd/aria/x86/avo")
	ConstraintExpr("amd64,gc,!purego")

	Generate()
	print("done")
}
