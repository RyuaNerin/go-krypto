package main

import (
	. "kryptosimd/lea/avo"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func leaEnc4NEON() {
	TEXT("leaEnc4NEON", NOSPLIT, "func(ctx *leaContext, dst []byte, src []byte)")

	ctx := GetCtx()
	dst := Mem{Base: Load(Param("dst").Base(), GP64())}
	src := Mem{Base: Load(Param("src").Base(), GP64())}

	XMM()

	/**
	__m128i x0, x1, x2, x3, tmp;
	__m128i tmp0, tmp1, tmp2, tmp3;
	*/
	x0 := XMM()
	x1 := XMM()
	x2 := XMM()
	x3 := XMM()

	MOVUPS(src.Offset(0x00), x0)
	MOVUPS(src.Offset(0x10), x1)
	MOVUPS(src.Offset(0x20), x2)
	MOVUPS(src.Offset(0x30), x3)
	x0, x1, x2, x3 = leaSSE2Swap(x0, x1, x2, x3)

	XAR3_SSE2(ctx.Rk, x3, x2, 4, 5)
	XAR5_SSE2(ctx.Rk, x2, x1, 2, 3)
	XAR9_SSE2(ctx.Rk, x1, x0, 0, 1)
	XAR3_SSE2(ctx.Rk, x0, x3, 10, 11)
	XAR5_SSE2(ctx.Rk, x3, x2, 8, 9)
	XAR9_SSE2(ctx.Rk, x2, x1, 6, 7)
	XAR3_SSE2(ctx.Rk, x1, x0, 16, 17)
	XAR5_SSE2(ctx.Rk, x0, x3, 14, 15)
	XAR9_SSE2(ctx.Rk, x3, x2, 12, 13)
	XAR3_SSE2(ctx.Rk, x2, x1, 22, 23)
	XAR5_SSE2(ctx.Rk, x1, x0, 20, 21)
	XAR9_SSE2(ctx.Rk, x0, x3, 18, 19)

	XAR3_SSE2(ctx.Rk, x3, x2, 28, 29)
	XAR5_SSE2(ctx.Rk, x2, x1, 26, 27)
	XAR9_SSE2(ctx.Rk, x1, x0, 24, 25)
	XAR3_SSE2(ctx.Rk, x0, x3, 34, 35)
	XAR5_SSE2(ctx.Rk, x3, x2, 32, 33)
	XAR9_SSE2(ctx.Rk, x2, x1, 30, 31)
	XAR3_SSE2(ctx.Rk, x1, x0, 40, 41)
	XAR5_SSE2(ctx.Rk, x0, x3, 38, 39)
	XAR9_SSE2(ctx.Rk, x3, x2, 36, 37)
	XAR3_SSE2(ctx.Rk, x2, x1, 46, 47)
	XAR5_SSE2(ctx.Rk, x1, x0, 44, 45)
	XAR9_SSE2(ctx.Rk, x0, x3, 42, 43)

	XAR3_SSE2(ctx.Rk, x3, x2, 52, 53)
	XAR5_SSE2(ctx.Rk, x2, x1, 50, 51)
	XAR9_SSE2(ctx.Rk, x1, x0, 48, 49)
	XAR3_SSE2(ctx.Rk, x0, x3, 58, 59)
	XAR5_SSE2(ctx.Rk, x3, x2, 56, 57)
	XAR9_SSE2(ctx.Rk, x2, x1, 54, 55)
	XAR3_SSE2(ctx.Rk, x1, x0, 64, 65)
	XAR5_SSE2(ctx.Rk, x0, x3, 62, 63)
	XAR9_SSE2(ctx.Rk, x3, x2, 60, 61)
	XAR3_SSE2(ctx.Rk, x2, x1, 70, 71)
	XAR5_SSE2(ctx.Rk, x1, x0, 68, 69)
	XAR9_SSE2(ctx.Rk, x0, x3, 66, 67)

	XAR3_SSE2(ctx.Rk, x3, x2, 76, 77)
	XAR5_SSE2(ctx.Rk, x2, x1, 74, 75)
	XAR9_SSE2(ctx.Rk, x1, x0, 72, 73)
	XAR3_SSE2(ctx.Rk, x0, x3, 82, 83)
	XAR5_SSE2(ctx.Rk, x3, x2, 80, 81)
	XAR9_SSE2(ctx.Rk, x2, x1, 78, 79)
	XAR3_SSE2(ctx.Rk, x1, x0, 88, 89)
	XAR5_SSE2(ctx.Rk, x0, x3, 86, 87)
	XAR9_SSE2(ctx.Rk, x3, x2, 84, 85)
	XAR3_SSE2(ctx.Rk, x2, x1, 94, 95)
	XAR5_SSE2(ctx.Rk, x1, x0, 92, 93)
	XAR9_SSE2(ctx.Rk, x0, x3, 90, 91)

	XAR3_SSE2(ctx.Rk, x3, x2, 100, 101)
	XAR5_SSE2(ctx.Rk, x2, x1, 98, 99)
	XAR9_SSE2(ctx.Rk, x1, x0, 96, 97)
	XAR3_SSE2(ctx.Rk, x0, x3, 106, 107)
	XAR5_SSE2(ctx.Rk, x3, x2, 104, 105)
	XAR9_SSE2(ctx.Rk, x2, x1, 102, 103)
	XAR3_SSE2(ctx.Rk, x1, x0, 112, 113)
	XAR5_SSE2(ctx.Rk, x0, x3, 110, 111)
	XAR9_SSE2(ctx.Rk, x3, x2, 108, 109)
	XAR3_SSE2(ctx.Rk, x2, x1, 118, 119)
	XAR5_SSE2(ctx.Rk, x1, x0, 116, 117)
	XAR9_SSE2(ctx.Rk, x0, x3, 114, 115)

	XAR3_SSE2(ctx.Rk, x3, x2, 124, 125)
	XAR5_SSE2(ctx.Rk, x2, x1, 122, 123)
	XAR9_SSE2(ctx.Rk, x1, x0, 120, 121)
	XAR3_SSE2(ctx.Rk, x0, x3, 130, 131)
	XAR5_SSE2(ctx.Rk, x3, x2, 128, 129)
	XAR9_SSE2(ctx.Rk, x2, x1, 126, 127)
	XAR3_SSE2(ctx.Rk, x1, x0, 136, 137)
	XAR5_SSE2(ctx.Rk, x0, x3, 134, 135)
	XAR9_SSE2(ctx.Rk, x3, x2, 132, 133)
	XAR3_SSE2(ctx.Rk, x2, x1, 142, 143)
	XAR5_SSE2(ctx.Rk, x1, x0, 140, 141)
	XAR9_SSE2(ctx.Rk, x0, x3, 138, 139)

	CMPB(ctx.Round, U8(24))
	JBE(LabelRef("OVER24_END"))
	XAR3_SSE2(ctx.Rk, x3, x2, 148, 149)
	XAR5_SSE2(ctx.Rk, x2, x1, 146, 147)
	XAR9_SSE2(ctx.Rk, x1, x0, 144, 145)
	XAR3_SSE2(ctx.Rk, x0, x3, 154, 155)
	XAR5_SSE2(ctx.Rk, x3, x2, 152, 153)
	XAR9_SSE2(ctx.Rk, x2, x1, 150, 151)
	XAR3_SSE2(ctx.Rk, x1, x0, 160, 161)
	XAR5_SSE2(ctx.Rk, x0, x3, 158, 159)
	XAR9_SSE2(ctx.Rk, x3, x2, 156, 157)
	XAR3_SSE2(ctx.Rk, x2, x1, 166, 167)
	XAR5_SSE2(ctx.Rk, x1, x0, 164, 165)
	XAR9_SSE2(ctx.Rk, x0, x3, 162, 163)
	Label("OVER24_END")

	CMPB(ctx.Round, U8(28))
	JBE(LabelRef("OVER28_END"))
	XAR3_SSE2(ctx.Rk, x3, x2, 172, 173)
	XAR5_SSE2(ctx.Rk, x2, x1, 170, 171)
	XAR9_SSE2(ctx.Rk, x1, x0, 168, 169)
	XAR3_SSE2(ctx.Rk, x0, x3, 178, 179)
	XAR5_SSE2(ctx.Rk, x3, x2, 176, 177)
	XAR9_SSE2(ctx.Rk, x2, x1, 174, 175)
	XAR3_SSE2(ctx.Rk, x1, x0, 184, 185)
	XAR5_SSE2(ctx.Rk, x0, x3, 182, 183)
	XAR9_SSE2(ctx.Rk, x3, x2, 180, 181)
	XAR3_SSE2(ctx.Rk, x2, x1, 190, 191)
	XAR5_SSE2(ctx.Rk, x1, x0, 188, 189)
	XAR9_SSE2(ctx.Rk, x0, x3, 186, 187)
	Label("OVER28_END")

	x0, x1, x2, x3 = leaSSE2Swap(x0, x1, x2, x3)
	MOVUPS(x0, dst.Offset(0x00))
	MOVUPS(x1, dst.Offset(0x10))
	MOVUPS(x2, dst.Offset(0x20))
	MOVUPS(x3, dst.Offset(0x30))

	/**
	return;
	*/
	RET()
}

func leaDec4NEON() {
	TEXT("leaDec4NEON", NOSPLIT, "func(ctx *leaContext, dst []byte, src []byte)")

	ctx := GetCtx()
	dst := Mem{Base: Load(Param("dst").Base(), GP64())}
	src := Mem{Base: Load(Param("src").Base(), GP64())}

	/**
	__m128i x0, x1, x2, x3;
	*/
	x0 := XMM()
	x1 := XMM()
	x2 := XMM()
	x3 := XMM()

	MOVUPS(src.Offset(0x00), x0)
	MOVUPS(src.Offset(0x10), x1)
	MOVUPS(src.Offset(0x20), x2)
	MOVUPS(src.Offset(0x30), x3)
	x0, x1, x2, x3 = leaSSE2Swap(x0, x1, x2, x3)

	CMPB(ctx.Round, U8(28))
	JBE(LabelRef("OVER28_END"))
	XSR9_SSE2(ctx.Rk, x0, x3, 186, 187)
	XSR5_SSE2(ctx.Rk, x1, x0, 188, 189)
	XSR3_SSE2(ctx.Rk, x2, x1, 190, 191)
	XSR9_SSE2(ctx.Rk, x3, x2, 180, 181)
	XSR5_SSE2(ctx.Rk, x0, x3, 182, 183)
	XSR3_SSE2(ctx.Rk, x1, x0, 184, 185)
	XSR9_SSE2(ctx.Rk, x2, x1, 174, 175)
	XSR5_SSE2(ctx.Rk, x3, x2, 176, 177)
	XSR3_SSE2(ctx.Rk, x0, x3, 178, 179)
	XSR9_SSE2(ctx.Rk, x1, x0, 168, 169)
	XSR5_SSE2(ctx.Rk, x2, x1, 170, 171)
	XSR3_SSE2(ctx.Rk, x3, x2, 172, 173)
	Label("OVER28_END")

	CMPB(ctx.Round, U8(24))
	JBE(LabelRef("OVER24_END"))
	XSR9_SSE2(ctx.Rk, x0, x3, 162, 163)
	XSR5_SSE2(ctx.Rk, x1, x0, 164, 165)
	XSR3_SSE2(ctx.Rk, x2, x1, 166, 167)
	XSR9_SSE2(ctx.Rk, x3, x2, 156, 157)
	XSR5_SSE2(ctx.Rk, x0, x3, 158, 159)
	XSR3_SSE2(ctx.Rk, x1, x0, 160, 161)
	XSR9_SSE2(ctx.Rk, x2, x1, 150, 151)
	XSR5_SSE2(ctx.Rk, x3, x2, 152, 153)
	XSR3_SSE2(ctx.Rk, x0, x3, 154, 155)
	XSR9_SSE2(ctx.Rk, x1, x0, 144, 145)
	XSR5_SSE2(ctx.Rk, x2, x1, 146, 147)
	XSR3_SSE2(ctx.Rk, x3, x2, 148, 149)
	Label("OVER24_END")

	XSR9_SSE2(ctx.Rk, x0, x3, 138, 139)
	XSR5_SSE2(ctx.Rk, x1, x0, 140, 141)
	XSR3_SSE2(ctx.Rk, x2, x1, 142, 143)
	XSR9_SSE2(ctx.Rk, x3, x2, 132, 133)
	XSR5_SSE2(ctx.Rk, x0, x3, 134, 135)
	XSR3_SSE2(ctx.Rk, x1, x0, 136, 137)
	XSR9_SSE2(ctx.Rk, x2, x1, 126, 127)
	XSR5_SSE2(ctx.Rk, x3, x2, 128, 129)
	XSR3_SSE2(ctx.Rk, x0, x3, 130, 131)
	XSR9_SSE2(ctx.Rk, x1, x0, 120, 121)
	XSR5_SSE2(ctx.Rk, x2, x1, 122, 123)
	XSR3_SSE2(ctx.Rk, x3, x2, 124, 125)

	XSR9_SSE2(ctx.Rk, x0, x3, 114, 115)
	XSR5_SSE2(ctx.Rk, x1, x0, 116, 117)
	XSR3_SSE2(ctx.Rk, x2, x1, 118, 119)
	XSR9_SSE2(ctx.Rk, x3, x2, 108, 109)
	XSR5_SSE2(ctx.Rk, x0, x3, 110, 111)
	XSR3_SSE2(ctx.Rk, x1, x0, 112, 113)
	XSR9_SSE2(ctx.Rk, x2, x1, 102, 103)
	XSR5_SSE2(ctx.Rk, x3, x2, 104, 105)
	XSR3_SSE2(ctx.Rk, x0, x3, 106, 107)
	XSR9_SSE2(ctx.Rk, x1, x0, 96, 97)
	XSR5_SSE2(ctx.Rk, x2, x1, 98, 99)
	XSR3_SSE2(ctx.Rk, x3, x2, 100, 101)

	XSR9_SSE2(ctx.Rk, x0, x3, 90, 91)
	XSR5_SSE2(ctx.Rk, x1, x0, 92, 93)
	XSR3_SSE2(ctx.Rk, x2, x1, 94, 95)
	XSR9_SSE2(ctx.Rk, x3, x2, 84, 85)
	XSR5_SSE2(ctx.Rk, x0, x3, 86, 87)
	XSR3_SSE2(ctx.Rk, x1, x0, 88, 89)
	XSR9_SSE2(ctx.Rk, x2, x1, 78, 79)
	XSR5_SSE2(ctx.Rk, x3, x2, 80, 81)
	XSR3_SSE2(ctx.Rk, x0, x3, 82, 83)
	XSR9_SSE2(ctx.Rk, x1, x0, 72, 73)
	XSR5_SSE2(ctx.Rk, x2, x1, 74, 75)
	XSR3_SSE2(ctx.Rk, x3, x2, 76, 77)

	XSR9_SSE2(ctx.Rk, x0, x3, 66, 67)
	XSR5_SSE2(ctx.Rk, x1, x0, 68, 69)
	XSR3_SSE2(ctx.Rk, x2, x1, 70, 71)
	XSR9_SSE2(ctx.Rk, x3, x2, 60, 61)
	XSR5_SSE2(ctx.Rk, x0, x3, 62, 63)
	XSR3_SSE2(ctx.Rk, x1, x0, 64, 65)
	XSR9_SSE2(ctx.Rk, x2, x1, 54, 55)
	XSR5_SSE2(ctx.Rk, x3, x2, 56, 57)
	XSR3_SSE2(ctx.Rk, x0, x3, 58, 59)
	XSR9_SSE2(ctx.Rk, x1, x0, 48, 49)
	XSR5_SSE2(ctx.Rk, x2, x1, 50, 51)
	XSR3_SSE2(ctx.Rk, x3, x2, 52, 53)

	XSR9_SSE2(ctx.Rk, x0, x3, 42, 43)
	XSR5_SSE2(ctx.Rk, x1, x0, 44, 45)
	XSR3_SSE2(ctx.Rk, x2, x1, 46, 47)
	XSR9_SSE2(ctx.Rk, x3, x2, 36, 37)
	XSR5_SSE2(ctx.Rk, x0, x3, 38, 39)
	XSR3_SSE2(ctx.Rk, x1, x0, 40, 41)
	XSR9_SSE2(ctx.Rk, x2, x1, 30, 31)
	XSR5_SSE2(ctx.Rk, x3, x2, 32, 33)
	XSR3_SSE2(ctx.Rk, x0, x3, 34, 35)
	XSR9_SSE2(ctx.Rk, x1, x0, 24, 25)
	XSR5_SSE2(ctx.Rk, x2, x1, 26, 27)
	XSR3_SSE2(ctx.Rk, x3, x2, 28, 29)

	XSR9_SSE2(ctx.Rk, x0, x3, 18, 19)
	XSR5_SSE2(ctx.Rk, x1, x0, 20, 21)
	XSR3_SSE2(ctx.Rk, x2, x1, 22, 23)
	XSR9_SSE2(ctx.Rk, x3, x2, 12, 13)
	XSR5_SSE2(ctx.Rk, x0, x3, 14, 15)
	XSR3_SSE2(ctx.Rk, x1, x0, 16, 17)
	XSR9_SSE2(ctx.Rk, x2, x1, 6, 7)
	XSR5_SSE2(ctx.Rk, x3, x2, 8, 9)
	XSR3_SSE2(ctx.Rk, x0, x3, 10, 11)
	XSR9_SSE2(ctx.Rk, x1, x0, 0, 1)
	XSR5_SSE2(ctx.Rk, x2, x1, 2, 3)
	XSR3_SSE2(ctx.Rk, x3, x2, 4, 5)

	x0, x1, x2, x3 = leaSSE2Swap(x0, x1, x2, x3)
	MOVUPS(x0, dst.Offset(0x00))
	MOVUPS(x1, dst.Offset(0x10))
	MOVUPS(x2, dst.Offset(0x20))
	MOVUPS(x3, dst.Offset(0x30))

	/**
	return;
	*/
	RET()
}

func leaSSE2Swap(x0, x1, x2, x3 VecVirtual) (VecVirtual, VecVirtual, VecVirtual, VecVirtual) {
	/**
	tmp0 = _mm_unpacklo_epi32(x0, x1);
	tmp1 = _mm_unpacklo_epi32(x2, x3);
	tmp2 = _mm_unpackhi_epi32(x0, x1);
	tmp3 = _mm_unpackhi_epi32(x2, x3);

	x0 = _mm_unpacklo_epi64(tmp0, tmp1);
	x1 = _mm_unpackhi_epi64(tmp0, tmp1);
	x2 = _mm_unpacklo_epi64(tmp2, tmp3);
	x3 = _mm_unpackhi_epi64(tmp2, tmp3);

	x86-64 clang 14.0.0
	-O3 -msse2

	movups     xmm0, xmmword ptr [rdi + 2]
	movups     xmm1, xmmword ptr [rdi + 20]
	movups     xmm2, xmmword ptr [rdi + 200]
	movups     xmm3, xmmword ptr [rdi + 2000]

	movaps     xmm4, xmm0
	unpcklps   xmm4, xmm1                      # xmm4 = xmm4[0],xmm1[0],xmm4[1],xmm1[1]
	movaps     xmm5, xmm2
	unpcklps   xmm5, xmm3                      # xmm5 = xmm5[0],xmm3[0],xmm5[1],xmm3[1]
	unpckhps   xmm0, xmm1                      # xmm0 = xmm0[2],xmm1[2],xmm0[3],xmm1[3]
	unpckhps   xmm2, xmm3                      # xmm2 = xmm2[2],xmm3[2],xmm2[3],xmm3[3]

	movaps     xmm1, xmm4
	movlhps    xmm1, xmm5                      # xmm1 = xmm1[0],xmm5[0]
	unpckhpd   xmm4, xmm5                      # xmm4 = xmm4[1],xmm5[1]
	movaps     xmm3, xmm0
	movlhps    xmm3, xmm2                      # xmm3 = xmm3[0],xmm2[0]
	unpckhpd   xmm0, xmm2                      # xmm0 = xmm0[1],xmm2[1]

	movups     xmmword ptr [rdi + 2], xmm1
	movups     xmmword ptr [rdi + 20], xmm4
	movups     xmmword ptr [rdi + 200], xmm3
	movups     xmmword ptr [rdi + 2000], xmm0
	*/
	xmm0 := x0
	xmm1 := x1
	xmm2 := x2
	xmm3 := x3
	xmm4 := XMM()
	xmm5 := XMM()

	MOVAPS(xmm0, xmm4)
	UNPCKLPS(xmm1, xmm4)

	MOVAPS(xmm2, xmm5)
	UNPCKLPS(xmm3, xmm5)

	UNPCKHPS(xmm1, xmm0)

	UNPCKHPS(xmm3, xmm2)

	MOVAPS(xmm4, xmm1)
	MOVLHPS(xmm5, xmm1)

	UNPCKHPD(xmm5, xmm4)

	MOVAPS(xmm0, xmm3)
	MOVLHPS(xmm2, xmm3)

	UNPCKHPD(xmm2, xmm0)

	return xmm1, xmm4, xmm3, xmm0
}

func XAR_SSE2(rk Mem, cur, pre VecVirtual, rk1, rk2, a, b int) {
	/**
	tmp = _mm_add_epi32(
		_mm_xor_si128(
			pre,
			_mm_set1_epi32(rk1)
		),
		_mm_xor_si128(
			cur,
			_mm_set1_epi32(rk2)
		)
	);
	cur = _mm_xor_si128(
		_mm_srli_epi32(tmp, 3),
		_mm_slli_epi32(tmp, 29)
	)
	*/

	tmp0 := XMM()
	tmp1 := XMM()

	{
		{
			MOVD(rk.Offset(4*rk1), tmp0)
			PSHUFD(U8(0), tmp0, tmp0)
		}
		PXOR(pre, tmp0)
	}
	{
		{
			MOVD(rk.Offset(4*rk2), tmp1)
			PSHUFD(U8(0), tmp1, tmp1)
		}
		PXOR(cur, tmp1)
	}
	PADDL(tmp1, tmp0) // paddd

	{
		MOVO(tmp0, cur)   // movdqa
		PSRLL(U8(a), cur) // psrld
	}
	{
		PSLLL(U8(b), tmp0) // pslld
	}
	PXOR(tmp0, cur)
}
func XAR3_SSE2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XAR_SSE2(rk, cur, pre, rk1, rk2, 3, 29)
}
func XAR5_SSE2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XAR_SSE2(rk, cur, pre, rk1, rk2, 5, 27)
}
func XAR9_SSE2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XAR_SSE2(rk, cur, pre, rk1, rk2, 23, 9)
}

func XSR_SSE2(rk Mem, cur, pre VecVirtual, rk1, rk2 int, a, b int) {
	/**
	cur = _mm_xor_si128(
		_mm_sub_epi32(
			_mm_xor_si128(
				_mm_srli_epi32(cur, a),
				_mm_slli_epi32(cur, b)
			),
			_mm_xor_si128(
				pre,
				_mm_set1_epi32(rk1)
			)
		),
		_mm_set1_epi32(rk2)
	);
	*/

	tmp := XMM()

	{

		{
			{
				MOVO(cur, tmp)
				PSRLL(U8(a), tmp) // psrld
			}
			{
				PSLLL(U8(b), cur) // pslld
			}
			PXOR(tmp, cur)
		}
		{
			{
				MOVD(rk.Offset(4*rk1), tmp)
				PSHUFD(U8(0), tmp, tmp)
			}
			PXOR(pre, tmp)
		}
		PSUBL(tmp, cur) // psubd
	}
	{
		MOVD(rk.Offset(4*rk2), tmp)
		PSHUFD(U8(0), tmp, tmp) // psubd
	}
	PXOR(tmp, cur)
}
func XSR9_SSE2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XSR_SSE2(rk, cur, pre, rk1, rk2, 9, 23)
}
func XSR5_SSE2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XSR_SSE2(rk, cur, pre, rk1, rk2, 27, 5)
}
func XSR3_SSE2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XSR_SSE2(rk, cur, pre, rk1, rk2, 29, 3)
}
