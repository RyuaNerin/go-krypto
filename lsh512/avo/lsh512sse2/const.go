package lsh512sse2

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

var (
	g_BytePermInfo Mem
)

func init() {
	f := func(name string, values ...uint64) Mem {
		m := GLOBL(name, NOPTR|RODATA)
		for i, v := range values {
			DATA(4*i, U32(v))
		}
		return m
	}

	g_BytePermInfo = f("g_BytePermInfo_sse2",
		0x00000000, 0x00000000, 0xffffffff, 0xffffffff,
		0xffffffff, 0xffffffff, 0x00000000, 0x00000000,
	)
}
