package lsh512avx2

import (
	. "github.com/mmcloughlin/avo/build"
	. "github.com/mmcloughlin/avo/operand"
	. "github.com/mmcloughlin/avo/reg"

	. "kryptosimd/avoutil"
	/*
		. "kryptosimd/avoutil/simd"
		. "kryptosimd/lsh/x86/sse2"
		. "kryptosimd/lsh512/avo/x86/lsh512avoconst"
		. "kryptosimd/lsh512/avo/x86/lsh512common"
	*/)

func UnitTest() {
	test_add_blk()
	test_mix_even()
	test_mix_odd()
	test_msg_add_even()
	test_msg_add_odd()
	test_msg_exp_even()
	test_msg_exp_odd()
	test_rotate_blk_even_alpha()
	test_rotate_blk_even_beta()
	test_rotate_blk_odd_alpha()
	test_rotate_blk_odd_beta()
	test_rotate_msg_gamma()
	test_word_perm()
	test_xor_with_const()
}

func loadState() LSH512AVX2_internal {
	i_state := LSH512AVX2_internal{
		submsg_e_l: []VecVirtual{YMM(), YMM()},
		submsg_e_r: []VecVirtual{YMM(), YMM()},
		submsg_o_l: []VecVirtual{YMM(), YMM()},
		submsg_o_r: []VecVirtual{YMM(), YMM()},
	}

	load_blk_mem2vec(i_state.submsg_e_l, Mem{Base: Load(Param("el").Base(), GP64())})
	load_blk_mem2vec(i_state.submsg_e_r, Mem{Base: Load(Param("er").Base(), GP64())})
	load_blk_mem2vec(i_state.submsg_o_l, Mem{Base: Load(Param("ol").Base(), GP64())})
	load_blk_mem2vec(i_state.submsg_o_r, Mem{Base: Load(Param("or").Base(), GP64())})

	return i_state
}
func saveState(i_state LSH512AVX2_internal) {
	store_blk(Mem{Base: Load(Param("el").Base(), GP64())}, i_state.submsg_e_l)
	store_blk(Mem{Base: Load(Param("er").Base(), GP64())}, i_state.submsg_e_r)
	store_blk(Mem{Base: Load(Param("ol").Base(), GP64())}, i_state.submsg_o_l)
	store_blk(Mem{Base: Load(Param("or").Base(), GP64())}, i_state.submsg_o_r)
}

func test_msg_exp_even() {
	TEXT("msg_exp_even", NOSPLIT, "func(el, er, ol, or []uint64)")

	i_state := loadState()

	msg_exp_even(i_state)

	saveState(i_state)

	RET()
}

func test_msg_exp_odd() {
	TEXT("msg_exp_odd", NOSPLIT, "func(el, er, ol, or []uint64)")

	i_state := loadState()

	msg_exp_odd(i_state)

	saveState(i_state)

	RET()
}

func test_msg_add_even() {
	TEXT("msg_add_even", NOSPLIT, "func(cv_l, cv_r []uint64, el, er, ol, or []uint64)")

	cv_l_mem := Mem{Base: Load(Param("cv_l").Base(), GP64())}
	cv_r_mem := Mem{Base: Load(Param("cv_r").Base(), GP64())}

	cv_r := []VecVirtual{YMM(), YMM()}
	cv_l := []VecVirtual{YMM(), YMM()}

	i_state := loadState()

	load_blk_mem2vec(cv_l, cv_l_mem)
	load_blk_mem2vec(cv_r, cv_r_mem)

	msg_add_even(cv_l, cv_r, i_state)

	store_blk(cv_l_mem, cv_l)
	store_blk(cv_r_mem, cv_r)

	saveState(i_state)

	RET()
}

func test_msg_add_odd() {
	TEXT("msg_add_odd", NOSPLIT, "func(cv_l, cv_r []uint64, el, er, ol, or []uint64)")

	cv_l_mem := Mem{Base: Load(Param("cv_l").Base(), GP64())}
	cv_r_mem := Mem{Base: Load(Param("cv_r").Base(), GP64())}

	cv_r := []VecVirtual{YMM(), YMM()}
	cv_l := []VecVirtual{YMM(), YMM()}

	i_state := loadState()

	load_blk_mem2vec(cv_l, cv_l_mem)
	load_blk_mem2vec(cv_r, cv_r_mem)

	msg_add_odd(cv_l, cv_r, i_state)

	store_blk(cv_l_mem, cv_l)
	store_blk(cv_r_mem, cv_r)

	saveState(i_state)

	RET()
}

func test_add_blk() {
	TEXT("add_blk", NOSPLIT, "func(cv_l, cv_r []uint64)")

	cv_l_mem := Mem{Base: Load(Param("cv_l").Base(), GP64())}
	cv_r_mem := Mem{Base: Load(Param("cv_r").Base(), GP64())}

	cv_l := []VecVirtual{YMM(), YMM()}
	cv_r := []VecVirtual{YMM(), YMM()}

	load_blk_mem2vec(cv_l, cv_l_mem)
	load_blk_mem2vec(cv_r, cv_r_mem)

	add_blk(cv_l, cv_r)

	store_blk(cv_l_mem, cv_l)
	store_blk(cv_r_mem, cv_r)

	RET()
}

func test_rotate_blk_even_alpha() {
	TEXT("rotate_blk_even_alpha", NOSPLIT, "func(cv []uint64)")

	cv_mem := Mem{Base: Load(Param("cv").Base(), GP64())}

	cv := []VecVirtual{YMM(), YMM()}

	load_blk_mem2vec(cv, cv_mem)

	rotate_blk_even_alpha(cv)

	store_blk(cv_mem, cv)

	RET()
}

func test_rotate_blk_even_beta() {
	TEXT("rotate_blk_even_beta", NOSPLIT, "func(cv []uint64)")

	cv_mem := Mem{Base: Load(Param("cv").Base(), GP64())}

	cv := []VecVirtual{YMM(), YMM()}

	load_blk_mem2vec(cv, cv_mem)

	rotate_blk_even_beta(cv)

	store_blk(cv_mem, cv)

	RET()
}

func test_rotate_blk_odd_alpha() {
	TEXT("rotate_blk_odd_alpha", NOSPLIT, "func(cv []uint64)")

	cv_mem := Mem{Base: Load(Param("cv").Base(), GP64())}

	cv := []VecVirtual{YMM(), YMM()}

	load_blk_mem2vec(cv, cv_mem)

	rotate_blk_odd_alpha(cv)

	store_blk(cv_mem, cv)

	RET()
}

func test_rotate_blk_odd_beta() {
	TEXT("rotate_blk_odd_beta", NOSPLIT, "func(cv []uint64)")

	cv_mem := Mem{Base: Load(Param("cv").Base(), GP64())}

	cv := []VecVirtual{YMM(), YMM()}

	load_blk_mem2vec(cv, cv_mem)

	rotate_blk_odd_beta(cv)

	store_blk(cv_mem, cv)

	RET()
}

func test_xor_with_const() {
	TEXT("xor_with_const", NOSPLIT, "func(cv_l, const_v []uint64)")

	cv_l_mem := Mem{Base: Load(Param("cv_l").Base(), GP64())}
	const_v_mem := Mem{Base: Load(Param("const_v").Base(), GP64())}

	cv_l := []VecVirtual{YMM(), YMM()}
	const_v := []VecVirtual{YMM(), YMM()}

	load_blk_mem2vec(cv_l, cv_l_mem)
	load_blk_mem2vec(const_v, const_v_mem)

	xor_with_const(cv_l, const_v)

	store_blk(cv_l_mem, cv_l)
	store_blk(const_v_mem, const_v)

	RET()
}

func test_rotate_msg_gamma() {
	TEXT("rotate_msg_gamma", NOSPLIT, "func(cv_r []uint64)")

	cv_r_mem := Mem{Base: Load(Param("cv_r").Base(), GP64())}

	cv_r := []VecVirtual{YMM(), YMM()}

	perm_step := []Mem{
		g_BytePermInfo.Offset(YmmSize * 0),
		g_BytePermInfo.Offset(YmmSize * 1),
	}

	load_blk_mem2vec(cv_r, cv_r_mem)

	rotate_msg_gamma(cv_r, perm_step)

	store_blk(cv_r_mem, cv_r)

	RET()
}

func test_word_perm() {
	TEXT("word_perm", NOSPLIT, "func(cv_l, cv_r []uint64)")

	cv_l_mem := Mem{Base: Load(Param("cv_l").Base(), GP64())}
	cv_r_mem := Mem{Base: Load(Param("cv_r").Base(), GP64())}

	cv_l := []VecVirtual{YMM(), YMM()}
	cv_r := []VecVirtual{YMM(), YMM()}

	load_blk_mem2vec(cv_l, cv_l_mem)
	load_blk_mem2vec(cv_r, cv_r_mem)

	word_perm(cv_l, cv_r)

	store_blk(cv_l_mem, cv_l)
	store_blk(cv_r_mem, cv_r)

	RET()
}

func test_mix_even() {
	TEXT("mix_even", NOSPLIT, "func(cv_l, cv_r, const_v []uint64)")

	cv_l_mem := Mem{Base: Load(Param("cv_l").Base(), GP64())}
	cv_r_mem := Mem{Base: Load(Param("cv_r").Base(), GP64())}
	const_v_mem := Mem{Base: Load(Param("const_v").Base(), GP64())}

	cv_l := []VecVirtual{YMM(), YMM()}
	cv_r := []VecVirtual{YMM(), YMM()}
	const_v := []VecVirtual{YMM(), YMM()}

	perm_step := []Mem{
		g_BytePermInfo.Offset(YmmSize * 0),
		g_BytePermInfo.Offset(YmmSize * 1),
	}

	load_blk_mem2vec(cv_l, cv_l_mem)
	load_blk_mem2vec(cv_r, cv_r_mem)
	load_blk_mem2vec(const_v, const_v_mem)

	mix_even(cv_l, cv_r, const_v, perm_step)

	store_blk(cv_l_mem, cv_l)
	store_blk(cv_r_mem, cv_r)
	store_blk(const_v_mem, const_v)

	RET()
}

func test_mix_odd() {
	TEXT("mix_odd", NOSPLIT, "func(cv_l, cv_r, const_v []uint64)")

	cv_l_mem := Mem{Base: Load(Param("cv_l").Base(), GP64())}
	cv_r_mem := Mem{Base: Load(Param("cv_r").Base(), GP64())}
	const_v_mem := Mem{Base: Load(Param("const_v").Base(), GP64())}

	cv_l := []VecVirtual{YMM(), YMM()}
	cv_r := []VecVirtual{YMM(), YMM()}
	const_v := []VecVirtual{YMM(), YMM()}

	perm_step := []Mem{
		g_BytePermInfo.Offset(YmmSize * 0),
		g_BytePermInfo.Offset(YmmSize * 1),
	}

	load_blk_mem2vec(cv_l, cv_l_mem)
	load_blk_mem2vec(cv_r, cv_r_mem)
	load_blk_mem2vec(const_v, const_v_mem)

	mix_odd(cv_l, cv_r, const_v, perm_step)

	store_blk(cv_l_mem, cv_l)
	store_blk(cv_r_mem, cv_r)
	store_blk(const_v_mem, const_v)

	RET()
}
