package lsh512sse2

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"

	. "kryptosimd/avoutil"
	. "kryptosimd/avoutil/simd"
	. "kryptosimd/lsh/sse2"
	. "kryptosimd/lsh512/avo/lsh512avoconst"
)

/* -------------------------------------------------------- */
// LSH: variables
/* -------------------------------------------------------- */

//	typedef struct LSH_ALIGNED_(32){
//		LSH_ALIGNED_(16) lsh_type algtype;
//		LSH_ALIGNED_(16) lsh_uint remain_databitlen;
//		LSH_ALIGNED_(32) __m128i cv_l[4];				// left chaining variable
//		LSH_ALIGNED_(32) __m128i cv_r[4];				// right chaining variable
//		LSH_ALIGNED_(32) lsh_u8 i_last_block[LSH512_MSG_BLK_BYTE_LEN];
//	} LSH512SSE2_Context;
type LSH512SSE2_Context struct {
	algtype           Mem // m32
	remain_databitlen Register
	cv_l              Mem
	cv_r              Mem
	i_last_block      Mem
}

//	typedef struct LSH_ALIGNED_(32) {
//		LSH_ALIGNED_(32) __m128i submsg_e_l[4];	/* even left sub-message */
//		LSH_ALIGNED_(32) __m128i submsg_e_r[4];	/* even right sub-message */
//		LSH_ALIGNED_(32) __m128i submsg_o_l[4];	/* odd left sub-message */
//		LSH_ALIGNED_(32) __m128i submsg_o_r[4];	/* odd right sub-message */
//	} LSH512SSE2_internal;
type LSH512SSE2_internal struct {
	submsg_e_l []VecVirtual
	submsg_e_r []VecVirtual
	submsg_o_l []VecVirtual
	submsg_o_r []VecVirtual
}

/* -------------------------------------------------------- */
// load a message block to register
/* -------------------------------------------------------- */

// static INLINE void load_blk(__m128i* dest, const void* src){
func load_blk_mem2vec(dst []VecVirtual, src Mem) {
	Comment("load_blk_mem2vec")

	//dest[0] = LOAD((const __m128i*)src);
	LOAD(dst[0], src)
	//dest[1] = LOAD((const __m128i*)src + 1);
	LOAD(dst[1], src.Offset(XmmSize))
	//dest[2] = LOAD((const __m128i*)src + 2);
	LOAD(dst[2], src.Offset(XmmSize*2))
	//dest[3] = LOAD((const __m128i*)src + 3);
	LOAD(dst[3], src.Offset(XmmSize*3))
}
func load_blk_vec2mem(dst Mem, src []VecVirtual) {
	Comment("load_blk_vec2mem")

	//dest[0] = LOAD((const __m128i*)src);
	LOAD(dst, src[0])
	//dest[1] = LOAD((const __m128i*)src + 1);
	LOAD(dst.Offset(XmmSize), src[1])
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
	LOAD(tmp, src.Offset(XmmSize))
	LOAD(dst.Offset(XmmSize), tmp)
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

// static INLINE void load_msg_blk(LSH512SSE2_internal * i_state, const lsh_u64* msgblk){
func load_msg_blk(i_state LSH512SSE2_internal, msgblk Mem /* uint32 */) {
	//load_blk(i_state->submsg_e_l, msgblk + 0);
	load_blk_mem2vec(i_state.submsg_e_l, msgblk.Offset(0*8))
	//load_blk(i_state->submsg_e_r, msgblk + 8);
	load_blk_mem2vec(i_state.submsg_e_r, msgblk.Offset(8*8))
	//load_blk(i_state->submsg_o_l, msgblk + 16);
	load_blk_mem2vec(i_state.submsg_o_l, msgblk.Offset(16*8))
	//load_blk(i_state->submsg_o_r, msgblk + 24);
	load_blk_mem2vec(i_state.submsg_o_r, msgblk.Offset(24*8))
}

// static INLINE void msg_exp_even(LSH512SSE2_internal * i_state){
func msg_exp_even(i_state LSH512SSE2_internal) {
	Comment("msg_exp_even")

	//__m128i temp;
	temp := XMM()
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
	MOVO_autoAU2(i_state.submsg_e_l[1], temp)
	//i_state->submsg_e_l[2] = _mm_unpacklo_epi64(i_state->submsg_e_l[3], i_state->submsg_e_l[2]);
	F_mm_unpacklo_epi64(i_state.submsg_e_l[2], i_state.submsg_e_l[3], i_state.submsg_e_l[2])
	//i_state->submsg_e_l[3] = _mm_unpackhi_epi64(temp, i_state->submsg_e_l[3]);
	F_mm_unpackhi_epi64(i_state.submsg_e_l[3], temp, i_state.submsg_e_l[3])
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

	ADD := ADD3

	//i_state->submsg_e_l[0] = ADD(i_state->submsg_o_l[0], i_state->submsg_e_l[0]);
	ADD(i_state.submsg_e_l[0], i_state.submsg_o_l[0], i_state.submsg_e_l[0])
	//i_state->submsg_e_l[1] = ADD(i_state->submsg_o_l[1], i_state->submsg_e_l[1]);
	ADD(i_state.submsg_e_l[1], i_state.submsg_o_l[1], i_state.submsg_e_l[1])
	//i_state->submsg_e_l[2] = ADD(i_state->submsg_o_l[2], i_state->submsg_e_l[2]);
	ADD(i_state.submsg_e_l[2], i_state.submsg_o_l[2], i_state.submsg_e_l[2])
	//i_state->submsg_e_l[3] = ADD(i_state->submsg_o_l[3], i_state->submsg_e_l[3]);
	ADD(i_state.submsg_e_l[3], i_state.submsg_o_l[3], i_state.submsg_e_l[3])
	//i_state->submsg_e_r[0] = ADD(i_state->submsg_o_r[0], i_state->submsg_e_r[0]);
	ADD(i_state.submsg_e_r[0], i_state.submsg_o_r[0], i_state.submsg_e_r[0])
	//i_state->submsg_e_r[1] = ADD(i_state->submsg_o_r[1], i_state->submsg_e_r[1]);
	ADD(i_state.submsg_e_r[1], i_state.submsg_o_r[1], i_state.submsg_e_r[1])
	//i_state->submsg_e_r[2] = ADD(i_state->submsg_o_r[2], i_state->submsg_e_r[2]);
	ADD(i_state.submsg_e_r[2], i_state.submsg_o_r[2], i_state.submsg_e_r[2])
	//i_state->submsg_e_r[3] = ADD(i_state->submsg_o_r[3], i_state->submsg_e_r[3]);
	ADD(i_state.submsg_e_r[3], i_state.submsg_o_r[3], i_state.submsg_e_r[3])
}

// static INLINE void msg_exp_odd(LSH512SSE2_internal * i_state){
func msg_exp_odd(i_state LSH512SSE2_internal) {
	Comment("msg_exp_odd")

	//__m128i temp;
	temp := XMM()
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

	ADD := ADD3
	//i_state->submsg_o_l[0] = ADD(i_state->submsg_e_l[0], i_state->submsg_o_l[0]);
	ADD(i_state.submsg_o_l[0], i_state.submsg_e_l[0], i_state.submsg_o_l[0])
	//i_state->submsg_o_l[1] = ADD(i_state->submsg_e_l[1], i_state->submsg_o_l[1]);
	ADD(i_state.submsg_o_l[1], i_state.submsg_e_l[1], i_state.submsg_o_l[1])
	//i_state->submsg_o_l[2] = ADD(i_state->submsg_e_l[2], i_state->submsg_o_l[2]);
	ADD(i_state.submsg_o_l[2], i_state.submsg_e_l[2], i_state.submsg_o_l[2])
	//i_state->submsg_o_l[3] = ADD(i_state->submsg_e_l[3], i_state->submsg_o_l[3]);
	ADD(i_state.submsg_o_l[3], i_state.submsg_e_l[3], i_state.submsg_o_l[3])
	//i_state->submsg_o_r[0] = ADD(i_state->submsg_e_r[0], i_state->submsg_o_r[0]);
	ADD(i_state.submsg_o_r[0], i_state.submsg_e_r[0], i_state.submsg_o_r[0])
	//i_state->submsg_o_r[1] = ADD(i_state->submsg_e_r[1], i_state->submsg_o_r[1]);
	ADD(i_state.submsg_o_r[1], i_state.submsg_e_r[1], i_state.submsg_o_r[1])
	//i_state->submsg_o_r[2] = ADD(i_state->submsg_e_r[2], i_state->submsg_o_r[2]);
	ADD(i_state.submsg_o_r[2], i_state.submsg_e_r[2], i_state.submsg_o_r[2])
	//i_state->submsg_o_r[3] = ADD(i_state->submsg_e_r[3], i_state->submsg_o_r[3]);
	ADD(i_state.submsg_o_r[3], i_state.submsg_e_r[3], i_state.submsg_o_r[3])
}

// static INLINE void load_sc(__m128i* const_v, lsh_uint i){
func load_sc(const_v []VecVirtual, i int) {
	Comment("load_sc")

	//load_blk(const_v, g_StepConstants + i);
	load_blk_mem2vec(const_v, G_StepConstants.Offset(i*8))

}

// static INLINE void msg_add_even(__m128i* cv_l, __m128i* cv_r, const LSH512SSE2_internal * i_state){
func msg_add_even(cv_l, cv_r []VecVirtual, i_state LSH512SSE2_internal) {
	Comment("msg_add_even")

	//cv_l[0] = XOR(cv_l[0], i_state->submsg_e_l[0]);
	XOR(cv_l[0], i_state.submsg_e_l[0])
	//cv_r[0] = XOR(cv_r[0], i_state->submsg_e_r[0]);
	XOR(cv_r[0], i_state.submsg_e_r[0])
	//cv_l[1] = XOR(cv_l[1], i_state->submsg_e_l[1]);
	XOR(cv_l[1], i_state.submsg_e_l[1])
	//cv_r[1] = XOR(cv_r[1], i_state->submsg_e_r[1]);
	XOR(cv_r[1], i_state.submsg_e_r[1])
	//cv_l[2] = XOR(cv_l[2], i_state->submsg_e_l[2]);
	XOR(cv_l[2], i_state.submsg_e_l[2])
	//cv_r[2] = XOR(cv_r[2], i_state->submsg_e_r[2]);
	XOR(cv_r[2], i_state.submsg_e_r[2])
	//cv_l[3] = XOR(cv_l[3], i_state->submsg_e_l[3]);
	XOR(cv_l[3], i_state.submsg_e_l[3])
	//cv_r[3] = XOR(cv_r[3], i_state->submsg_e_r[3]);
	XOR(cv_r[3], i_state.submsg_e_r[3])
}

// static INLINE void msg_add_odd(__m128i* cv_l, __m128i* cv_r, const LSH512SSE2_internal * i_state){
func msg_add_odd(cv_l, cv_r []VecVirtual, i_state LSH512SSE2_internal) {
	Comment("msg_add_odd")

	//cv_l[0] = XOR(cv_l[0], i_state->submsg_o_l[0]);
	XOR(cv_l[0], i_state.submsg_o_l[0])
	//cv_r[0] = XOR(cv_r[0], i_state->submsg_o_r[0]);
	XOR(cv_r[0], i_state.submsg_o_r[0])
	//cv_l[1] = XOR(cv_l[1], i_state->submsg_o_l[1]);
	XOR(cv_l[1], i_state.submsg_o_l[1])
	//cv_r[1] = XOR(cv_r[1], i_state->submsg_o_r[1]);
	XOR(cv_r[1], i_state.submsg_o_r[1])
	//cv_l[2] = XOR(cv_l[2], i_state->submsg_o_l[2]);
	XOR(cv_l[2], i_state.submsg_o_l[2])
	//cv_r[2] = XOR(cv_r[2], i_state->submsg_o_r[2]);
	XOR(cv_r[2], i_state.submsg_o_r[2])
	//cv_l[3] = XOR(cv_l[3], i_state->submsg_o_l[3]);
	XOR(cv_l[3], i_state.submsg_o_l[3])
	//cv_r[3] = XOR(cv_r[3], i_state->submsg_o_r[3]);
	XOR(cv_r[3], i_state.submsg_o_r[3])
}

// static INLINE void add_blk(__m128i* cv_l, const __m128i* cv_r){
func add_blk(cv_l, cv_r []VecVirtual) {
	Comment("add_blk")

	//cv_l[0] = ADD(cv_l[0], cv_r[0]);
	ADD(cv_l[0], cv_r[0])
	//cv_l[1] = ADD(cv_l[1], cv_r[1]);
	ADD(cv_l[1], cv_r[1])
	//cv_l[2] = ADD(cv_l[2], cv_r[2]);
	ADD(cv_l[2], cv_r[2])
	//cv_l[3] = ADD(cv_l[3], cv_r[3]);
	ADD(cv_l[3], cv_r[3])
}

// dst = OR(SHIFT_L(dst, ROT_EVEN_ALPHA), SHIFT_R(dst, WORD_BIT_LEN - ROT_EVEN_ALPHA))
func rotate_blk(dst VecVirtual, v int) {
	tmpXmm := XMM()

	// dst = OR(SHIFT_L(dst, ROT_EVEN_ALPHA), SHIFT_R(dst, WORD_BIT_LEN - ROT_EVEN_ALPHA))

	MOVO_autoAU2(tmpXmm, dst)
	SHIFT_L(tmpXmm, U8(v))

	SHIFT_R(dst, U8(WORD_BIT_LEN-v))
	OR(dst, tmpXmm)
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

// static INLINE void rotate_msg_gamma(__m128i* cv_r){
func rotate_msg_gamma(cv_r []VecVirtual) {
	Comment("rotate_msg_gamma")

	//__m128i temp;
	//temp = _mm_and_si128(cv_r[0], _mm_set_epi32(0xffffffff, 0xffffffff, 0x0, 0x0));\
	//cv_r[0] = _mm_and_si128(cv_r[0], _mm_set_epi32(0x0, 0x0, 0xffffffff, 0xffffffff));\
	//temp = _mm_xor_si128(_mm_slli_epi64(temp, 16), _mm_srli_epi64(temp, 48));\
	//cv_r[0] = _mm_xor_si128(cv_r[0], temp);\
	//cv_r[1] = _mm_xor_si128(_mm_slli_epi64(cv_r[1], 32), _mm_srli_epi64(cv_r[1], 32));\
	//temp = _mm_and_si128(cv_r[1], _mm_set_epi32(0xffffffff, 0xffffffff, 0x0, 0x0));\
	//cv_r[1] = _mm_and_si128(cv_r[1], _mm_set_epi32(0x0, 0x0, 0xffffffff, 0xffffffff));\
	//temp = _mm_xor_si128(_mm_slli_epi64(temp, 16), _mm_srli_epi64(temp, 48));\
	//cv_r[1] = _mm_xor_si128(cv_r[1], temp);\
	//cv_r[2] = _mm_xor_si128(_mm_slli_epi64(cv_r[2], 8), _mm_srli_epi64(cv_r[2], 56));\
	//temp = _mm_and_si128(cv_r[2], _mm_set_epi32(0xffffffff, 0xffffffff, 0x0, 0x0));\
	//cv_r[2] = _mm_and_si128(cv_r[2], _mm_set_epi32(0x0, 0x0, 0xffffffff, 0xffffffff));\
	//temp = _mm_xor_si128(_mm_slli_epi64(temp, 16), _mm_srli_epi64(temp, 48));\
	//cv_r[2] = _mm_xor_si128(cv_r[2], temp);\
	//cv_r[3] = _mm_xor_si128(_mm_slli_epi64(cv_r[3], 40), _mm_srli_epi64(cv_r[3], 24));\
	//temp= _mm_and_si128(cv_r[3], _mm_set_epi32(0xffffffff, 0xffffffff, 0x0, 0x0));\
	//cv_r[3] = _mm_and_si128(cv_r[3], _mm_set_epi32(0x0, 0x0, 0xffffffff, 0xffffffff));\
	//temp = _mm_xor_si128(_mm_slli_epi64(temp, 16), _mm_srli_epi64(temp, 48));\
	//cv_r[3] = _mm_xor_si128(cv_r[3], temp);
} //

// static INLINE void word_perm(__m128i* cv_l, __m128i* cv_r){
func word_perm(cv_l, cv_r []VecVirtual) {
	//__m128i temp[2];
	temp := []Op{XMM(), XMM()} // AllocLocal(XmmSize)

	//temp[0] = cv_l[0];
	MOVO_autoAU(temp[0], cv_l[0])
	//cv_l[0] = _mm_unpacklo_epi64(cv_l[1], cv_l[0]);
	F_mm_unpacklo_epi64(cv_l[0], cv_l[1], cv_l[0])
	//cv_l[1] = _mm_unpackhi_epi64(temp[0], cv_l[1]);
	F_mm_unpackhi_epi64(cv_l[1], temp[0], cv_l[1])
	//temp[0] = cv_l[2];
	MOVO_autoAU(temp[0], cv_l[2])
	//cv_l[2] = _mm_unpacklo_epi64(cv_l[3], cv_l[2]);
	F_mm_unpacklo_epi64(cv_l[2], cv_l[3], cv_l[2])
	//cv_l[3] = _mm_unpackhi_epi64(temp[0], cv_l[3]);
	F_mm_unpackhi_epi64(cv_l[3], temp[0], cv_l[3])
	//cv_r[1] = _mm_shuffle_epi32(cv_r[1], 0x4e);
	F_mm_shuffle_epi32(cv_r[1], cv_r[1], U8(0x4e))
	//temp[0] = cv_r[0];
	MOVO_autoAU(temp[0], cv_r[0])
	//cv_r[0] = _mm_unpacklo_epi64(cv_r[0], cv_r[1]);
	F_mm_unpacklo_epi64(cv_r[0], cv_r[0], cv_r[1])
	//cv_r[1] = _mm_unpackhi_epi64(cv_r[1], temp[0]);
	F_mm_unpackhi_epi64(cv_r[1], cv_r[1], temp[0])
	//cv_r[3] = _mm_shuffle_epi32(cv_r[3], 0x4e);
	F_mm_shuffle_epi32(cv_r[3], cv_r[3], U8(0x4e))
	//temp[0] = cv_r[2];
	MOVO_autoAU(temp[0], cv_r[2])
	//cv_r[2] = _mm_unpacklo_epi64(cv_r[2], cv_r[3]);
	F_mm_unpacklo_epi64(cv_r[2], cv_r[2], cv_r[3])
	//cv_r[3] = _mm_unpackhi_epi64(cv_r[3], temp[0]);
	F_mm_unpackhi_epi64(cv_r[3], cv_r[3], temp[0])
	//temp[0] = cv_l[0];
	MOVO_autoAU(temp[0], cv_l[0])
	//temp[1] = cv_l[1];
	MOVO_autoAU(temp[1], cv_l[1])
	//cv_l[0] = cv_l[2];
	MOVO_autoAU(cv_l[0], cv_l[2])
	//cv_l[1] = cv_l[3];
	MOVO_autoAU(cv_l[1], cv_l[3])
	//cv_l[2] = cv_r[2];
	MOVO_autoAU(cv_l[2], cv_r[2])
	//cv_l[3] = cv_r[3];
	MOVO_autoAU(cv_l[3], cv_r[3])
	//cv_r[2] = cv_r[0];
	MOVO_autoAU(cv_r[2], cv_r[0])
	//cv_r[3] = cv_r[1];
	MOVO_autoAU(cv_r[3], cv_r[1])
	//cv_r[0] = temp[0];
	MOVO_autoAU(cv_r[0], temp[0])
	//cv_r[1] = temp[1];
	MOVO_autoAU(cv_r[1], temp[1])
}

/* -------------------------------------------------------- */
// step function
/* -------------------------------------------------------- */

// static INLINE void mix_even(__m128i* cv_l, __m128i* cv_r, const __m128i* const_v){
func mix_even(cv_l, cv_r []VecVirtual, const_v []VecVirtual) {
	Comment("mix_even")

	add_blk(cv_l, cv_r)
	rotate_blk_even_alpha(cv_l)
	xor_with_const(cv_l, const_v)
	add_blk(cv_r, cv_l)
	rotate_blk_even_beta(cv_r)
	add_blk(cv_l, cv_r)
	rotate_msg_gamma(cv_r)
}

// static INLINE void mix_odd(__m128i* cv_l, __m128i* cv_r, const __m128i* const_v){
func mix_odd(cv_l, cv_r []VecVirtual, const_v []VecVirtual) {
	Comment("mix_odd")

	add_blk(cv_l, cv_r)
	rotate_blk_odd_alpha(cv_l)
	xor_with_const(cv_l, const_v)
	add_blk(cv_r, cv_l)
	rotate_blk_odd_beta(cv_r)
	add_blk(cv_l, cv_r)
	rotate_msg_gamma(cv_r)
}

/* -------------------------------------------------------- */
// compression function
/* -------------------------------------------------------- */

// static INLINE void compress(__m128i* cv_l, __m128i* cv_r, const lsh_u64 pdMsgBlk[MSG_BLK_WORD_LEN])
func compress(cv_l, cv_r []VecVirtual, pdMsgBlk Mem) {
	Comment("compress")
	//__m128i const_v[4];			// step function constant
	const_v := []VecVirtual{XMM(), XMM(), XMM(), XMM()}
	//LSH512SSE2_internal i_state[1];
	i_state := LSH512SSE2_internal{
		submsg_e_l: []VecVirtual{XMM(), XMM(), XMM(), XMM()},
		submsg_e_r: []VecVirtual{XMM(), XMM(), XMM(), XMM()},
		submsg_o_l: []VecVirtual{XMM(), XMM(), XMM(), XMM()},
		submsg_o_r: []VecVirtual{XMM(), XMM(), XMM(), XMM()},
	}
	//int i;

	load_msg_blk(i_state, pdMsgBlk)

	msg_add_even(cv_l, cv_r, i_state)
	load_sc(const_v, 0)
	mix_even(cv_l, cv_r, const_v)
	word_perm(cv_l, cv_r)

	msg_add_odd(cv_l, cv_r, i_state)
	load_sc(const_v, 8)
	mix_odd(cv_l, cv_r, const_v)
	word_perm(cv_l, cv_r)

	//for (i = 1; i < NUM_STEPS / 2; i++){
	for i := 1; i < NUM_STEPS/2; i++ {
		msg_exp_even(i_state)
		msg_add_even(cv_l, cv_r, i_state)
		load_sc(const_v, i*16)
		mix_even(cv_l, cv_r, const_v)
		word_perm(cv_l, cv_r)

		msg_exp_odd(i_state)
		msg_add_odd(cv_l, cv_r, i_state)
		load_sc(const_v, i*16+8)
		mix_odd(cv_l, cv_r, const_v)
		word_perm(cv_l, cv_r)
	}

	msg_exp_even(i_state)
	msg_add_even(cv_l, cv_r, i_state)
}

/* -------------------------------------------------------- */
//static INLINE void init224(LSH512SSE2_Context* state)
func init224(state *LSH512SSE2_Context) {
	Comment("init224")

	//load_blk(state->cv_l, g_IV224);
	load_blk_mem2mem(state.cv_l, G_IV224)
	//load_blk(state->cv_r, g_IV224 + 8);
	load_blk_mem2mem(state.cv_l, G_IV224)
}

// static INLINE void init256(LSH512SSE2_Context* state)
func init256(state *LSH512SSE2_Context) {
	Comment("init256")

	//load_blk(state->cv_l, g_IV256);
	load_blk_mem2mem(state.cv_l, G_IV256)
	//load_blk(state->cv_r, g_IV256 + 8);
	load_blk_mem2mem(state.cv_l, G_IV256)
}

// static INLINE void init384(LSH512SSE2_Context* state)
func init384(state *LSH512SSE2_Context) {
	Comment("init384")

	//load_blk(state->cv_l, g_IV384);
	load_blk_mem2mem(state.cv_l, G_IV384)
	//load_blk(state->cv_r, g_IV384 + 8);
	load_blk_mem2mem(state.cv_l, G_IV384)
}

// static INLINE void init512(LSH512SSE2_Context* state)
func init512(state *LSH512SSE2_Context) {
	Comment("init512")

	//load_blk(state->cv_l, g_IV512);
	load_blk_mem2mem(state.cv_l, G_IV512)
	//load_blk(state->cv_r, g_IV512 + 8);
	load_blk_mem2mem(state.cv_l, G_IV512)
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
	hash_val_byte_len := GP32()
	LSH_GET_HASHBYTE(hash_val_byte_len, algtype)
	//lsh_uint hash_val_bit_len = LSH_GET_SMALL_HASHBIT(algtype);
	hash_val_bit_len := GP32()
	LSH_GET_SMALL_HASHBIT(hash_val_bit_len, algtype)

	//store_blk((__m128i*)hash_val, cv_l);
	store_blk(hash_val, cv_l)
	//memcpy(pbHashVal, hash_val, sizeof(lsh_u8) * hash_val_byte_len);
	//if (hash_val_bit_len){
	CMPL(hash_val_bit_len, U32(0))
	JE(LabelRef("get_hash_if_end"))
	{
		//pbHashVal[hash_val_byte_len-1] &= (((lsh_u8)0xff) << hash_val_bit_len);
		tmp8 := GP8()
		MOVB(U8(0xff), tmp8)
		MOVB(hash_val_bit_len.As8(), CL)
		SHLB(CL, tmp8)

		addr := GP32()
		MOVL(hash_val_byte_len, addr.As32())
		SUBL(U8(1), addr)

		MOVB(tmp8, pbHashVal.Idx(addr, 1))
	}
	Label("get_hash_if_end")
}

/* -------------------------------------------------------- */

// lsh_err lsh512_sse2_init(struct LSH512_Context * _ctx, const lsh_type algtype){
func lsh512_sse2_init(ctx *LSH512SSE2_Context) {
	Comment("lsh512_sse2_init")

	//LSH512SSE2_Context* ctx = (LSH512SSE2_Context*)_ctx;
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
	CMPL(ctx.algtype, U32(LSH_TYPE_512_512))
	JNE(LabelRef("lsh512_sse2_init_if0_end"))
	{
		//	init512(ctx);
		init256(ctx)
		//	return LSH_SUCCESS;
		JMP(LabelRef("lsh512_sse2_init_ret"))
	}
	Label("lsh512_sse2_init_if0_end")
	//case LSH_TYPE_512_384:
	CMPL(ctx.algtype, U32(LSH_TYPE_512_384))
	JNE(LabelRef("lsh512_sse2_init_if1_end"))
	{
		//	init384(ctx);
		init384(ctx)
		//	return LSH_SUCCESS;
		JMP(LabelRef("lsh512_sse2_init_ret"))
	}
	Label("lsh512_sse2_init_if1_end")
	//case LSH_TYPE_512_256:
	CMPL(ctx.algtype, U32(LSH_TYPE_512_256))
	JNE(LabelRef("lsh512_sse2_init_if2_end"))
	{
		//	init256(ctx);
		init256(ctx)
		//	return LSH_SUCCESS;
		JMP(LabelRef("lsh512_sse2_init_ret"))
	}
	Label("lsh512_sse2_init_if2_end")
	//case LSH_TYPE_512_224:
	{
		//	init224(ctx);
		init224(ctx)
		//	return LSH_SUCCESS;
		JMP(LabelRef("lsh512_sse2_init_ret"))
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

	//for (i = 0; i < NUM_STEPS / 2; i++)
	//{
	//	//Mix
	//	load_sc(const_v, i * 16);
	//	mix_even(cv_l, cv_r, const_v);
	//	word_perm(cv_l, cv_r);

	//	load_sc(const_v, i * 16 + 8);
	//	mix_odd(cv_l, cv_r, const_v);
	//	word_perm(cv_l, cv_r);
	//}

	//store_blk(ctx->cv_l, cv_l);
	//store_blk(ctx->cv_r, cv_r);

	//return LSH_SUCCESS;
	Label("lsh512_sse2_init_ret")
}

// lsh_err lsh512_sse2_update(struct LSH512_Context * _ctx, const lsh_u8 * data, size_t databitlen){
func lsh512_sse2_update(ctx *LSH512SSE2_Context, data Mem, databitlen Register) {
	Comment("lsh512_sse2_update")

	//__m128i cv_l[4];
	cv_l := []VecVirtual{XMM(), XMM(), XMM(), XMM()}
	//__m128i cv_r[4];
	cv_r := []VecVirtual{XMM(), XMM(), XMM(), XMM()}
	//size_t databytelen = databitlen >> 3;
	databytelen := GP32()
	MOVL(databitlen, databytelen)
	SHRL(U8(3), databytelen)
	//lsh_u32 pos2 = databitlen & 0x7;

	//LSH512SSE2_Context* ctx = (LSH512SSE2_Context*)_ctx;
	//lsh_uint remain_msg_byte;
	remain_msg_byte := GP32()
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
	MOVL(ctx.remain_databitlen, remain_msg_byte)
	SHRL(U8(3), remain_msg_byte)
	//remain_msg_bit = ctx->remain_databitlen & 7;
	//if (remain_msg_byte >= LSH512_MSG_BLK_BYTE_LEN){
	//	return LSH_ERR_INVALID_STATE;
	//}
	//if (remain_msg_bit > 0){
	//	return LSH_ERR_INVALID_DATABITLEN;
	//}

	//if (databytelen + remain_msg_byte < LSH512_MSG_BLK_BYTE_LEN){
	tmp32 := GP32()
	MOVL(databytelen, tmp32)
	ADDL(remain_msg_byte, tmp32)
	CMPL(tmp32, U32(LSH512_MSG_BLK_BYTE_LEN))
	JGE(LabelRef("lsh512_sse2_update_if0_end"))
	{
		//memcpy(ctx->i_last_block + remain_msg_byte, data, databytelen);
		Memcpy(ctx.i_last_block.Idx(remain_msg_byte, 1), data, databytelen, false)
		//ctx->remain_databitlen += (lsh_uint)databitlen;
		ADDL(databitlen, ctx.remain_databitlen)
		//remain_msg_byte += (lsh_uint)databytelen;
		ADDL(databytelen, remain_msg_byte)
		//if (pos2){
		//	ctx->i_last_block[remain_msg_byte] = data[databytelen] & ((0xff >> pos2) ^ 0xff);
		//}
		//return LSH_SUCCESS;
		JMP(LabelRef("lsh512_sse2_update_ret"))
	}
	Label("lsh512_sse2_update_if0_end")

	//load_blk(cv_l, ctx->cv_l);
	load_blk_mem2vec(cv_l, ctx.cv_l)
	//load_blk(cv_r, ctx->cv_r);
	load_blk_mem2vec(cv_r, ctx.cv_r)

	//if (remain_msg_byte > 0){
	CMPL(remain_msg_byte, U32(0))
	JE(LabelRef("lsh512_sse2_update_if1_end"))
	{
		//size_t more_BYTE = LSH512_MSG_BLK_BYTE_LEN - remain_msg_byte;
		more_BYTE := GP32()
		MOVL(U32(LSH512_MSG_BLK_BYTE_LEN), more_BYTE)
		SUBL(remain_msg_byte, more_BYTE)
		//memcpy(ctx->i_last_block + remain_msg_byte, data, more_BYTE);
		Memcpy(ctx.i_last_block.Idx(remain_msg_byte, 1), data, more_BYTE, false)
		//compress(cv_l, cv_r, (lsh_u64*)ctx->i_last_block);
		compress(cv_l, cv_r, ctx.i_last_block)
		//data += more_BYTE;
		ADDQ(more_BYTE.As64(), data.Base)
		//databytelen -= more_BYTE;
		SUBL(more_BYTE, databytelen)
		//remain_msg_byte = 0;
		MOVL(U32(0), remain_msg_byte)
		//ctx->remain_databitlen = 0;
		MOVL(U32(0), ctx.remain_databitlen)
	}
	Label("lsh512_sse2_update_if1_end")

	//while (databytelen >= LSH512_MSG_BLK_BYTE_LEN)
	Label("lsh512_sse2_update_while_start")
	CMPL(databytelen, U32(LSH512_MSG_BLK_BYTE_LEN))
	JL(LabelRef("lsh512_sse2_update_while_end"))
	{
		//compress(cv_l, cv_r, (lsh_u64*)data);
		compress(cv_l, cv_r, data)
		//data += LSH512_MSG_BLK_BYTE_LEN;
		ADDQ(U32(LSH512_MSG_BLK_BYTE_LEN), data.Base)
		//databytelen -= LSH512_MSG_BLK_BYTE_LEN;
		SUBL(U32(LSH512_MSG_BLK_BYTE_LEN), databytelen)

		JMP(LabelRef("lsh512_sse2_update_while_start"))
	}
	Label("lsh512_sse2_update_while_end")

	//store_blk(ctx->cv_l, cv_l);
	store_blk(ctx.cv_l, cv_l)
	//store_blk(ctx->cv_r, cv_r);
	store_blk(ctx.cv_r, cv_r)

	//if (databytelen > 0){
	CMPL(remain_msg_byte, U32(0))
	JE(LabelRef("lsh512_sse2_update_if3_end"))
	{
		//memcpy(ctx->i_last_block, data, databytelen);
		Memcpy(ctx.i_last_block, data, databytelen, false)
		//ctx->remain_databitlen = (lsh_uint)(databytelen << 3);
		MOVL(databytelen, ctx.remain_databitlen)
		SHLL(U8(3), ctx.remain_databitlen)
	}
	Label("lsh512_sse2_update_if3_end")

	//if (pos2){
	//	ctx->i_last_block[databytelen] = data[databytelen] & ((0xff >> pos2) ^ 0xff);
	//	ctx->remain_databitlen += pos2;
	//}
	//return LSH_SUCCESS;

	Label("lsh512_sse2_update_ret")
}

// lsh_err lsh512_sse2_final(struct LSH512_Context * _ctx, lsh_u8 * hashval){
func lsh512_sse2_final(ctx *LSH512SSE2_Context, hashval Mem) {
	Comment("lsh512_sse2_final")

	tmp32 := GP32()

	//__m128i cv_l[4];
	cv_l := []VecVirtual{XMM(), XMM(), XMM(), XMM()}
	//__m128i cv_r[4];
	cv_r := []VecVirtual{XMM(), XMM(), XMM(), XMM()}
	//LSH512SSE2_Context* ctx = (LSH512SSE2_Context*)_ctx;
	//lsh_uint remain_msg_byte;
	remain_msg_byte := GP32()
	//lsh_uint remain_msg_bit;

	//if (ctx == NULL || hashval == NULL){
	//	return LSH_ERR_NULL_PTR;
	//}
	//if (ctx->algtype == 0 || LSH_GET_HASHBYTE(ctx->algtype) > LSH512_HASH_VAL_MAX_BYTE_LEN){
	//	return LSH_ERR_INVALID_STATE;
	//}

	//remain_msg_byte = ctx->remain_databitlen >> 3;
	MOVL(ctx.remain_databitlen, remain_msg_byte)
	SHRL(U8(3), remain_msg_byte)
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
		MOVB(U8(0x80), ctx.i_last_block.Idx(remain_msg_byte, 1))
	}
	//memset(ctx->i_last_block + remain_msg_byte + 1, 0, LSH512_MSG_BLK_BYTE_LEN - remain_msg_byte - 1);
	MOVL(remain_msg_byte, tmp32)
	ADDL(U8(1), tmp32)
	arg2 := GP32()
	MOVL(U32(LSH512_MSG_BLK_BYTE_LEN-1), arg2)
	SUBL(remain_msg_byte, arg2)
	Memset(ctx.i_last_block.Idx(tmp32, 1), 0, arg2, false)

	//load_blk(cv_l, ctx->cv_l);
	load_blk_mem2vec(cv_l, ctx.cv_l)
	//load_blk(cv_r, ctx->cv_r);
	load_blk_mem2vec(cv_r, ctx.cv_r)

	//compress(cv_l, cv_r, (lsh_u64*)ctx->i_last_block);
	compress(cv_l, cv_r, ctx.i_last_block)

	//fin(cv_l, cv_r);
	fin(cv_l, cv_r)
	//get_hash(cv_l, hashval, ctx->algtype);
	get_hash(cv_l, hashval, ctx.algtype)

	//memset(ctx, 0, sizeof(struct LSH512_Context));

	//return LSH_SUCCESS;
}

func getCtx() *LSH512SSE2_Context {
	ctx := Dereference(Param("ctx"))

	algtype, err := ctx.Field("algtype").Resolve()
	if err != nil {
		panic(err)
	}
	cv_l, err := ctx.Field("cv_l").Index(0).Resolve()
	if err != nil {
		panic(err)
	}
	cv_r, err := ctx.Field("cv_r").Index(0).Resolve()
	if err != nil {
		panic(err)
	}
	i_last_block, err := ctx.Field("i_last_block").Index(0).Resolve()
	if err != nil {
		panic(err)
	}

	return &LSH512SSE2_Context{
		algtype:           algtype.Addr,
		remain_databitlen: Load(ctx.Field("remain_databitlen"), GP32()),
		cv_l:              cv_l.Addr,
		cv_r:              cv_r.Addr,
		i_last_block:      i_last_block.Addr,
	}
}

func LSH512InitSSE2() {
	TEXT("lsh512InitSSE2", NOSPLIT, "func(ctx *lsh512ContextAsmData)")

	ctx := getCtx()
	lsh512_sse2_init(ctx)
	Store(ctx.remain_databitlen, Dereference(Param("ctx")).Field("remain_databitlen"))

	RET()
}

func LSH512UpdateSSE2() {
	TEXT("lsh512UpdateSSE2", NOSPLIT, "func(ctx *lsh512ContextAsmData, data []byte, databitlen uint32)")

	ctx := getCtx()
	data := Mem{Base: Load(Param("data").Base(), GP64())}
	databitlen := Load(Param("databitlen"), GP32())

	lsh512_sse2_update(ctx, data, databitlen)
	Store(ctx.remain_databitlen, Dereference(Param("ctx")).Field("remain_databitlen"))

	RET()
}

func LSH512FinalSSE2() {
	TEXT("lsh512FinalSSE2", NOSPLIT, "func(ctx *lsh512ContextAsmData, hashval []byte)")

	ctx := getCtx()
	hashval := Mem{Base: Load(Param("hashval").Base(), GP64())}

	lsh512_sse2_final(ctx, hashval)
	Store(ctx.remain_databitlen, Dereference(Param("ctx")).Field("remain_databitlen"))

	RET()
}
