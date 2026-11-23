package unaligned

import (
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func TestCopyPreservesNaturalAlignment(t *testing.T) {
	source := new(atomic.Uint64)
	source.Store(42)
	const alignment = uintptr(32)
	naturalAlignment := unsafe.Alignof(*source)

	for residue := uintptr(0); residue < alignment; residue++ {
		if residue%naturalAlignment != 0 {
			require.PanicsWithValue(t, msgNaturalAlignmentRequired, func() {
				Copy(source, alignment, residue)
			})
			continue
		}

		context := Copy(source, alignment, residue)
		require.Equal(t, residue, uintptr(unsafe.Pointer(context.Value))%alignment)
		require.Equal(t, uint64(42), context.Value.Load())
		context.Value.Store(99)
		require.Equal(t, uint64(42), source.Load(), "copy aliases source storage")
		context.KeepAlive()
	}

	require.PanicsWithValue(t, msgNaturalAlignmentRequired, func() {
		Copy(source, naturalAlignment+1, 0)
	})
}

func TestNaturallyAlignedOffsets(t *testing.T) {
	source := new(atomic.Uint64)
	Test(t, 32, source, func(residue uintptr) {
		require.NotZero(t, residue%32, "context must remain unaligned to the SIMD boundary")
		require.Zero(t, residue%unsafe.Alignof(*source))
		context := Copy(source, 32, residue)
		context.Value.Store(42)
		require.Equal(t, uint64(42), context.Value.Load())
		context.KeepAlive()
	})
}
