package lsh256avx2

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"

	. "kryptosimd/avoutil"
	. "kryptosimd/avoutil/simd"
	. "kryptosimd/lsh/avx2"
	. "kryptosimd/lsh256/avo/lsh256avoconst"
	. "kryptosimd/lsh256/avo/lsh256common"
)

//	typedef struct LSH_ALIGNED_(32){
//		__m256i submsg_e_l[1];
//		__m256i submsg_e_r[1];
//		__m256i submsg_o_l[1];
//		__m256i submsg_o_r[1];
//	} LSH256AVX2_internal;
type LSH256AVX2_internal struct {
	submsg_e_l []VecVirtual
	submsg_e_r []VecVirtual
	submsg_o_l []VecVirtual
	submsg_o_r []VecVirtual
}

// static INLINE void load_blk(__m256i* dest, const void* src){
func load_blk_mem2vec(dest []VecVirtual, src Mem) {
	Comment("load_blk_mem2vec")

	//dest[0] = LOAD(src);
	LOAD(dest[0], src)
}
func load_blk_vec2mem(dest Mem, src []VecVirtual) {
	Comment("load_blk_vec2mem")

	//dest[0] = LOAD(src);
	LOAD(dest, src[0])
}
func load_blk_mem2mem(dest, src Mem) {
	Comment("load_blk_mem2mem")

	tmp := YMM()
	//dest[0] = LOAD(src);
	LOAD(tmp, src)
	LOAD(dest, tmp)
}

// static INLINE void store_blk(__m256i* dest, const __m256i* src){
func store_blk(dest Mem, src []VecVirtual) {
	Comment("store_blk")

	//STORE(dest, src[0]);
	STORE(dest, src[0])
}

// static INLINE void load_msg_blk(LSH256AVX2_internal * i_state, const lsh_u32* msgblk){
func load_msg_blk(i_state LSH256AVX2_internal, msgblk Mem) {
	Comment("load_msg_blk")

	//load_blk(i_state->submsg_e_l, msgblk + 0);
	load_blk_mem2vec(i_state.submsg_e_l, msgblk.Offset(0*4))
	//load_blk(i_state->submsg_e_r, msgblk + 8);
	load_blk_mem2vec(i_state.submsg_e_r, msgblk.Offset(8*4))
	//load_blk(i_state->submsg_o_l, msgblk + 16);
	load_blk_mem2vec(i_state.submsg_o_l, msgblk.Offset(16*4))
	//load_blk(i_state->submsg_o_r, msgblk + 24);
	load_blk_mem2vec(i_state.submsg_o_r, msgblk.Offset(24*4))
}

// static INLINE void msg_exp_even(LSH256AVX2_internal * i_state, const __m256i perm_step){
func msg_exp_even(i_state LSH256AVX2_internal, perm_step Op) {
	Comment("msg_exp_even")

	//i_state->submsg_e_l[0] = ADD(i_state->submsg_o_l[0], SHUFFLE8(i_state->submsg_e_l[0], perm_step));
	SHUFFLE8(i_state.submsg_e_l[0], i_state.submsg_e_l[0], perm_step)
	ADD32(i_state.submsg_e_l[0], i_state.submsg_e_l[0], i_state.submsg_o_l[0])
	//i_state->submsg_e_r[0] = ADD(i_state->submsg_o_r[0], SHUFFLE8(i_state->submsg_e_r[0], perm_step));
	SHUFFLE8(i_state.submsg_e_r[0], i_state.submsg_e_r[0], perm_step)
	ADD32(i_state.submsg_e_r[0], i_state.submsg_e_r[0], i_state.submsg_o_r[0])
}

// static INLINE void msg_exp_odd(LSH256AVX2_internal * i_state, const __m256i perm_step){
func msg_exp_odd(i_state LSH256AVX2_internal, perm_step Op) {
	Comment("msg_exp_odd")

	//i_state->submsg_o_l[0] = ADD(i_state->submsg_e_l[0], SHUFFLE8(i_state->submsg_o_l[0], perm_step));
	SHUFFLE8(i_state.submsg_o_l[0], i_state.submsg_o_l[0], perm_step)
	ADD32(i_state.submsg_o_l[0], i_state.submsg_o_l[0], i_state.submsg_e_l[0])
	//i_state->submsg_o_r[0] = ADD(i_state->submsg_e_r[0], SHUFFLE8(i_state->submsg_o_r[0], perm_step));
	SHUFFLE8(i_state.submsg_o_r[0], i_state.submsg_o_r[0], perm_step)
	ADD32(i_state.submsg_o_r[0], i_state.submsg_o_r[0], i_state.submsg_e_r[0])
}

// static INLINE void load_sc(__m256i* const_v, lsh_uint i){
func load_sc(const_v []VecVirtual, i int) {
	Comment("load_sc")

	//load_blk(const_v, g_StepConstants + i);
	load_blk_mem2vec(const_v, G_StepConstants.Offset(i*4))
	//const_v[0] = G_StepConstants.Offset(i * 4)
}

// static INLINE void msg_add_even(__m256i* cv_l, __m256i* cv_r, const LSH256AVX2_internal * i_state){
func msg_add_even(cv_l, cv_r []VecVirtual, i_state LSH256AVX2_internal) {
	Comment("msg_add_even")

	//*cv_l = XOR(*cv_l, i_state->submsg_e_l[0]);
	XOR(cv_l[0], cv_l[0], i_state.submsg_e_l[0])
	//*cv_r = XOR(*cv_r, i_state->submsg_e_r[0]);
	XOR(cv_r[0], cv_r[0], i_state.submsg_e_r[0])
}

// static INLINE void msg_add_odd(__m256i* cv_l, __m256i* cv_r, const LSH256AVX2_internal * i_state){
func msg_add_odd(cv_l, cv_r []VecVirtual, i_state LSH256AVX2_internal) {
	Comment("msg_add_odd")

	//*cv_l = XOR(*cv_l, i_state->submsg_o_l[0]);
	XOR(cv_l[0], cv_l[0], i_state.submsg_o_l[0])
	//*cv_r = XOR(*cv_r, i_state->submsg_o_r[0]);
	XOR(cv_r[0], cv_r[0], i_state.submsg_o_r[0])
}

// static INLINE void add_blk(__m256i* cv_l, const __m256i* cv_r){
func add_blk(cv_l, cv_r []VecVirtual) {
	Comment("add_blk")

	//*cv_l = ADD(*cv_l, *cv_r);
	ADD32(cv_l[0], cv_l[0], cv_r[0])
}

// static INLINE void rotate_blk_even_alpha(__m256i* cv){
func rotate_blk_even_alpha(cv []VecVirtual) {
	Comment("rotate_blk_even_alpha")

	tmpYmm := YMM()
	//*cv = OR(SHIFT_L(*cv, ROT_EVEN_ALPHA), SHIFT_R(*cv, WORD_BIT_LEN - ROT_EVEN_ALPHA));
	SHIFT_L32(tmpYmm, cv[0], U8(ROT_EVEN_ALPHA))
	SHIFT_R32(cv[0], cv[0], U8(WORD_BIT_LEN-ROT_EVEN_ALPHA))
	OR(cv[0], cv[0], tmpYmm)
}

// static INLINE void rotate_blk_even_beta(__m256i* cv){
func rotate_blk_even_beta(cv []VecVirtual) {
	Comment("rotate_blk_even_beta")

	tmpYmm := YMM()
	//*cv = OR(SHIFT_L(*cv, ROT_EVEN_BETA), SHIFT_R(*cv, WORD_BIT_LEN - ROT_EVEN_BETA));
	SHIFT_L32(tmpYmm, cv[0], U8(ROT_EVEN_BETA))
	SHIFT_R32(cv[0], cv[0], U8(WORD_BIT_LEN-ROT_EVEN_BETA))
	OR(cv[0], cv[0], tmpYmm)
}

// static INLINE void rotate_blk_odd_alpha(__m256i* cv){
func rotate_blk_odd_alpha(cv []VecVirtual) {
	Comment("rotate_blk_odd_alpha")

	tmpYmm := YMM()
	//*cv = OR(SHIFT_L(*cv, ROT_ODD_ALPHA), SHIFT_R(*cv, WORD_BIT_LEN - ROT_ODD_ALPHA));
	SHIFT_L32(tmpYmm, cv[0], U8(ROT_ODD_ALPHA))
	SHIFT_R32(cv[0], cv[0], U8(WORD_BIT_LEN-ROT_ODD_ALPHA))
	OR(cv[0], cv[0], tmpYmm)
}

// static INLINE void rotate_blk_odd_beta(__m256i* cv){
func rotate_blk_odd_beta(cv []VecVirtual) {
	Comment("rotate_blk_odd_beta")

	tmpYmm := YMM()
	//*cv = OR(SHIFT_L(*cv, ROT_ODD_BETA), SHIFT_R(*cv, WORD_BIT_LEN - ROT_ODD_BETA));
	SHIFT_L32(tmpYmm, cv[0], U8(ROT_ODD_BETA))
	SHIFT_R32(cv[0], cv[0], U8(WORD_BIT_LEN-ROT_ODD_BETA))
	OR(cv[0], cv[0], tmpYmm)
}

// static INLINE void xor_with_const(__m256i* cv_l, const __m256i* const_v){
func xor_with_const(cv_l []VecVirtual, const_v []VecVirtual) {
	Comment("xor_with_const")

	//*cv_l = XOR(*cv_l, *const_v);
	XOR(cv_l[0], cv_l[0], const_v[0])
}

// static INLINE void rotate_msg_gamma(__m256i* cv_r, const __m256i byte_perm_step){
func rotate_msg_gamma(cv_r []VecVirtual, byte_perm_step VecVirtual) {
	Comment("rotate_msg_gamma")

	//*cv_r = SHUFFLE8(*cv_r, byte_perm_step);
	SHUFFLE8(cv_r[0], cv_r[0], byte_perm_step)
}

// static INLINE void word_perm(__m256i* cv_l, __m256i* cv_r){
func word_perm(cv_l, cv_r []VecVirtual) {
	Comment("word_perm")

	//__m256i temp;
	temp := YMM()
	//temp = _mm256_shuffle_epi32(*cv_l, 0xd2);
	F_mm256_shuffle_epi32(temp, cv_l[0], U8(0xd2))
	//*cv_r = _mm256_shuffle_epi32(*cv_r, 0x6c);
	F_mm256_shuffle_epi32(cv_r[0], cv_r[0], U8(0x6c))
	//*cv_l = _mm256_permute2x128_si256(temp, *cv_r, 0x31);
	F_mm256_permute2x128_si256(cv_l[0], temp, cv_r[0], U8(0x31))
	//*cv_r = _mm256_permute2x128_si256(temp, *cv_r, 0x20);
	F_mm256_permute2x128_si256(cv_r[0], temp, cv_r[0], U8(0x20))
}

///* -------------------------------------------------------- *
//*  step function
//*  -------------------------------------------------------- */

// static INLINE void mix_even(__m256i* cv_l, __m256i* cv_r, const __m256i* const_v, const __m256i byte_perm_step){
func mix_even(cv_l, cv_r []VecVirtual, const_v []VecVirtual, byte_perm_step VecVirtual) {
	Comment("mix_even")

	add_blk(cv_l, cv_r)
	rotate_blk_even_alpha(cv_l)
	xor_with_const(cv_l, const_v)
	add_blk(cv_r, cv_l)
	rotate_blk_even_beta(cv_r)
	add_blk(cv_l, cv_r)
	rotate_msg_gamma(cv_r, byte_perm_step)
}

// static INLINE void mix_odd(__m256i* cv_l, __m256i* cv_r, const __m256i* const_v, const __m256i byte_perm_step){
func mix_odd(cv_l, cv_r []VecVirtual, const_v []VecVirtual, byte_perm_step VecVirtual) {
	Comment("mix_odd")

	add_blk(cv_l, cv_r)
	rotate_blk_odd_alpha(cv_l)
	xor_with_const(cv_l, const_v)
	add_blk(cv_r, cv_l)
	rotate_blk_odd_beta(cv_r)
	add_blk(cv_l, cv_r)
	rotate_msg_gamma(cv_r, byte_perm_step)
}

///* -------------------------------------------------------- *
//*  compression function
//*  -------------------------------------------------------- */

// static INLINE void compress(__m256i* cv_l, __m256i* cv_r, const lsh_u32 pdMsgBlk[MSG_BLK_WORD_LEN])
func compress(cv_l, cv_r []VecVirtual, pdMsgBlk Mem) {
	Comment("compress")

	//__m256i const_v[1];				// step function constant
	const_v := []VecVirtual{YMM()}
	//__m256i byte_perm_step;		// byte permutation info
	byte_perm_step := YMM()
	//__m256i word_perm_step;	// msg_word permutation info
	word_perm_step := YMM()
	//LSH256AVX2_internal i_state[1];
	i_state := LSH256AVX2_internal{
		submsg_e_l: []VecVirtual{YMM()},
		submsg_e_r: []VecVirtual{YMM()},
		submsg_o_l: []VecVirtual{YMM()},
		submsg_o_r: []VecVirtual{YMM()},
	}
	//int i;

	//byte_perm_step = LOAD(g_BytePermInfo);
	LOAD(byte_perm_step, g_BytePermInfo)
	//word_perm_step = LOAD(g_MsgWordPermInfo);
	LOAD(word_perm_step, g_MsgWordPermInfo)

	load_msg_blk(i_state, pdMsgBlk)

	msg_add_even(cv_l, cv_r, i_state)
	load_sc(const_v, 0)
	mix_even(cv_l, cv_r, const_v, byte_perm_step)
	word_perm(cv_l, cv_r)

	msg_add_odd(cv_l, cv_r, i_state)
	load_sc(const_v, 8)
	mix_odd(cv_l, cv_r, const_v, byte_perm_step)
	word_perm(cv_l, cv_r)

	//for (i = 1; i < NUM_STEPS / 2; i++){
	for i := 1; i < NUM_STEPS/2; i++ {
		msg_exp_even(i_state, word_perm_step)
		msg_add_even(cv_l, cv_r, i_state)
		load_sc(const_v, 16*i)
		mix_even(cv_l, cv_r, const_v, byte_perm_step)
		word_perm(cv_l, cv_r)

		msg_exp_odd(i_state, word_perm_step)
		msg_add_odd(cv_l, cv_r, i_state)
		load_sc(const_v, 16*i+8)
		mix_odd(cv_l, cv_r, const_v, byte_perm_step)
		word_perm(cv_l, cv_r)
	}

	msg_exp_even(i_state, word_perm_step)
	msg_add_even(cv_l, cv_r, i_state)
	//}
}

///* -------------------------------------------------------- */

// static INLINE void init224(LSH256AVX2_Context * state)
func init224(state *LSH256_Context) {
	Comment("init224")

	//load_blk(state->cv_l, g_IV224);
	load_blk_mem2mem(state.Cv_l, G_IV224)
	//load_blk(state->cv_r, g_IV224 + 8);
	load_blk_mem2mem(state.Cv_r, G_IV224.Offset(8*4))
}

// static INLINE void init256(LSH256AVX2_Context * state)
func init256(state *LSH256_Context) {
	Comment("init256")

	//load_blk(state->cv_l, g_IV256);
	load_blk_mem2mem(state.Cv_l, G_IV256)
	//load_blk(state->cv_r, g_IV256 + 8);
	load_blk_mem2mem(state.Cv_r, G_IV256.Offset(8*4))
}

///* -------------------------------------------------------- */

// static INLINE void fin(__m256i* cv_l, const __m256i* cv_r)
func fin(cv_l, cv_r []VecVirtual) {
	Comment("fin")

	//cv_l[0] = XOR(cv_l[0], cv_r[0]);
	XOR(cv_l[0], cv_l[0], cv_r[0])
}

///* -------------------------------------------------------- */

// static INLINE void get_hash(__m256i* cv_l, lsh_u8 * pbHashVal, const lsh_type algtype)
func get_hash(cv_l []VecVirtual, pbHashVal Mem, algtype Op) {
	Comment("get_hash")

	//lsh_u8 hash_val[LSH256_HASH_VAL_MAX_BYTE_LEN] = { 0x0, };
	hash_val := pbHashVal
	//lsh_uint hash_val_byte_len = LSH_GET_HASHBYTE(algtype);
	hash_val_byte_len := GP32()
	LSH_GET_HASHBYTE(hash_val_byte_len, algtype)
	//lsh_uint hash_val_bit_len = LSH_GET_SMALL_HASHBIT(algtype);
	hash_val_bit_len := GP32()
	LSH_GET_SMALL_HASHBIT(hash_val_bit_len, algtype)

	//STORE(hash_val, cv_l[0]);
	STORE(hash_val, cv_l[0])
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

		MOVB(tmp8, pbHashVal.Offset(-1).Idx(hash_val_byte_len, 1))
	}
	Label("get_hash_if_end")
}

///* -------------------------------------------------------- */

// lsh_err lsh256_avx2_init(struct LSH256_Context * _ctx, const lsh_type algtype){
func Lsh256_avx2_init(ctx *LSH256_Context) {
	Comment("lsh256_avx2_init")

	//__m256i cv_l[1];s
	//__m256i cv_r[1];
	//__m256i const_v[1];
	//__m256i byte_perm_step;
	//LSH256AVX2_Context* ctx = (LSH256AVX2_Context*)_ctx;
	//lsh_uint i;

	//if (ctx == NULL){
	//	return LSH_ERR_NULL_PTR;
	//}

	//ctx->algtype = algtype;
	//ctx->remain_databitlen = 0;

	//if (!LSH_IS_LSH256(algtype)){
	//	return LSH_ERR_INVALID_ALGTYPE;
	//}

	//if (LSH_GET_HASHBYTE(algtype) > LSH256_HASH_VAL_MAX_BYTE_LEN || LSH_GET_HASHBYTE(algtype) == 0){
	//	return LSH_ERR_INVALID_ALGTYPE;
	//}

	//switch (algtype){
	//case LSH_TYPE_256_256:
	CMPL(ctx.Algtype, U32(LSH_TYPE_256_256))
	JNE(LabelRef("lsh256_avx2_init_if0_end"))
	//	init256(ctx);
	init256(ctx)
	//	return LSH_SUCCESS;
	JMP(LabelRef("lsh256_avx2_init_ret"))
	//case LSH_TYPE_256_224:
	Label("lsh256_avx2_init_if0_end")
	//	init224(ctx);
	init224(ctx)
	//	return LSH_SUCCESS;
	JMP(LabelRef("lsh256_avx2_init_ret"))
	//default:
	//	break;
	//}

	//*cv_l = _mm256_set_epi32(0, 0, 0, 0, 0, 0, LSH_GET_HASHBIT(algtype), LSH256_HASH_VAL_MAX_BYTE_LEN);
	//*cv_r = _mm256_setzero_si256();
	//byte_perm_step = LOAD(g_BytePermInfo);

	//for (i = 0; i < NUM_STEPS / 2; i++)
	//{
	//	//Mix
	//	load_sc(const_v, i * 16);
	//	mix_even(cv_l, cv_r, const_v, byte_perm_step);
	//	word_perm(cv_l, cv_r);
	//	load_sc(const_v, i * 16 + 8);
	//	mix_odd(cv_l, cv_r, const_v, byte_perm_step);
	//	word_perm(cv_l, cv_r);
	//}

	//ctx->cv_l[0] = cv_l[0];
	//ctx->cv_r[0] = cv_r[0];

	Label("lsh256_avx2_init_ret")
	//return LSH_SUCCESS;
}

// lsh_err lsh256_avx2_update(struct LSH256_Context * _ctx, const lsh_u8 * data, size_t databitlen){
func Lsh256_avx2_update(ctx *LSH256_Context, data Mem, databytelen Register) {
	Comment("lsh256_avx2_update")

	//__m256i cv_l[1];
	cv_l := []VecVirtual{YMM()}
	//__m256i cv_r[1];
	cv_r := []VecVirtual{YMM()}
	//size_t databytelen = databitlen >> 3;

	//lsh_uint pos2 = databitlen & 0x7;

	//LSH256AVX2_Context* ctx = (LSH256AVX2_Context*)_ctx;
	//lsh_uint remain_msg_byte;
	remain_msg_byte := GP64()
	//lsh_uint remain_msg_bit;

	//if (ctx == NULL || data == NULL){
	//	return LSH_ERR_NULL_PTR;
	//}
	//if (ctx->algtype == 0 || LSH_GET_HASHBYTE(ctx->algtype) > LSH256_HASH_VAL_MAX_BYTE_LEN){
	//	return LSH_ERR_INVALID_STATE;
	//}
	//if (databitlen == 0){
	//	return LSH_SUCCESS;
	//}

	//remain_msg_byte = ctx->remain_databitlen >> 3;
	MOVQ(ctx.Remain_databytelen, remain_msg_byte)

	//remain_msg_bit = ctx->remain_databitlen & 7;
	//if (remain_msg_byte >= LSH256_MSG_BLK_BYTE_LEN){
	//	return LSH_ERR_INVALID_STATE;
	//}
	//if (remain_msg_bit > 0){
	//	return LSH_ERR_INVALID_DATABITLEN;
	//}

	//if (databytelen + remain_msg_byte < LSH256_MSG_BLK_BYTE_LEN){
	tmp32 := GP64()
	MOVQ(databytelen, tmp32)
	ADDQ(remain_msg_byte, tmp32)
	CMPQ(tmp32, U32(LSH256_MSG_BLK_BYTE_LEN))
	JGE(LabelRef("lsh256_avx2_update_if0_end"))
	{
		//memcpy(ctx->last_block + remain_msg_byte, data, databytelen);
		Memcpy(ctx.Last_block.Idx(remain_msg_byte, 1), data, databytelen, true)
		//ctx->remain_databitlen += (lsh_uint)databitlen;
		ADDQ(databytelen, ctx.Remain_databytelen)
		//remain_msg_byte += (lsh_uint)databytelen;
		//if (pos2){
		//	ctx->last_block[remain_msg_byte] = data[databytelen] & ((0xff >> pos2) ^ 0xff);
		//}

		//return LSH_SUCCESS;
		JMP(LabelRef("lsh256_avx2_update_ret"))
	}
	Label("lsh256_avx2_update_if0_end")

	//load_blk(cv_l, ctx->cv_l);
	load_blk_mem2vec(cv_l, ctx.Cv_l)
	//load_blk(cv_r, ctx->cv_r);
	load_blk_mem2vec(cv_r, ctx.Cv_r)

	//if (remain_msg_byte > 0){
	CMPQ(remain_msg_byte, U32(0))
	JE(LabelRef("lsh256_avx2_update_if2_end"))
	{
		//lsh_uint more_byte = LSH256_MSG_BLK_BYTE_LEN - remain_msg_byte;
		more_byte := GP64()
		MOVQ(U32(LSH256_MSG_BLK_BYTE_LEN), more_byte)
		SUBQ(remain_msg_byte, more_byte)
		//memcpy(ctx->last_block + remain_msg_byte, data, more_byte);
		Memcpy(ctx.Last_block.Idx(remain_msg_byte, 1), data, more_byte, true)
		//compress(cv_l, cv_r, (lsh_u32*)ctx->last_block);
		compress(cv_l, cv_r, ctx.Last_block)
		//data += more_byte;
		ADDQ(more_byte.As64(), data.Base)
		//databytelen -= more_byte;
		SUBQ(more_byte, databytelen)
		//remain_msg_byte = 0;
		MOVQ(U32(0), remain_msg_byte)
		//ctx->remain_databitlen = 0;
		MOVQ(U32(0), ctx.Remain_databytelen)
	}
	Label("lsh256_avx2_update_if2_end")

	//while (databytelen >= LSH256_MSG_BLK_BYTE_LEN)
	Label("lsh256_avx2_update_while_start")
	CMPQ(databytelen, U32(LSH256_MSG_BLK_BYTE_LEN))
	JL(LabelRef("lsh256_avx2_update_while_end"))
	{
		//compress(cv_l, cv_r, (lsh_u32*)data);
		compress(cv_l, cv_r, data)
		//data += LSH256_MSG_BLK_BYTE_LEN;
		ADDQ(U32(LSH256_MSG_BLK_BYTE_LEN), data.Base)
		//databytelen -= LSH256_MSG_BLK_BYTE_LEN;
		SUBQ(U32(LSH256_MSG_BLK_BYTE_LEN), databytelen)

		JMP(LabelRef("lsh256_avx2_update_while_start"))
		//	}
	}
	Label("lsh256_avx2_update_while_end")

	//store_blk(ctx->cv_l, cv_l);
	store_blk(ctx.Cv_l, cv_l)
	//store_blk(ctx->cv_r, cv_r);
	store_blk(ctx.Cv_r, cv_r)

	//if (databytelen > 0){
	CMPQ(remain_msg_byte, U32(0))
	JE(LabelRef("lsh256_avx2_update_if3_end"))
	{
		//memcpy(ctx->last_block, data, databytelen);
		Memcpy(ctx.Last_block, data, databytelen, true)
		//ctx->remain_databitlen = (lsh_uint)(databytelen << 3);
		MOVQ(databytelen, ctx.Remain_databytelen)
		//}
	}
	Label("lsh256_avx2_update_if3_end")

	//if (pos2){
	//	ctx->last_block[databytelen] = data[databytelen] & ((0xff >> pos2) ^ 0xff);
	//	ctx->remain_databitlen += pos2;
	//}

	//return LSH_SUCCESS;
	Label("lsh256_avx2_update_ret")
}

// lsh_err lsh256_avx2_final(struct LSH256_Context * _ctx, lsh_u8 * hashval){
func Lsh256_avx2_final(ctx *LSH256_Context, hashval Mem) {
	Comment("lsh256_avx2_final")

	//__m256i cv_l[1];
	cv_l := []VecVirtual{YMM()}
	//__m256i cv_r[1];
	cv_r := []VecVirtual{YMM()}
	//LSH256AVX2_Context* ctx = (LSH256AVX2_Context*)_ctx;
	//lsh_uint remain_msg_byte;
	remain_msg_byte := GP64()
	//lsh_uint remain_msg_bit;

	//if (ctx == NULL || hashval == NULL){
	//	return LSH_ERR_NULL_PTR;
	//}
	//if (ctx->algtype == 0 || LSH_GET_HASHBYTE(ctx->algtype) > LSH256_HASH_VAL_MAX_BYTE_LEN){
	//	return LSH_ERR_INVALID_STATE;
	//}

	//remain_msg_byte = ctx->remain_databitlen >> 3;
	MOVQ(ctx.Remain_databytelen, remain_msg_byte)
	//remain_msg_bit = ctx->remain_databitlen & 7;

	//if (remain_msg_byte >= LSH256_MSG_BLK_BYTE_LEN){
	//	return LSH_ERR_INVALID_STATE;
	//}

	//if (remain_msg_bit){
	//	ctx->last_block[remain_msg_byte] |= (0x1 << (7 - remain_msg_bit));
	//}
	//else{
	//	ctx->last_block[remain_msg_byte] = 0x80;
	MOVB(U8(0x80), ctx.Last_block.Idx(remain_msg_byte, 1))
	//}
	//memset(ctx->last_block + remain_msg_byte + 1, 0, LSH256_MSG_BLK_BYTE_LEN - remain_msg_byte - 1);
	arg2 := GP64()
	MOVQ(U64(LSH256_MSG_BLK_BYTE_LEN-1), arg2)
	SUBQ(remain_msg_byte, arg2)
	Memset(ctx.Last_block.Offset(1).Idx(remain_msg_byte, 1), 0, arg2, false)

	//load_blk(cv_l, ctx->cv_l);
	load_blk_mem2vec(cv_l, ctx.Cv_l)
	//load_blk(cv_r, ctx->cv_r);
	load_blk_mem2vec(cv_r, ctx.Cv_r)

	//compress(cv_l, cv_r, (lsh_u32*)ctx->last_block);
	compress(cv_l, cv_r, ctx.Last_block)

	//fin(cv_l, cv_r);
	fin(cv_l, cv_r)
	//get_hash(cv_l, hashval, ctx->algtype);
	get_hash(cv_l, hashval, ctx.Algtype)

	//memset(ctx, 0, sizeof(struct LSH256_Context));

	//return LSH_SUCCESS;
}
