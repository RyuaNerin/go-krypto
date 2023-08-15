package main

import (
	//. "kryptosimd/lea/avo"

	. "github.com/mmcloughlin/avo/build"
	//. "github.com/mmcloughlin/avo/operand"
	//. "github.com/mmcloughlin/avo/reg"
)

func leaEnc4NEON() {
	TEXT("leaEnc4NEON", NOSPLIT, "func(ctx *leaContext, dst []byte, src []byte)")

	RET()
}

func leaDec4NEON() {
	TEXT("leaDec4NEON", NOSPLIT, "func(ctx *leaContext, dst []byte, src []byte)")

	RET()
}
