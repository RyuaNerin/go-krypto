package avx2

import (
	. "github.com/mmcloughlin/avo/operand"

	. "kryptosimd/avoutil/simd"
)

/* -------------------------------------------------------- *
*  LSH: functions
*  -------------------------------------------------------- */

// #define LOAD(x) _mm256_loadu_si256((__m256i*)x)
func LOAD(dst, src Op) Op { return F_mm256_loadu_si256(dst, src) }

// #define STORE(x,y) _mm256_storeu_si256((__m256i*)x, y)
func STORE(dst, src Op) { F_mm256_storeu_si256(dst, src) }

// #define XOR(x,y) _mm256_xor_si256(x,y)
func XOR(dst, x, y Op) Op { return F_mm256_xor_si256(dst, x, y) }

// #define OR(x,y) _mm256_or_si256(x,y)
func OR(dst, x, y Op) Op { return F_mm256_or_si256(dst, x, y) }

// #define AND(x,y) _mm256_and_si256(x,y)
func AND(dst, x, y Op) Op { return F_mm256_and_si256(dst, x, y) }

// #define SHUFFLE8(x,y) _mm256_shuffle_epi8(x,y)
func SHUFFLE8(dst, x, y Op) Op { return F_mm256_shuffle_epi8(dst, x, y) }

// #define ADD(x,y) _mm256_add_epi32(x,y)
func ADD(dst, x, y Op) Op { return F_mm256_add_epi32(dst, x, y) }

// #define SHIFT_L(x,r) _mm256_slli_epi32(x,r)
func SHIFT_L(dst, x, y Op) Op { return F_mm256_slli_epi32(dst, x, y) }

// #define SHIFT_R(x,r) _mm256_srli_epi32(x,r)
func SHIFT_R(dst, x, y Op) Op { return F_mm256_srli_epi32(dst, x, y) }
