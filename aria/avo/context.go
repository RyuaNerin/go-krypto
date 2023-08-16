package avo

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

type ariaContext struct {
	round   uint32
	ecb     bool
	enc_key [68]uint32
	dec_key [68]uint32
}

type AriaContext struct {
	Round   Register // U32
	Enc_key Mem
	Dec_key Mem
}

func GetCtx() AriaContext {
	ctx := Dereference(Param("ctx"))

	round := Load(ctx.Field("round"), GP32())
	enc_key, err := ctx.Field("enc_key").Index(0).Resolve()
	if err != nil {
		panic(err)
	}
	dec_key, err := ctx.Field("dec_key").Index(0).Resolve()
	if err != nil {
		panic(err)
	}

	return AriaContext{
		Round:   round,
		Enc_key: enc_key.Addr,
		Dec_key: dec_key.Addr,
	}
}
