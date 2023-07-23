package lsh256avoconst

type lsh256ContextAsmData struct {
	// 16 aligned
	algtype uint32
	_pad0   [16 - 4]byte
	// 16 aligned
	remain_databitlen uint32
	_pad1             [16 - 4]byte

	cv_l       [32 / 4]uint32
	cv_r       [32 / 4]uint32
	last_block [128]byte
}
