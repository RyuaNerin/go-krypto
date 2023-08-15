package avo

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

type leaContext struct {
	round uint8
	rk    [192]uint32
	ecb   bool
}

type LeaContext struct {
	Round Register // U8
	Rk    Mem
}

func GetCtx() LeaContext {
	ctx := Dereference(Param("ctx"))

	round := Load(ctx.Field("round"), GP8())
	rk, err := ctx.Field("rk").Index(0).Resolve()
	if err != nil {
		panic(err)
	}

	return LeaContext{
		Round: round,
		Rk:    rk.Addr,
	}
}
