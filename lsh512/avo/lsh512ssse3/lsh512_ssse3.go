package lsh512ssse3

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"

	. "kryptosimd/avoutil"
	. "kryptosimd/avoutil/simd"
	. "kryptosimd/lsh/ssse3"
	. "kryptosimd/lsh512/avo/lsh512avoconst"
	. "kryptosimd/lsh512/avo/lsh512common"
)

/* -------------------------------------------------------- */
// LSH: variables
/* -------------------------------------------------------- */

//	typedef struct LSH_ALIGNED_(32) {
//		LSH_ALIGNED_(32) __m128i submsg_e_l[4];	/* even left sub-message */
//		LSH_ALIGNED_(32) __m128i submsg_e_r[4];	/* even right sub-message */
//		LSH_ALIGNED_(32) __m128i submsg_o_l[4];	/* odd left sub-message */
//		LSH_ALIGNED_(32) __m128i submsg_o_r[4];	/* odd right sub-message */
//	} LSH512SSSE3_internal;
type LSH512SSSE3_internal struct {
	submsg_e_l []VecVirtual
	submsg_e_r []VecVirtual
	submsg_o_l []VecVirtual
	submsg_o_r []VecVirtual

	submsg_e_l_Mem Mem // xmmSize * 4
	submsg_e_r_Mem Mem // xmmSize * 4
	submsg_o_l_Mem Mem // xmmSize * 4
	submsg_o_r_Mem Mem // xmmSize * 4

	tempVec0 []VecVirtual
	tempVec1 []VecVirtual
}

func (ctx *LSH512SSSE3_internal) load(v []VecVirtual, m Mem) {
	Comment("i_state_load___start")
	load_blk_mem2vec(v, m)
	Comment("i_state_load___end")
}
func (ctx *LSH512SSSE3_internal) save(v []VecVirtual, m Mem) {
	Comment("i_state_save___start")
	load_blk_vec2mem(m, v)
	Comment("i_state_save___end")
}

/* -------------------------------------------------------- */
// load a message block to register
/* -------------------------------------------------------- */

// static INLINE void load_blk(__m128i* dest, const void* src){
func load_blk_mem2vec(dst []VecVirtual, src Mem) {
	Comment("load_blk_mem2vec")

	//dest[0] = LOAD((const __m128i*)src);
	LOAD(dst[0], src.Offset(XmmSize*0))
	//dest[1] = LOAD((const __m128i*)src + 1);
	LOAD(dst[1], src.Offset(XmmSize*1))
	//dest[2] = LOAD((const __m128i*)src + 2);
	LOAD(dst[2], src.Offset(XmmSize*2))
	//dest[3] = LOAD((const __m128i*)src + 3);
	LOAD(dst[3], src.Offset(XmmSize*3))
}
func load_blk_vec2mem(dst Mem, src []VecVirtual) {
	Comment("load_blk_vec2mem")

	//dest[0] = LOAD((const __m128i*)src);
	LOAD(dst.Offset(XmmSize*0), src[0])
	//dest[1] = LOAD((const __m128i*)src + 1);
	LOAD(dst.Offset(XmmSize*1), src[1])
	//dest[2] = LOAD((const __m128i*)src + 2);
	LOAD(dst.Offset(XmmSize*2), src[2])
	//dest[3] = LOAD((const __m128i*)src + 3);
	LOAD(dst.Offset(XmmSize*3), src[3])
}
func load_blk_mem2mem(dst Mem, src Mem) {
	Comment("load_blk_mem2mem")

	tmp := XMM()
	//dest[0] = LOAD((const __m128i*)src);
	LOAD(tmp, src)
	LOAD(dst, tmp)
	//dest[1] = LOAD((const __m128i*)src + 1);
	LOAD(tmp, src.Offset(XmmSize*1))
	LOAD(dst.Offset(XmmSize*1), tmp)
	//dest[2] = LOAD((const __m128i*)src + 2);
	LOAD(tmp, src.Offset(XmmSize*2))
	LOAD(dst.Offset(XmmSize*2), tmp)
	//dest[3] = LOAD((const __m128i*)src + 3);
	LOAD(tmp, src.Offset(XmmSize*3))
	LOAD(dst.Offset(XmmSize*3), tmp)
}

// static INLINE void store_blk(__m128i* dest, const __m128i* src){
func store_blk(dst Mem, src []VecVirtual) {
	Comment("store_blk")

	//STORE(dest, src[0]);
	STORE(dst, src[0])
	//STORE(dest + 1, src[1]);
	STORE(dst.Offset(XmmSize), src[1])
	//STORE(dest + 2, src[2]);
	STORE(dst.Offset(XmmSize*2), src[2])
	//STORE(dest + 3, src[3]);
	STORE(dst.Offset(XmmSize*3), src[3])
}

// static INLINE void load_msg_blk(LSH512SSSE3_internal * i_state, const lsh_u64* msgblk){
func load_msg_blk(i_state LSH512SSSE3_internal, msgblk Mem /* uint32 */) {
	//load_blk(i_state->submsg_e_l, msgblk + 0);
	load_blk_mem2mem(i_state.submsg_e_l_Mem, msgblk.Offset(0*8))
	//load_blk(i_state->submsg_e_r, msgblk + 8);
	load_blk_mem2mem(i_state.submsg_e_r_Mem, msgblk.Offset(8*8))
	//load_blk(i_state->submsg_o_l, msgblk + 16);
	load_blk_mem2mem(i_state.submsg_o_l_Mem, msgblk.Offset(16*8))
	//load_blk(i_state->submsg_o_r, msgblk + 24);
	load_blk_mem2mem(i_state.submsg_o_r_Mem, msgblk.Offset(24*8))
}

// static INLINE void msg_exp_even(LSH512SSSE3_internal * i_state){
func msg_exp_even(i_state LSH512SSSE3_internal) {
	if i_state.tempVec0 == nil || i_state.tempVec1 == nil {
		panic(`tempVec is required`)
	}
	Comment("msg_exp_even")

	ADD := ADD64_

	//__m128i temp;
	temp := XMM()

	////////////////////////////////////////////////////////////////////////////////////////////////////

	i_state.submsg_e_l = i_state.tempVec0
	i_state.load(i_state.submsg_e_l, i_state.submsg_e_l_Mem)

	i_state.submsg_o_l = i_state.tempVec1
	i_state.load(i_state.submsg_o_l, i_state.submsg_o_l_Mem)

	//i_state->submsg_e_l[1] = _mm_shuffle_epi32(i_state->submsg_e_l[1], 0x4e);
	F_mm_shuffle_epi32(i_state.submsg_e_l[1], i_state.submsg_e_l[1], U8(0x4e))
	//temp = i_state->submsg_e_l[0];
	MOVO_autoAU2(temp, i_state.submsg_e_l[0])
	//i_state->submsg_e_l[0] = i_state->submsg_e_l[1];
	MOVO_autoAU2(i_state.submsg_e_l[0], i_state.submsg_e_l[1])
	//i_state->submsg_e_l[1] = temp;
	MOVO_autoAU2(i_state.submsg_e_l[1], temp)
	//i_state->submsg_e_l[3] = _mm_shuffle_epi32(i_state->submsg_e_l[3], 0x4e);
	F_mm_shuffle_epi32(i_state.submsg_e_l[3], i_state.submsg_e_l[3], U8(0x4e))
	//temp = i_state->submsg_e_l[2];
	MOVO_autoAU2(temp, i_state.submsg_e_l[2])
	//i_state->submsg_e_l[2] = _mm_unpacklo_epi64(i_state->submsg_e_l[3], i_state->submsg_e_l[2]);
	F_mm_unpacklo_epi64(i_state.submsg_e_l[2], i_state.submsg_e_l[3], i_state.submsg_e_l[2])
	//i_state->submsg_e_l[3] = _mm_unpackhi_epi64(temp, i_state->submsg_e_l[3]);
	F_mm_unpackhi_epi64(i_state.submsg_e_l[3], temp, i_state.submsg_e_l[3])

	//i_state->submsg_e_l[0] = ADD(i_state->submsg_o_l[0], i_state->submsg_e_l[0]);
	ADD(i_state.submsg_e_l[0], i_state.submsg_o_l[0], i_state.submsg_e_l[0])
	//i_state->submsg_e_l[1] = ADD(i_state->submsg_o_l[1], i_state->submsg_e_l[1]);
	ADD(i_state.submsg_e_l[1], i_state.submsg_o_l[1], i_state.submsg_e_l[1])
	//i_state->submsg_e_l[2] = ADD(i_state->submsg_o_l[2], i_state->submsg_e_l[2]);
	ADD(i_state.submsg_e_l[2], i_state.submsg_o_l[2], i_state.submsg_e_l[2])
	//i_state->submsg_e_l[3] = ADD(i_state->submsg_o_l[3], i_state->submsg_e_l[3]);
	ADD(i_state.submsg_e_l[3], i_state.submsg_o_l[3], i_state.submsg_e_l[3])

	i_state.save(i_state.submsg_e_l, i_state.submsg_e_l_Mem)

	////////////////////////////////////////////////////////////////////////////////////////////////////

	i_state.submsg_e_r = i_state.tempVec0
	i_state.load(i_state.submsg_e_r, i_state.submsg_e_r_Mem)

	i_state.submsg_o_r = i_state.tempVec1
	i_state.load(i_state.submsg_o_r, i_state.submsg_o_r_Mem)

	//i_state->submsg_e_r[1] = _mm_shuffle_epi32(i_state->submsg_e_r[1], 0x4e);
	F_mm_shuffle_epi32(i_state.submsg_e_r[1], i_state.submsg_e_r[1], U8(0x4e))
	//temp = i_state->submsg_e_r[0];
	MOVO_autoAU2(temp, i_state.submsg_e_r[0])
	//i_state->submsg_e_r[0] = i_state->submsg_e_r[1];
	MOVO_autoAU2(i_state.submsg_e_r[0], i_state.submsg_e_r[1])
	//i_state->submsg_e_r[1] = temp;
	MOVO_autoAU2(i_state.submsg_e_r[1], temp)
	//i_state->submsg_e_r[3] = _mm_shuffle_epi32(i_state->submsg_e_r[3], 0x4e);
	F_mm_shuffle_epi32(i_state.submsg_e_r[3], i_state.submsg_e_r[3], U8(0x4e))
	//temp = i_state->submsg_e_r[2];
	MOVO_autoAU2(temp, i_state.submsg_e_r[2])
	//i_state->submsg_e_r[2] = _mm_unpacklo_epi64(i_state->submsg_e_r[3], i_state->submsg_e_r[2]);
	F_mm_unpacklo_epi64(i_state.submsg_e_r[2], i_state.submsg_e_r[3], i_state.submsg_e_r[2])
	//i_state->submsg_e_r[3] = _mm_unpackhi_epi64(temp, i_state->submsg_e_r[3]);
	F_mm_unpackhi_epi64(i_state.submsg_e_r[3], temp, i_state.submsg_e_r[3])

	//i_state->submsg_e_r[0] = ADD(i_state->submsg_o_r[0], i_state->submsg_e_r[0]);
	ADD(i_state.submsg_e_r[0], i_state.submsg_o_r[0], i_state.submsg_e_r[0])
	//i_state->submsg_e_r[1] = ADD(i_state->submsg_o_r[1], i_state->submsg_e_r[1]);
	ADD(i_state.submsg_e_r[1], i_state.submsg_o_r[1], i_state.submsg_e_r[1])
	//i_state->submsg_e_r[2] = ADD(i_state->submsg_o_r[2], i_state->submsg_e_r[2]);
	ADD(i_state.submsg_e_r[2], i_state.submsg_o_r[2], i_state.submsg_e_r[2])
	//i_state->submsg_e_r[3] = ADD(i_state->submsg_o_r[3], i_state->submsg_e_r[3]);
	ADD(i_state.submsg_e_r[3], i_state.submsg_o_r[3], i_state.submsg_e_r[3])

	i_state.save(i_state.submsg_e_r, i_state.submsg_e_r_Mem)
}

// static INLINE void msg_exp_odd(LSH512SSSE3_internal * i_state){
func msg_exp_odd(i_state LSH512SSSE3_internal) {
	if i_state.tempVec0 == nil || i_state.tempVec1 == nil {
		panic(`tempVec is required`)
	}
	Comment("msg_exp_odd")

	ADD := ADD64_

	//__m128i temp;
	temp := XMM()

	////////////////////////////////////////////////////////////////////////////////////////////////////

	i_state.submsg_o_l = i_state.tempVec0
	i_state.load(i_state.submsg_o_l, i_state.submsg_o_l_Mem)

	i_state.submsg_e_l = i_state.tempVec1
	i_state.load(i_state.submsg_e_l, i_state.submsg_e_l_Mem)

	//i_state->submsg_o_l[1] = _mm_shuffle_epi32(i_state->submsg_o_l[1], 0x4e);
	F_mm_shuffle_epi32(i_state.submsg_o_l[1], i_state.submsg_o_l[1], U8(0x4e))
	//temp = i_state->submsg_o_l[0];
	MOVO_autoAU2(temp, i_state.submsg_o_l[0])
	//i_state->submsg_o_l[0] = i_state->submsg_o_l[1];
	MOVO_autoAU2(i_state.submsg_o_l[0], i_state.submsg_o_l[1])
	//i_state->submsg_o_l[1] = temp;
	MOVO_autoAU2(i_state.submsg_o_l[1], temp)
	//i_state->submsg_o_l[3] = _mm_shuffle_epi32(i_state->submsg_o_l[3], 0x4e);
	F_mm_shuffle_epi32(i_state.submsg_o_l[3], i_state.submsg_o_l[3], U8(0x4e))
	//temp = i_state->submsg_o_l[2];
	MOVO_autoAU2(temp, i_state.submsg_o_l[2])
	//i_state->submsg_o_l[2] = _mm_unpacklo_epi64(i_state->submsg_o_l[3], i_state->submsg_o_l[2]);
	F_mm_unpacklo_epi64(i_state.submsg_o_l[2], i_state.submsg_o_l[3], i_state.submsg_o_l[2])
	//i_state->submsg_o_l[3] = _mm_unpackhi_epi64(temp, i_state->submsg_o_l[3]);
	F_mm_unpackhi_epi64(i_state.submsg_o_l[3], temp, i_state.submsg_o_l[3])

	//i_state->submsg_o_l[0] = ADD(i_state->submsg_e_l[0], i_state->submsg_o_l[0]);
	//i_state->submsg_o_l[1] = ADD(i_state->submsg_e_l[1], i_state->submsg_o_l[1]);
	//i_state->submsg_o_l[2] = ADD(i_state->submsg_e_l[2], i_state->submsg_o_l[2]);
	//i_state->submsg_o_l[3] = ADD(i_state->submsg_e_l[3], i_state->submsg_o_l[3]);
	ADD(i_state.submsg_o_l[0], i_state.submsg_e_l[0], i_state.submsg_o_l[0])
	ADD(i_state.submsg_o_l[1], i_state.submsg_e_l[1], i_state.submsg_o_l[1])
	ADD(i_state.submsg_o_l[2], i_state.submsg_e_l[2], i_state.submsg_o_l[2])
	ADD(i_state.submsg_o_l[3], i_state.submsg_e_l[3], i_state.submsg_o_l[3])

	i_state.save(i_state.submsg_o_l, i_state.submsg_o_l_Mem)

	////////////////////////////////////////////////////////////////////////////////////////////////////

	i_state.submsg_o_r = i_state.tempVec0
	i_state.load(i_state.submsg_o_r, i_state.submsg_o_r_Mem)

	i_state.submsg_e_r = i_state.tempVec1
	i_state.load(i_state.submsg_e_r, i_state.submsg_e_r_Mem)

	//i_state->submsg_o_r[1] = _mm_shuffle_epi32(i_state->submsg_o_r[1], 0x4e);
	F_mm_shuffle_epi32(i_state.submsg_o_r[1], i_state.submsg_o_r[1], U8(0x4e))
	//temp = i_state->submsg_o_r[0];
	MOVO_autoAU2(temp, i_state.submsg_o_r[0])
	//i_state->submsg_o_r[0] = i_state->submsg_o_r[1];
	MOVO_autoAU2(i_state.submsg_o_r[0], i_state.submsg_o_r[1])
	//i_state->submsg_o_r[1] = temp;
	MOVO_autoAU2(i_state.submsg_o_r[1], temp)
	//i_state->submsg_o_r[3] = _mm_shuffle_epi32(i_state->submsg_o_r[3], 0x4e);
	F_mm_shuffle_epi32(i_state.submsg_o_r[3], i_state.submsg_o_r[3], U8(0x4e))
	//temp = i_state->submsg_o_r[2];
	MOVO_autoAU2(temp, i_state.submsg_o_r[2])
	//i_state->submsg_o_r[2] = _mm_unpacklo_epi64(i_state->submsg_o_r[3], i_state->submsg_o_r[2]);
	F_mm_unpacklo_epi64(i_state.submsg_o_r[2], i_state.submsg_o_r[3], i_state.submsg_o_r[2])
	//i_state->submsg_o_r[3] = _mm_unpackhi_epi64(temp, i_state->submsg_o_r[3]);
	F_mm_unpackhi_epi64(i_state.submsg_o_r[3], temp, i_state.submsg_o_r[3])

	//i_state->submsg_o_r[0] = ADD(i_state->submsg_e_r[0], i_state->submsg_o_r[0]);
	//i_state->submsg_o_r[1] = ADD(i_state->submsg_e_r[1], i_state->submsg_o_r[1]);
	//i_state->submsg_o_r[2] = ADD(i_state->submsg_e_r[2], i_state->submsg_o_r[2]);
	//i_state->submsg_o_r[3] = ADD(i_state->submsg_e_r[3], i_state->submsg_o_r[3]);
	ADD(i_state.submsg_o_r[0], i_state.submsg_e_r[0], i_state.submsg_o_r[0])
	ADD(i_state.submsg_o_r[1], i_state.submsg_e_r[1], i_state.submsg_o_r[1])
	ADD(i_state.submsg_o_r[2], i_state.submsg_e_r[2], i_state.submsg_o_r[2])
	ADD(i_state.submsg_o_r[3], i_state.submsg_e_r[3], i_state.submsg_o_r[3])

	i_state.save(i_state.submsg_o_r, i_state.submsg_o_r_Mem)
}

// static INLINE void load_sc(__m128i* const_v, lsh_uint i){
func load_sc(const_v []VecVirtual, i int) {
	Comment("load_sc")

	//load_blk(const_v, g_StepConstants + i);
	load_blk_mem2vec(const_v, G_StepConstants.Offset(i*8))

}

// static INLINE void msg_add_even(__m128i* cv_l, __m128i* cv_r, const LSH512SSSE3_internal * i_state){
func msg_add_even(cv_l, cv_r []VecVirtual, i_state LSH512SSSE3_internal, tempVec []VecVirtual) {
	Comment("msg_add_even")

	//cv_l[0] = XOR(cv_l[0], i_state->submsg_e_l[0]);
	//cv_r[0] = XOR(cv_r[0], i_state->submsg_e_r[0]);
	//cv_l[1] = XOR(cv_l[1], i_state->submsg_e_l[1]);
	//cv_r[1] = XOR(cv_r[1], i_state->submsg_e_r[1]);
	//cv_l[2] = XOR(cv_l[2], i_state->submsg_e_l[2]);
	//cv_r[2] = XOR(cv_r[2], i_state->submsg_e_r[2]);
	//cv_l[3] = XOR(cv_l[3], i_state->submsg_e_l[3]);
	//cv_r[3] = XOR(cv_r[3], i_state->submsg_e_r[3]);

	load_blk_mem2vec(tempVec, i_state.submsg_e_l_Mem)
	XOR(cv_l[0], tempVec[0])
	XOR(cv_l[1], tempVec[1])
	XOR(cv_l[2], tempVec[2])
	XOR(cv_l[3], tempVec[3])

	load_blk_mem2vec(tempVec, i_state.submsg_e_r_Mem)
	XOR(cv_r[0], tempVec[0])
	XOR(cv_r[1], tempVec[1])
	XOR(cv_r[2], tempVec[2])
	XOR(cv_r[3], tempVec[3])
}

// static INLINE void msg_add_odd(__m128i* cv_l, __m128i* cv_r, const LSH512SSSE3_internal * i_state){
func msg_add_odd(cv_l, cv_r []VecVirtual, i_state LSH512SSSE3_internal, tempVec []VecVirtual) {
	Comment("msg_add_odd")

	//cv_l[0] = XOR(cv_l[0], i_state->submsg_o_l[0]);
	//cv_r[0] = XOR(cv_r[0], i_state->submsg_o_r[0]);
	//cv_l[1] = XOR(cv_l[1], i_state->submsg_o_l[1]);
	//cv_r[1] = XOR(cv_r[1], i_state->submsg_o_r[1]);
	//cv_l[2] = XOR(cv_l[2], i_state->submsg_o_l[2]);
	//cv_r[2] = XOR(cv_r[2], i_state->submsg_o_r[2]);
	//cv_l[3] = XOR(cv_l[3], i_state->submsg_o_l[3]);
	//cv_r[3] = XOR(cv_r[3], i_state->submsg_o_r[3]);

	load_blk_mem2vec(tempVec, i_state.submsg_o_l_Mem)
	XOR(cv_l[0], tempVec[0])
	XOR(cv_l[1], tempVec[1])
	XOR(cv_l[2], tempVec[2])
	XOR(cv_l[3], tempVec[3])

	load_blk_mem2vec(tempVec, i_state.submsg_o_r_Mem)
	XOR(cv_r[0], tempVec[0])
	XOR(cv_r[1], tempVec[1])
	XOR(cv_r[2], tempVec[2])
	XOR(cv_r[3], tempVec[3])
}

// static INLINE void add_blk(__m128i* cv_l, const __m128i* cv_r){
func add_blk(cv_l, cv_r []VecVirtual) {
	Comment("add_blk")

	//cv_l[0] = ADD(cv_l[0], cv_r[0]);
	ADD64(cv_l[0], cv_r[0])
	//cv_l[1] = ADD(cv_l[1], cv_r[1]);
	ADD64(cv_l[1], cv_r[1])
	//cv_l[2] = ADD(cv_l[2], cv_r[2]);
	ADD64(cv_l[2], cv_r[2])
	//cv_l[3] = ADD(cv_l[3], cv_r[3]);
	ADD64(cv_l[3], cv_r[3])
}

func rotate_blk(dst VecVirtual, v int) {
	tmpXmm := XMM()

	// dst = OR(SHIFT_L(dst, ROT_EVEN_ALPHA), SHIFT_R(dst, WORD_BIT_LEN - ROT_EVEN_ALPHA))
	MOVO_autoAU2(tmpXmm, dst)          //         tmpXmm = dst
	SHIFT_L64(tmpXmm, U8(v))           // tmpXmm = SHIFT_L(dst, ROT_EVEN_ALPHA)
	SHIFT_R64(dst, U8(WORD_BIT_LEN-v)) //                                  dst = SHIFT_R(dst, WORD_BIT_LEN - ROT_EVEN_ALPHA)
	OR(dst, tmpXmm)                    // dst = OR(SHIFT_L(dst, ROT_EVEN_ALPHA), SHIFT_R(dst, WORD_BIT_LEN - ROT_EVEN_ALPHA))
}

// static INLINE void rotate_blk_even_alpha(__m128i* cv){
func rotate_blk_even_alpha(cv []VecVirtual) {
	Comment("rotate_blk_even_alpha")

	//cv[0] = OR(SHIFT_L(cv[0], ROT_EVEN_ALPHA), SHIFT_R(cv[0], WORD_BIT_LEN - ROT_EVEN_ALPHA));
	//cv[1] = OR(SHIFT_L(cv[1], ROT_EVEN_ALPHA), SHIFT_R(cv[1], WORD_BIT_LEN - ROT_EVEN_ALPHA));
	//cv[2] = OR(SHIFT_L(cv[2], ROT_EVEN_ALPHA), SHIFT_R(cv[2], WORD_BIT_LEN - ROT_EVEN_ALPHA));
	//cv[3] = OR(SHIFT_L(cv[3], ROT_EVEN_ALPHA), SHIFT_R(cv[3], WORD_BIT_LEN - ROT_EVEN_ALPHA));
	rotate_blk(cv[0], ROT_EVEN_ALPHA)
	rotate_blk(cv[1], ROT_EVEN_ALPHA)
	rotate_blk(cv[2], ROT_EVEN_ALPHA)
	rotate_blk(cv[3], ROT_EVEN_ALPHA)
}

// static INLINE void rotate_blk_even_beta(__m128i* cv){
func rotate_blk_even_beta(cv []VecVirtual) {
	Comment("rotate_blk_even_beta")

	//cv[0] = OR(SHIFT_L(cv[0], ROT_EVEN_BETA), SHIFT_R(cv[0], WORD_BIT_LEN - ROT_EVEN_BETA));
	//cv[1] = OR(SHIFT_L(cv[1], ROT_EVEN_BETA), SHIFT_R(cv[1], WORD_BIT_LEN - ROT_EVEN_BETA));
	//cv[2] = OR(SHIFT_L(cv[2], ROT_EVEN_BETA), SHIFT_R(cv[2], WORD_BIT_LEN - ROT_EVEN_BETA));
	//cv[3] = OR(SHIFT_L(cv[3], ROT_EVEN_BETA), SHIFT_R(cv[3], WORD_BIT_LEN - ROT_EVEN_BETA));
	rotate_blk(cv[0], ROT_EVEN_BETA)
	rotate_blk(cv[1], ROT_EVEN_BETA)
	rotate_blk(cv[2], ROT_EVEN_BETA)
	rotate_blk(cv[3], ROT_EVEN_BETA)
}

// static INLINE void rotate_blk_odd_alpha(__m128i* cv){
func rotate_blk_odd_alpha(cv []VecVirtual) {
	Comment("rotate_blk_odd_alpha")

	//cv[0] = OR(SHIFT_L(cv[0], ROT_ODD_ALPHA), SHIFT_R(cv[0], WORD_BIT_LEN - ROT_ODD_ALPHA));
	//cv[1] = OR(SHIFT_L(cv[1], ROT_ODD_ALPHA), SHIFT_R(cv[1], WORD_BIT_LEN - ROT_ODD_ALPHA));
	//cv[2] = OR(SHIFT_L(cv[2], ROT_ODD_ALPHA), SHIFT_R(cv[2], WORD_BIT_LEN - ROT_ODD_ALPHA));
	//cv[3] = OR(SHIFT_L(cv[3], ROT_ODD_ALPHA), SHIFT_R(cv[3], WORD_BIT_LEN - ROT_ODD_ALPHA));
	rotate_blk(cv[0], ROT_ODD_ALPHA)
	rotate_blk(cv[1], ROT_ODD_ALPHA)
	rotate_blk(cv[2], ROT_ODD_ALPHA)
	rotate_blk(cv[3], ROT_ODD_ALPHA)
}

// static INLINE void rotate_blk_odd_beta(__m128i* cv){
func rotate_blk_odd_beta(cv []VecVirtual) {
	Comment("rotate_blk_odd_beta")

	//cv[0] = OR(SHIFT_L(cv[0], ROT_ODD_BETA), SHIFT_R(cv[0], WORD_BIT_LEN - ROT_ODD_BETA));
	//cv[1] = OR(SHIFT_L(cv[1], ROT_ODD_BETA), SHIFT_R(cv[1], WORD_BIT_LEN - ROT_ODD_BETA));
	//cv[2] = OR(SHIFT_L(cv[2], ROT_ODD_BETA), SHIFT_R(cv[2], WORD_BIT_LEN - ROT_ODD_BETA));
	//cv[3] = OR(SHIFT_L(cv[3], ROT_ODD_BETA), SHIFT_R(cv[3], WORD_BIT_LEN - ROT_ODD_BETA));
	rotate_blk(cv[0], ROT_ODD_BETA)
	rotate_blk(cv[1], ROT_ODD_BETA)
	rotate_blk(cv[2], ROT_ODD_BETA)
	rotate_blk(cv[3], ROT_ODD_BETA)
}

// static INLINE void xor_with_const(__m128i* cv_l, const __m128i* const_v){
func xor_with_const(cv_l []VecVirtual, const_v []VecVirtual) {
	//cv_l[0] = XOR(cv_l[0], const_v[0]);
	XOR(cv_l[0], const_v[0])
	//cv_l[1] = XOR(cv_l[1], const_v[1]);
	XOR(cv_l[1], const_v[1])
	//cv_l[2] = XOR(cv_l[2], const_v[2]);
	XOR(cv_l[2], const_v[2])
	//cv_l[3] = XOR(cv_l[3], const_v[3]);
	XOR(cv_l[3], const_v[3])
}

// static INLINE void rotate_msg_gamma(__m128i* cv_r, const __m128i * perm_step){
func rotate_msg_gamma(cv_r []VecVirtual, perm_step []Mem) {
	Comment("rotate_msg_gamma")

	//cv_r[0] = SHUFFLE8(cv_r[0], perm_step[0]);
	SHUFFLE8(cv_r[0], perm_step[0])
	//cv_r[1] = SHUFFLE8(cv_r[1], perm_step[1]);
	SHUFFLE8(cv_r[1], perm_step[1])
	//cv_r[2] = SHUFFLE8(cv_r[2], perm_step[2]);
	SHUFFLE8(cv_r[2], perm_step[2])
	//cv_r[3] = SHUFFLE8(cv_r[3], perm_step[3]);
	SHUFFLE8(cv_r[3], perm_step[3])
}

// static INLINE void word_perm(__m128i* cv_l, __m128i* cv_r){
func word_perm(cv_l, cv_r []VecVirtual, temp []VecVirtual) {
	Comment("word_perm")
	//__m128i temp[2];

	//temp[0] = cv_l[0];
	MOVO_autoAU2(temp[0], cv_l[0])
	//cv_l[0] = _mm_unpacklo_epi64(cv_l[1], cv_l[0]);
	F_mm_unpacklo_epi64(cv_l[0], cv_l[1], cv_l[0])
	//cv_l[1] = _mm_unpackhi_epi64(temp[0], cv_l[1]);
	F_mm_unpackhi_epi64(cv_l[1], temp[0], cv_l[1])
	//temp[0] = cv_l[2];
	MOVO_autoAU2(temp[0], cv_l[2])
	//cv_l[2] = _mm_unpacklo_epi64(cv_l[3], cv_l[2]);
	F_mm_unpacklo_epi64(cv_l[2], cv_l[3], cv_l[2])
	//cv_l[3] = _mm_unpackhi_epi64(temp[0], cv_l[3]);
	F_mm_unpackhi_epi64(cv_l[3], temp[0], cv_l[3])
	//cv_r[1] = _mm_shuffle_epi32(cv_r[1], 0x4e);
	F_mm_shuffle_epi32(cv_r[1], cv_r[1], U8(0x4e))
	//temp[0] = cv_r[0];
	MOVO_autoAU2(temp[0], cv_r[0])
	//cv_r[0] = _mm_unpacklo_epi64(cv_r[0], cv_r[1]);
	F_mm_unpacklo_epi64(cv_r[0], cv_r[0], cv_r[1])
	//cv_r[1] = _mm_unpackhi_epi64(cv_r[1], temp[0]);
	F_mm_unpackhi_epi64(cv_r[1], cv_r[1], temp[0])
	//cv_r[3] = _mm_shuffle_epi32(cv_r[3], 0x4e);
	F_mm_shuffle_epi32(cv_r[3], cv_r[3], U8(0x4e))
	//temp[0] = cv_r[2];
	MOVO_autoAU2(temp[0], cv_r[2])
	//cv_r[2] = _mm_unpacklo_epi64(cv_r[2], cv_r[3]);
	F_mm_unpacklo_epi64(cv_r[2], cv_r[2], cv_r[3])
	//cv_r[3] = _mm_unpackhi_epi64(cv_r[3], temp[0]);
	F_mm_unpackhi_epi64(cv_r[3], cv_r[3], temp[0])
	//temp[0] = cv_l[0];
	MOVO_autoAU2(temp[0], cv_l[0])
	//temp[1] = cv_l[1];
	MOVO_autoAU2(temp[1], cv_l[1])
	//cv_l[0] = cv_l[2];
	MOVO_autoAU2(cv_l[0], cv_l[2])
	//cv_l[1] = cv_l[3];
	MOVO_autoAU2(cv_l[1], cv_l[3])
	//cv_l[2] = cv_r[2];
	MOVO_autoAU2(cv_l[2], cv_r[2])
	//cv_l[3] = cv_r[3];
	MOVO_autoAU2(cv_l[3], cv_r[3])
	//cv_r[2] = cv_r[0];
	MOVO_autoAU2(cv_r[2], cv_r[0])
	//cv_r[3] = cv_r[1];
	MOVO_autoAU2(cv_r[3], cv_r[1])
	//cv_r[0] = temp[0];
	MOVO_autoAU2(cv_r[0], temp[0])
	//cv_r[1] = temp[1];
	MOVO_autoAU2(cv_r[1], temp[1])
}

/* -------------------------------------------------------- */
// step function
/* -------------------------------------------------------- */

// static INLINE void mix_even(__m128i* cv_l, __m128i* cv_r, const __m128i* const_v, const __m128i * perm_step){
func mix_even(cv_l, cv_r []VecVirtual, const_v []VecVirtual, perm_step []Mem) {
	Comment("mix_even")

	add_blk(cv_l, cv_r)
	rotate_blk_even_alpha(cv_l)
	xor_with_const(cv_l, const_v)
	add_blk(cv_r, cv_l)
	rotate_blk_even_beta(cv_r)
	add_blk(cv_l, cv_r)
	rotate_msg_gamma(cv_r, perm_step)
}

// static INLINE void mix_odd(__m128i* cv_l, __m128i* cv_r, const __m128i* const_v, const __m128i * perm_step){
func mix_odd(cv_l, cv_r []VecVirtual, const_v []VecVirtual, perm_step []Mem) {
	Comment("mix_odd")

	add_blk(cv_l, cv_r)
	rotate_blk_odd_alpha(cv_l)
	xor_with_const(cv_l, const_v)
	add_blk(cv_r, cv_l)
	rotate_blk_odd_beta(cv_r)
	add_blk(cv_l, cv_r)
	rotate_msg_gamma(cv_r, perm_step)
}

/* -------------------------------------------------------- */
// compression function
/* -------------------------------------------------------- */

// static INLINE void compress(__m128i* cv_l, __m128i* cv_r, const lsh_u64 pdMsgBlk[MSG_BLK_WORD_LEN])
func compress(cv_l, cv_r []VecVirtual, pdMsgBlk Mem, i_state_alloc, cv_l_mem, cv_r_mem Mem) {
	Comment("compress")
	xmmTemp := []VecVirtual{XMM(), XMM(), XMM(), XMM()}

	//__m128i const_v[4];			// step function constant
	const_v := xmmTemp
	// __m128i perm_step[4];	// byte permutation info
	perm_step := make([]Mem, 4)
	//LSH512SSSE3_internal i_state[1];
	i_state := LSH512SSSE3_internal{
		tempVec0:       cv_l,
		tempVec1:       cv_r,
		submsg_e_l_Mem: i_state_alloc.Offset(XmmSize * 4 * 0),
		submsg_e_r_Mem: i_state_alloc.Offset(XmmSize * 4 * 1),
		submsg_o_l_Mem: i_state_alloc.Offset(XmmSize * 4 * 2),
		submsg_o_r_Mem: i_state_alloc.Offset(XmmSize * 4 * 3),
	}
	//int i;

	//perm_step[0] = LOAD(g_BytePermInfo[0]);
	perm_step[0] = g_BytePermInfo.Offset(XmmSize * 0)
	//perm_step[1] = LOAD(g_BytePermInfo[1]);
	perm_step[1] = g_BytePermInfo.Offset(XmmSize * 1)
	//perm_step[2] = LOAD(g_BytePermInfo[2]);
	perm_step[2] = g_BytePermInfo.Offset(XmmSize * 2)
	//perm_step[3] = LOAD(g_BytePermInfo[3]);
	perm_step[3] = g_BytePermInfo.Offset(XmmSize * 3)

	save := func() {
		Comment("save___start")
		load_blk_vec2mem(cv_l_mem, cv_l)
		load_blk_vec2mem(cv_r_mem, cv_r)
		Comment("save___end")
	}
	load := func() {
		Comment("load___start")
		load_blk_mem2vec(cv_l, cv_l_mem)
		load_blk_mem2vec(cv_r, cv_r_mem)
		Comment("load___end")
	}

	load_msg_blk(i_state, pdMsgBlk)

	msg_add_even(cv_l, cv_r, i_state, xmmTemp)
	load_sc(const_v, 0)
	mix_even(cv_l, cv_r, const_v, perm_step)
	word_perm(cv_l, cv_r, xmmTemp)

	msg_add_odd(cv_l, cv_r, i_state, xmmTemp)
	load_sc(const_v, 8)
	mix_odd(cv_l, cv_r, const_v, perm_step)
	word_perm(cv_l, cv_r, xmmTemp)

	//for (i = 1; i < NUM_STEPS / 2; i++){
	for i := 1; i < NUM_STEPS/2; i++ {
		save()
		msg_exp_even(i_state)
		load()
		msg_add_even(cv_l, cv_r, i_state, xmmTemp)
		load_sc(const_v, i*16)
		mix_even(cv_l, cv_r, const_v, perm_step)
		word_perm(cv_l, cv_r, xmmTemp)

		save()
		msg_exp_odd(i_state)
		load()
		msg_add_odd(cv_l, cv_r, i_state, xmmTemp)
		load_sc(const_v, i*16+8)
		mix_odd(cv_l, cv_r, const_v, perm_step)
		word_perm(cv_l, cv_r, xmmTemp)
	}

	save()
	msg_exp_even(i_state)
	load()
	msg_add_even(cv_l, cv_r, i_state, xmmTemp)
}

/* -------------------------------------------------------- */
//static INLINE void init224(LSH512SSSE3_Context* state)
func init224(state *LSH512_Context) {
	Comment("init224")

	//load_blk(state->cv_l, g_IV224);
	load_blk_mem2mem(state.Cv_l, G_IV224)
	//load_blk(state->cv_r, g_IV224 + 8);
	load_blk_mem2mem(state.Cv_r, G_IV224.Offset(8*8))
}

// static INLINE void init256(LSH512SSSE3_Context* state)
func init256(state *LSH512_Context) {
	Comment("init256")

	//load_blk(state->cv_l, g_IV256);
	load_blk_mem2mem(state.Cv_l, G_IV256)
	//load_blk(state->cv_r, g_IV256 + 8);
	load_blk_mem2mem(state.Cv_r, G_IV256.Offset(8*8))
}

// static INLINE void init384(LSH512SSSE3_Context* state)
func init384(state *LSH512_Context) {
	Comment("init384")

	//load_blk(state->cv_l, g_IV384);
	load_blk_mem2mem(state.Cv_l, G_IV384)
	//load_blk(state->cv_r, g_IV384 + 8);
	load_blk_mem2mem(state.Cv_r, G_IV384.Offset(8*8))
}

// static INLINE void init512(LSH512SSSE3_Context* state)
func init512(state *LSH512_Context) {
	Comment("init512")

	//load_blk(state->cv_l, g_IV512);
	load_blk_mem2mem(state.Cv_l, G_IV512)
	//load_blk(state->cv_r, g_IV512 + 8);
	load_blk_mem2mem(state.Cv_r, G_IV512.Offset(8*8))
}

/* -------------------------------------------------------- */

// static INLINE void fin(__m128i *cv_l, const __m128i *cv_r)
func fin(cv_l, cv_r []VecVirtual) {
	Comment("fin")

	//cv_l[0] = XOR(cv_l[0], cv_r[0]);
	XOR(cv_l[0], cv_r[0])
	//cv_l[1] = XOR(cv_l[1], cv_r[1]);
	XOR(cv_l[1], cv_r[1])
	//cv_l[2] = XOR(cv_l[2], cv_r[2]);
	XOR(cv_l[2], cv_r[2])
	//cv_l[3] = XOR(cv_l[3], cv_r[3]);
	XOR(cv_l[3], cv_r[3])
}

/* -------------------------------------------------------- */

// static INLINE void get_hash(__m128i *cv_l, lsh_u8 * pbHashVal, const lsh_type algtype)
func get_hash(cv_l []VecVirtual, pbHashVal Mem, algtype Op) {
	Comment("get_hash")

	//lsh_u8 hash_val[LSH512_HASH_VAL_MAX_BYTE_LEN] = { 0x0, };
	hash_val := pbHashVal
	//lsh_uint hash_val_byte_len = LSH_GET_HASHBYTE(algtype);
	//lsh_uint hash_val_bit_len = LSH_GET_SMALL_HASHBIT(algtype);

	//store_blk((__m128i*)hash_val, cv_l);
	store_blk(hash_val, cv_l)
	//memcpy(pbHashVal, hash_val, sizeof(lsh_u8) * hash_val_byte_len);
	//if (hash_val_bit_len){
	//	pbHashVal[hash_val_byte_len-1] &= (((lsh_u8)0xff) << hash_val_bit_len);
	//}
}

/* -------------------------------------------------------- */

// lsh_err lsh512_ssse3_init(struct LSH512_Context * _ctx, const lsh_type algtype){
func Lsh512_ssse3_init(ctx *LSH512_Context) {
	Comment("lsh512_ssse3_init")

	//LSH512SSSE3_Context* ctx = (LSH512SSSE3_Context*)_ctx;
	//__m128i cv_l[4];
	//__m128i cv_r[4];
	//__m128i const_v[4];
	//lsh_uint i;

	//if (ctx == NULL){
	//	return LSH_ERR_NULL_PTR;
	//}

	//ctx->algtype = algtype;
	//ctx->remain_databitlen = 0;

	//if (!LSH_IS_LSH512(algtype)){
	//	return LSH_ERR_INVALID_ALGTYPE;
	//}

	//if (LSH_GET_HASHBYTE(algtype) > LSH512_HASH_VAL_MAX_BYTE_LEN || LSH_GET_HASHBYTE(algtype) == 0){
	//	return LSH_ERR_INVALID_ALGTYPE;
	//}

	//switch (algtype){
	//case LSH_TYPE_512_512:
	CMPL(ctx.Algtype, U32(LSH_TYPE_512_512))
	JNE(LabelRef("lsh512_ssse3_init_if0_end"))
	{
		//	init512(ctx);
		init512(ctx)
		//	return LSH_SUCCESS;
		JMP(LabelRef("lsh512_ssse3_init_ret"))
	}
	Label("lsh512_ssse3_init_if0_end")
	//case LSH_TYPE_512_384:
	CMPL(ctx.Algtype, U32(LSH_TYPE_512_384))
	JNE(LabelRef("lsh512_ssse3_init_if1_end"))
	{
		//	init384(ctx);
		init384(ctx)
		//	return LSH_SUCCESS;
		JMP(LabelRef("lsh512_ssse3_init_ret"))
	}
	Label("lsh512_ssse3_init_if1_end")
	//case LSH_TYPE_512_256:
	CMPL(ctx.Algtype, U32(LSH_TYPE_512_256))
	JNE(LabelRef("lsh512_ssse3_init_if2_end"))
	{
		//	init256(ctx);
		init256(ctx)
		//	return LSH_SUCCESS;
		JMP(LabelRef("lsh512_ssse3_init_ret"))
	}
	Label("lsh512_ssse3_init_if2_end")
	//case LSH_TYPE_512_224:
	{
		//	init224(ctx);
		init224(ctx)
		//	return LSH_SUCCESS;
		JMP(LabelRef("lsh512_ssse3_init_ret"))
	}
	//default:
	//	break;
	//}

	//cv_l[0] = _mm_set_epi32(0, LSH_GET_HASHBIT(algtype), 0, LSH512_HASH_VAL_MAX_BYTE_LEN);
	//cv_l[1] = _mm_setzero_si128();
	//cv_l[2] = _mm_setzero_si128();
	//cv_l[3] = _mm_setzero_si128();
	//cv_r[0] = _mm_setzero_si128();
	//cv_r[1] = _mm_setzero_si128();
	//cv_r[2] = _mm_setzero_si128();
	//cv_r[3] = _mm_setzero_si128();
	//perm_step[0] = LOAD(g_BytePermInfo[0]);
	//perm_step[1] = LOAD(g_BytePermInfo[1]);
	//perm_step[2] = LOAD(g_BytePermInfo[2]);
	//perm_step[3] = LOAD(g_BytePermInfo[3]);

	//for (i = 0; i < NUM_STEPS / 2; i++)
	//{
	//	//Mix
	//	load_sc(const_v, i * 16);
	//	mix_even(cv_l, cv_r, const_v, perm_step);
	//	word_perm(cv_l, cv_r);

	//	load_sc(const_v, i * 16 + 8);
	//	mix_odd(cv_l, cv_r, const_v, perm_step);
	//	word_perm(cv_l, cv_r);
	//}

	//store_blk(ctx->cv_l, cv_l);
	//store_blk(ctx->cv_r, cv_r);

	//return LSH_SUCCESS;
	Label("lsh512_ssse3_init_ret")
}

// lsh_err lsh512_ssse3_update(struct LSH512_Context * _ctx, const lsh_u8 * data, size_t databitlen){
func Lsh512_ssse3_update(ctx *LSH512_Context, data Mem, databytelen Register) {
	Comment("lsh512_ssse3_update")

	i_state_alloc := AllocLocal(XmmSize * 4 * 4)

	//__m128i cv_l[4];
	cv_l := []VecVirtual{XMM(), XMM(), XMM(), XMM()}
	//__m128i cv_r[4];
	cv_r := []VecVirtual{XMM(), XMM(), XMM(), XMM()}
	//size_t databytelen = databitlen >> 3;
	//lsh_u32 pos2 = databitlen & 0x7;

	//LSH512SSSE3_Context* ctx = (LSH512SSSE3_Context*)_ctx;
	//lsh_uint remain_msg_byte;
	remain_msg_byte := GP64()
	//lsh_uint remain_msg_bit;

	//if (ctx == NULL || data == NULL){
	//	return LSH_ERR_NULL_PTR;
	//}
	//if (ctx->algtype == 0 || LSH_GET_HASHBYTE(ctx->algtype) > LSH512_HASH_VAL_MAX_BYTE_LEN){
	//	return LSH_ERR_INVALID_STATE;
	//}
	//if (databitlen == 0){
	//	return LSH_SUCCESS;
	//}

	//remain_msg_byte = ctx->remain_databitlen >> 3;
	MOVQ(ctx.Remain_databytelen, remain_msg_byte)
	//remain_msg_bit = ctx->remain_databitlen & 7;
	//if (remain_msg_byte >= LSH512_MSG_BLK_BYTE_LEN){
	//	return LSH_ERR_INVALID_STATE;
	//}
	//if (remain_msg_bit > 0){
	//	return LSH_ERR_INVALID_DATABITLEN;
	//}

	//if (databytelen + remain_msg_byte < LSH512_MSG_BLK_BYTE_LEN){
	tmp32 := GP64()
	MOVQ(databytelen, tmp32)
	ADDQ(remain_msg_byte, tmp32)
	CMPQ(tmp32, U32(LSH512_MSG_BLK_BYTE_LEN))
	JGE(LabelRef("lsh512_ssse3_update_if0_end"))
	{
		//memcpy(ctx->i_last_block + remain_msg_byte, data, databytelen);
		Memcpy(ctx.I_last_block.Idx(remain_msg_byte, 1), data, databytelen, false)
		//ctx->remain_databitlen += (lsh_uint)databitlen;
		ADDQ(databytelen, ctx.Remain_databytelen)
		//remain_msg_byte += (lsh_uint)databytelen;
		ADDQ(databytelen, remain_msg_byte)
		//if (pos2){
		//	ctx->i_last_block[remain_msg_byte] = data[databytelen] & ((0xff >> pos2) ^ 0xff);
		//}
		//return LSH_SUCCESS;
		JMP(LabelRef("lsh512_ssse3_update_ret"))
	}
	Label("lsh512_ssse3_update_if0_end")

	//load_blk(cv_l, ctx->cv_l);
	load_blk_mem2vec(cv_l, ctx.Cv_l)
	//load_blk(cv_r, ctx->cv_r);
	load_blk_mem2vec(cv_r, ctx.Cv_r)

	//if (remain_msg_byte > 0){
	CMPQ(remain_msg_byte, U32(0))
	JE(LabelRef("lsh512_ssse3_update_if1_end"))
	{
		//size_t more_BYTE = LSH512_MSG_BLK_BYTE_LEN - remain_msg_byte;
		more_BYTE := GP64()
		MOVQ(U32(LSH512_MSG_BLK_BYTE_LEN), more_BYTE)
		SUBQ(remain_msg_byte, more_BYTE)
		//memcpy(ctx->i_last_block + remain_msg_byte, data, more_BYTE);
		Memcpy(ctx.I_last_block.Idx(remain_msg_byte, 1), data, more_BYTE, false)
		//compress(cv_l, cv_r, (lsh_u64*)ctx->i_last_block);
		compress(cv_l, cv_r, ctx.I_last_block, i_state_alloc, ctx.Cv_l, ctx.Cv_r)
		//data += more_BYTE;
		ADDQ(more_BYTE, data.Base)
		//databytelen -= more_BYTE;
		SUBQ(more_BYTE, databytelen)
		//remain_msg_byte = 0;
		MOVQ(U32(0), remain_msg_byte)
		//ctx->remain_databitlen = 0;
		MOVQ(U32(0), ctx.Remain_databytelen)
	}
	Label("lsh512_ssse3_update_if1_end")

	//while (databytelen >= LSH512_MSG_BLK_BYTE_LEN)
	Label("lsh512_ssse3_update_while_start")
	CMPQ(databytelen, U32(LSH512_MSG_BLK_BYTE_LEN))
	JL(LabelRef("lsh512_ssse3_update_while_end"))
	{
		//compress(cv_l, cv_r, (lsh_u64*)data);
		compress(cv_l, cv_r, data, i_state_alloc, ctx.Cv_l, ctx.Cv_r)
		//data += LSH512_MSG_BLK_BYTE_LEN;
		ADDQ(U32(LSH512_MSG_BLK_BYTE_LEN), data.Base)
		//databytelen -= LSH512_MSG_BLK_BYTE_LEN;
		SUBQ(U32(LSH512_MSG_BLK_BYTE_LEN), databytelen)

		JMP(LabelRef("lsh512_ssse3_update_while_start"))
	}
	Label("lsh512_ssse3_update_while_end")

	//store_blk(ctx->cv_l, cv_l);
	store_blk(ctx.Cv_l, cv_l)
	//store_blk(ctx->cv_r, cv_r);
	store_blk(ctx.Cv_r, cv_r)

	//if (databytelen > 0){
	CMPQ(remain_msg_byte, U32(0))
	JE(LabelRef("lsh512_ssse3_update_if3_end"))
	{
		//memcpy(ctx->i_last_block, data, databytelen);
		Memcpy(ctx.I_last_block, data, databytelen, false)
		//ctx->remain_databitlen = (lsh_uint)(databytelen << 3);
		MOVQ(databytelen, ctx.Remain_databytelen)
	}
	Label("lsh512_ssse3_update_if3_end")

	//if (pos2){
	//	ctx->i_last_block[databytelen] = data[databytelen] & ((0xff >> pos2) ^ 0xff);
	//	ctx->remain_databitlen += pos2;
	//}
	//return LSH_SUCCESS;

	Label("lsh512_ssse3_update_ret")
}

// lsh_err lsh512_ssse3_final(struct LSH512_Context * _ctx, lsh_u8 * hashval){
func Lsh512_ssse3_final(ctx *LSH512_Context, hashval Mem) {
	Comment("lsh512_ssse3_final")

	i_state_alloc := AllocLocal(XmmSize * 4 * 4)

	//__m128i cv_l[4];
	cv_l := []VecVirtual{XMM(), XMM(), XMM(), XMM()}
	//__m128i cv_r[4];
	cv_r := []VecVirtual{XMM(), XMM(), XMM(), XMM()}
	//LSH512SSSE3_Context* ctx = (LSH512SSSE3_Context*)_ctx;
	//lsh_uint remain_msg_byte;
	remain_msg_byte := GP64()
	//lsh_uint remain_msg_bit;

	//if (ctx == NULL || hashval == NULL){
	//	return LSH_ERR_NULL_PTR;
	//}
	//if (ctx->algtype == 0 || LSH_GET_HASHBYTE(ctx->algtype) > LSH512_HASH_VAL_MAX_BYTE_LEN){
	//	return LSH_ERR_INVALID_STATE;
	//}

	//remain_msg_byte = ctx->remain_databitlen >> 3;
	MOVQ(ctx.Remain_databytelen, remain_msg_byte)
	//remain_msg_bit = ctx->remain_databitlen & 7;

	//if (remain_msg_byte >= LSH512_MSG_BLK_BYTE_LEN){
	//	return LSH_ERR_INVALID_STATE;
	//}

	//if (remain_msg_bit){
	//	ctx->i_last_block[remain_msg_byte] |= (0x1 << (7 - remain_msg_bit));
	//}
	//else
	{
		//ctx->i_last_block[remain_msg_byte] = 0x80;
		MOVB(U8(0x80), ctx.I_last_block.Idx(remain_msg_byte, 1))
	}
	//memset(ctx->i_last_block + remain_msg_byte + 1, 0, LSH512_MSG_BLK_BYTE_LEN - remain_msg_byte - 1);
	size := GP64()
	MOVQ(U32(LSH512_MSG_BLK_BYTE_LEN-1), size)
	SUBQ(remain_msg_byte, size)
	Memset(ctx.I_last_block.Offset(1).Idx(remain_msg_byte, 1), 0, size, false)

	//load_blk(cv_l, ctx->cv_l);
	load_blk_mem2vec(cv_l, ctx.Cv_l)
	//load_blk(cv_r, ctx->cv_r);
	load_blk_mem2vec(cv_r, ctx.Cv_r)

	//compress(cv_l, cv_r, (lsh_u64*)ctx->i_last_block);
	compress(cv_l, cv_r, ctx.I_last_block, i_state_alloc, ctx.Cv_l, ctx.Cv_r)

	//fin(cv_l, cv_r);
	fin(cv_l, cv_r)
	//get_hash(cv_l, hashval, ctx->algtype);
	get_hash(cv_l, hashval, ctx.Algtype)

	//memset(ctx, 0, sizeof(struct LSH512_Context));

	//return LSH_SUCCESS;
}
