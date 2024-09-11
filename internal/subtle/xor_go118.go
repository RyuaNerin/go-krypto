// https://github.com/golang/go/blob/release-branch.go1.21/src/crypto/subtle/xor_arm64.go

//go:build gc && !go1.20 && !purego && go1.18 && (amd64 || arm64 || ppc64 || ppc64le)
// +build gc
// +build !go1.20
// +build !purego
// +build go1.18
// +build amd64 arm64 ppc64 ppc64le

package subtle

// XORBytes sets dst[i] = x[i] ^ y[i] for all i < n = min(len(x), len(y)),
// returning n, the number of bytes written to dst.
// If dst does not have length at least n,
// XORBytes panics without writing anything to dst.
func XORBytes(dst, x, y []byte) int {
	n := len(x)
	if len(y) < n {
		n = len(y)
	}
	if n == 0 {
		return 0
	}
	if n > len(dst) {
		panic("subtle.XORBytes: dst too short")
	}
	xorBytes(&dst[0], &x[0], &y[0], n) // arch-specific
	return n
}
