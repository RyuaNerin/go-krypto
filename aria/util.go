//go:build !go1.19
// +build !go1.19

package aria

import (
	"reflect"
	"unsafe"
)

func toUint32Array(b []byte) []uint32 {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	header.Len /= 4
	header.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(&header))
}

func toBytePtr(b []byte) *byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	return (*byte)(unsafe.Pointer(header.Data))
}
