package memory

import (
	"unsafe"

	"github.com/RyuaNerin/go-krypto/internal/zeroize"
)

// ZeroizeUnlessSubsliceOf zeroizes target only if it is not fully contained
// within the memory region of src.
// If target is a subslice of src, zeroizing src covers it as well, so the
// redundant wipe is skipped.
func ZeroizeUnlessSubsliceOf[T any](target, src []T) {
	if len(target) == 0 {
		return
	}

	var zero T
	size := unsafe.Sizeof(zero)

	tStart := uintptr(unsafe.Pointer(unsafe.SliceData(target)))
	tEnd := tStart + uintptr(len(target))*size

	sStart := uintptr(unsafe.Pointer(unsafe.SliceData(src)))
	sEnd := sStart + uintptr(len(src))*size

	// Skip if [tStart, tEnd) is contained in [sStart, sEnd).
	if len(src) > 0 && tStart >= sStart && tEnd <= sEnd {
		return
	}

	zeroize.Zeroize(target)
}
