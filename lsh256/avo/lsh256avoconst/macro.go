package avoutil

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

// #define LSH_GET_SMALL_HASHBIT(lsh_type_val)		((lsh_type_val)>>24)
func LSH_GET_SMALL_HASHBIT(dst Register, lsh_type_val Op) Register {
	if dst != lsh_type_val {
		MOVL(lsh_type_val, dst)
	}
	SHRL(U8(24), dst)
	return dst
}

// #define LSH_GET_HASHBYTE(lsh_type_val)			((lsh_type_val) & 0xffff)
func LSH_GET_HASHBYTE(dst Register, lsh_type_val Op) Register {
	if dst != lsh_type_val {
		MOVL(lsh_type_val, dst)
	}
	ANDL(U32(0xffff), dst)
	return dst
}
