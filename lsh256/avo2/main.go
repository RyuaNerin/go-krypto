package main

import (
	. "github.com/mmcloughlin/avo/build"
)

type lsh256ContextAsmData struct {
	algtype           uint32
	____p0            [16 - 4]byte
	remain_databitlen uint32
	____p1            [16 - 4]byte
	cv_l              [8]uint32
	cv_r              [8]uint32
	last_block        [128]byte
}

func main() {
	Package("kryptosimd/lsh256/avo2")

	lsh256_sse2_init()
	//lsh256_sse2_update()

	Generate()
	print("done")
}
