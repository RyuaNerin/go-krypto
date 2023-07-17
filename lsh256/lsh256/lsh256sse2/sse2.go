package lsh256sse2

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"

	"github.com/RyuaNerin/go-krypto/lsh256/lsh256/lsh256avoconst"
)

const (
	LSH_TYPE_256_256 = 0x0000020
	LSH_TYPE_256_224 = 0x000001C

	MSG_BLK_WORD_LEN      = 32
	CV_WORD_LEN           = 16
	CONST_WORD_LEN        = 8
	HASH_VAL_MAX_WORD_LEN = 8

	WORD_BIT_LEN = 32

	LSH256_MSG_BLK_BYTE_LEN      = 128
	LSH256_MSG_BLK_BIT_LEN       = 1024
	LSH256_CV_BYTE_LEN           = 64
	LSH256_HASH_VAL_MAX_BYTE_LEN = 32

	/* -------------------------------------------------------- */

	NUM_STEPS = 26

	ROT_EVEN_ALPHA = 29
	ROT_EVEN_BETA  = 1
	ROT_ODD_ALPHA  = 5
	ROT_ODD_BETA   = 17
)

func LSH_GET_SMALL_HASHBIT(lsh_type_val int) int { return ((lsh_type_val) >> 24) }
func LSH_GET_HASHBYTE(lsh_type_val int) int      { return ((lsh_type_val) & 0xffff) }
func LSH_GET_HASHBIT(lsh_type_val int) int {
	return ((LSH_GET_HASHBYTE(lsh_type_val) << 3) - LSH_GET_SMALL_HASHBIT(lsh_type_val))
}

/**
#define LOAD(x)      _mm_loadu_si128((__m128i*)x)
#define STORE(x,y)   _mm_storeu_si128((__m128i*)x, y)
#define XOR(x,y)     _mm_xor_si128(x,y)
#define OR(x,y)      _mm_or_si128(x,y)
#define AND(x,y)     _mm_and_si128(x,y)

#define ADD(x,y)     _mm_add_epi32(x,y)
#define SHIFT_L(x,r) _mm_slli_epi32(x,r)
#define SHIFT_R(x,r) _mm_srli_epi32(x,r)
*/

type LSH256SSE2_internal struct {
	submsg_e_l []VecVirtual
	submsg_e_r []VecVirtual
	submsg_o_l []VecVirtual
	submsg_o_r []VecVirtual
}

func load_blk(dest []VecVirtual, src Mem) {
	Comment("load_blk")

	/**
	dest[0] = LOAD((const __m128i*)src);
	dest[1] = LOAD((const __m128i*)src + 1);
	*/
	MOVOU(src.Offset(0x00), dest[0]) // movdqu
	MOVOU(src.Offset(0x10), dest[1]) // movdqu
}

func store_blk(dest Mem, src []VecVirtual) {
	Comment("store_blk")

	/**
	STORE(dest, src[0]);
	STORE(dest + 1, src[1]);
	*/
	MOVOU(src[0], dest.Offset(0x00)) // movdqu
	MOVOU(src[1], dest.Offset(0x10)) // movdqu
}

func load_msg_blk(i_state LSH256SSE2_internal, msgblk Mem /* uint32 */) {
	Comment("load_msg_blk")

	/**
	load_blk(i_state->submsg_e_l, msgblk + 0);
	load_blk(i_state->submsg_e_r, msgblk + 8);
	load_blk(i_state->submsg_o_l, msgblk + 16);
	load_blk(i_state->submsg_o_r, msgblk + 24);
	*/
	load_blk(i_state.submsg_e_l, msgblk.Offset(0x00*4))
	load_blk(i_state.submsg_e_r, msgblk.Offset(0x08*4))
	load_blk(i_state.submsg_o_l, msgblk.Offset(0x10*4))
	load_blk(i_state.submsg_o_r, msgblk.Offset(0x18*4))
}

func msg_exp_even(i_state LSH256SSE2_internal) {
	Comment("msg_exp_even")

	/**
	i_state->submsg_e_l[0] = ADD(i_state->submsg_o_l[0], _mm_shuffle_epi32(i_state->submsg_e_l[0], 0x4b));
	i_state->submsg_e_l[1] = ADD(i_state->submsg_o_l[1], _mm_shuffle_epi32(i_state->submsg_e_l[1], 0x93));
	i_state->submsg_e_r[0] = ADD(i_state->submsg_o_r[0], _mm_shuffle_epi32(i_state->submsg_e_r[0], 0x4b));
	i_state->submsg_e_r[1] = ADD(i_state->submsg_o_r[1], _mm_shuffle_epi32(i_state->submsg_e_r[1], 0x93));
	*/
	PSHUFD(U8(0x4b), i_state.submsg_e_l[0], i_state.submsg_e_l[0])
	PAND(i_state.submsg_o_l[0], i_state.submsg_e_l[0])

	PSHUFD(U8(0x93), i_state.submsg_e_l[1], i_state.submsg_e_l[1])
	PAND(i_state.submsg_o_l[1], i_state.submsg_e_l[1])

	PSHUFD(U8(0x4b), i_state.submsg_e_r[0], i_state.submsg_e_r[0])
	PAND(i_state.submsg_o_r[0], i_state.submsg_e_r[0])

	PSHUFD(U8(0x93), i_state.submsg_e_r[1], i_state.submsg_e_r[1])
	PAND(i_state.submsg_o_r[1], i_state.submsg_e_r[1])
}

func msg_exp_odd(i_state LSH256SSE2_internal) {
	Comment("msg_exp_odd")

	/**
	i_state->submsg_o_l[0] = ADD(i_state->submsg_e_l[0], _mm_shuffle_epi32(i_state->submsg_o_l[0], 0x4b));
	i_state->submsg_o_l[1] = ADD(i_state->submsg_e_l[1], _mm_shuffle_epi32(i_state->submsg_o_l[1], 0x93));
	i_state->submsg_o_r[0] = ADD(i_state->submsg_e_r[0], _mm_shuffle_epi32(i_state->submsg_o_r[0], 0x4b));
	i_state->submsg_o_r[1] = ADD(i_state->submsg_e_r[1], _mm_shuffle_epi32(i_state->submsg_o_r[1], 0x93));
	*/

	PSHUFD(U8(0x4b), i_state.submsg_o_l[0], i_state.submsg_o_l[0])
	PAND(i_state.submsg_e_l[0], i_state.submsg_o_l[0])

	PSHUFD(U8(0x93), i_state.submsg_o_l[1], i_state.submsg_o_l[1])
	PAND(i_state.submsg_e_l[1], i_state.submsg_o_l[1])

	PSHUFD(U8(0x4b), i_state.submsg_o_r[0], i_state.submsg_o_r[0])
	PAND(i_state.submsg_e_r[0], i_state.submsg_o_r[0])

	PSHUFD(U8(0x93), i_state.submsg_o_r[1], i_state.submsg_o_r[1])
	PAND(i_state.submsg_e_r[1], i_state.submsg_o_r[1])
}

func load_sc(const_v []VecVirtual, i int) {
	Comment("load_sc")

	/**
	load_blk(const_v, g_StepConstants + i);
	*/
	load_blk(const_v, lsh256avoconst.G_StepConstants.Offset(4*i))
}

func load_sc_const_v(const_v []Mem, i int) {
	const_v[0] = lsh256avoconst.G_StepConstants.Offset(4*i + 0x00)
	const_v[1] = lsh256avoconst.G_StepConstants.Offset(4*i + 0x10)
}

func msg_add_even(cv_l []VecVirtual, cv_r []VecVirtual, i_state LSH256SSE2_internal) {
	Comment("msg_add_even")

	/**
	cv_l[0] = XOR(cv_l[0], i_state->submsg_e_l[0]);
	cv_r[0] = XOR(cv_r[0], i_state->submsg_e_r[0]);
	cv_l[1] = XOR(cv_l[1], i_state->submsg_e_l[1]);
	cv_r[1] = XOR(cv_r[1], i_state->submsg_e_r[1]);
	*/
	PXOR(i_state.submsg_e_l[0], cv_l[0])
	PXOR(i_state.submsg_e_r[0], cv_r[0])
	PXOR(i_state.submsg_e_l[1], cv_l[1])
	PXOR(i_state.submsg_e_r[1], cv_r[1])
}

func msg_add_odd(cv_l []VecVirtual, cv_r []VecVirtual, i_state LSH256SSE2_internal) {
	Comment("msg_add_odd")

	/**
	cv_l[0] = XOR(cv_l[0], i_state->submsg_o_l[0]);
	cv_r[0] = XOR(cv_r[0], i_state->submsg_o_r[0]);
	cv_l[1] = XOR(cv_l[1], i_state->submsg_o_l[1]);
	cv_r[1] = XOR(cv_r[1], i_state->submsg_o_r[1]);
	*/
	PXOR(i_state.submsg_o_l[0], cv_l[0])
	PXOR(i_state.submsg_o_r[0], cv_r[0])
	PXOR(i_state.submsg_o_l[1], cv_l[1])
	PXOR(i_state.submsg_o_r[1], cv_r[1])
}

func add_blk(cv_l []VecVirtual, cv_r []VecVirtual) {
	Comment("add_blk")

	/**
	cv_l[0] = ADD(cv_l[0], cv_r[0]);
	cv_l[1] = ADD(cv_l[1], cv_r[1]);
	*/
	PADDD(cv_r[0], cv_l[0])
	PADDD(cv_r[1], cv_l[1])
}

func rotate_blk_even_alpha(cv []VecVirtual) {
	Comment("rotate_blk_even_alpha")

	/**
	cv[0] = OR(SHIFT_L(cv[0], ROT_EVEN_ALPHA), SHIFT_R(cv[0], WORD_BIT_LEN - ROT_EVEN_ALPHA));
	cv[1] = OR(SHIFT_L(cv[1], ROT_EVEN_ALPHA), SHIFT_R(cv[1], WORD_BIT_LEN - ROT_EVEN_ALPHA));
	*/

	tmp := XMM()

	MOVO(cv[0], tmp)                              // movdqa
	PSLLL(U8(ROT_EVEN_ALPHA), tmp)                // pslld
	PSRLL(U8(WORD_BIT_LEN-ROT_EVEN_ALPHA), cv[0]) // psrld
	POR(tmp, cv[0])

	MOVO(cv[1], tmp)                              // movdqa
	PSLLL(U8(ROT_EVEN_ALPHA), tmp)                // pslld
	PSRLL(U8(WORD_BIT_LEN-ROT_EVEN_ALPHA), cv[1]) // psrld
	POR(tmp, cv[1])
}

func rotate_blk_even_beta(cv []VecVirtual) {
	Comment("rotate_blk_even_beta")

	/**
	cv[0] = OR(SHIFT_L(cv[0], ROT_EVEN_BETA), SHIFT_R(cv[0], WORD_BIT_LEN - ROT_EVEN_BETA));
	cv[1] = OR(SHIFT_L(cv[1], ROT_EVEN_BETA), SHIFT_R(cv[1], WORD_BIT_LEN - ROT_EVEN_BETA));
	*/

	tmp := XMM()

	MOVO(cv[0], tmp)                             // movdqa
	PSLLL(U8(ROT_EVEN_BETA), tmp)                // pslld
	PSRLL(U8(WORD_BIT_LEN-ROT_EVEN_BETA), cv[0]) // psrld
	POR(tmp, cv[0])

	MOVO(cv[1], tmp)                             // movdqa
	PSLLL(U8(ROT_EVEN_BETA), tmp)                // pslld
	PSRLL(U8(WORD_BIT_LEN-ROT_EVEN_BETA), cv[1]) // psrld
	POR(tmp, cv[1])
}

func rotate_blk_odd_alpha(cv []VecVirtual) {
	Comment("rotate_blk_odd_alpha")

	/**
	cv[0] = OR(SHIFT_L(cv[0], ROT_ODD_ALPHA), SHIFT_R(cv[0], WORD_BIT_LEN - ROT_ODD_ALPHA));
	cv[1] = OR(SHIFT_L(cv[1], ROT_ODD_ALPHA), SHIFT_R(cv[1], WORD_BIT_LEN - ROT_ODD_ALPHA));
	*/

	tmp := XMM()

	MOVO(cv[0], tmp)                             // movdqa
	PSLLL(U8(ROT_ODD_ALPHA), tmp)                // pslld
	PSRLL(U8(WORD_BIT_LEN-ROT_ODD_ALPHA), cv[0]) // psrld
	POR(tmp, cv[0])

	MOVO(cv[1], tmp)                             // movdqa
	PSLLL(U8(ROT_ODD_ALPHA), tmp)                // pslld
	PSRLL(U8(WORD_BIT_LEN-ROT_ODD_ALPHA), cv[1]) // psrld
	POR(tmp, cv[1])
}
func rotate_blk_odd_beta(cv []VecVirtual) {
	Comment("rotate_blk_odd_beta")

	/**
	cv[0] = OR(SHIFT_L(cv[0], ROT_ODD_BETA), SHIFT_R(cv[0], WORD_BIT_LEN - ROT_ODD_BETA));
	cv[1] = OR(SHIFT_L(cv[1], ROT_ODD_BETA), SHIFT_R(cv[1], WORD_BIT_LEN - ROT_ODD_BETA));
	*/

	tmp := XMM()

	MOVO(cv[0], tmp)                            // movdqa
	PSLLL(U8(ROT_ODD_BETA), tmp)                // pslld
	PSRLL(U8(WORD_BIT_LEN-ROT_ODD_BETA), cv[0]) // psrld
	POR(tmp, cv[0])

	MOVO(cv[1], tmp)                            // movdqa
	PSLLL(U8(ROT_ODD_BETA), tmp)                // pslld
	PSRLL(U8(WORD_BIT_LEN-ROT_ODD_BETA), cv[1]) // psrld
	POR(tmp, cv[1])
}

func xor_with_const(cv_l []VecVirtual, const_v []Mem) {
	Comment("xor_with_const")

	/**
	cv_l[0] = XOR(cv_l[0], const_v[0])
	cv_l[1] = XOR(cv_l[1], const_v[1])
	*/
	PXOR(const_v[0], cv_l[0])
	PXOR(const_v[1], cv_l[1])
}

var (
	g_BytePermInfo_L [6]Mem
	g_BytePermInfo_R [6]Mem
)

func init() {
	f := func(name string, values ...uint64) Mem {
		m := GLOBL(name, NOPTR|RODATA)
		for i, v := range values {
			DATA(4*i, U32(v))
		}
		return m
	}

	g_BytePermInfo_L[0] = f("g_BytePermInfo_L_sse2_0", 0xffffffff, 0xffffffff, 0xffffffff, 0x00000000)
	g_BytePermInfo_L[1] = f("g_BytePermInfo_L_sse2_1", 0x00000000, 0x00000000, 0x00000000, 0xffffffff)
	g_BytePermInfo_L[2] = f("g_BytePermInfo_L_sse2_2", 0xffffffff, 0xffffffff, 0x00000000, 0x00000000)
	g_BytePermInfo_L[3] = f("g_BytePermInfo_L_sse2_3", 0x00000000, 0x00000000, 0xffffffff, 0xffffffff)
	g_BytePermInfo_L[4] = f("g_BytePermInfo_L_sse2_4", 0xffffffff, 0x00000000, 0x00000000, 0x00000000)
	g_BytePermInfo_L[5] = f("g_BytePermInfo_L_sse2_5", 0x00000000, 0xffffffff, 0xffffffff, 0xffffffff)

	g_BytePermInfo_R[0] = f("g_BytePermInfo_R_sse2_0", 0x00000000, 0xffffffff, 0xffffffff, 0xffffffff)
	g_BytePermInfo_R[1] = f("g_BytePermInfo_R_sse2_1", 0xffffffff, 0x00000000, 0x00000000, 0x00000000)
	g_BytePermInfo_R[2] = f("g_BytePermInfo_R_sse2_2", 0x00000000, 0x00000000, 0xffffffff, 0xffffffff)
	g_BytePermInfo_R[3] = f("g_BytePermInfo_R_sse2_3", 0xffffffff, 0xffffffff, 0x00000000, 0x00000000)
	g_BytePermInfo_R[4] = f("g_BytePermInfo_R_sse2_4", 0x00000000, 0x00000000, 0x00000000, 0xffffffff)
	g_BytePermInfo_R[5] = f("g_BytePermInfo_R_sse2_5", 0xffffffff, 0xffffffff, 0xffffffff, 0x00000000)
}

func rotate_msg_gamma(cv_r []VecVirtual) {
	Comment("rotate_msg_gamma")

	/**
	__m128i temp;

	temp = AND(cv_r[0], _mm_set_epi32(0xffffffff, 0xffffffff, 0xffffffff, 0x0));
	cv_r[0] = AND(cv_r[0], _mm_set_epi32(0x0, 0x0, 0x0, 0xffffffff));
	temp = XOR(SHIFT_L(temp, 8), SHIFT_R(temp, 24));
	cv_r[0] = XOR(cv_r[0], temp);
	temp = AND(cv_r[0], _mm_set_epi32(0xffffffff, 0xffffffff, 0x0, 0x0));
	cv_r[0] = AND(cv_r[0], _mm_set_epi32(0x0, 0x0, 0xffffffff, 0xffffffff));
	temp = XOR(SHIFT_L(temp, 8), SHIFT_R(temp, 24));
	cv_r[0] = XOR(cv_r[0], temp);
	temp = AND(cv_r[0], _mm_set_epi32(0xffffffff, 0x0, 0x0, 0x0));
	cv_r[0] = AND(cv_r[0], _mm_set_epi32(0x0, 0xffffffff, 0xffffffff, 0xffffffff));
	temp = XOR(SHIFT_L(temp, 8), SHIFT_R(temp, 24));
	cv_r[0] = XOR(cv_r[0], temp);

	temp = AND(cv_r[1], _mm_set_epi32(0x0, 0xffffffff, 0xffffffff, 0xffffffff));
	cv_r[1] = AND(cv_r[1], _mm_set_epi32(0xffffffff, 0x0, 0x0, 0x0));
	temp = XOR(SHIFT_L(temp, 8), SHIFT_R(temp, 24));
	cv_r[1] = XOR(cv_r[1], temp);
	temp = AND(cv_r[1], _mm_set_epi32(0x0, 0x0, 0xffffffff, 0xffffffff));
	cv_r[1] = AND(cv_r[1], _mm_set_epi32(0xffffffff, 0xffffffff, 0x0, 0x0));
	temp = XOR(SHIFT_L(temp, 8), SHIFT_R(temp, 24));
	cv_r[1] = XOR(cv_r[1], temp);
	temp = AND(cv_r[1], _mm_set_epi32(0x0, 0x0, 0x0, 0xffffffff));
	cv_r[1] = AND(cv_r[1], _mm_set_epi32(0xffffffff, 0xffffffff, 0xffffffff, 0x0));
	temp = XOR(SHIFT_L(temp, 8), SHIFT_R(temp, 24));
	cv_r[1] = XOR(cv_r[1], temp);
	*/

	{
		// gcc -O3 -msse2
		xmm0 := cv_r[0]
		xmm1 := XMM()
		xmm2 := XMM()

		MOVO(g_BytePermInfo_L[0], xmm1)
		PAND(xmm0, xmm1)
		PAND(g_BytePermInfo_L[1], xmm0)
		MOVO(xmm1, xmm2)
		PSLLL(U8(8), xmm1)
		PSRLL(U8(24), xmm2)
		PXOR(xmm2, xmm1)
		PXOR(xmm0, xmm1)
		MOVO(g_BytePermInfo_L[2], xmm0)
		PAND(xmm1, xmm0)
		PAND(g_BytePermInfo_L[3], xmm1)
		MOVO(xmm0, xmm2)
		PSLLL(U8(8), xmm0)
		PSRLL(U8(24), xmm2)
		PXOR(xmm2, xmm0)
		PXOR(xmm1, xmm0)
		MOVO(g_BytePermInfo_L[4], xmm1)
		PAND(xmm0, xmm1)
		PAND(g_BytePermInfo_L[5], xmm0)
		MOVO(xmm1, xmm2)
		PSLLL(U8(8), xmm1)
		PSRLL(U8(24), xmm2)
		PXOR(xmm2, xmm1)
		PXOR(xmm1, xmm0)
	}
	{
		// gcc -O3 -msse2
		xmm0 := cv_r[1]
		xmm1 := XMM()
		xmm2 := XMM()

		MOVO(g_BytePermInfo_R[0], xmm1)
		PAND(xmm0, xmm1)
		PAND(g_BytePermInfo_R[1], xmm0)
		MOVO(xmm1, xmm2)
		PSLLL(U8(8), xmm1)
		PSRLL(U8(24), xmm2)
		PXOR(xmm2, xmm1)
		PXOR(xmm0, xmm1)
		MOVO(g_BytePermInfo_R[2], xmm0)
		PAND(xmm1, xmm0)
		PAND(g_BytePermInfo_R[3], xmm1)
		MOVO(xmm0, xmm2)
		PSLLL(U8(8), xmm0)
		PSRLL(U8(24), xmm2)
		PXOR(xmm2, xmm0)
		PXOR(xmm1, xmm0)
		MOVO(g_BytePermInfo_R[4], xmm1)
		PAND(xmm0, xmm1)
		PAND(g_BytePermInfo_R[5], xmm0)
		MOVO(xmm1, xmm2)
		PSLLL(U8(8), xmm1)
		PSRLL(U8(24), xmm2)
		PXOR(xmm2, xmm1)
		PXOR(xmm1, xmm0)
	}
}

func word_perm(cv_l, cv_r []VecVirtual) {
	Comment("word_perm")

	/**
	__m128i temp;
	cv_l[0] = _mm_shuffle_epi32(cv_l[0], 0xd2);
	cv_l[1] = _mm_shuffle_epi32(cv_l[1], 0xd2);
	cv_r[0] = _mm_shuffle_epi32(cv_r[0], 0x6c);
	cv_r[1] = _mm_shuffle_epi32(cv_r[1], 0x6c);
	temp = cv_l[0];
	cv_l[0] = cv_l[1];
	cv_l[1] = cv_r[1];
	cv_r[1] = cv_r[0];
	cv_r[0] = temp;
	*/

	PSHUFD(U8(0xd2), cv_l[0], cv_l[0])
	PSHUFD(U8(0xd2), cv_l[1], cv_l[1])
	PSHUFD(U8(0x6c), cv_r[0], cv_r[0])
	PSHUFD(U8(0x6c), cv_r[1], cv_r[1])

	temp := cv_l[0]
	cv_l[0] = cv_l[1]
	cv_l[1] = cv_r[1]
	cv_r[1] = cv_r[0]
	cv_r[0] = temp
}

func mix_even(cv_l, cv_r []VecVirtual, const_v []Mem) {
	Comment("mix_even")

	add_blk(cv_l, cv_r)
	rotate_blk_even_alpha(cv_l)
	xor_with_const(cv_l, const_v)
	add_blk(cv_r, cv_l)
	rotate_blk_even_beta(cv_r)
	add_blk(cv_l, cv_r)
	rotate_msg_gamma(cv_r)
}

func mix_odd(cv_l, cv_r []VecVirtual, const_v []Mem) {
	Comment("mix_odd")

	add_blk(cv_l, cv_r)
	rotate_blk_odd_alpha(cv_l)
	xor_with_const(cv_l, const_v)
	add_blk(cv_r, cv_l)
	rotate_blk_odd_beta(cv_r)
	add_blk(cv_l, cv_r)
	rotate_msg_gamma(cv_r)
}

func compress(cv_l, cv_r []VecVirtual, pdMsgBlk Mem /* uint32 */) {
	Comment("compress")

	/**
	__m128i const_v[2];			// step function constant
	LSH256SSE2_internal i_state[1];
	int i;
	*/
	const_v := make([]Mem, 2)

	i_state := LSH256SSE2_internal{
		submsg_e_l: []VecVirtual{XMM(), XMM()},
		submsg_e_r: []VecVirtual{XMM(), XMM()},
		submsg_o_l: []VecVirtual{XMM(), XMM()},
		submsg_o_r: []VecVirtual{XMM(), XMM()},
	}

	load_msg_blk(i_state, pdMsgBlk)

	msg_add_even(cv_l, cv_r, i_state)
	load_sc_const_v(const_v, 0)
	mix_even(cv_l, cv_r, const_v)
	word_perm(cv_l, cv_r)

	msg_add_odd(cv_l, cv_r, i_state)
	load_sc_const_v(const_v, 8)
	mix_odd(cv_l, cv_r, const_v)
	word_perm(cv_l, cv_r)

	for i := 1; i < NUM_STEPS/2; i++ {
		msg_exp_even(i_state)
		msg_add_even(cv_l, cv_r, i_state)
		load_sc_const_v(const_v, 16*i)
		mix_even(cv_l, cv_r, const_v)
		word_perm(cv_l, cv_r)

		msg_exp_odd(i_state)
		msg_add_odd(cv_l, cv_r, i_state)
		load_sc_const_v(const_v, 16*i+8)
		mix_odd(cv_l, cv_r, const_v)
		word_perm(cv_l, cv_r)
	}

	msg_exp_even(i_state)
	msg_add_even(cv_l, cv_r, i_state)
}

func fin(cv_l, cv_r []VecVirtual) {
	Comment("fin")

	/**
	cv_l[0] = XOR(cv_l[0], cv_r[0]);
	cv_l[1] = XOR(cv_l[1], cv_r[1]);
	*/
	PXOR(cv_r[0], cv_l[0])
	PXOR(cv_r[1], cv_l[1])
}

func get_hash(cv_l []VecVirtual, pbHashVal Mem, algtype Register) {
	Comment("get_hash")

	/**
	lsh_u8 hash_val[LSH256_HASH_VAL_MAX_BYTE_LEN] = { 0x0, };
	*/

	/**
	#define LSH_GET_SMALL_HASHBIT(lsh_type_val)		((lsh_type_val)>>24)
	#define LSH_GET_HASHBYTE(lsh_type_val)			((lsh_type_val) & 0xffff)
	#define LSH_GET_HASHBIT(lsh_type_val)			((LSH_GET_HASHBYTE(lsh_type_val)<<3)-LSH_GET_SMALL_HASHBIT(lsh_type_val))


	lsh_uint hash_val_byte_len = LSH_GET_HASHBYTE(algtype);
	lsh_uint hash_val_bit_len = LSH_GET_SMALL_HASHBIT(algtype);
	*/
	hash_val_byte_len := GP64()
	MOVQ(algtype, hash_val_byte_len)
	ANDQ(U32(0xFFFF), hash_val_byte_len)

	/**
	STORE(hash_val, cv_l[0]);
	STORE((hash_val + 16), cv_l[1]);
	*/
	MOVOU(cv_l[0], pbHashVal.Offset(0x00))
	MOVOU(cv_l[1], pbHashVal.Offset(0x10))

	//memcpy(pbHashVal, hash_val, sizeof(lsh_u8) * hash_val_byte_len);

	Label("get_hash_ret")
}
func LSH256InitSSE2() {
	TEXT("lsh256InitSSE2", NOSPLIT, "func(ctx *lsh256ContextAsmData)")

	ctx := Dereference(Param("ctx"))
	ctx_algtype := Load(ctx.Field("algtype"), GP64())
	ctx_cv_l := Mem{Base: Load(ctx.Field("cv_l").Base(), GP64())}
	ctx_cv_r := Mem{Base: Load(ctx.Field("cv_r").Base(), GP64())}

	copyIV := func(iv Mem) {
		tmp := XMM()

		MOVO(iv.Offset(0x00), tmp)
		MOVAPS(tmp, ctx_cv_l.Offset(0x00))
		MOVO(iv.Offset(0x10), tmp)
		MOVAPS(tmp, ctx_cv_l.Offset(0x10))

		MOVO(iv.Offset(0x20), tmp)
		MOVAPS(tmp, ctx_cv_r.Offset(0x00))
		MOVO(iv.Offset(0x30), tmp)
		MOVAPS(tmp, ctx_cv_r.Offset(0x10))
	}

	CMPQ(ctx_algtype, U32(LSH_TYPE_256_256))
	JNE(LabelRef("not_256"))
	{
		copyIV(lsh256avoconst.G_IV256)
		JMP(LabelRef("ret"))
	}
	Label("not_256")
	{
		copyIV(lsh256avoconst.G_IV224)
		JMP(LabelRef("ret"))
	}

	Label("ret")
	RET()
}

func LSH256UpdateSSE2() {
	TEXT("lsh256UpdateSSE2", NOSPLIT, "func(ctx *lsh256ContextAsmData, data []byte, remain_msg_byte int) int")
	Comment("return databytelen")

	ctx := Dereference(Param("ctx"))
	//ctx_algtype := Load(ctx.Field("algtype"), GP64())
	ctx_remain_databytelen := Load(ctx.Field("remain_databytelen"), GP64())
	ctx_cv_l := Mem{Base: Load(ctx.Field("cv_l").Base(), GP64())}
	ctx_cv_r := Mem{Base: Load(ctx.Field("cv_r").Base(), GP64())}
	ctx_last_block := Mem{Base: Load(ctx.Field("last_block").Base(), GP64())}

	data := Mem{Base: Load(Param("data").Base(), GP64())}
	databytelen := Load(Param("data").Len(), GP64())

	remain_msg_byte := Load(Param("remain_msg_byte"), GP64())

	/**
	__m128i cv_l[2];
	__m128i cv_r[2];
	*/
	cv_l := []VecVirtual{XMM(), XMM()}
	cv_r := []VecVirtual{XMM(), XMM()}

	load_blk(cv_l, ctx_cv_l)
	load_blk(cv_r, ctx_cv_r)

	/**
	if (remain_msg_byte > 0){
	*/
	CMPQ(remain_msg_byte, U32(0))
	JLE(LabelRef("end_if1"))
	{
		/**
		lsh_uint more_byte = LSH256_MSG_BLK_BYTE_LEN - remain_msg_byte;
		*/
		more_byte := GP64()
		MOVQ(U64(LSH256_MSG_BLK_BYTE_LEN), more_byte)
		SUBQ(remain_msg_byte, more_byte)

		/**
		memcpy(ctx->last_block + remain_msg_byte, data, more_byte);
		이 부분은 go에서 먼저 처리하고 들어와야 함.
		*/

		compress(cv_l, cv_r, ctx_last_block)

		/**
		data += more_byte;
		*/
		Comment("data += more_byte;")
		ADDQ(more_byte, data.Base)

		/**
		databytelen -= more_byte;
		*/
		SUBQ(more_byte, databytelen)

		/**
		remain_msg_byte = 0;
		ctx->remain_databytelen = 0;
		*/
		MOVQ(U64(0), remain_msg_byte)
		MOVQ(U64(0), ctx_remain_databytelen)
	}
	Label("end_if1")

	/**
	while (databytelen >= LSH256_MSG_BLK_BYTE_LEN)
	*/
	Label("while_0_loop")
	CMPQ(databytelen, U32(LSH256_MSG_BLK_BYTE_LEN))
	JL(LabelRef("while_0_done"))
	{
		compress(cv_l, cv_r, data)

		/**
		data += LSH256_MSG_BLK_BYTE_LEN;
		databytelen -= LSH256_MSG_BLK_BYTE_LEN;
		*/
		Comment("data += LSH256_MSG_BLK_BYTE_LEN")
		ADDQ(U32(LSH256_MSG_BLK_BYTE_LEN), data.Base)
		SUBQ(U8(LSH256_MSG_BLK_BYTE_LEN), databytelen)

		JMP(LabelRef("while_0_loop"))
	}
	Label("while_0_done")

	store_blk(ctx_cv_l, cv_l)
	store_blk(ctx_cv_r, cv_r)

	Comment("return value")
	Store(databytelen, ReturnIndex(0))

	RET()
}

func LSH256FinalSSE2() {
	TEXT("lsh256FinalSSE2", NOSPLIT, "func(ctx *lsh256ContextAsmData, hashval []byte)")

	ctx := Dereference(Param("ctx"))
	ctx_algtype := Load(ctx.Field("algtype"), GP64())
	//ctx_remain_databytelen := Load(ctx.Field("remain_databytelen"), GP64())
	ctx_cv_l := Mem{Base: Load(ctx.Field("cv_l").Base(), GP64())}
	ctx_cv_r := Mem{Base: Load(ctx.Field("cv_r").Base(), GP64())}
	ctx_last_block := Mem{Base: Load(ctx.Field("last_block").Base(), GP64())}

	hashval := Mem{Base: Load(Param("hashval").Base(), GP64())}

	/**
	__m128i cv_l[2];
	__m128i cv_r[2];
	*/
	cv_l := []VecVirtual{XMM(), XMM()}
	cv_r := []VecVirtual{XMM(), XMM()}

	load_blk(cv_l, ctx_cv_l)
	load_blk(cv_r, ctx_cv_r)

	compress(cv_l, cv_r, ctx_last_block)

	fin(cv_l, cv_r)
	get_hash(cv_l, hashval, ctx_algtype)

	RET()
}
