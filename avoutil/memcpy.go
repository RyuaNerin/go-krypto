package avoutil

import (
	"fmt"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

const memcpyName = "memcpy"

var memcpyN = 0

// TODO
// size: uint64
func Memcpy(dst, src Op, size Register, avx2 bool) {
	if size.Size() != 8 {
		panic("wong arguments")
	}

	Comment("Memcpy")
	memcpyN++

	op256 := YMM()
	op128 := op256.AsX()
	op64 := GP64()
	op32 := op64.As32()
	op16 := op64.As16()
	op8 := op64.As8()

	//////////////////////////////

	dstAddr := GP64()
	srcAddr := GP64()
	remain := GP64()

	LEAQ(dst, dstAddr)
	LEAQ(src, srcAddr)
	MOVQ(size, remain)

	//////////////////////////////

	cpy := func(sz int, tmp Op, mov func(a, b Op)) {
		labelStart := fmt.Sprintf("memcpy_%d_sz%d_start", memcpyN, sz)
		labelEnd := fmt.Sprintf("memcpy_%d_sz%d_end", memcpyN, sz)

		Label(labelStart)
		CMPQ(remain, U32(sz))
		JL(LabelRef(labelEnd))
		{
			mov(Mem{Base: srcAddr}, tmp)
			mov(tmp, Mem{Base: dstAddr})

			ADDQ(U32(sz), srcAddr)
			ADDQ(U32(sz), dstAddr)
			SUBQ(U32(sz), remain)

			JMP(LabelRef(labelStart))
		}
		Label(labelEnd)
	}

	if enableXYZ {
		if avx2 {
			cpy(32, op256, VMOVDQ_autoAU)
		}
		cpy(16, op128, MOVO_autoAU)

	}
	cpy(8, op64, MOVQ)
	cpy(4, op32, MOVL)
	cpy(2, op16, MOVW)
	cpy(1, op8, MOVB)
}

func MemcpyStatic(dst, src Mem, size int, avx2 bool) {
	Comment("MemcpyStatic")

	op256 := YMM()
	op128 := XMM()
	op64 := GP64()
	op32 := op64.As32()
	op16 := op64.As16()
	op8 := op64.As8()

	//////////////////////////////

	idx := 0
	for size > 0 {
		sz := 1

		if enableXYZ && avx2 && size >= 32 {
			VMOVDQ_autoAU(src.Offset(idx), op256)
			VMOVDQ_autoAU(op256, dst.Offset(idx))
			sz = 32
		} else if enableXYZ && size >= 16 {
			MOVO_autoAU(src.Offset(idx), op128)
			MOVO_autoAU(op128, dst.Offset(idx))
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
