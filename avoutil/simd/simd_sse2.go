package simd

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"

	. "kryptosimd/avoutil"
)

/*
Synopsis

	__m128i _mm_loadu_si128 (__m128i const* mem_addr)
	#include <emmintrin.h>
	Instruction: movdqu xmm, m128
	CPUID Flags: SSE2

Description

	Load 128-bits of integer data from memory into dst. mem_addr does not need to be aligned on any particular boundary.

Operation

	dst[127:0] := MEM[mem_addr+127:mem_addr]
*/
func F_mm_loadu_si128(dst VecVirtual, src Op) VecVirtual {
	MOVO_autoAU2(dst, src)
	return dst
}

/*
Synopsis

	void _mm_storeu_si128 (__m128i* mem_addr, __m128i a)
	#include <emmintrin.h>
	Instruction: movdqu m128, xmm
	CPUID Flags: SSE2

Description

	Store 128-bits of integer data from a into memory. mem_addr does not need to be aligned on any particular boundary.

Operation

	MEM[mem_addr+127:mem_addr] := a[127:0]
*/
func F_mm_storeu_si128(dst, src Op) {
	MOVO_autoAU2(dst, src)
}

/*
Synopsis

	__m128i _mm_xor_si128 (__m128i a, __m128i b)
	#include <emmintrin.h>
	Instruction: pxor xmm, xmm
	CPUID Flags: SSE2

Description

	Compute the bitwise XOR of 128 bits (representing integer data) in a and b, and store the result in dst.

Operation

	dst[127:0] := (a[127:0] XOR b[127:0])
*/
func F_mm_xor_si128(dst VecVirtual, a, b Op) VecVirtual {
	switch {
	case dst == a:
		CheckType(
			`
			//	PXOR m128 xmm
			//	PXOR xmm  xmm
			`,
			b, dst,
		)
		PXOR(b, dst)
	case dst == b:
		CheckType(
			`
			//	PXOR m128 xmm
			//	PXOR xmm  xmm
			`,
			b, dst,
		)
		PXOR(a, dst)
	default:
		CheckType(
			`
			//	PXOR m128 xmm
			//	PXOR xmm  xmm
			`,
			b, dst,
		)

		MOVO_autoAU2(dst, a)
		PXOR(b, dst)
	}

	return dst
}

/*
Synopsis

	__m128i _mm_or_si128 (__m128i a, __m128i b)
	#include <emmintrin.h>
	Instruction: por xmm, xmm
	CPUID Flags: SSE2

Description

	Compute the bitwise OR of 128 bits (representing integer data) in a and b, and store the result in dst.

Operation

	dst[127:0] := (a[127:0] OR b[127:0])
*/
func F_mm_or_si128(dst VecVirtual, a, b Op) VecVirtual {
	switch {
	case dst == a:
		CheckType(
			`
			//	POR m128 xmm
			//	POR xmm  xmm
			`,
			b, dst,
		)
		POR(b, dst)
	case dst == b:
		CheckType(
			`
			//	POR m128 xmm
			//	POR xmm  xmm
			`,
			b, dst,
		)
		POR(a, dst)
	default:
		CheckType(
			`
			//	POR m128 xmm
			//	POR xmm  xmm
			`,
			b, dst,
		)

		MOVO_autoAU2(dst, a)
		POR(b, dst)
	}
	return dst
}

/*
Synopsis

	__m128i _mm_and_si128 (__m128i a, __m128i b)
	#include <emmintrin.h>
	Instruction: pand xmm, xmm
	CPUID Flags: SSE2

Description

	Compute the bitwise AND of 128 bits (representing integer data) in a and b, and store the result in dst.

Operation

	dst[127:0] := (a[127:0] AND b[127:0])
*/
func F_mm_and_si128(dst VecVirtual, a, b Op) VecVirtual {
	switch {
	case dst == a:
		CheckType(
			`
			//	PAND m128 xmm
			//	PAND xmm  xmm
			`,
			b, dst,
		)
		PAND(b, dst)
	case dst == b:
		CheckType(
			`
			//	PAND m128 xmm
			//	PAND xmm  xmm
			`,
			b, dst,
		)
		PAND(a, dst)
	default:
		CheckType(
			`
			//	PAND m128 xmm
			//	PAND xmm  xmm
			`,
			b, dst,
		)

		MOVO_autoAU2(dst, a)
		PAND(b, dst)
	}
	return dst
}

/*
Synopsis

	__m128i _mm_add_epi32 (__m128i a, __m128i b)
	#include <emmintrin.h>
	Instruction: paddd xmm, xmm
	CPUID Flags: SSE2

Description

	Add packed 32-bit integers in a and b, and store the results in dst.

Operation

	FOR j := 0 to 3
		i := j*32
		dst[i+31:i] := a[i+31:i] + b[i+31:i]
	ENDFOR
*/
func F_mm_add_epi32(dst VecVirtual, a, b Op) VecVirtual {
	switch {
	case dst == a:
		CheckType(
			`
			//	PADDD m128 xmm
			//	PADDD xmm  xmm
			`,
			b, dst,
		)
		PADDD(b, dst)
	case dst == b:
		CheckType(
			`
			//	PADDD m128 xmm
			//	PADDD xmm  xmm
			`,
			b, dst,
		)
		PADDD(a, dst)
	default:
		CheckType(
			`
			//	PADDD m128 xmm
			//	PADDD xmm  xmm
			`,
			b, dst,
		)

		MOVO_autoAU2(dst, a)
		PADDD(b, dst)
	}
	return dst
}

/*
Synopsis

	__m128i _mm_slli_epi32 (__m128i a, int imm8)
	#include <emmintrin.h>
	Instruction: pslld xmm, imm8
	CPUID Flags: SSE2

Description

	Shift packed 32-bit integers in a left by imm8 while shifting in zeros, and store the results in dst.

Operation

	FOR j := 0 to 3
		i := j*32
		IF imm8[7:0] > 31
			dst[i+31:i] := 0
		ELSE
			dst[i+31:i] := ZeroExtend32(a[i+31:i] << imm8[7:0])
		FI
	ENDFOR
*/
func F_mm_slli_epi32(dst VecVirtual, a, imm8 Op) VecVirtual {
	if dst != a {
		MOVO_autoAU2(dst, a)
	}

	CheckType(
		`
		//	PSLLL imm8 xmm
		//	PSLLL m128 xmm
		//	PSLLL xmm  xmm
		`,
		imm8, dst,
	)

	PSLLL(imm8, dst)
	return dst
}

/*
Synopsis

	__m128i _mm_srli_epi32 (__m128i a, int imm8)
	#include <emmintrin.h>
	Instruction: psrld xmm, imm8
	CPUID Flags: SSE2

Description

	Shift packed 32-bit integers in a right by imm8 while shifting in zeros, and store the results in dst.

Operation

	FOR j := 0 to 3
		i := j*32
		IF imm8[7:0] > 31
			dst[i+31:i] := 0
		ELSE
			dst[i+31:i] := ZeroExtend32(a[i+31:i] >> imm8[7:0])
		FI
	ENDFOR
*/
func F_mm_srli_epi32(dst VecVirtual, a, imm8 Op) VecVirtual {
	if dst != a {
		MOVO_autoAU2(dst, a)
	}

	CheckType(
		`
		//	PSRLL imm8 xmm
		//	PSRLL m128 xmm
		//	PSRLL xmm  xmm
		`,
		imm8, dst,
	)

	PSRLL(imm8, dst)
	return dst
}

/*
Synopsis

	__m128i _mm_shuffle_epi32 (__m128i a, int imm8)
	#include <emmintrin.h>
	Instruction: pshufd xmm, xmm, imm8
	CPUID Flags: SSE2

Description

	Shuffle 32-bit integers in a using the control in imm8, and store the results in dst.

Operation

	DEFINE SELECT4(src, control) {
		CASE(control[1:0]) OF
		0:	tmp[31:0] := src[31:0]
		1:	tmp[31:0] := src[63:32]
		2:	tmp[31:0] := src[95:64]
		3:	tmp[31:0] := src[127:96]
		ESAC
		RETURN tmp[31:0]
	}
	dst[31:0] := SELECT4(a[127:0], imm8[1:0])
	dst[63:32] := SELECT4(a[127:0], imm8[3:2])
	dst[95:64] := SELECT4(a[127:0], imm8[5:4])
	dst[127:96] := SELECT4(a[127:0], imm8[7:6])
*/
func F_mm_shuffle_epi32(dst VecVirtual, a, imm8 Op) VecVirtual {
	if dst != a {
		MOVO_autoAU2(dst, a)
	}

	CheckType(
		`
		//	PSHUFD imm8 m128 xmm
		//	PSHUFD imm8 xmm  xmm
		`,
		imm8, a, dst,
	)

	PSHUFD(imm8, a, dst)
	return dst
}

/*
*
Synopsis

	__m128i _mm_unpacklo_epi64 (__m128i a, __m128i b)
	#include <emmintrin.h>
	Instruction: punpcklqdq xmm, xmm
	CPUID Flags: SSE2

Description

	Unpack and interleave 64-bit integers from the low half of a and b, and store the results in dst.

Operation

	DEFINE INTERLEAVE_QWORDS(src1[127:0], src2[127:0]) {
		dst[63:0] := src1[63:0]
		dst[127:64] := src2[63:0]
		RETURN dst[127:0]
	}
	dst[127:0] := INTERLEAVE_QWORDS(a[127:0], b[127:0])
*/
func F_mm_unpacklo_epi64(dst VecVirtual, a, b Op) VecVirtual {
	switch {
	case dst == a:
		CheckType(
			`
			//	PUNPCKLQDQ m128 xmm
			//	PUNPCKLQDQ xmm  xmm
			`,
			b, dst,
		)

		PUNPCKLQDQ(b, dst)

	case dst == b:
		CheckType(
			`
			//	PUNPCKLQDQ m128 xmm
			//	PUNPCKLQDQ xmm  xmm
			`,
			a, dst,
		)

		tmp := XMM()
		MOVO_autoAU2(tmp, b)
		MOVO_autoAU2(dst, a)
		PUNPCKLQDQ(tmp, dst)

	default:
		CheckType(
			`
			//	PUNPCKLQDQ m128 xmm
			//	PUNPCKLQDQ xmm  xmm
			`,
			b, dst,
		)

		MOVO_autoAU2(dst, a)
		PUNPCKLQDQ(b, dst)
	}

	return dst
}

/*
*
Synopsis

	__m128i _mm_unpackhi_epi64 (__m128i a, __m128i b)
	#include <emmintrin.h>
	Instruction: punpckhqdq xmm, xmm
	CPUID Flags: SSE2

Description

	Unpack and interleave 64-bit integers from the high half of a and b, and store the results in dst.

Operation

	DEFINE INTERLEAVE_HIGH_QWORDS(src1[127:0], src2[127:0]) {
		dst[63:0] := src1[127:64]
		dst[127:64] := src2[127:64]
		RETURN dst[127:0]
	}
	dst[127:0] := INTERLEAVE_HIGH_QWORDS(a[127:0], b[127:0])
*/
func F_mm_unpackhi_epi64(dst VecVirtual, a, b Op) VecVirtual {
	switch {
	case dst == a:
		CheckType(
			`
			//	PUNPCKHQDQ m128 xmm
			//	PUNPCKHQDQ xmm  xmm
			`,
			b, dst,
		)

		PUNPCKHQDQ(b, dst)

	case dst == b:
		CheckType(
			`
			//	PUNPCKHQDQ m128 xmm
			//	PUNPCKHQDQ xmm  xmm
			`,
			a, dst,
		)

		tmp := XMM()
		MOVO_autoAU2(tmp, b)
		MOVO_autoAU2(dst, a)
		PUNPCKHQDQ(tmp, dst)

	default:
		CheckType(
			`
			//	PUNPCKHQDQ m128 xmm
			//	PUNPCKHQDQ xmm  xmm
			`,
			b, dst,
		)

		MOVO_autoAU2(dst, a)
		PUNPCKHQDQ(b, dst)
	}

	return dst
}

/*
*
Synopsis

	__m128i _mm_add_epi64 (__m128i a, __m128i b)
	#include <emmintrin.h>
	Instruction: paddq xmm, xmm
	CPUID Flags: SSE2

Description

	Add packed 64-bit integers in a and b, and store the results in dst.

Operation

	FOR j := 0 to 1
		i := j*64
		dst[i+63:i] := a[i+63:i] + b[i+63:i]
	ENDFOR
*/
func F_mm_add_epi64(dst VecVirtual, a, b Op) VecVirtual {
	switch {
	case dst == a:
		CheckType(
			`
			//	PADDQ m128 xmm
			//	PADDQ xmm  xmm
			`,
			b, dst,
		)

		PADDQ(b, dst)

	case dst == b:
		CheckType(
			`
			//	PADDQ m128 xmm
			//	PADDQ xmm  xmm
			`,
			a, dst,
		)
		PADDQ(a, dst)

	default:
		CheckType(
			`
			//	PADDQ m128 xmm
			//	PADDQ xmm  xmm
			`,
			b, dst,
		)

		MOVO_autoAU2(dst, a)
		PADDQ(b, dst)
	}
	return dst
}

/*
*
Synopsis

	__m128i _mm_srli_epi64 (__m128i a, int imm8)
	#include <emmintrin.h>
	Instruction: psrlq xmm, imm8
	CPUID Flags: SSE2

Description

	Shift packed 64-bit integers in a right by imm8 while shifting in zeros, and store the results in dst.

Operation

	FOR j := 0 to 1
		i := j*64
		IF imm8[7:0] > 63
			dst[i+63:i] := 0
		ELSE
			dst[i+63:i] := ZeroExtend64(a[i+63:i] >> imm8[7:0])
		FI
	ENDFOR
*/
func F_mm_srli_epi64(dst VecVirtual, a, imm8 Op) VecVirtual {
	if dst != a {
		MOVO_autoAU2(dst, a)
	}

	CheckType(
		`
		//	PSRLQ imm8 xmm
		//	PSRLQ m128 xmm
		//	PSRLQ xmm  xmm
		`,
		imm8, dst,
	)

	PSRLQ(imm8, dst)
	return dst
}

/*
*
Synopsis

	__m128i _mm_slli_epi64 (__m128i a, int imm8)
	#include <emmintrin.h>
	Instruction: psllq xmm, imm8
	CPUID Flags: SSE2

Description

	Shift packed 64-bit integers in a left by imm8 while shifting in zeros, and store the results in dst.

Operation

	FOR j := 0 to 1
		i := j*64
		IF imm8[7:0] > 63
			dst[i+63:i] := 0
		ELSE
			dst[i+63:i] := ZeroExtend64(a[i+63:i] << imm8[7:0])
		FI
	ENDFOR
*/
func F_mm_slli_epi64(dst VecVirtual, a, imm8 Op) VecVirtual {
	if dst != a {
		MOVO_autoAU2(dst, a)
	}

	CheckType(
		`
		//	PSLLQ imm8 xmm
		//	PSLLQ m128 xmm
		//	PSLLQ xmm  xmm
		`,
		imm8, dst,
	)

	PSLLQ(imm8, dst)
	return dst
}

/*
*
Synopsis

	__m128i _mm_load_si128 (__m128i const* mem_addr)
	#include <emmintrin.h>
	Instruction: movdqa xmm, m128
	CPUID Flags: SSE2

Description

	Load 128-bits of integer data from memory into dst. mem_addr must be aligned on a 16-byte boundary or a general-protection exception may be generated.

Operation

	dst[127:0] := MEM[mem_addr+127:mem_addr]
*/
func F_mm_load_si128(dst VecVirtual, src Mem) VecVirtual {
	CheckType(
		`
		//	MOVOA m128 xmm
		`,
		src, dst,
	)
	MOVOA(dst, src)
	return dst
}

/*
*
Synopsis

	__m128i _mm_srli_epi16 (__m128i a, int imm8)
	#include <emmintrin.h>
	Instruction: psrlw xmm, imm8
	CPUID Flags: SSE2

Description

	Shift packed 16-bit integers in a right by imm8 while shifting in zeros, and store the results in dst.

Operation

	FOR j := 0 to 7
		i := j*16
		IF imm8[7:0] > 15
			dst[i+15:i] := 0
		ELSE
			dst[i+15:i] := ZeroExtend16(a[i+15:i] >> imm8[7:0])
		FI
	ENDFOR
*/
func F_mm_srli_epi16(dst VecVirtual, a, imm8 Op) VecVirtual {
	if dst != a {
		MOVO_autoAU2(dst, a)
	}

	CheckType(
		`
		//	PSRLW imm8 xmm
		`,
		imm8, dst,
	)

	PSRLW(imm8, dst)
	return dst
}
