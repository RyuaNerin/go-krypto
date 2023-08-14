package avoutil

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

func Alloc32(name string, values ...uint32) Mem {
	mem := GLOBL(name, NOPTR|RODATA)

	for idx, v := range values {
		DATA(4*idx, U32(v))
	}
	return mem
}

func Alloc64(name string, values ...uint64) Mem {
	mem := GLOBL(name, NOPTR|RODATA)

	for idx, v := range values {
		DATA(8*idx, U64(v))
	}
	return mem
}
