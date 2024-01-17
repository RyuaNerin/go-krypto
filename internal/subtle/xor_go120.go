//go:build go1.20
// +build go1.20

package subtle

import (
	go_subtle "crypto/subtle"
)

func XORBytes(dst, x, y []byte) int {
	return go_subtle.XORBytes(dst, x, y)
}
