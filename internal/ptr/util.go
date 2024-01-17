//go:build !go1.19
// +build !go1.19

package ptr

import (
	"reflect"
	"unsafe"
)

func ByteToUint32Array(b []byte) []uint32 {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	header.Len /= 4
	header.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(&header))
}

func BytePtr(b []byte) *byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	return (*byte)(unsafe.Pointer(header.Data))
}
