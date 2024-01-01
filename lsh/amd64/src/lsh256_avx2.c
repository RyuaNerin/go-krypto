/*
 * Copyright (c) 2016 NSR (National Security Research Institute)
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy 
 * of this software and associated documentation files (the "Software"), to deal 
 * in the Software without restriction, including without limitation the rights 
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell 
 * copies of the Software, and to permit persons to whom the Software is 
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in 
 * all copies or substantial portions of the Software.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR 
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, 
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE 
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER 
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, 
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN 
 * THE SOFTWARE.
 */

#include "../../string.h"
#include "../../lsh_local.h"


#if defined(_MSC_VER)
#include "intrin.h"
#else
#include "emmintrin.h"
#include "xmmintrin.h"
#include "immintrin.h"
#include "x86intrin.h"
#endif

/* -------------------------------------------------------- *
*  LSH: parameters
*  -------------------------------------------------------- */
#define MSG_BLK_WORD_LEN		32
#define CV_WORD_LEN				16
#define CONST_WORD_LEN			8
#define HASH_VAL_MAX_WORD_LEN	8

#define WORD_BIT_LEN			32

/* -------------------------------------------------------- */

#define NUM_STEPS				26

#define ROT_EVEN_ALPHA			29
#define ROT_EVEN_BETA			1
#define ROT_ODD_ALPHA			5
#define ROT_ODD_BETA			17

/* -------------------------------------------------------- */

#define _LSH256_
#define _VER_256_BIT_REG_

/* -------------------------------------------------------- *
*  LSH: variables
*  -------------------------------------------------------- */

typedef struct LSH_ALIGNED_(32){
	LSH_ALIGNED_(16) lsh_type algtype;
	LSH_ALIGNED_(16) lsh_uint remain_databytelen;
	LSH_ALIGNED_(32) __m256i cv_l[1];				// left chaining variable
	LSH_ALIGNED_(32) __m256i cv_r[1];				// right chaining variable
	LSH_ALIGNED_(32) lsh_u8 last_block[LSH256_MSG_BLK_BYTE_LEN];
} LSH256AVX2_Context;

typedef struct LSH_ALIGNED_(32){
	__m256i submsg_e_l[1];
	__m256i submsg_e_r[1];
	__m256i submsg_o_l[1];
	__m256i submsg_o_r[1];
} LSH256AVX2_internal;

/* -------------------------------------------------------- *
*  LSH: iv
*  -------------------------------------------------------- */
static const LSH_ALIGNED_(32) lsh_u32 g_IV224[CV_WORD_LEN] = {
	0x068608D3, 0x62D8F7A7, 0xD76652AB, 0x4C600A43, 0xBDC40AA8, 0x1ECA0B68, 0xDA1A89BE, 0x3147D354,
	0x707EB4F9, 0xF65B3862, 0x6B0B2ABE, 0x56B8EC0A, 0xCF237286, 0xEE0D1727, 0x33636595, 0x8BB8D05F,
};

static const LSH_ALIGNED_(32) lsh_u32 g_IV256[CV_WORD_LEN] = {
	0x46a10f1f, 0xfddce486, 0xb41443a8, 0x198e6b9d, 0x3304388d, 0xb0f5a3c7, 0xb36061c4, 0x7adbd553,
	0x105d5378, 0x2f74de54, 0x5c2f2d95, 0xf2553fbe, 0x8051357a, 0x138668c8, 0x47aa4484, 0xe01afb41
};

/* -------------------------------------------------------- *
*  LSH: step constants
*  -------------------------------------------------------- */

static const LSH_ALIGNED_(32) lsh_u32 g_StepConstants[CONST_WORD_LEN * NUM_STEPS] = {
	0x917caf90, 0x6c1b10a2, 0x6f352943, 0xcf778243, 0x2ceb7472, 0x29e96ff2, 0x8a9ba428, 0x2eeb2642,
	0x0e2c4021, 0x872bb30e, 0xa45e6cb2, 0x46f9c612, 0x185fe69e, 0x1359621b, 0x263fccb2, 0x1a116870,
	0x3a6c612f, 0xb2dec195, 0x02cb1f56, 0x40bfd858, 0x784684b6, 0x6cbb7d2e, 0x660c7ed8, 0x2b79d88a,
	0xa6cd9069, 0x91a05747, 0xcdea7558, 0x00983098, 0xbecb3b2e, 0x2838ab9a, 0x728b573e, 0xa55262b5,
	0x745dfa0f, 0x31f79ed8, 0xb85fce25, 0x98c8c898, 0x8a0669ec, 0x60e445c2, 0xfde295b0, 0xf7b5185a,
	0xd2580983, 0x29967709, 0x182df3dd, 0x61916130, 0x90705676, 0x452a0822, 0xe07846ad, 0xaccd7351,
	0x2a618d55, 0xc00d8032, 0x4621d0f5, 0xf2f29191, 0x00c6cd06, 0x6f322a67, 0x58bef48d, 0x7a40c4fd,
	0x8beee27f, 0xcd8db2f2, 0x67f2c63b, 0xe5842383, 0xc793d306, 0xa15c91d6, 0x17b381e5, 0xbb05c277,
	0x7ad1620a, 0x5b40a5bf, 0x5ab901a2, 0x69a7a768, 0x5b66d9cd, 0xfdee6877, 0xcb3566fc, 0xc0c83a32,
	0x4c336c84, 0x9be6651a, 0x13baa3fc, 0x114f0fd1, 0xc240a728, 0xec56e074, 0x009c63c7, 0x89026cf2,
	0x7f9ff0d0, 0x824b7fb5, 0xce5ea00f, 0x605ee0e2, 0x02e7cfea, 0x43375560, 0x9d002ac7, 0x8b6f5f7b,
	0x1f90c14f, 0xcdcb3537, 0x2cfeafdd, 0xbf3fc342, 0xeab7b9ec, 0x7a8cb5a3, 0x9d2af264, 0xfacedb06,
	0xb052106e, 0x99006d04, 0x2bae8d09, 0xff030601, 0xa271a6d6, 0x0742591d, 0xc81d5701, 0xc9a9e200,
	0x02627f1e, 0x996d719d, 0xda3b9634, 0x02090800, 0x14187d78, 0x499b7624, 0xe57458c9, 0x738be2c9,
	0x64e19d20, 0x06df0f36, 0x15d1cb0e, 0x0b110802, 0x2c95f58c, 0xe5119a6d, 0x59cd22ae, 0xff6eac3c,
	0x467ebd84, 0xe5ee453c, 0xe79cd923, 0x1c190a0d, 0xc28b81b8, 0xf6ac0852, 0x26efd107, 0x6e1ae93b,
	0xc53c41ca, 0xd4338221, 0x8475fd0a, 0x35231729, 0x4e0d3a7a, 0xa2b45b48, 0x16c0d82d, 0x890424a9,
	0x017e0c8f, 0x07b5a3f5, 0xfa73078e, 0x583a405e, 0x5b47b4c8, 0x570fa3ea, 0xd7990543, 0x8d28ce32,
	0x7f8a9b90, 0xbd5998fc, 0x6d7a9688, 0x927a9eb6, 0xa2fc7d23, 0x66b38e41, 0x709e491a, 0xb5f700bf,
	0x0a262c0f, 0x16f295b9, 0xe8111ef5, 0x0d195548, 0x9f79a0c5, 0x1a41cfa7, 0x0ee7638a, 0xacf7c074,
	0x30523b19, 0x09884ecf, 0xf93014dd, 0x266e9d55, 0x191a6664, 0x5c1176c1, 0xf64aed98, 0xa4b83520,
	0x828d5449, 0x91d71dd8, 0x2944f2d6, 0x950bf27b, 0x3380ca7d, 0x6d88381d, 0x4138868e, 0x5ced55c4,
	0x0fe19dcb, 0x68f4f669, 0x6e37c8ff, 0xa0fe6e10, 0xb44b47b0, 0xf5c0558a, 0x79bf14cf, 0x4a431a20,
	0xf17f68da, 0x5deb5fd1, 0xa600c86d, 0x9f6c7eb0, 0xff92f864, 0xb615e07f, 0x38d3e448, 0x8d5d3a6a,
	0x70e843cb, 0x494b312e, 0xa6c93613, 0x0beb2f4f, 0x928b5d63, 0xcbf66035, 0x0cb82c80, 0xea97a4f7,
	0x592c0f3b, 0x947c5f77, 0x6fff49b9, 0xf71a7e5a, 0x1de8c0f5, 0xc2569600, 0xc4e4ac8c, 0x823c9ce1
};


/* -------------------------------------------------------- *
*  LSH : permutation information
*  -------------------------------------------------------- */

static const LSH_ALIGNED_(32) lsh_u32 g_BytePermInfo[8] = {
	0x03020100, 0x06050407, 0x09080b0a, 0x0c0f0e0d, 0x10131211, 0x15141716, 0x1a19181b, 0x1f1e1d1c };
static const LSH_ALIGNED_(32) lsh_u32 g_MsgWordPermInfo[8] = {
	0x0f0e0d0c, 0x0b0a0908, 0x03020100, 0x07060504, 0x1f1e1d1c, 0x13121110, 0x17161514, 0x1b1a1918
};


/* -------------------------------------------------------- *
*  LSH: functions
*  -------------------------------------------------------- */

#define LOAD(x) _mm256_loadu_si256((__m256i*)x)
#define STORE(x,y) _mm256_storeu_si256((__m256i*)x, y)
#define XOR(x,y) _mm256_xor_si256(x,y)
#define OR(x,y) _mm256_or_si256(x,y)
#define AND(x,y) _mm256_and_si256(x,y)
#define SHUFFLE8(x,y) _mm256_shuffle_epi8(x,y)

#define ADD(x,y) _mm256_add_epi32(x,y)
#define SHIFT_L(x,r) _mm256_slli_epi32(x,r)
#define SHIFT_R(x,r) _mm256_srli_epi32(x,r)

static INLINE void load_blk(__m256i* dest, const void* src){
	dest[0] = LOAD(src);
}

static INLINE void store_blk(__m256i* dest, const __m256i* src){
	STORE(dest, src[0]);
}

static INLINE void load_msg_blk(LSH256AVX2_internal * i_state, const lsh_u32* msgblk){
	load_blk(i_state->submsg_e_l, msgblk + 0);
	load_blk(i_state->submsg_e_r, msgblk + 8);
	load_blk(i_state->submsg_o_l, msgblk + 16);
	load_blk(i_state->submsg_o_r, msgblk + 24);
}
static INLINE void msg_exp_even(LSH256AVX2_internal * i_state, const __m256i perm_step){
	i_state->submsg_e_l[0] = ADD(i_state->submsg_o_l[0], SHUFFLE8(i_state->submsg_e_l[0], perm_step));
	i_state->submsg_e_r[0] = ADD(i_state->submsg_o_r[0], SHUFFLE8(i_state->submsg_e_r[0], perm_step));
}
static INLINE void msg_exp_odd(LSH256AVX2_internal * i_state, const __m256i perm_step){
	i_state->submsg_o_l[0] = ADD(i_state->submsg_e_l[0], SHUFFLE8(i_state->submsg_o_l[0], perm_step));
	i_state->submsg_o_r[0] = ADD(i_state->submsg_e_r[0], SHUFFLE8(i_state->submsg_o_r[0], perm_step));
}
static INLINE void load_sc(__m256i* const_v, lsh_uint i){
	load_blk(const_v, g_StepConstants + i);
}
static INLINE void msg_add_even(__m256i* cv_l, __m256i* cv_r, const LSH256AVX2_internal * i_state){
	*cv_l = XOR(*cv_l, i_state->submsg_e_l[0]);
	*cv_r = XOR(*cv_r, i_state->submsg_e_r[0]);
}
static INLINE void msg_add_odd(__m256i* cv_l, __m256i* cv_r, const LSH256AVX2_internal * i_state){
	*cv_l = XOR(*cv_l, i_state->submsg_o_l[0]);
	*cv_r = XOR(*cv_r, i_state->submsg_o_r[0]);
}
static INLINE void add_blk(__m256i* cv_l, const __m256i* cv_r){
	*cv_l = ADD(*cv_l, *cv_r);
}

static INLINE void rotate_blk_even_alpha(__m256i* cv){
	*cv = OR(SHIFT_L(*cv, ROT_EVEN_ALPHA), SHIFT_R(*cv, WORD_BIT_LEN - ROT_EVEN_ALPHA));
}

static INLINE void rotate_blk_even_beta(__m256i* cv){
	*cv = OR(SHIFT_L(*cv, ROT_EVEN_BETA), SHIFT_R(*cv, WORD_BIT_LEN - ROT_EVEN_BETA));
}

static INLINE void rotate_blk_odd_alpha(__m256i* cv){
	*cv = OR(SHIFT_L(*cv, ROT_ODD_ALPHA), SHIFT_R(*cv, WORD_BIT_LEN - ROT_ODD_ALPHA));
}

static INLINE void rotate_blk_odd_beta(__m256i* cv){
	*cv = OR(SHIFT_L(*cv, ROT_ODD_BETA), SHIFT_R(*cv, WORD_BIT_LEN - ROT_ODD_BETA));
}

static INLINE void xor_with_const(__m256i* cv_l, const __m256i* const_v){
	*cv_l = XOR(*cv_l, *const_v);
}
static INLINE void rotate_msg_gamma(__m256i* cv_r, const __m256i byte_perm_step){
	*cv_r = SHUFFLE8(*cv_r, byte_perm_step);
}
static INLINE void word_perm(__m256i* cv_l, __m256i* cv_r){
	__m256i temp;
	temp = _mm256_shuffle_epi32(*cv_l, 0xd2);
	*cv_r = _mm256_shuffle_epi32(*cv_r, 0x6c);
	*cv_l = _mm256_permute2x128_si256(temp, *cv_r, 0x31);
	*cv_r = _mm256_permute2x128_si256(temp, *cv_r, 0x20);
};

/* -------------------------------------------------------- *
*  step function
*  -------------------------------------------------------- */

static INLINE void mix_even(__m256i* cv_l, __m256i* cv_r, const __m256i* const_v, const __m256i byte_perm_step){
	add_blk(cv_l, cv_r);
	rotate_blk_even_alpha(cv_l);
	xor_with_const(cv_l, const_v);
	add_blk(cv_r, cv_l);
	rotate_blk_even_beta(cv_r);
	add_blk(cv_l, cv_r);
	rotate_msg_gamma(cv_r, byte_perm_step);
}

static INLINE void mix_odd(__m256i* cv_l, __m256i* cv_r, const __m256i* const_v, const __m256i byte_perm_step){
	add_blk(cv_l, cv_r);
	rotate_blk_odd_alpha(cv_l);
	xor_with_const(cv_l, const_v);
	add_blk(cv_r, cv_l);
	rotate_blk_odd_beta(cv_r);
	add_blk(cv_l, cv_r);
	rotate_msg_gamma(cv_r, byte_perm_step);
}

/* -------------------------------------------------------- *
*  compression function
*  -------------------------------------------------------- */

static INLINE void compress(__m256i* cv_l, __m256i* cv_r, const lsh_u32 pdMsgBlk[MSG_BLK_WORD_LEN])
{
	__m256i const_v[1];				// step function constant
	__m256i byte_perm_step;		// byte permutation info
	__m256i word_perm_step;	// msg_word permutation info
	LSH256AVX2_internal i_state[1];
	int i;

	

	byte_perm_step = LOAD(g_BytePermInfo);
	word_perm_step = LOAD(g_MsgWordPermInfo);

	load_msg_blk(i_state, pdMsgBlk);
 
	msg_add_even(cv_l, cv_r, i_state);
	load_sc(const_v, 0);
	mix_even(cv_l, cv_r, const_v, byte_perm_step);
	word_perm(cv_l, cv_r);

	msg_add_odd(cv_l, cv_r, i_state); 
	load_sc(const_v, 8);
	mix_odd(cv_l, cv_r, const_v, byte_perm_step);
	word_perm(cv_l, cv_r);

	for (i = 1; i < NUM_STEPS / 2; i++){
		msg_exp_even(i_state, word_perm_step); 
		msg_add_even(cv_l, cv_r, i_state); 
		load_sc(const_v, 16 * i);
		mix_even(cv_l, cv_r, const_v, byte_perm_step);
		word_perm(cv_l, cv_r);

		msg_exp_odd(i_state, word_perm_step); 
		msg_add_odd(cv_l, cv_r, i_state); 
		load_sc(const_v, 16 * i + 8);
		mix_odd(cv_l, cv_r, const_v, byte_perm_step);
		word_perm(cv_l, cv_r);
	}

	msg_exp_even(i_state, word_perm_step); 
	msg_add_even(cv_l, cv_r, i_state);
}


/* -------------------------------------------------------- */

static INLINE void init224(LSH256AVX2_Context * state)
{
	load_blk(state->cv_l, g_IV224);
	load_blk(state->cv_r, g_IV224 + 8);
}

static INLINE void init256(LSH256AVX2_Context * state)
{
	load_blk(state->cv_l, g_IV256);
	load_blk(state->cv_r, g_IV256 + 8);
}

/* -------------------------------------------------------- */

static INLINE void fin(__m256i* cv_l, const __m256i* cv_r)
{
	cv_l[0] = XOR(cv_l[0], cv_r[0]);
}

/* -------------------------------------------------------- */

static INLINE void get_hash(__m256i* cv_l, lsh_u8 * pbHashVal, const lsh_type algtype)
{
	STORE(pbHashVal, cv_l[0]);
}

/* -------------------------------------------------------- */

void lsh256_avx2_init(struct LSH256_Context * _ctx, const lsh_type algtype){
	LSH256AVX2_Context* ctx = (LSH256AVX2_Context*)_ctx;
	ctx->algtype = algtype;
	ctx->remain_databytelen = 0;

	if (algtype == LSH_TYPE_256_256) {
		init256(ctx);
	} else {
		init224(ctx);
	}
}

void lsh256_avx2_update(struct LSH256_Context * _ctx, const lsh_u8 * data, size_t databytelen){
	__m256i cv_l[1];
	__m256i cv_r[1];

	LSH256AVX2_Context* ctx = (LSH256AVX2_Context*)_ctx;
	lsh_uint remain_msg_byte;

	remain_msg_byte = ctx->remain_databytelen;

	if (databytelen + remain_msg_byte < LSH256_MSG_BLK_BYTE_LEN){
		memcpy(ctx->last_block + remain_msg_byte, data, databytelen);
		ctx->remain_databytelen += databytelen;
		remain_msg_byte += (lsh_uint)databytelen;
		return;
	}

	load_blk(cv_l, ctx->cv_l);
	load_blk(cv_r, ctx->cv_r);

	if (remain_msg_byte > 0){
		lsh_uint more_byte = LSH256_MSG_BLK_BYTE_LEN - remain_msg_byte;
		memcpy(ctx->last_block + remain_msg_byte, data, more_byte);
		compress(cv_l, cv_r, (lsh_u32*)ctx->last_block);
		data += more_byte;
		databytelen -= more_byte;
		remain_msg_byte = 0;
		ctx->remain_databytelen = 0;
	}

	while (databytelen >= LSH256_MSG_BLK_BYTE_LEN)
	{
		compress(cv_l, cv_r, (lsh_u32*)data);
		data += LSH256_MSG_BLK_BYTE_LEN;
		databytelen -= LSH256_MSG_BLK_BYTE_LEN;
	}

	store_blk(ctx->cv_l, cv_l);
	store_blk(ctx->cv_r, cv_r);

	if (databytelen > 0){
		memcpy(ctx->last_block, data, databytelen);
		ctx->remain_databytelen = databytelen;
	}
}

void lsh256_avx2_final(struct LSH256_Context * _ctx, lsh_u8 * hashval){
	__m256i cv_l[1];
	__m256i cv_r[1];
	LSH256AVX2_Context* ctx = (LSH256AVX2_Context*)_ctx;
	lsh_uint remain_msg_byte;

	remain_msg_byte = ctx->remain_databytelen;

	ctx->last_block[remain_msg_byte] = 0x80;
	memset(ctx->last_block + remain_msg_byte + 1, 0, LSH256_MSG_BLK_BYTE_LEN - remain_msg_byte - 1);

	load_blk(cv_l, ctx->cv_l);
	load_blk(cv_r, ctx->cv_r);

	compress(cv_l, cv_r, (lsh_u32*)ctx->last_block);

	fin(cv_l, cv_r);
	get_hash(cv_l, hashval, ctx->algtype);
}
