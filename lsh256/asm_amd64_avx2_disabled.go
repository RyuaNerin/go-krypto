//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh256 && amd64 && !krypto_enable_lsh256_avx2

package lsh256

const enableAVX2 = false
