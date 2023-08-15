package lsh256ssse3

import (
	"kryptosimd/avoutil"
)

var (
	g_BytePermInfo = avoutil.Alloc32(
		"g_BytePermInfo_ssse3",
		0x03020100, 0x06050407, 0x09080b0a, 0x0c0f0e0d,
		0x00030201, 0x05040706, 0x0a09080b, 0x0f0e0d0c,
	)
	/**
	g_BytePermInfo_L = avoutil.Alloc32(
		"g_BytePermInfo_L_ssse3",
		0x03020100, 0x06050407, 0x09080b0a, 0x0c0f0e0d,
	)
	g_BytePermInfo_R = avoutil.Alloc32(
		"g_BytePermInfo_R_ssse3",
	)
	g_MsgWordPermInfo = avoutil.Alloc32(
		"g_MsgWordPermInfo_ssse3",
		0x0f0e0d0c, 0x0b0a0908, 0x03020100, 0x07060504,
		0x1f1e1d1c, 0x13121110, 0x17161514, 0x1b1a1918,
	)
	*/
)
