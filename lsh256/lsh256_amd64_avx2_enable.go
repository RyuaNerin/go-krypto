//go:build krypto_lsh256_avx2

package lsh256

func init() {
	useAVX2 = true
}
