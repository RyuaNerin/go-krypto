package lsh512avoconst

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"
)

// lsh_def.h
const (
	LSH_TYPE_512_512 = 0x0000040 // 0x0010040
	LSH_TYPE_512_384 = 0x0000030 // 0x0010030
	LSH_TYPE_512_256 = 0x0000020 // 0x0010020
	LSH_TYPE_512_224 = 0x000001C // 0x001001C

	LSH512_MSG_BLK_BYTE_LEN      = 256
	LSH512_MSG_BLK_BIT_LEN       = 2048
	LSH512_CV_BYTE_LEN           = 128
	LSH512_HASH_VAL_MAX_BYTE_LEN = 64
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
