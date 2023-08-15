package lsh512avx2

import (
	"kryptosimd/avoutil"
)

var (
	g_BytePermInfo = avoutil.Alloc64(
		"g_BytePermInfo_avx2",
		0x0706050403020100, 0x0d0c0b0a09080f0e, 0x1312111017161514, 0x19181f1e1d1c1b1a,
		0x0605040302010007, 0x0c0b0a09080f0e0d, 0x1211101716151413, 0x181f1e1d1c1b1a19,
	)
	g_MsgWordPermInfo = avoutil.Alloc64(
		"g_MsgWordPermInfo_avx2",
		0x0706050403020100, 0x0f0e0d0c0b0a0908, 0x1716151413121110, 0x1f1e1d1c1b1a1918,
	)
)
