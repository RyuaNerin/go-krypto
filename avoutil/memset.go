package avoutil

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
	if size.Size() != 8 {
		panic("size must be gp64")
	}

	Comment("memset")

	memsetN++

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
	remain := GP64()

	LEAQ(dst, dstAddr)
	MOVQ(size, remain)

	step := func(sz int, tmp Op, read func(a, b Op), mov func(imr, emr Op)) {
		if sz >= 16 && !enableXYZ {
			return
		}

		labelStart := fmt.Sprintf("memset_%d_sz%d_start", memsetN, sz)
		labelEnd := fmt.Sprintf("memset_%d_sz%d_end", memsetN, sz)

		CMPQ(remain, U32(sz))
		JL(LabelRef(labelEnd))

		read(constVal, tmp)

		Label(labelStart)
		{
			mov(tmp, Mem{Base: dstAddr})

			SUBQ(U32(sz), remain)
			ADDQ(U32(sz), dstAddr)

			CMPQ(remain, U32(sz))
			JL(LabelRef(labelEnd))

			JMP(LabelRef(labelStart))
		}
		Label(labelEnd)
	}

	tmp := GP64()

	if enableXYZ {
		reg := YMM()
		if useAVX2 {
			step(32, reg, MOVO, VMOVDQa)
		}
		step(16, reg.AsX(), MOVO, MOVOa)
	}
	step(8, tmp.As64(), MOVQ, MOVQ)
	step(4, tmp.As32(), MOVL, MOVL)
	step(2, tmp.As16(), MOVW, MOVW)

	{
		labelStart := fmt.Sprintf("memset_%d_1_start", memsetN)
		labelEnd := fmt.Sprintf("memset_%d_1_end", memsetN)

		Label(labelStart)
		CMPQ(remain, U32(0))
		JE(LabelRef(labelEnd))
		{
			MOVB(U8(val), Mem{Base: dstAddr})

			SUBQ(U32(1), remain)
			ADDQ(U32(1), dstAddr)

			JMP(LabelRef(labelStart))
		}
		Label(labelEnd)
	}
}
