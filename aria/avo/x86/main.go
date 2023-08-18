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
	tmp := XMM()
	F_mm_storeu_si128(
		dst,
		F_mm_xor_si128(
			tmp,
			F_mm_loadu_si128(tmp, rk),
			F_mm_loadu_si128(XMM(), t),
		),
	)

	RET()
}
