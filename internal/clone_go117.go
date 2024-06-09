//go:build !go1.18
// +build !go1.18

package internal

import (
	"unsafe"
)

func StringClone(s string) string {
	if len(s) == 0 {
		return ""
	}
	b := make([]byte, len(s))
	copy(b, s)

	return *(*string)(unsafe.Pointer(&b))
}
