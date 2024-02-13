//go:build go1.20
// +build go1.20

package ptr

import (
	"unsafe"
)

func PUint32(b []byte) []uint32 {
	sd := unsafe.SliceData(b)
	return unsafe.Slice((*uint32)(unsafe.Pointer(sd)), len(b)/4)
}

func PByte(b []byte) *byte {
	sd := unsafe.SliceData(b)
	return sd
}
