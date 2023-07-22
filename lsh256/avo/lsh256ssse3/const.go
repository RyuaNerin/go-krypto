package lsh256ssse3

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
)

var (
	g_BytePermInfo_L  Mem
	g_BytePermInfo_R  Mem
	g_MsgWordPermInfo Mem
)

func init() {
	f := func(name string, values ...uint64) Mem {
		m := GLOBL(name, NOPTR|RODATA)
		for i, v := range values {
			DATA(4*i, U32(v))
		}
		return m
	}

	g_BytePermInfo_L = f("g_BytePermInfo_L_ssse3",
		0x03020100, 0x06050407, 0x09080b0a, 0x0c0f0e0d,
	)
	g_BytePermInfo_R = f("g_BytePermInfo_R_ssse3",
		0x00030201, 0x05040706, 0x0a09080b, 0x0f0e0d0c,
	)
	g_MsgWordPermInfo = f("g_MsgWordPermInfo_ssse3",
		0x0f0e0d0c, 0x0b0a0908, 0x03020100, 0x07060504,
		0x1f1e1d1c, 0x13121110, 0x17161514, 0x1b1a1918,
	)
}
