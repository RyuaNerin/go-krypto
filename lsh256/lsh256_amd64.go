package lsh256

type lsh256ContextAsmData struct {
	algtype           uint32
	remain_databitlen uint32
	cv_l              []uint32
	cv_r              []uint32
	last_block        []byte
}
