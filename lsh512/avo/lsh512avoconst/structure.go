package lsh512avoconst

type lsh512ContextAsmData struct {
	// 16 aligned
	algtype            uint32
	_                  [4]byte
	remain_databytelen uint64

	cv_l         [8]uint64
	cv_r         [8]uint64
	i_last_block [256]byte
}
