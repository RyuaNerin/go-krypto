package main

import (
	. "github.com/mmcloughlin/avo/build"
)

func main() {
	Package("kryptosimd/lea/avo")

	leaEnc4NEON()
	leaDec4NEON()

	Generate()
	print("done")
}
