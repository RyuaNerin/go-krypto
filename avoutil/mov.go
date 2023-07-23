package avoutil

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

const (
	YmmSize = 256 / 8
	XmmSize = 128 / 8
)

func isAligned(alignedByte int, args ...Op) bool {
	for _, v := range args {
		_, isGP := v.(GPVirtual)
		_, isVec := v.(VecVirtual)
		mem, isMem := v.(Mem)

		switch {
		case isGP:
			continue
		case isVec:
			continue
		case isMem:
			if mem.Symbol.Static && mem.Symbol.Name != "" && mem.Disp%alignedByte == 0 && mem.Index == nil {
				continue
			}
		}
		return false
	}

	return true
}

func VMOVDQ_autoAU(mxy, mxy1 Op) {
	if isAligned(YmmSize, mxy, mxy1) {
		CheckType(
			`
			//	VMOVDQA m128 xmm
			//	VMOVDQA m256 ymm
			//	VMOVDQA xmm  m128
			//	VMOVDQA xmm  xmm
			//	VMOVDQA ymm  m256
			//	VMOVDQA ymm  ymm
			`,
			mxy, mxy1,
		)

		VMOVDQA(mxy, mxy1)
	} else {
		CheckType(
			`
			//	VMOVDQU m128 xmm
			//	VMOVDQU m256 ymm
			//	VMOVDQU xmm  m128
			//	VMOVDQU xmm  xmm
			//	VMOVDQU ymm  m256
			//	VMOVDQU ymm  ymm
			`,
			mxy, mxy1,
		)

		VMOVDQU(mxy, mxy1)
	}
}

func MOVO_autoAU2(dst, src Op) {
	MOVO_autoAU(src, dst)
}
func MOVO_autoAU(mx, mx1 Op) {
	if isAligned(XmmSize, mx, mx1) {
		CheckType(
			`
			//	MOVOA m128 xmm
			//	MOVOA xmm  m128
			//	MOVOA xmm  xmm
			`,
			mx, mx1,
		)

		MOVOA(mx, mx1)
	} else {
		CheckType(
			`
			//	MOVOU m128 xmm
			//	MOVOU xmm  m128
			//	MOVOU xmm  xmm
			`,
			mx, mx1,
		)

		MOVOU(mx, mx1)
	}
}
