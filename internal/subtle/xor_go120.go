//go:build go1.20
// +build go1.20

package subtle

import (
	go_subtle "crypto/subtle"
)

// XORBytes sets dst[i] = x[i] ^ y[i] for all i < n = min(len(x), len(y)),
// returning n, the number of bytes written to dst.
// If dst does not have length at least n,
// XORBytes panics without writing anything to dst.
func XORBytes(dst, x, y []byte) int {
	return go_subtle.XORBytes(dst, x, y)
}
