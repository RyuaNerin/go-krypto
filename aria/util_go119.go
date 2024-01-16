//go:build go1.19
// +build go1.19

package aria

import (
	"unsafe"
)

func toUint32Array(b []byte) []uint32 {
	sd := unsafe.SliceData(b)
	return unsafe.Slice((*uint32)(unsafe.Pointer(sd)), len(b)/4)
}
