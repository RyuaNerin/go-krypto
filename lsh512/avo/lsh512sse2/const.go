package lsh512sse2

import (
	"kryptosimd/avoutil"
)

var (
	g_BytePermInfo = avoutil.Alloc32(
		"g_BytePermInfo_sse2",
		0x00000000, 0x00000000, 0xffffffff, 0xffffffff,
		0xffffffff, 0xffffffff, 0x00000000, 0x00000000,
	)
)
