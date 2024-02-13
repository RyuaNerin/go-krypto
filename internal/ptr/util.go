//go:build !go1.20
// +build !go1.20

package ptr

import (
	"reflect"
	"unsafe"
)

func PUint32(b []byte) []uint32 {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	header.Len /= 4
	header.Cap /= 4
	return *(*[]uint32)(unsafe.Pointer(&header))
}

func PByte(b []byte) *byte {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	return (*byte)(unsafe.Pointer(header.Data))
}
