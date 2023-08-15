package main

import (
	. "kryptosimd/lea/avo"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

func leaEnc8AVX2() {
	TEXT("leaEnc8AVX2", NOSPLIT, "func(ctx *leaContext, dst []byte, src []byte)")

	ctx := GetCtx()
	dst := Mem{Base: Load(Param("dst").Base(), GP64())}
	src := Mem{Base: Load(Param("src").Base(), GP64())}

	/**
	__m256i x0, x1, x2, x3, tmp;
	__m128i tmp128;
	*/
	x0 := YMM()
	x1 := YMM()
	x2 := YMM()
	x3 := YMM()

	/**
	x0 = _mm256_setr_epi32(
		*((unsigned int *)pt + 0x00), *((unsigned int *)pt + 0x04),
		*((unsigned int *)pt + 0x08), *((unsigned int *)pt + 0x0c),
		*((unsigned int *)pt + 0x10), *((unsigned int *)pt + 0x14),
		*((unsigned int *)pt + 0x18), *((unsigned int *)pt + 0x1c)
	);
	x1 = _mm256_setr_epi32(
		*((unsigned int *)pt + 0x01), *((unsigned int *)pt + 0x05),
		*((unsigned int *)pt + 0x09), *((unsigned int *)pt + 0x0d),
		*((unsigned int *)pt + 0x11), *((unsigned int *)pt + 0x15),
		*((unsigned int *)pt + 0x19), *((unsigned int *)pt + 0x1d)
	);
	x2 = _mm256_setr_epi32(
		*((unsigned int *)pt + 0x02), *((unsigned int *)pt + 0x06),
		*((unsigned int *)pt + 0x0a), *((unsigned int *)pt + 0x0e),
		*((unsigned int *)pt + 0x12), *((unsigned int *)pt + 0x16),
		*((unsigned int *)pt + 0x1a), *((unsigned int *)pt + 0x1e)
	);
	x3 = _mm256_setr_epi32(
		*((unsigned int *)pt + 0x03), *((unsigned int *)pt + 0x07),
		*((unsigned int *)pt + 0x0b), *((unsigned int *)pt + 0x0f),
		*((unsigned int *)pt + 0x13), *((unsigned int *)pt + 0x17),
		*((unsigned int *)pt + 0x1b), *((unsigned int *)pt + 0x1f)
	);
	*/
	leaAVX2Int2Ymm(x0, src, 0x00, 0x04, 0x08, 0x0c, 0x10, 0x14, 0x18, 0x1c)
	leaAVX2Int2Ymm(x1, src, 0x01, 0x05, 0x09, 0x0d, 0x11, 0x15, 0x19, 0x1d)
	leaAVX2Int2Ymm(x2, src, 0x02, 0x06, 0x0a, 0x0e, 0x12, 0x16, 0x1a, 0x1e)
	leaAVX2Int2Ymm(x3, src, 0x03, 0x07, 0x0b, 0x0f, 0x13, 0x17, 0x1b, 0x1f)

	XAR3_AVX2(ctx.Rk, x3, x2, 4, 5)
	XAR5_AVX2(ctx.Rk, x2, x1, 2, 3)
	XAR9_AVX2(ctx.Rk, x1, x0, 0, 1)
	XAR3_AVX2(ctx.Rk, x0, x3, 10, 11)
	XAR5_AVX2(ctx.Rk, x3, x2, 8, 9)
	XAR9_AVX2(ctx.Rk, x2, x1, 6, 7)
	XAR3_AVX2(ctx.Rk, x1, x0, 16, 17)
	XAR5_AVX2(ctx.Rk, x0, x3, 14, 15)
	XAR9_AVX2(ctx.Rk, x3, x2, 12, 13)
	XAR3_AVX2(ctx.Rk, x2, x1, 22, 23)
	XAR5_AVX2(ctx.Rk, x1, x0, 20, 21)
	XAR9_AVX2(ctx.Rk, x0, x3, 18, 19)

	XAR3_AVX2(ctx.Rk, x3, x2, 28, 29)
	XAR5_AVX2(ctx.Rk, x2, x1, 26, 27)
	XAR9_AVX2(ctx.Rk, x1, x0, 24, 25)
	XAR3_AVX2(ctx.Rk, x0, x3, 34, 35)
	XAR5_AVX2(ctx.Rk, x3, x2, 32, 33)
	XAR9_AVX2(ctx.Rk, x2, x1, 30, 31)
	XAR3_AVX2(ctx.Rk, x1, x0, 40, 41)
	XAR5_AVX2(ctx.Rk, x0, x3, 38, 39)
	XAR9_AVX2(ctx.Rk, x3, x2, 36, 37)
	XAR3_AVX2(ctx.Rk, x2, x1, 46, 47)
	XAR5_AVX2(ctx.Rk, x1, x0, 44, 45)
	XAR9_AVX2(ctx.Rk, x0, x3, 42, 43)

	XAR3_AVX2(ctx.Rk, x3, x2, 52, 53)
	XAR5_AVX2(ctx.Rk, x2, x1, 50, 51)
	XAR9_AVX2(ctx.Rk, x1, x0, 48, 49)
	XAR3_AVX2(ctx.Rk, x0, x3, 58, 59)
	XAR5_AVX2(ctx.Rk, x3, x2, 56, 57)
	XAR9_AVX2(ctx.Rk, x2, x1, 54, 55)
	XAR3_AVX2(ctx.Rk, x1, x0, 64, 65)
	XAR5_AVX2(ctx.Rk, x0, x3, 62, 63)
	XAR9_AVX2(ctx.Rk, x3, x2, 60, 61)
	XAR3_AVX2(ctx.Rk, x2, x1, 70, 71)
	XAR5_AVX2(ctx.Rk, x1, x0, 68, 69)
	XAR9_AVX2(ctx.Rk, x0, x3, 66, 67)

	XAR3_AVX2(ctx.Rk, x3, x2, 76, 77)
	XAR5_AVX2(ctx.Rk, x2, x1, 74, 75)
	XAR9_AVX2(ctx.Rk, x1, x0, 72, 73)
	XAR3_AVX2(ctx.Rk, x0, x3, 82, 83)
	XAR5_AVX2(ctx.Rk, x3, x2, 80, 81)
	XAR9_AVX2(ctx.Rk, x2, x1, 78, 79)
	XAR3_AVX2(ctx.Rk, x1, x0, 88, 89)
	XAR5_AVX2(ctx.Rk, x0, x3, 86, 87)
	XAR9_AVX2(ctx.Rk, x3, x2, 84, 85)
	XAR3_AVX2(ctx.Rk, x2, x1, 94, 95)
	XAR5_AVX2(ctx.Rk, x1, x0, 92, 93)
	XAR9_AVX2(ctx.Rk, x0, x3, 90, 91)

	XAR3_AVX2(ctx.Rk, x3, x2, 100, 101)
	XAR5_AVX2(ctx.Rk, x2, x1, 98, 99)
	XAR9_AVX2(ctx.Rk, x1, x0, 96, 97)
	XAR3_AVX2(ctx.Rk, x0, x3, 106, 107)
	XAR5_AVX2(ctx.Rk, x3, x2, 104, 105)
	XAR9_AVX2(ctx.Rk, x2, x1, 102, 103)
	XAR3_AVX2(ctx.Rk, x1, x0, 112, 113)
	XAR5_AVX2(ctx.Rk, x0, x3, 110, 111)
	XAR9_AVX2(ctx.Rk, x3, x2, 108, 109)
	XAR3_AVX2(ctx.Rk, x2, x1, 118, 119)
	XAR5_AVX2(ctx.Rk, x1, x0, 116, 117)
	XAR9_AVX2(ctx.Rk, x0, x3, 114, 115)

	XAR3_AVX2(ctx.Rk, x3, x2, 124, 125)
	XAR5_AVX2(ctx.Rk, x2, x1, 122, 123)
	XAR9_AVX2(ctx.Rk, x1, x0, 120, 121)
	XAR3_AVX2(ctx.Rk, x0, x3, 130, 131)
	XAR5_AVX2(ctx.Rk, x3, x2, 128, 129)
	XAR9_AVX2(ctx.Rk, x2, x1, 126, 127)
	XAR3_AVX2(ctx.Rk, x1, x0, 136, 137)
	XAR5_AVX2(ctx.Rk, x0, x3, 134, 135)
	XAR9_AVX2(ctx.Rk, x3, x2, 132, 133)
	XAR3_AVX2(ctx.Rk, x2, x1, 142, 143)
	XAR5_AVX2(ctx.Rk, x1, x0, 140, 141)
	XAR9_AVX2(ctx.Rk, x0, x3, 138, 139)

	CMPB(ctx.Round, U8(24))
	JBE(LabelRef("OVER24_END"))
	XAR3_AVX2(ctx.Rk, x3, x2, 148, 149)
	XAR5_AVX2(ctx.Rk, x2, x1, 146, 147)
	XAR9_AVX2(ctx.Rk, x1, x0, 144, 145)
	XAR3_AVX2(ctx.Rk, x0, x3, 154, 155)
	XAR5_AVX2(ctx.Rk, x3, x2, 152, 153)
	XAR9_AVX2(ctx.Rk, x2, x1, 150, 151)
	XAR3_AVX2(ctx.Rk, x1, x0, 160, 161)
	XAR5_AVX2(ctx.Rk, x0, x3, 158, 159)
	XAR9_AVX2(ctx.Rk, x3, x2, 156, 157)
	XAR3_AVX2(ctx.Rk, x2, x1, 166, 167)
	XAR5_AVX2(ctx.Rk, x1, x0, 164, 165)
	XAR9_AVX2(ctx.Rk, x0, x3, 162, 163)
	Label("OVER24_END")

	CMPB(ctx.Round, U8(28))
	JBE(LabelRef("OVER28_END"))
	XAR3_AVX2(ctx.Rk, x3, x2, 172, 173)
	XAR5_AVX2(ctx.Rk, x2, x1, 170, 171)
	XAR9_AVX2(ctx.Rk, x1, x0, 168, 169)
	XAR3_AVX2(ctx.Rk, x0, x3, 178, 179)
	XAR5_AVX2(ctx.Rk, x3, x2, 176, 177)
	XAR9_AVX2(ctx.Rk, x2, x1, 174, 175)
	XAR3_AVX2(ctx.Rk, x1, x0, 184, 185)
	XAR5_AVX2(ctx.Rk, x0, x3, 182, 183)
	XAR9_AVX2(ctx.Rk, x3, x2, 180, 181)
	XAR3_AVX2(ctx.Rk, x2, x1, 190, 191)
	XAR5_AVX2(ctx.Rk, x1, x0, 188, 189)
	XAR9_AVX2(ctx.Rk, x0, x3, 186, 187)
	Label("OVER28_END")

	/**
	tmp128 = _mm256_extractf128_si256(x0, 0);
	*((unsigned int *)ct + 0x00) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)ct + 0x04) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)ct + 0x08) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)ct + 0x0c) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x0, 1);
	*((unsigned int *)ct + 0x10) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)ct + 0x14) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)ct + 0x18) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)ct + 0x1c) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x1, 0);
	*((unsigned int *)ct + 0x01) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)ct + 0x05) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)ct + 0x09) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)ct + 0x0d) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x1, 1);
	*((unsigned int *)ct + 0x11) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)ct + 0x15) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)ct + 0x19) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)ct + 0x1d) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x2, 0);
	*((unsigned int *)ct + 0x02) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)ct + 0x06) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)ct + 0x0a) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)ct + 0x0e) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x2, 1);
	*((unsigned int *)ct + 0x12) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)ct + 0x16) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)ct + 0x1a) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)ct + 0x1e) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x3, 0);
	*((unsigned int *)ct + 0x03) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)ct + 0x07) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)ct + 0x0b) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)ct + 0x0f) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x3, 1);
	*((unsigned int *)ct + 0x13) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)ct + 0x17) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)ct + 0x1b) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)ct + 0x1f) = _mm_extract_epi32(tmp128, 3);
	*/
	leaAVX2Ymm2Int(x0, dst, 0x00, 0x04, 0x08, 0x0c, 0x10, 0x14, 0x18, 0x1c)
	leaAVX2Ymm2Int(x1, dst, 0x01, 0x05, 0x09, 0x0d, 0x11, 0x15, 0x19, 0x1d)
	leaAVX2Ymm2Int(x2, dst, 0x02, 0x06, 0x0a, 0x0e, 0x12, 0x16, 0x1a, 0x1e)
	leaAVX2Ymm2Int(x3, dst, 0x03, 0x07, 0x0b, 0x0f, 0x13, 0x17, 0x1b, 0x1f)

	/**
	return
	*/
	RET()
}

func leaDec8AVX2() {
	TEXT("leaDec8AVX2", NOSPLIT, "func(ctx *leaContext, dst []byte, src []byte)")

	ctx := GetCtx()
	dst := Mem{Base: Load(Param("dst").Base(), GP64())}
	src := Mem{Base: Load(Param("src").Base(), GP64())}

	/**
	__m256i x0, x1, x2, x3;
	__m128i tmp128;
	*/
	x0 := YMM()
	x1 := YMM()
	x2 := YMM()
	x3 := YMM()

	/**
	x0 = _mm256_setr_epi32(
		*((unsigned int *)ct + 0x00), *((unsigned int *)ct + 0x04),
		*((unsigned int *)ct + 0x08), *((unsigned int *)ct + 0x0c),
		*((unsigned int *)ct + 0x10), *((unsigned int *)ct + 0x14),
		*((unsigned int *)ct + 0x18), *((unsigned int *)ct + 0x1c)
	);
	x1 = _mm256_setr_epi32(
		*((unsigned int *)ct + 0x01), *((unsigned int *)ct + 0x05),
		*((unsigned int *)ct + 0x09), *((unsigned int *)ct + 0x0d),
		*((unsigned int *)ct + 0x11), *((unsigned int *)ct + 0x15),
		*((unsigned int *)ct + 0x19), *((unsigned int *)ct + 0x1d)
	);
	x2 = _mm256_setr_epi32(
		*((unsigned int *)ct + 0x02), *((unsigned int *)ct + 0x06),
		*((unsigned int *)ct + 0x0a), *((unsigned int *)ct + 0x0e),
		*((unsigned int *)ct + 0x12), *((unsigned int *)ct + 0x16),
		*((unsigned int *)ct + 0x1a), *((unsigned int *)ct + 0x1e)
	);
	x3 = _mm256_setr_epi32(
		*((unsigned int *)ct + 0x03), *((unsigned int *)ct + 0x07),
		*((unsigned int *)ct + 0x0b), *((unsigned int *)ct + 0x0f),
		*((unsigned int *)ct + 0x13), *((unsigned int *)ct + 0x17),
		*((unsigned int *)ct + 0x1b), *((unsigned int *)ct + 0x1f)
	);
	*/
	leaAVX2Int2Ymm(x0, src, 0x00, 0x04, 0x08, 0x0c, 0x10, 0x14, 0x18, 0x1c)
	leaAVX2Int2Ymm(x1, src, 0x01, 0x05, 0x09, 0x0d, 0x11, 0x15, 0x19, 0x1d)
	leaAVX2Int2Ymm(x2, src, 0x02, 0x06, 0x0a, 0x0e, 0x12, 0x16, 0x1a, 0x1e)
	leaAVX2Int2Ymm(x3, src, 0x03, 0x07, 0x0b, 0x0f, 0x13, 0x17, 0x1b, 0x1f)

	CMPB(ctx.Round, U8(28))
	JBE(LabelRef("OVER28_END"))
	XSR9_AVX2(ctx.Rk, x0, x3, 186, 187)
	XSR5_AVX2(ctx.Rk, x1, x0, 188, 189)
	XSR3_AVX2(ctx.Rk, x2, x1, 190, 191)
	XSR9_AVX2(ctx.Rk, x3, x2, 180, 181)
	XSR5_AVX2(ctx.Rk, x0, x3, 182, 183)
	XSR3_AVX2(ctx.Rk, x1, x0, 184, 185)
	XSR9_AVX2(ctx.Rk, x2, x1, 174, 175)
	XSR5_AVX2(ctx.Rk, x3, x2, 176, 177)
	XSR3_AVX2(ctx.Rk, x0, x3, 178, 179)
	XSR9_AVX2(ctx.Rk, x1, x0, 168, 169)
	XSR5_AVX2(ctx.Rk, x2, x1, 170, 171)
	XSR3_AVX2(ctx.Rk, x3, x2, 172, 173)
	Label("OVER28_END")

	CMPB(ctx.Round, U8(24))
	JBE(LabelRef("OVER24_END"))
	XSR9_AVX2(ctx.Rk, x0, x3, 162, 163)
	XSR5_AVX2(ctx.Rk, x1, x0, 164, 165)
	XSR3_AVX2(ctx.Rk, x2, x1, 166, 167)
	XSR9_AVX2(ctx.Rk, x3, x2, 156, 157)
	XSR5_AVX2(ctx.Rk, x0, x3, 158, 159)
	XSR3_AVX2(ctx.Rk, x1, x0, 160, 161)
	XSR9_AVX2(ctx.Rk, x2, x1, 150, 151)
	XSR5_AVX2(ctx.Rk, x3, x2, 152, 153)
	XSR3_AVX2(ctx.Rk, x0, x3, 154, 155)
	XSR9_AVX2(ctx.Rk, x1, x0, 144, 145)
	XSR5_AVX2(ctx.Rk, x2, x1, 146, 147)
	XSR3_AVX2(ctx.Rk, x3, x2, 148, 149)
	Label("OVER24_END")

	XSR9_AVX2(ctx.Rk, x0, x3, 138, 139)
	XSR5_AVX2(ctx.Rk, x1, x0, 140, 141)
	XSR3_AVX2(ctx.Rk, x2, x1, 142, 143)
	XSR9_AVX2(ctx.Rk, x3, x2, 132, 133)
	XSR5_AVX2(ctx.Rk, x0, x3, 134, 135)
	XSR3_AVX2(ctx.Rk, x1, x0, 136, 137)
	XSR9_AVX2(ctx.Rk, x2, x1, 126, 127)
	XSR5_AVX2(ctx.Rk, x3, x2, 128, 129)
	XSR3_AVX2(ctx.Rk, x0, x3, 130, 131)
	XSR9_AVX2(ctx.Rk, x1, x0, 120, 121)
	XSR5_AVX2(ctx.Rk, x2, x1, 122, 123)
	XSR3_AVX2(ctx.Rk, x3, x2, 124, 125)

	XSR9_AVX2(ctx.Rk, x0, x3, 114, 115)
	XSR5_AVX2(ctx.Rk, x1, x0, 116, 117)
	XSR3_AVX2(ctx.Rk, x2, x1, 118, 119)
	XSR9_AVX2(ctx.Rk, x3, x2, 108, 109)
	XSR5_AVX2(ctx.Rk, x0, x3, 110, 111)
	XSR3_AVX2(ctx.Rk, x1, x0, 112, 113)
	XSR9_AVX2(ctx.Rk, x2, x1, 102, 103)
	XSR5_AVX2(ctx.Rk, x3, x2, 104, 105)
	XSR3_AVX2(ctx.Rk, x0, x3, 106, 107)
	XSR9_AVX2(ctx.Rk, x1, x0, 96, 97)
	XSR5_AVX2(ctx.Rk, x2, x1, 98, 99)
	XSR3_AVX2(ctx.Rk, x3, x2, 100, 101)

	XSR9_AVX2(ctx.Rk, x0, x3, 90, 91)
	XSR5_AVX2(ctx.Rk, x1, x0, 92, 93)
	XSR3_AVX2(ctx.Rk, x2, x1, 94, 95)
	XSR9_AVX2(ctx.Rk, x3, x2, 84, 85)
	XSR5_AVX2(ctx.Rk, x0, x3, 86, 87)
	XSR3_AVX2(ctx.Rk, x1, x0, 88, 89)
	XSR9_AVX2(ctx.Rk, x2, x1, 78, 79)
	XSR5_AVX2(ctx.Rk, x3, x2, 80, 81)
	XSR3_AVX2(ctx.Rk, x0, x3, 82, 83)
	XSR9_AVX2(ctx.Rk, x1, x0, 72, 73)
	XSR5_AVX2(ctx.Rk, x2, x1, 74, 75)
	XSR3_AVX2(ctx.Rk, x3, x2, 76, 77)

	XSR9_AVX2(ctx.Rk, x0, x3, 66, 67)
	XSR5_AVX2(ctx.Rk, x1, x0, 68, 69)
	XSR3_AVX2(ctx.Rk, x2, x1, 70, 71)
	XSR9_AVX2(ctx.Rk, x3, x2, 60, 61)
	XSR5_AVX2(ctx.Rk, x0, x3, 62, 63)
	XSR3_AVX2(ctx.Rk, x1, x0, 64, 65)
	XSR9_AVX2(ctx.Rk, x2, x1, 54, 55)
	XSR5_AVX2(ctx.Rk, x3, x2, 56, 57)
	XSR3_AVX2(ctx.Rk, x0, x3, 58, 59)
	XSR9_AVX2(ctx.Rk, x1, x0, 48, 49)
	XSR5_AVX2(ctx.Rk, x2, x1, 50, 51)
	XSR3_AVX2(ctx.Rk, x3, x2, 52, 53)

	XSR9_AVX2(ctx.Rk, x0, x3, 42, 43)
	XSR5_AVX2(ctx.Rk, x1, x0, 44, 45)
	XSR3_AVX2(ctx.Rk, x2, x1, 46, 47)
	XSR9_AVX2(ctx.Rk, x3, x2, 36, 37)
	XSR5_AVX2(ctx.Rk, x0, x3, 38, 39)
	XSR3_AVX2(ctx.Rk, x1, x0, 40, 41)
	XSR9_AVX2(ctx.Rk, x2, x1, 30, 31)
	XSR5_AVX2(ctx.Rk, x3, x2, 32, 33)
	XSR3_AVX2(ctx.Rk, x0, x3, 34, 35)
	XSR9_AVX2(ctx.Rk, x1, x0, 24, 25)
	XSR5_AVX2(ctx.Rk, x2, x1, 26, 27)
	XSR3_AVX2(ctx.Rk, x3, x2, 28, 29)

	XSR9_AVX2(ctx.Rk, x0, x3, 18, 19)
	XSR5_AVX2(ctx.Rk, x1, x0, 20, 21)
	XSR3_AVX2(ctx.Rk, x2, x1, 22, 23)
	XSR9_AVX2(ctx.Rk, x3, x2, 12, 13)
	XSR5_AVX2(ctx.Rk, x0, x3, 14, 15)
	XSR3_AVX2(ctx.Rk, x1, x0, 16, 17)
	XSR9_AVX2(ctx.Rk, x2, x1, 6, 7)
	XSR5_AVX2(ctx.Rk, x3, x2, 8, 9)
	XSR3_AVX2(ctx.Rk, x0, x3, 10, 11)
	XSR9_AVX2(ctx.Rk, x1, x0, 0, 1)
	XSR5_AVX2(ctx.Rk, x2, x1, 2, 3)
	XSR3_AVX2(ctx.Rk, x3, x2, 4, 5)

	/**
	tmp128 = _mm256_extractf128_si256(x0, 0);
	*((unsigned int *)pt + 0x00) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)pt + 0x04) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)pt + 0x08) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)pt + 0x0c) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x0, 1);
	*((unsigned int *)pt + 0x10) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)pt + 0x14) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)pt + 0x18) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)pt + 0x1c) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x1, 0);
	*((unsigned int *)pt + 0x01) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)pt + 0x05) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)pt + 0x09) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)pt + 0x0d) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x1, 1);
	*((unsigned int *)pt + 0x11) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)pt + 0x15) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)pt + 0x19) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)pt + 0x1d) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x2, 0);
	*((unsigned int *)pt + 0x02) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)pt + 0x06) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)pt + 0x0a) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)pt + 0x0e) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x2, 1);
	*((unsigned int *)pt + 0x12) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)pt + 0x16) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)pt + 0x1a) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)pt + 0x1e) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x3, 0);
	*((unsigned int *)pt + 0x03) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)pt + 0x07) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)pt + 0x0b) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)pt + 0x0f) = _mm_extract_epi32(tmp128, 3);
	tmp128 = _mm256_extractf128_si256(x3, 1);
	*((unsigned int *)pt + 0x13) = _mm_extract_epi32(tmp128, 0); *((unsigned int *)pt + 0x17) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)pt + 0x1b) = _mm_extract_epi32(tmp128, 2); *((unsigned int *)pt + 0x1f) = _mm_extract_epi32(tmp128, 3);
	*/
	leaAVX2Ymm2Int(x0, dst, 0x00, 0x04, 0x08, 0x0c, 0x10, 0x14, 0x18, 0x1c)
	leaAVX2Ymm2Int(x1, dst, 0x01, 0x05, 0x09, 0x0d, 0x11, 0x15, 0x19, 0x1d)
	leaAVX2Ymm2Int(x2, dst, 0x02, 0x06, 0x0a, 0x0e, 0x12, 0x16, 0x1a, 0x1e)
	leaAVX2Ymm2Int(x3, dst, 0x03, 0x07, 0x0b, 0x0f, 0x13, 0x17, 0x1b, 0x1f)

	/**
	return
	*/
	RET()
}

func leaAVX2Int2Ymm(dst VecVirtual, src Mem, r0, r1, r2, r3, r4, r5, r6, r7 int) {
	/**
	vmovd        xmm1,       dword ptr [rsi      ]    # xmm1 = mem[0],zero,zero,zero
	vpinsrd      xmm1, xmm1, dword ptr [rsi +  16], 1
	vpinsrd      xmm1, xmm1, dword ptr [rsi +  32], 2
	vpinsrd      xmm1, xmm1, dword ptr [rsi +  48], 3

	vmovd        xmm0,       dword ptr [rsi +  64]    # xmm0 = mem[0],zero,zero,zero
	vpinsrd      xmm0, xmm0, dword ptr [rsi +  80], 1
	vpinsrd      xmm0, xmm0, dword ptr [rsi +  96], 2
	vpinsrd      xmm0, xmm0, dword ptr [rsi + 112], 3

	vinserti128  ymm0, ymm1, xmm0, 1
	*/

	ymm0 := dst
	xmm0 := dst.AsX()

	ymm1 := YMM()
	xmm1 := ymm1.AsX()

	VMOVD(src.Offset(4*r0), xmm0)
	VPINSRD(U8(1), src.Offset(4*r1), xmm0, xmm0)
	VPINSRD(U8(2), src.Offset(4*r2), xmm0, xmm0)
	VPINSRD(U8(3), src.Offset(4*r3), xmm0, xmm0)

	VMOVD(src.Offset(4*r4), xmm1)
	VPINSRD(U8(1), src.Offset(4*r5), xmm1, xmm1)
	VPINSRD(U8(2), src.Offset(4*r6), xmm1, xmm1)
	VPINSRD(U8(3), src.Offset(4*r7), xmm1, xmm1)

	VINSERTI128(U8(1), xmm1, ymm0, ymm0)
}

func leaAVX2Ymm2Int(y VecVirtual, dst Mem, d0, d1, d2, d3, d4, d5, d6, d7 int) {
	/**
	tmp128 = _mm256_extractf128_si256(x0, 0);
	*((unsigned int *)ct) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)ct + 0x04) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)ct + 0x08) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)ct + 0x0c) = _mm_extract_epi32(tmp128, 3);

	tmp128 = _mm256_extractf128_si256(x0, 1);
	*((unsigned int *)ct + 0x10) = _mm_extract_epi32(tmp128, 0);
	*((unsigned int *)ct + 0x14) = _mm_extract_epi32(tmp128, 1);
	*((unsigned int *)ct + 0x18) = _mm_extract_epi32(tmp128, 2);
	*((unsigned int *)ct + 0x1c) = _mm_extract_epi32(tmp128, 3);

	    0                                           128                                         256
	xmm | 00000000 | 00000000 | 00000000 | 00000000 |                                           |
	ymm | 00000000 | 00000000 | 00000000 | 00000000 | 00000000 | 00000000 | 00000000 | 00000000 |
	    |    x0         x1         x2         x3    |    x4         x5         x6         x7    |
	*/

	ymm0 := y
	xmm0 := y.AsX()

	VEXTRACTPS(U8(0), xmm0, dst.Offset(4*d0))
	VEXTRACTPS(U8(1), xmm0, dst.Offset(4*d1))
	VEXTRACTPS(U8(2), xmm0, dst.Offset(4*d2))
	VEXTRACTPS(U8(3), xmm0, dst.Offset(4*d3))

	VEXTRACTF128(U8(1), ymm0, xmm0)
	VEXTRACTPS(U8(0), xmm0, dst.Offset(4*d4))
	VEXTRACTPS(U8(1), xmm0, dst.Offset(4*d5))
	VEXTRACTPS(U8(2), xmm0, dst.Offset(4*d6))
	VEXTRACTPS(U8(3), xmm0, dst.Offset(4*d7))
}

func XAR_AVX2(rk Mem, cur, pre VecVirtual, rk1, rk2, a, b int) {
	/**
	tmp = _mm256_add_epi32(              ----> tmp0
		_mm256_xor_si256(                ------> tmp0
			pre,
			_mm256_set1_epi32(ctx.Rk1)       --------> tmp0
		),
		_mm256_xor_si256(                ------> tmp1
			cur,
			_mm256_set1_epi32(ctx.Rk2)       --------> tmp1
		)
	);
	cur = _mm256_xor_si256(              ----> cur
		_mm256_srli_epi32(tmp, a),       ------> cur
		_mm256_slli_epi32(tmp, b)        ------> tmp0
	);
	*/

	tmp0 := YMM()
	tmp1 := YMM()

	VPBROADCASTD(rk.Offset(4*rk1), tmp0)
	VPBROADCASTD(rk.Offset(4*rk2), tmp1)

	VPXOR(pre, tmp0, tmp0)
	VPXOR(cur, tmp1, tmp1)
	VPADDD(tmp1, tmp0, tmp0)

	VPSRLD(U8(a), tmp0, cur)
	VPSLLD(U8(b), tmp0, tmp0)
	VPXOR(tmp0, cur, cur)
}
func XAR3_AVX2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XAR_AVX2(rk, cur, pre, rk1, rk2, 3, 29)
}
func XAR5_AVX2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XAR_AVX2(rk, cur, pre, rk1, rk2, 5, 27)
}
func XAR9_AVX2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XAR_AVX2(rk, cur, pre, rk1, rk2, 23, 9)
}

func XSR_AVX2(rk Mem, cur, pre VecVirtual, rk1, rk2, a, b int) {
	/**
	cur = _mm256_xor_si256(                ----> cur
		_mm256_sub_epi32(                  ------> cur
			_mm256_xor_si256(              --------> cur
				_mm256_srli_epi32(cur, a), ----------> tmp
				_mm256_slli_epi32(cur, b)  ----------> cur
			),
			_mm256_xor_si256(              --------> tmp
				pre,
				_mm256_set1_epi32(ctx.Rk1)     ----------> tmp
			)
		),
		_mm256_set1_epi32(ctx.Rk2)             ------> tmp
	);
	*/

	tmp := YMM()

	VPSRLD(U8(a), cur, tmp)

	VPSLLD(U8(b), cur, cur)
	VPXOR(tmp, cur, cur)

	VPBROADCASTD(rk.Offset(4*rk1), tmp)
	VPXOR(pre, tmp, tmp)

	VPSUBD(tmp, cur, cur)

	VPBROADCASTD(rk.Offset(4*rk2), tmp)

	VPXOR(tmp, cur, cur)
}

func XSR9_AVX2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XSR_AVX2(rk, cur, pre, rk1, rk2, 9, 23)
}
func XSR5_AVX2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XSR_AVX2(rk, cur, pre, rk1, rk2, 27, 5)
}
func XSR3_AVX2(rk Mem, cur, pre VecVirtual, rk1, rk2 int) {
	XSR_AVX2(rk, cur, pre, rk1, rk2, 29, 3)
}
