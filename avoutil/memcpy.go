package lsh256avoconst

import (
	"fmt"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

const memcpyName = "memcpy"

var memcpyN = 0

// TODO
// size: int32
func Memcpy(dst, src Op, size Register) {
	Comment("Memcpy")

	memcpyN++

	tmp8 := GP8()

	dstAddr := GP64()
	srcAddr := GP64()
	size2 := GP32()

	LEAQ(dst, dstAddr)
	LEAQ(src, srcAddr)
	MOVL(size, size2)

	labelStart := fmt.Sprintf("memcpy_%d_while_start", memcpyN)
	labelEnd := fmt.Sprintf("memcpy_%d_while_end", memcpyN)

	Label(labelStart)
	CMPL(size2, U32(0))
	JE(LabelRef(labelEnd))
	{
		MOVB(Mem{Base: srcAddr}, tmp8)
		MOVB(tmp8, Mem{Base: dstAddr})

		ADDQ(U32(1), srcAddr)
		ADDQ(U32(1), dstAddr)
		SUBL(U32(1), size2)

		JMP(LabelRef(labelStart))
	}
	Label(labelEnd)
}

func MemcpyStatic(dstMem, srcMem Mem, size int, avx2 bool, xmmTmp VecVirtual, ymmTmp VecVirtual) {
	Comment("MemcpyStatic")

	op256 := ymmTmp
	op128 := xmmTmp
	op64 := GP64()
	op32 := op64.As32()
	op16 := op64.As16()
	op8 := op64.As8()

	if xmmTmp == nil && size >= 16 {
		op128 = XMM()
	}
	if ymmTmp == nil && size >= 32 && avx2 {
		op256 = YMM()
	}

	dst := Mem{Base: GP64()}
	src := Mem{Base: GP64()}

	LEAQ(dstMem, dst.Base)
	LEAQ(srcMem, src.Base)

	idx := 0
	for size > 0 {
		sz := 1

		if size >= 32 && avx2 {
			VMOVDQU(src.Offset(idx), op256)
			VMOVDQU(op256, dst.Offset(idx))
			sz = 32
		} else if size >= 16 {
			MOVOU(src.Offset(idx), op128)
			MOVOU(op128, dst.Offset(idx))
			sz = 16
		} else if size >= 8 {
			MOVQ(src.Offset(idx), op64)
			MOVQ(op64, dst.Offset(idx))
			sz = 8
		} else if size >= 4 {
			MOVL(src.Offset(idx), op32)
			MOVL(op32, dst.Offset(idx))
			sz = 4
		} else if size >= 2 {
			MOVW(src.Offset(idx), op16)
			MOVW(op16, dst.Offset(idx))
			sz = 2
		} else {
			MOVB(src.Offset(idx), op8)
			MOVB(op8, dst.Offset(idx))
		}

		size -= sz
		idx += sz
	}
}
