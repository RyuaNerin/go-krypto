package lsh512avoconst

type lsh512ContextAsmData struct {
	// 16 aligned
	algtype uint32
	_pad0   [16 - 4]byte
	// 16 aligned
	remain_databitlen uint32
	_pad1             [16 - 4]byte

	cv_l         [8]uint64
	cv_r         [8]uint64
	i_last_block [256]byte
}
