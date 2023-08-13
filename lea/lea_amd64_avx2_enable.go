//go:build krypto_lea_avx2

package lea

func init() {
	useAVX2 = true
}
