package memory

import (
	"unsafe"
)

func IsAlignedPtr(p unsafe.Pointer, align uintptr) bool {
	return uintptr(p)%align == 0
}

func IsAligned[T any](ptr *T, align uintptr) bool {
	return uintptr(unsafe.Pointer(ptr))%align == 0
}

type Aligned16[T any] struct {
	data      T
	alignment [15]byte // alignment
}

func (a *Aligned16[T]) Zeroize() {
	var zero T
	a.data = zero
	clear(a.alignment[:])
}

func (a *Aligned16[T]) Aligned() *T {
	base := unsafe.Pointer(&a.data)
	return (*T)(unsafe.Add(base, (-uintptr(base))&15))
}
