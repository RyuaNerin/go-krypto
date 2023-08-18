package main

import (
	. "kryptosimd/avoutil/simd"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func main() {
	Package("kryptosimd/aria/avo/x86")
	ConstraintExpr("amd64,gc,!purego")

	processFinSSE2()

	Generate()
	print("done")
}

func processFinSSE2() {
	TEXT("processFinSSE2", NOSPLIT, "func(dst []byte, rk []byte, t []byte)")

	dst := Mem{Base: Load(Param("dst").Base(), GP64())}
	rk := Mem{Base: Load(Param("rk").Base(), GP64())}
	t := Mem{Base: Load(Param("t").Base(), GP64())}

	/**
	for j := 0; j < 16; j++ {
		dst[j] = rk[j] ^ t[j]
	}
	*/
	F_mm_storeu_si128(
		dst,
		F_mm_xor_si128(
			A_mm_loadu_si128(rk),
			A_mm_loadu_si128(t),
		),
	)

	RET()
}
