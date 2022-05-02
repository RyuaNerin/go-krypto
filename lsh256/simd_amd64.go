//go:build amd64

package lsh256

var (
	simdSetSSE2 = simdSet{
		init:   lsh256InitSSE2,
		update: lsh256UpdateSSE2,
		final:  lsh256FinalSSE2,
	}
	simdSetSSSE3 = simdSet{
		/**
		init:   lsh256InitSSSE3,
		update: lsh256UpdateSSSE3,
		final:  lsh256FinalSSSE3,
		*/
	}
	simdSetAVX2 = simdSet{
		/**
		init:   lsh256InitAVX2,
		update: lsh256UpdateAVX2,
		final:  lsh256FinalAVX2,
		*/
	}
)

func init() {
	simdSetDefault = simdSetSSE2

	/**
	switch {
	case cpu.X86.HasSSSE3:
		simdSetDefault = simdSetSSSE3
	case cpu.X86.HasAVX2:
		simdSetDefault = simdSetAVX2
	}
	*/
}
