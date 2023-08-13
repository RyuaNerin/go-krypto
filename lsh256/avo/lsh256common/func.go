package lsh256common

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

type LSH256_Context struct {
	Algtype            Mem // m32
	Remain_databytelen Register
	Cv_l               Mem
	Cv_r               Mem
	Last_block         Mem
}

func getCtx() *LSH256_Context {
	ctx := Dereference(Param("ctx"))

	algtype, err := ctx.Field("algtype").Resolve()
	if err != nil {
		panic(err)
	}
	cv_l, err := ctx.Field("cv_l").Index(0).Resolve()
	if err != nil {
		panic(err)
	}
	cv_r, err := ctx.Field("cv_r").Index(0).Resolve()
	if err != nil {
		panic(err)
	}
	last_block, err := ctx.Field("last_block").Index(0).Resolve()
	if err != nil {
		panic(err)
	}

	return &LSH256_Context{
		Algtype:            algtype.Addr,
		Remain_databytelen: Load(ctx.Field("remain_databytelen"), GP64()),
		Cv_l:               cv_l.Addr,
		Cv_r:               cv_r.Addr,
		Last_block:         last_block.Addr,
	}
}

func LSH256Init(simd string, body func(ctx *LSH256_Context)) {
	TEXT("lsh256Init"+simd, NOSPLIT, "func(ctx *lsh256ContextAsmData)")

	ctx := getCtx()
	body(ctx)
	Store(ctx.Remain_databytelen, Dereference(Param("ctx")).Field("remain_databytelen"))

	RET()
}

func LSH256Update(simd string, body func(ctx *LSH256_Context, data Mem, databytelen Register)) {
	TEXT("lsh256Update"+simd, NOSPLIT, "func(ctx *lsh256ContextAsmData, data []byte)")

	ctx := getCtx()
	data := Mem{Base: Load(Param("data").Base(), GP64())}
	databytelen := Load(Param("data").Len(), GP64())

	body(ctx, data, databytelen)
	Store(ctx.Remain_databytelen, Dereference(Param("ctx")).Field("remain_databytelen"))

	RET()
}

func LSH256Final(simd string, body func(ctx *LSH256_Context, hashval Mem)) {
	TEXT("lsh256Final"+simd, NOSPLIT, "func(ctx *lsh256ContextAsmData, hashval []byte)")

	ctx := getCtx()
	hashval := Mem{Base: Load(Param("hashval").Base(), GP64())}

	body(ctx, hashval)
	Store(ctx.Remain_databytelen, Dereference(Param("ctx")).Field("remain_databytelen"))

	RET()
}
