package lsh256avoconst

type lsh256ContextAsmData struct {
	// 16 aligned
	algtype uint32
	_pad0   [16 - 4]byte
	// 16 aligned
	remain_databytelen uint64
	_pad1              [16 - 8]byte

	cv_l       [32]byte
	cv_r       [32]byte
	last_block [128]byte
}
