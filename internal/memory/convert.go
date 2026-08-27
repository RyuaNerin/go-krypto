package memory

import (
	"unsafe"
)

func Cast[T any](b []byte) []T {
	var zero T
	sd := unsafe.SliceData(b)
	return unsafe.Slice((*T)(unsafe.Pointer(sd)), uintptr(len(b))/unsafe.Sizeof(zero))
}

func CastPtr[TTo any](b []byte) *TTo {
	sd := unsafe.SliceData(b)
	return (*TTo)(unsafe.Pointer(sd))
}

func Pointer[S ~[]E, E any](b S) unsafe.Pointer {
	return unsafe.Pointer(unsafe.SliceData(b))
}
