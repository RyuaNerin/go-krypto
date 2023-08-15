package lsh512common

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

type LSH512_Context struct {
	Algtype            Register // m32
	Remain_databytelen Register
	Cv_l               Mem
	Cv_r               Mem
	I_last_block       Mem
}

func getCtx() *LSH512_Context {
	ctx := Dereference(Param("ctx"))

	cv_l, err := ctx.Field("cv_l").Index(0).Resolve()
	if err != nil {
		panic(err)
	}
	cv_r, err := ctx.Field("cv_r").Index(0).Resolve()
	if err != nil {
		panic(err)
	}
	last_block, err := ctx.Field("i_last_block").Index(0).Resolve()
	if err != nil {
		panic(err)
	}

	return &LSH512_Context{
		Algtype:            Load(ctx.Field("algtype"), GP32()),
		Remain_databytelen: Load(ctx.Field("remain_databytelen"), GP64()),
		Cv_l:               cv_l.Addr,
		Cv_r:               cv_r.Addr,
		I_last_block:       last_block.Addr,
	}
}

func LSH512Init(simd string, body func(ctx *LSH512_Context)) {
	TEXT("lsh512Init"+simd, NOSPLIT, "func(ctx *lsh512ContextAsmData)")

	ctx := getCtx()
	body(ctx)
	Store(ctx.Remain_databytelen, Dereference(Param("ctx")).Field("remain_databytelen"))

	RET()
}

func LSH512Update(simd string, body func(ctx *LSH512_Context, data Mem, databytelen Register)) {
	TEXT("lsh512Update"+simd, NOSPLIT, "func(ctx *lsh512ContextAsmData, data []byte)")

	ctx := getCtx()
	data := Mem{Base: Load(Param("data").Base(), GP64())}
	databytelen := Load(Param("data").Len(), GP64())

	body(ctx, data, databytelen)
	Store(ctx.Remain_databytelen, Dereference(Param("ctx")).Field("remain_databytelen"))

	RET()
}

func LSH512Final(simd string, body func(ctx *LSH512_Context, hashval Mem)) {
	TEXT("lsh512Final"+simd, NOSPLIT, "func(ctx *lsh512ContextAsmData, hashval []byte)")

	ctx := getCtx()
	hashval := Mem{Base: Load(Param("hashval").Base(), GP64())}

	body(ctx, hashval)
	Store(ctx.Remain_databytelen, Dereference(Param("ctx")).Field("remain_databytelen"))

	RET()
}
