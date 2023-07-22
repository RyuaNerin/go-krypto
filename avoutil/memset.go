package lsh256avoconst

import (
	"fmt"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

const memsetName = "memset"

var memsetN = 0

var (
	memsetConstVal = make(map[uint8]Mem)
)

// TODO
// size : gp32
func Memset(dst Op, val uint8, size Register, useAVX2 bool) {
	Comment("memset")

	memsetN++

	ZMM()
	constVal, ok := memsetConstVal[val]
	if !ok {
		constVal = GLOBL(fmt.Sprintf("memset_value_%d", val), NOPTR|RODATA)
		DATA(8*0, U64(uint64(val)*0x_0101_0101_0101_0101))
		DATA(8*1, U64(uint64(val)*0x_0101_0101_0101_0101)) // 128 XMM
		DATA(8*2, U64(uint64(val)*0x_0101_0101_0101_0101))
		DATA(8*3, U64(uint64(val)*0x_0101_0101_0101_0101)) // 256 YMM
		memsetConstVal[val] = constVal
	}

	dstAddr := GP64()
	size2 := GP32()
	idx := GP32()

	LEAQ(dst, dstAddr)
	MOVL(size, size2)
	MOVL(U32(0), idx)

	step := func(sz int, tmp Op, read func(a, b Op), mov func(imr, emr Op)) {
		labelStart := fmt.Sprintf("memset_%d_sz%d_start", memsetN, sz)
		labelEnd := fmt.Sprintf("memset_%d_sz%d_end", memsetN, sz)

		CMPL(size2, U32(sz))
		JL(LabelRef(labelEnd))

		read(constVal, tmp)

		Label(labelStart)
		{
			mov(tmp, Mem{Base: dstAddr}.Idx(idx, 1))

			SUBL(U32(sz), size2)
			ADDL(U32(sz), idx)

			CMPL(size2, U32(sz))
			JL(LabelRef(labelEnd))

			JMP(LabelRef(labelStart))
		}
		Label(labelEnd)
	}

	tmp := GP64()

	if useAVX2 {
		step(32, YMM(), MOVO, VMOVDQ_autoAU)
	}
	step(16, XMM(), MOVO, MOVO_autoAU)
	step(8, tmp.As64(), MOVQ, MOVQ)
	step(4, tmp.As32(), MOVL, MOVL)
	step(2, tmp.As16(), MOVW, MOVW)

	{
		labelStart := fmt.Sprintf("memset_%d_1_start", memsetN)
		labelEnd := fmt.Sprintf("memset_%d_1_end", memsetN)

		Label(labelStart)
		CMPL(size2, U32(0))
		JE(LabelRef(labelEnd))
		{
			MOVB(U8(val), Mem{Base: dstAddr}.Idx(idx, 1))

			SUBL(U32(1), size2)
			ADDL(U32(1), idx)

			JMP(LabelRef(labelStart))
		}
		Label(labelEnd)
	}
}
