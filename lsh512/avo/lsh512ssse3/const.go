package lsh512ssse3

import (
	"kryptosimd/avoutil"
)

var (
	g_BytePermInfo = avoutil.Alloc64(
		"g_BytePermInfo_ssse3",
		0x0706050403020100, 0x0d0c0b0a09080f0e,
		0x0302010007060504, 0x09080f0e0d0c0b0a,
		0x0605040302010007, 0x0c0b0a09080f0e0d,
		0x0201000706050403, 0x080f0e0d0c0b0a09,
	)
	g_MsgWordPermInfo = avoutil.Alloc64(
		"g_MsgWordPermInfo_ssse3",
		0x0706050403020100, 0x0f0e0d0c0b0a0908, 0x1716151413121110, 0x1f1e1d1c1b1a1918,
	)
)
