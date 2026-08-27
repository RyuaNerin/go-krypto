package has160

import (
	"encoding/binary"
	"math/bits"
)

func (ctx *has160ContextData) stepBlock(block *[BlockSize]byte) {
	x := [20]uint32{
		binary.LittleEndian.Uint32(block[0:]),
		binary.LittleEndian.Uint32(block[4:]),
		binary.LittleEndian.Uint32(block[8:]),
		binary.LittleEndian.Uint32(block[12:]),
		binary.LittleEndian.Uint32(block[16:]),
		binary.LittleEndian.Uint32(block[20:]),
		binary.LittleEndian.Uint32(block[24:]),
		binary.LittleEndian.Uint32(block[28:]),
		binary.LittleEndian.Uint32(block[32:]),
		binary.LittleEndian.Uint32(block[36:]),
		binary.LittleEndian.Uint32(block[40:]),
		binary.LittleEndian.Uint32(block[44:]),
		binary.LittleEndian.Uint32(block[48:]),
		binary.LittleEndian.Uint32(block[52:]),
		binary.LittleEndian.Uint32(block[56:]),
		binary.LittleEndian.Uint32(block[60:]),
	}

	a, b, c, d, e := ctx.h[0], ctx.h[1], ctx.h[2], ctx.h[3], ctx.h[4]

	// 00 ~ 19
	// Round 1
	x[16] = x[0] ^ x[1] ^ x[2] ^ x[3]
	x[17] = x[4] ^ x[5] ^ x[6] ^ x[7]
	x[18] = x[8] ^ x[9] ^ x[10] ^ x[11]
	x[19] = x[12] ^ x[13] ^ x[14] ^ x[15]
	e, b = f1(a, b, c, d, e, x[18], 5)  // j = 0
	d, a = f1(e, a, b, c, d, x[0], 11)  // j = 1
	c, e = f1(d, e, a, b, c, x[1], 7)   // j = 2
	b, d = f1(c, d, e, a, b, x[2], 15)  // j = 3
	a, c = f1(b, c, d, e, a, x[3], 6)   // j = 4
	e, b = f1(a, b, c, d, e, x[19], 13) // j = 5
	d, a = f1(e, a, b, c, d, x[4], 8)   // j = 6
	c, e = f1(d, e, a, b, c, x[5], 14)  // j = 7
	b, d = f1(c, d, e, a, b, x[6], 7)   // j = 8
	a, c = f1(b, c, d, e, a, x[7], 12)  // j = 9
	e, b = f1(a, b, c, d, e, x[16], 9)  // j = 10
	d, a = f1(e, a, b, c, d, x[8], 11)  // j = 11
	c, e = f1(d, e, a, b, c, x[9], 8)   // j = 12
	b, d = f1(c, d, e, a, b, x[10], 15) // j = 13
	a, c = f1(b, c, d, e, a, x[11], 6)  // j = 14
	e, b = f1(a, b, c, d, e, x[17], 12) // j = 15
	d, a = f1(e, a, b, c, d, x[12], 9)  // j = 16
	c, e = f1(d, e, a, b, c, x[13], 14) // j = 17
	b, d = f1(c, d, e, a, b, x[14], 5)  // j = 18
	a, c = f1(b, c, d, e, a, x[15], 13) // j = 19

	// 20 ~ 39
	// Round 2
	x[16] = x[3] ^ x[6] ^ x[9] ^ x[12]
	x[17] = x[2] ^ x[5] ^ x[8] ^ x[15]
	x[18] = x[1] ^ x[4] ^ x[11] ^ x[14]
	x[19] = x[0] ^ x[7] ^ x[10] ^ x[13]
	e, b = f2(a, b, c, d, e, x[18], 5)  // j = 20
	d, a = f2(e, a, b, c, d, x[3], 11)  // j = 21
	c, e = f2(d, e, a, b, c, x[6], 7)   // j = 22
	b, d = f2(c, d, e, a, b, x[9], 15)  // j = 23
	a, c = f2(b, c, d, e, a, x[12], 6)  // j = 24
	e, b = f2(a, b, c, d, e, x[19], 13) // j = 25
	d, a = f2(e, a, b, c, d, x[15], 8)  // j = 26
	c, e = f2(d, e, a, b, c, x[2], 14)  // j = 27
	b, d = f2(c, d, e, a, b, x[5], 7)   // j = 28
	a, c = f2(b, c, d, e, a, x[8], 12)  // j = 29
	e, b = f2(a, b, c, d, e, x[16], 9)  // j = 30
	d, a = f2(e, a, b, c, d, x[11], 11) // j = 31
	c, e = f2(d, e, a, b, c, x[14], 8)  // j = 32
	b, d = f2(c, d, e, a, b, x[1], 15)  // j = 33
	a, c = f2(b, c, d, e, a, x[4], 6)   // j = 34
	e, b = f2(a, b, c, d, e, x[17], 12) // j = 35
	d, a = f2(e, a, b, c, d, x[7], 9)   // j = 36
	c, e = f2(d, e, a, b, c, x[10], 14) // j = 37
	b, d = f2(c, d, e, a, b, x[13], 5)  // j = 38
	a, c = f2(b, c, d, e, a, x[0], 13)  // j = 39

	// 40 ~ 59
	// Round 3
	x[16] = x[5] ^ x[7] ^ x[12] ^ x[14]
	x[17] = x[0] ^ x[2] ^ x[9] ^ x[11]
	x[18] = x[4] ^ x[6] ^ x[13] ^ x[15]
	x[19] = x[1] ^ x[3] ^ x[8] ^ x[10]
	e, b = f3(a, b, c, d, e, x[18], 5)  // j = 40
	d, a = f3(e, a, b, c, d, x[12], 11) // j = 41
	c, e = f3(d, e, a, b, c, x[5], 7)   // j = 42
	b, d = f3(c, d, e, a, b, x[14], 15) // j = 43
	a, c = f3(b, c, d, e, a, x[7], 6)   // j = 44
	e, b = f3(a, b, c, d, e, x[19], 13) // j = 45
	d, a = f3(e, a, b, c, d, x[0], 8)   // j = 46
	c, e = f3(d, e, a, b, c, x[9], 14)  // j = 47
	b, d = f3(c, d, e, a, b, x[2], 7)   // j = 48
	a, c = f3(b, c, d, e, a, x[11], 12) // j = 49
	e, b = f3(a, b, c, d, e, x[16], 9)  // j = 50
	d, a = f3(e, a, b, c, d, x[4], 11)  // j = 51
	c, e = f3(d, e, a, b, c, x[13], 8)  // j = 52
	b, d = f3(c, d, e, a, b, x[6], 15)  // j = 53
	a, c = f3(b, c, d, e, a, x[15], 6)  // j = 54
	e, b = f3(a, b, c, d, e, x[17], 12) // j = 55
	d, a = f3(e, a, b, c, d, x[8], 9)   // j = 56
	c, e = f3(d, e, a, b, c, x[1], 14)  // j = 57
	b, d = f3(c, d, e, a, b, x[10], 5)  // j = 58
	a, c = f3(b, c, d, e, a, x[3], 13)  // j = 59

	// 60 ~ 79
	// Round 4
	x[16] = x[2] ^ x[7] ^ x[8] ^ x[13]
	x[17] = x[3] ^ x[4] ^ x[9] ^ x[14]
	x[18] = x[0] ^ x[5] ^ x[10] ^ x[15]
	x[19] = x[1] ^ x[6] ^ x[11] ^ x[12]
	e, b = f4(a, b, c, d, e, x[18], 5)  // j = 60
	d, a = f4(e, a, b, c, d, x[7], 11)  // j = 61
	c, e = f4(d, e, a, b, c, x[2], 7)   // j = 62
	b, d = f4(c, d, e, a, b, x[13], 15) // j = 63
	a, c = f4(b, c, d, e, a, x[8], 6)   // j = 64
	e, b = f4(a, b, c, d, e, x[19], 13) // j = 65
	d, a = f4(e, a, b, c, d, x[3], 8)   // j = 66
	c, e = f4(d, e, a, b, c, x[14], 14) // j = 67
	b, d = f4(c, d, e, a, b, x[9], 7)   // j = 68
	a, c = f4(b, c, d, e, a, x[4], 12)  // j = 69
	e, b = f4(a, b, c, d, e, x[16], 9)  // j = 70
	d, a = f4(e, a, b, c, d, x[15], 11) // j = 71
	c, e = f4(d, e, a, b, c, x[10], 8)  // j = 72
	b, d = f4(c, d, e, a, b, x[5], 15)  // j = 73
	a, c = f4(b, c, d, e, a, x[0], 6)   // j = 74
	e, b = f4(a, b, c, d, e, x[17], 12) // j = 75
	d, a = f4(e, a, b, c, d, x[11], 9)  // j = 76
	c, e = f4(d, e, a, b, c, x[6], 14)  // j = 77
	b, d = f4(c, d, e, a, b, x[1], 5)   // j = 78
	a, c = f4(b, c, d, e, a, x[12], 13) // j = 79

	ctx.h[0] += a
	ctx.h[1] += b
	ctx.h[2] += c
	ctx.h[3] += d
	ctx.h[4] += e
}

func f1(a, b, c, d, e, x uint32, s int) (he, hb uint32) {
	he = e + bits.RotateLeft32(a, s) + x + ((b & c) | (^b & d))
	hb = bits.RotateLeft32(b, 10)
	return
}

func f2(a, b, c, d, e, x uint32, s int) (he, hb uint32) {
	he = e + bits.RotateLeft32(a, s) + x + 0x5a827999 + (b ^ c ^ d)
	hb = bits.RotateLeft32(b, 17)
	return
}

func f3(a, b, c, d, e, x uint32, s int) (he, hb uint32) {
	he = e + bits.RotateLeft32(a, s) + x + 0x6ed9eba1 + (c ^ (b | ^d))
	hb = bits.RotateLeft32(b, 25)
	return
}

func f4(a, b, c, d, e, x uint32, s int) (he, hb uint32) {
	he = e + bits.RotateLeft32(a, s) + x + 0x8f1bbcdc + (b ^ c ^ d)
	hb = bits.RotateLeft32(b, 30)
	return
}
