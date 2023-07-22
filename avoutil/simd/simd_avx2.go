package simd

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"

	. "kryptosimd/avoutil"
)

/*
Synopsis

	__m256i _mm256_loadu_si256 (__m256i const * mem_addr)
	#include <immintrin.h>
	Instruction: vmovdqu ymm, m256
	CPUID Flags: AVX

Description

	Load 256-bits of integer data from memory into dst. mem_addr does not need to be aligned on any particular boundary.

Operation

	dst[255:0] := MEM[mem_addr+255:mem_addr]
	dst[MAX:256] := 0
*/
func F_mm256_loadu_si256(dst, src Op) Op {
	VMOVDQ_autoAU(src, dst)
	return dst
}

/*
Synopsis

	void _mm256_storeu_si256 (__m256i * mem_addr, __m256i a)
	#include <immintrin.h>
	Instruction: vmovdqu m256, ymm
	CPUID Flags: AVX

Description

	Store 256-bits of integer data from a into memory. mem_addr does not need to be aligned on any particular boundary.

Operation

	MEM[mem_addr+255:mem_addr] := a[255:0]
*/
func F_mm256_storeu_si256(dst, src Op) Op {
	VMOVDQ_autoAU(src, dst)
	return dst
}

/*
Synopsis

	__m256i _mm256_xor_si256 (__m256i a, __m256i b)
	#include <immintrin.h>
	Instruction: vpxor ymm, ymm, ymm
	CPUID Flags: AVX2

Description

	Compute the bitwise XOR of 256 bits (representing integer data) in a and b, and store the result in dst.

Operation

	dst[255:0] := (a[255:0] XOR b[255:0])
	dst[MAX:256] := 0
*/
func F_mm256_xor_si256(dst, a, b Op) Op {
	CheckType(
		`
		//	VPXOR m256 ymm ymm
		//	VPXOR ymm  ymm ymm
		//	VPXOR m128 xmm xmm
		//	VPXOR xmm  xmm xmm
		`,
		b, a, dst,
	)

	VPXOR(b, a, dst)
	return dst
}

/*
Synopsis

	__m256i _mm256_or_si256 (__m256i a, __m256i b)
	#include <immintrin.h>
	Instruction: vpor ymm, ymm, ymm
	CPUID Flags: AVX2

Description

	Compute the bitwise OR of 256 bits (representing integer data) in a and b, and store the result in dst.

Operation

	dst[255:0] := (a[255:0] OR b[255:0])
	dst[MAX:256] := 0
*/
func F_mm256_or_si256(dst, a, b Op) Op {
	CheckType(
		`
		//	VPOR m256 ymm ymm
		//	VPOR ymm  ymm ymm
		//	VPOR m128 xmm xmm
		//	VPOR xmm  xmm xmm
		`,
		b, a, dst,
	)

	VPOR(b, a, dst)
	return dst
}

/*
Synopsis

	__m256i _mm256_and_si256 (__m256i a, __m256i b)
	#include <immintrin.h>
	Instruction: vpand ymm, ymm, ymm
	CPUID Flags: AVX2

Description

	Compute the bitwise AND of 256 bits (representing integer data) in a and b, and store the result in dst.

Operation

	dst[255:0] := (a[255:0] AND b[255:0])
	dst[MAX:256] := 0
*/
func F_mm256_and_si256(dst, a, b Op) Op {
	CheckType(
		`
		//	VPAND m256 ymm ymm
		//	VPAND ymm  ymm ymm
		//	VPAND m128 xmm xmm
		//	VPAND xmm  xmm xmm
		`,
		b, a, dst,
	)

	VPAND(b, a, dst)
	return dst
}

/*
Synopsis

	__m256i _mm256_shuffle_epi8 (__m256i a, __m256i b)
	#include <immintrin.h>
	Instruction: vpshufb ymm, ymm, ymm
	CPUID Flags: AVX2

Description

	Shuffle 8-bit integers in a within 128-bit lanes according to shuffle control mask in the corresponding 8-bit element of b, and store the results in dst.

Operation

	FOR j := 0 to 15
		i := j*8
		IF b[i+7] == 1
			dst[i+7:i] := 0
		ELSE
			index[3:0] := b[i+3:i]
			dst[i+7:i] := a[index*8+7:index*8]
		FI
		IF b[128+i+7] == 1
			dst[128+i+7:128+i] := 0
		ELSE
			index[3:0] := b[128+i+3:128+i]
			dst[128+i+7:128+i] := a[128+index*8+7:128+index*8]
		FI
	ENDFOR
	dst[MAX:256] := 0
*/
func F_mm256_shuffle_epi8(dst, x, y Op) Op {
	CheckType(
		`
		//	VPSHUFB m256 ymm ymm
		//	VPSHUFB ymm  ymm ymm
		//	VPSHUFB m128 xmm xmm
		//	VPSHUFB xmm  xmm xmm
		//	VPSHUFB m128 xmm k xmm
		//	VPSHUFB m256 ymm k ymm
		//	VPSHUFB xmm  xmm k xmm
		//	VPSHUFB ymm  ymm k ymm
		//	VPSHUFB m512 zmm k zmm
		//	VPSHUFB m512 zmm zmm
		//	VPSHUFB zmm  zmm k zmm
		//	VPSHUFB zmm  zmm zmm
		`,
		y, x, dst,
	)

	VPSHUFB(y, x, dst)
	return dst
}

/*
Synopsis

	__m256i _mm256_add_epi32 (__m256i a, __m256i b)
	#include <immintrin.h>
	Instruction: vpaddd ymm, ymm, ymm
	CPUID Flags: AVX2

Description

	Add packed 32-bit integers in a and b, and store the results in dst.

Operation

	FOR j := 0 to 7
		i := j*32
		dst[i+31:i] := a[i+31:i] + b[i+31:i]
	ENDFOR
	dst[MAX:256] := 0
*/
func F_mm256_add_epi32(dst, a, b Op) Op {
	CheckType(
		`
		//	VPADDD m256 ymm ymm
		//	VPADDD ymm  ymm ymm
		//	VPADDD m128 xmm xmm
		//	VPADDD xmm  xmm xmm
		//	VPADDD m128 xmm k xmm
		//	VPADDD m256 ymm k ymm
		//	VPADDD xmm  xmm k xmm
		//	VPADDD ymm  ymm k ymm
		//	VPADDD m512 zmm k zmm
		//	VPADDD m512 zmm zmm
		//	VPADDD zmm  zmm k zmm
		//	VPADDD zmm  zmm zmm
		`,
		b, a, dst,
	)

	VPADDD(b, a, dst)
	return dst
}

/*
Synopsis

	__m256i _mm256_slli_epi32 (__m256i a, int imm8)
	#include <immintrin.h>
	Instruction: vpslld ymm, ymm, imm8
	CPUID Flags: AVX2

Description

	Shift packed 32-bit integers in a left by imm8 while shifting in zeros, and store the results in dst.

Operation

	FOR j := 0 to 7
		i := j*32
		IF imm8[7:0] > 31
			dst[i+31:i] := 0
		ELSE
			dst[i+31:i] := ZeroExtend32(a[i+31:i] << imm8[7:0])
		FI
	ENDFOR
	dst[MAX:256] := 0
*/
func F_mm256_slli_epi32(dst, a, imm8 Op) Op {
	CheckType(
		`
		//	VPSLLD imm8 ymm  ymm
		//	VPSLLD m128 ymm  ymm
		//	VPSLLD xmm  ymm  ymm
		//	VPSLLD imm8 xmm  xmm
		//	VPSLLD m128 xmm  xmm
		//	VPSLLD xmm  xmm  xmm
		//	VPSLLD imm8 m128 k xmm
		//	VPSLLD imm8 m128 xmm
		//	VPSLLD imm8 m256 k ymm
		//	VPSLLD imm8 m256 ymm
		//	VPSLLD imm8 xmm  k xmm
		//	VPSLLD imm8 ymm  k ymm
		//	VPSLLD m128 xmm  k xmm
		//	VPSLLD m128 ymm  k ymm
		//	VPSLLD xmm  xmm  k xmm
		//	VPSLLD xmm  ymm  k ymm
		//	VPSLLD imm8 m512 k zmm
		//	VPSLLD imm8 m512 zmm
		//	VPSLLD imm8 zmm  k zmm
		//	VPSLLD imm8 zmm  zmm
		//	VPSLLD m128 zmm  k zmm
		//	VPSLLD m128 zmm  zmm
		//	VPSLLD xmm  zmm  k zmm
		//	VPSLLD xmm  zmm  zmm
		`,
		imm8, a, dst,
	)

	VPSLLD(imm8, a, dst)
	return dst
}

/*
Synopsis

	__m256i _mm256_srli_epi32 (__m256i a, int imm8)
	#include <immintrin.h>
	Instruction: vpsrld ymm, ymm, imm8
	CPUID Flags: AVX2

Description

	Shift packed 32-bit integers in a right by imm8 while shifting in zeros, and store the results in dst.

Operation

	FOR j := 0 to 7
		i := j*32
		IF imm8[7:0] > 31
			dst[i+31:i] := 0
		ELSE
			dst[i+31:i] := ZeroExtend32(a[i+31:i] >> imm8[7:0])
		FI
	ENDFOR
	dst[MAX:256] := 0
*/
func F_mm256_srli_epi32(dst, a, imm8 Op) Op {
	CheckType(
		`
		//	VPSRLD imm8 ymm  ymm
		//	VPSRLD m128 ymm  ymm
		//	VPSRLD xmm  ymm  ymm
		//	VPSRLD imm8 xmm  xmm
		//	VPSRLD m128 xmm  xmm
		//	VPSRLD xmm  xmm  xmm
		//	VPSRLD imm8 m128 k xmm
		//	VPSRLD imm8 m128 xmm
		//	VPSRLD imm8 m256 k ymm
		//	VPSRLD imm8 m256 ymm
		//	VPSRLD imm8 xmm  k xmm
		//	VPSRLD imm8 ymm  k ymm
		//	VPSRLD m128 xmm  k xmm
		//	VPSRLD m128 ymm  k ymm
		//	VPSRLD xmm  xmm  k xmm
		//	VPSRLD xmm  ymm  k ymm
		//	VPSRLD imm8 m512 k zmm
		//	VPSRLD imm8 m512 zmm
		//	VPSRLD imm8 zmm  k zmm
		//	VPSRLD imm8 zmm  zmm
		//	VPSRLD m128 zmm  k zmm
		//	VPSRLD m128 zmm  zmm
		//	VPSRLD xmm  zmm  k zmm
		//	VPSRLD xmm  zmm  zmm
		`,
		imm8, a, dst,
	)

	VPSRLD(imm8, a, dst)
	return dst
}

/*
Synopsis

	__m256i _mm256_shuffle_epi32 (__m256i a, const int imm8)
	#include <immintrin.h>
	Instruction: vpshufd ymm, ymm, imm8
	CPUID Flags: AVX2

Description

	Shuffle 32-bit integers in a within 128-bit lanes using the control in imm8, and store the results in dst.

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
	dst[159:128] := SELECT4(a[255:128], imm8[1:0])
	dst[191:160] := SELECT4(a[255:128], imm8[3:2])
	dst[223:192] := SELECT4(a[255:128], imm8[5:4])
	dst[255:224] := SELECT4(a[255:128], imm8[7:6])
	dst[MAX:256] := 0
*/
func F_mm256_shuffle_epi32(dst, a, imm8 Op) Op {
	CheckType(
		`
		//	VPSHUFD imm8 m256 ymm
		//	VPSHUFD imm8 ymm  ymm
		//	VPSHUFD imm8 m128 xmm
		//	VPSHUFD imm8 xmm  xmm
		//	VPSHUFD imm8 m128 k xmm
		//	VPSHUFD imm8 m256 k ymm
		//	VPSHUFD imm8 xmm  k xmm
		//	VPSHUFD imm8 ymm  k ymm
		//	VPSHUFD imm8 m512 k zmm
		//	VPSHUFD imm8 m512 zmm
		//	VPSHUFD imm8 zmm  k zmm
		//	VPSHUFD imm8 zmm  zmm
		`,
		imm8, a, dst,
	)

	VPSHUFD(imm8, a, dst)
	return dst
}

/*
Synopsis

	__m256i _mm256_permute2x128_si256 (__m256i a, __m256i b, const int imm8)
	#include <immintrin.h>
	Instruction: vperm2i128 ymm, ymm, ymm, imm8
	CPUID Flags: AVX2

Description

	Shuffle 128-bits (composed of integer data) selected by imm8 from a and b, and store the results in dst.

Operation

	DEFINE SELECT4(src1, src2, control) {
		CASE(control[1:0]) OF
		0:	tmp[127:0] := src1[127:0]
		1:	tmp[127:0] := src1[255:128]
		2:	tmp[127:0] := src2[127:0]
		3:	tmp[127:0] := src2[255:128]
		ESAC
		IF control[3]
			tmp[127:0] := 0
		FI
		RETURN tmp[127:0]
	}
	dst[127:0] := SELECT4(a[255:0], b[255:0], imm8[3:0])
	dst[255:128] := SELECT4(a[255:0], b[255:0], imm8[7:4])
	dst[MAX:256] := 0
*/
func F_mm256_permute2x128_si256(dst, a, b, imm8 Op) Op {
	CheckType(
		`
		//	VPERM2I128 imm8 m256 ymm ymm
		//	VPERM2I128 imm8 ymm  ymm ymm
		`,
		imm8, b, a, dst,
	)

	VPERM2I128(imm8, b, a, dst)
	return dst
}
