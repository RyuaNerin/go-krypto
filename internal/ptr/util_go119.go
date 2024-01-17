//go:build go1.19
// +build go1.19

package ptr

import (
	"unsafe"
)

func ByteToUint32Array(b []byte) []uint32 {
	sd := unsafe.SliceData(b)
	return unsafe.Slice((*uint32)(unsafe.Pointer(sd)), len(b)/4)
}

func BytePtr(b []byte) *byte {
	sd := unsafe.SliceData(b)
	return sd
}
