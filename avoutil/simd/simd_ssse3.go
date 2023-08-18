package simd

import (
	. "kryptosimd/avoutil"

	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

/**
Synopsis
	__m128i _mm_shuffle_epi8 (__m128i a, __m128i b)
	#include <tmmintrin.h>
	Instruction: pshufb xmm, xmm
	CPUID Flags: SSSE3
Description
	Shuffle packed 8-bit integers in a according to shuffle control mask in the corresponding 8-bit element of b, and store the results in dst.
Operation
	FOR j := 0 to 15
		i := j*8
		IF b[i+7] == 1
			dst[i+7:i] := 0
		ELSE
			index[3:0] := b[i+3:i]
			dst[i+7:i] := a[index*8+7:index*8]
		FI
	ENDFOR

dst = a
*/

func F_mm_shuffle_epi8(dst VecVirtual, a, b Op) VecVirtual {
	if dst != a {
		MOVO_autoAU2(dst, a)
	}

	CheckType(
		`
		//	PSHUFB m128 xmm
		//	PSHUFB xmm  xmm
		`,
		b, dst,
	)

	PSHUFB(b, dst)

	return dst
}
