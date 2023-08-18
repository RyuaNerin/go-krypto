package sse2

import (
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"

	. "kryptosimd/avoutil/simd"
)

/* -------------------------------------------------------- */
// LSH: functions
/* -------------------------------------------------------- */
/* -------------------------------------------------------- */
// register functions macro
/* -------------------------------------------------------- */

// #define LOAD(x) _mm_loadu_si128((__m128i*)x)
func LOAD(dst VecVirtual, x Op) Op { return F_mm_loadu_si128(dst, x) }

// #define STORE(x,y) _mm_storeu_si128((__m128i*)x, y)
func STORE(dst, y Op) { F_mm_storeu_si128(dst, y) }

// #define XOR(x,y) _mm_xor_si128(x,y)
func XOR(dst VecVirtual, x, y Op) VecVirtual { return F_mm_xor_si128(dst, x, y) }

// #define OR(x,y) _mm_or_si128(x,y)
func OR(dst VecVirtual, x, y Op) VecVirtual { return F_mm_or_si128(dst, x, y) }

// #define AND(x,y) _mm_and_si128(x,y)
func AND(dst VecVirtual, x, y Op) VecVirtual { return F_mm_and_si128(dst, x, y) }

func ADD32(dst VecVirtual, x, y Op) VecVirtual { return F_mm_add_epi32(dst, x, y) } // #define ADD(x,y) _mm_add_epi32(x,y)
func ADD64(dst VecVirtual, x, y Op) VecVirtual { return F_mm_add_epi64(dst, x, y) } // #define ADD(x,y) _mm_add_epi64(x,y)

func SHIFT_L32(dst VecVirtual, x, r Op) VecVirtual { return F_mm_slli_epi32(dst, x, r) } // #define SHIFT_L(x,r) _mm_slli_epi32(x,r)
func SHIFT_R32(dst VecVirtual, x, r Op) VecVirtual { return F_mm_srli_epi32(dst, x, r) } // #define SHIFT_R(x,r) _mm_srli_epi32(x,r)

func SHIFT_L64(dst VecVirtual, x, r Op) VecVirtual { return F_mm_slli_epi64(dst, x, r) } // #define SHIFT_L(x,r) _mm_slli_epi64(x,r)
func SHIFT_R64(dst VecVirtual, x, r Op) VecVirtual { return F_mm_srli_epi64(dst, x, r) } // #define SHIFT_R(x,r) _mm_srli_epi64(x,r)
