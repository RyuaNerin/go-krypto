package lsh512avoconst

type lsh512ContextAsmData struct {
	// 16 aligned
	algtype            uint32
	_                  [4]byte
	remain_databytelen uint64

	cv_l         [64]byte
	cv_r         [64]byte
	i_last_block [256]byte
}
