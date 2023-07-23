package sse2

import (
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"

	. "kryptosimd/avoutil"
	. "kryptosimd/avoutil/simd"
)

/* -------------------------------------------------------- */
// LSH: functions
/* -------------------------------------------------------- */
/* -------------------------------------------------------- */
// register functions macro
/* -------------------------------------------------------- */

// #define LOAD(x) _mm_loadu_si128((__m128i*)x)
func LOAD(dst, x Op) Op { return F_mm_loadu_si128(dst, x) }

// #define STORE(x,y) _mm_storeu_si128((__m128i*)x, y)
func STORE(dst, y Op) { F_mm_storeu_si128(dst, y) }

// #define XOR(x,y) _mm_xor_si128(x,y)
func XOR(dst Op, y VecVirtual) Op { return F_mm_xor_si128(dst, y) }

// #define OR(x,y) _mm_or_si128(x,y)
func OR(dst Op, y VecVirtual) Op { return F_mm_or_si128(dst, y) }

// #define AND(x,y) _mm_and_si128(x,y)
func AND(dst Op, y VecVirtual) Op { return F_mm_and_si128(dst, y) }

// #define ADD(x,y) _mm_add_epi32(x,y)
func ADD(dst VecVirtual, y Op) Op { return F_mm_add_epi32(dst, y) }
func ADD3(dst VecVirtual, a, b Op) Op {
	MOVO_autoAU2(dst, a)
	ADD(dst, b)
	return dst
}

// #define SHIFT_L(x,r) _mm_slli_epi32(x,r)
func SHIFT_L(dst VecVirtual, r Op) Op { return F_mm_slli_epi32(dst, r) }

// #define SHIFT_R(x,r) _mm_srli_epi32(x,r)
func SHIFT_R(dst VecVirtual, r Op) Op { return F_mm_srli_epi32(dst, r) }
