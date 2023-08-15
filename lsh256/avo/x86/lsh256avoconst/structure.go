package lsh256avoconst

type lsh256ContextAsmData struct {
	// 16 aligned
	algtype            uint32
	_                  [4]byte
	remain_databytelen uint64

	cv_l       [32]byte
	cv_r       [32]byte
	last_block [128]byte
}
