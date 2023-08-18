package simd

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"

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
func F_mm256_loadu_si256(dst VecVirtual, src Op) VecVirtual {
	VMOVDQ_autoAU2(dst, src)
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
func F_mm256_storeu_si256(dst, src Op) {
	VMOVDQ_autoAU2(dst, src)
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
func F_mm256_xor_si256(dst VecVirtual, a, b Op) VecVirtual {
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
func F_mm256_or_si256(dst VecVirtual, a, b Op) VecVirtual {
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
func F_mm256_and_si256(dst VecVirtual, a, b Op) VecVirtual {
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
func F_mm256_shuffle_epi8(dst VecVirtual, x, y Op) VecVirtual {
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
func F_mm256_add_epi32(dst VecVirtual, a, b Op) VecVirtual {
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
func F_mm256_slli_epi32(dst VecVirtual, a, imm8 Op) VecVirtual {
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
func F_mm256_srli_epi32(dst VecVirtual, a, imm8 Op) VecVirtual {
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
func F_mm256_shuffle_epi32(dst VecVirtual, a, imm8 Op) VecVirtual {
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
func F_mm256_permute2x128_si256(dst VecVirtual, a, b, imm8 Op) VecVirtual {
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

/*
*
Synopsis

	__m256i _mm256_add_epi64 (__m256i a, __m256i b)
	#include <immintrin.h>
	Instruction: vpaddq ymm, ymm, ymm
	CPUID Flags: AVX2

Description

	Add packed 64-bit integers in a and b, and store the results in dst.

Operation

	FOR j := 0 to 3
		i := j*64
		dst[i+63:i] := a[i+63:i] + b[i+63:i]
	ENDFOR
	dst[MAX:256] := 0
*/
func F_mm256_add_epi64(dst VecVirtual, x, y Op) VecVirtual {
	CheckType(
		`
		//	VPADDQ m256 ymm ymm
		//	VPADDQ ymm  ymm ymm
		//	VPADDQ m128 xmm xmm
		//	VPADDQ xmm  xmm xmm
		//	VPADDQ m128 xmm k xmm
		//	VPADDQ m256 ymm k ymm
		//	VPADDQ xmm  xmm k xmm
		//	VPADDQ ymm  ymm k ymm
		//	VPADDQ m512 zmm k zmm
		//	VPADDQ m512 zmm zmm
		//	VPADDQ zmm  zmm k zmm
		//	VPADDQ zmm  zmm zmm
		`,
		y, x, dst,
	)

	VPADDQ(y, x, dst)
	return dst
}

/*
*
Synopsis

	__m256i _mm256_slli_epi64 (__m256i a, int imm8)
	#include <immintrin.h>
	Instruction: vpsllq ymm, ymm, imm8
	CPUID Flags: AVX2

Description

	Shift packed 64-bit integers in a left by imm8 while shifting in zeros, and store the results in dst.

Operation

	FOR j := 0 to 3
		i := j*64
		IF imm8[7:0] > 63
			dst[i+63:i] := 0
		ELSE
			dst[i+63:i] := ZeroExtend64(a[i+63:i] << imm8[7:0])
		FI
	ENDFOR
	dst[MAX:256] := 0
*/
func F_mm256_slli_epi64(dst VecVirtual, x, r Op) VecVirtual {
	CheckType(
		`
		//	VPSLLQ imm8 ymm  ymm
		//	VPSLLQ m128 ymm  ymm
		//	VPSLLQ xmm  ymm  ymm
		//	VPSLLQ imm8 xmm  xmm
		//	VPSLLQ m128 xmm  xmm
		//	VPSLLQ xmm  xmm  xmm
		//	VPSLLQ imm8 m128 k xmm
		//	VPSLLQ imm8 m128 xmm
		//	VPSLLQ imm8 m256 k ymm
		//	VPSLLQ imm8 m256 ymm
		//	VPSLLQ imm8 xmm  k xmm
		//	VPSLLQ imm8 ymm  k ymm
		//	VPSLLQ m128 xmm  k xmm
		//	VPSLLQ m128 ymm  k ymm
		//	VPSLLQ xmm  xmm  k xmm
		//	VPSLLQ xmm  ymm  k ymm
		//	VPSLLQ imm8 m512 k zmm
		//	VPSLLQ imm8 m512 zmm
		//	VPSLLQ imm8 zmm  k zmm
		//	VPSLLQ imm8 zmm  zmm
		//	VPSLLQ m128 zmm  k zmm
		//	VPSLLQ m128 zmm  zmm
		//	VPSLLQ xmm  zmm  k zmm
		//	VPSLLQ xmm  zmm  zmm
		`,
		r, x, dst,
	)

	VPSLLQ(r, x, dst)
	return dst
}

/*
*
Synopsis

	__m256i _mm256_srli_epi64 (__m256i a, int imm8)
	#include <immintrin.h>
	Instruction: vpsrlq ymm, ymm, imm8
	CPUID Flags: AVX2

Description

	Shift packed 64-bit integers in a right by imm8 while shifting in zeros, and store the results in dst.

Operation

	FOR j := 0 to 3
		i := j*64
		IF imm8[7:0] > 63
			dst[i+63:i] := 0
		ELSE
			dst[i+63:i] := ZeroExtend64(a[i+63:i] >> imm8[7:0])
		FI
	ENDFOR
	dst[MAX:256] := 0
*/
func F_mm256_srli_epi64(dst VecVirtual, x, r Op) VecVirtual {
	CheckType(
		`
		//	VPSRLQ imm8 ymm  ymm
		//	VPSRLQ m128 ymm  ymm
		//	VPSRLQ xmm  ymm  ymm
		//	VPSRLQ imm8 xmm  xmm
		//	VPSRLQ m128 xmm  xmm
		//	VPSRLQ xmm  xmm  xmm
		//	VPSRLQ imm8 m128 k xmm
		//	VPSRLQ imm8 m128 xmm
		//	VPSRLQ imm8 m256 k ymm
		//	VPSRLQ imm8 m256 ymm
		//	VPSRLQ imm8 xmm  k xmm
		//	VPSRLQ imm8 ymm  k ymm
		//	VPSRLQ m128 xmm  k xmm
		//	VPSRLQ m128 ymm  k ymm
		//	VPSRLQ xmm  xmm  k xmm
		//	VPSRLQ xmm  ymm  k ymm
		//	VPSRLQ imm8 m512 k zmm
		//	VPSRLQ imm8 m512 zmm
		//	VPSRLQ imm8 zmm  k zmm
		//	VPSRLQ imm8 zmm  zmm
		//	VPSRLQ m128 zmm  k zmm
		//	VPSRLQ m128 zmm  zmm
		//	VPSRLQ xmm  zmm  k zmm
		//	VPSRLQ xmm  zmm  zmm
		`,
		r, x, dst,
	)

	VPSRLQ(r, x, dst)
	return dst
}

/*
*
Synopsis

	__m256i _mm256_permute4x64_epi64 (__m256i a, const int imm8)
	#include <immintrin.h>
	Instruction: vpermq ymm, ymm, imm8
	CPUID Flags: AVX2

Description

	Shuffle 64-bit integers in a across lanes using the control in imm8, and store the results in dst.

Operation

	DEFINE SELECT4(src, control) {
		CASE(control[1:0]) OF
		0:	tmp[63:0] := src[63:0]
		1:	tmp[63:0] := src[127:64]
		2:	tmp[63:0] := src[191:128]
		3:	tmp[63:0] := src[255:192]
		ESAC
		RETURN tmp[63:0]
	}
	dst[63:0] := SELECT4(a[255:0], imm8[1:0])
	dst[127:64] := SELECT4(a[255:0], imm8[3:2])
	dst[191:128] := SELECT4(a[255:0], imm8[5:4])
	dst[255:192] := SELECT4(a[255:0], imm8[7:6])
	dst[MAX:256] := 0
*/
func F_mm256_permute4x64_epi64(dst VecVirtual, a, imm8 Op) VecVirtual {
	CheckType(
		`
		//	VPERMQ imm8 m256 ymm
		//	VPERMQ imm8 ymm  ymm
		//	VPERMQ imm8 m256 k ymm
		//	VPERMQ imm8 ymm  k ymm
		//	VPERMQ m256 ymm  k ymm
		//	VPERMQ m256 ymm  ymm
		//	VPERMQ ymm  ymm  k ymm
		//	VPERMQ ymm  ymm  ymm
		//	VPERMQ imm8 m512 k zmm
		//	VPERMQ imm8 m512 zmm
		//	VPERMQ imm8 zmm  k zmm
		//	VPERMQ imm8 zmm  zmm
		//	VPERMQ m512 zmm  k zmm
		//	VPERMQ m512 zmm  zmm
		//	VPERMQ zmm  zmm  k zmm
		//	VPERMQ zmm  zmm  zmm
		`,
		imm8, a, dst,
	)

	VPERMQ(imm8, a, dst)
	return dst
}
