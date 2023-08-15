package lsh256avx2

import (
	"kryptosimd/avoutil"
)

var (
	g_BytePermInfo = avoutil.Alloc32(
		"g_BytePermInfo_avx2",
		0x03020100, 0x06050407, 0x09080b0a, 0x0c0f0e0d,
		0x10131211, 0x15141716, 0x1a19181b, 0x1f1e1d1c,
	)
	g_MsgWordPermInfo = avoutil.Alloc32(
		"g_MsgWordPermInfo_avx2",
		0x0f0e0d0c, 0x0b0a0908, 0x03020100, 0x07060504,
		0x1f1e1d1c, 0x13121110, 0x17161514, 0x1b1a1918,
	)
)
