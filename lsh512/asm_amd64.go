//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh512 && amd64

package lsh512

import (
	"unsafe"

	"golang.org/x/sys/cpu"
)

var (
	hasSSSE3  = cpu.X86.HasSSE2 && cpu.X86.HasSSSE3
	hasAVX2   = cpu.X86.HasAVX2 && cpu.X86.HasAVX && hasSSSE3
	hasAVX512 = cpu.X86.HasAVX512F && cpu.X86.HasAVX512VL && hasAVX2
)

func init() {
	switch {
	case hasAVX512:
		defaultSimd = simdAVX512{}

	case hasAVX2:
		defaultSimd = simdAVX2{}

	case hasSSSE3:
		defaultSimd = simdSSSE3{}

	default:
		defaultSimd = simdSSSE2{}
	}
}

type (
	simdSSSE2  struct{}
	simdSSSE3  struct{}
	simdAVX2   struct{}
	simdAVX512 struct{}
)

var (
	_ simdSet = (*simdSSSE2)(nil)
	_ simdSet = (*simdSSSE3)(nil)
	_ simdSet = (*simdAVX2)(nil)
	_ simdSet = (*simdAVX512)(nil)
)

func (s simdSSSE2) update(c, d unsafe.Pointer, n int64) { __lsh512_sse2_update(c, d, n) }
func (s simdSSSE2) final(c *lsh512ContextData) (sum [Size]byte) {
	__lsh512_sse2_final(unsafe.Pointer(c), unsafe.Pointer(&sum[0]))
	return sum
}

func (s simdSSSE3) update(c, d unsafe.Pointer, n int64) { __lsh512_ssse3_update(c, d, n) }
func (s simdSSSE3) final(c *lsh512ContextData) (sum [Size]byte) {
	__lsh512_sse2_final(unsafe.Pointer(c), unsafe.Pointer(&sum[0]))
	return sum
}

func (s simdAVX2) update(c, d unsafe.Pointer, n int64) { __lsh512_avx2_update(c, d, n) }
func (s simdAVX2) final(c *lsh512ContextData) (sum [Size]byte) {
	__lsh512_sse2_final(unsafe.Pointer(c), unsafe.Pointer(&sum[0]))
	return sum
}

func (s simdAVX512) update(c, d unsafe.Pointer, n int64) { __lsh512_avx512_update(c, d, n) }
func (s simdAVX512) final(c *lsh512ContextData) (sum [Size]byte) {
	__lsh512_avx512_final(unsafe.Pointer(c), unsafe.Pointer(&sum[0]))
	return sum
}

// constants for SIMD implementation

var (
	bytePermInfoSSSE3 = [...]uint64{
		0x0706050403020100, 0x0d0c0b0a09080f0e, 0x0302010007060504, 0x09080f0e0d0c0b0a,
		0x0605040302010007, 0x0c0b0a09080f0e0d, 0x0201000706050403, 0x080f0e0d0c0b0a09,
	}
	bytePermInfoAVX2 = [...]uint64{
		0x0706050403020100, 0x0d0c0b0a09080f0e, 0x1312111017161514, 0x19181f1e1d1c1b1a,
		0x0605040302010007, 0x0c0b0a09080f0e0d, 0x1211101716151413, 0x181f1e1d1c1b1a19,
	}
	msgWordPermInfo = [...]uint64{
		// ssse3, avx2
		0x0706050403020100, 0x0f0e0d0c0b0a0908, 0x1716151413121110, 0x1f1e1d1c1b1a1918,
	}
)

func init() {
	_ = bytePermInfoSSSE3
	_ = bytePermInfoAVX2
	_ = msgWordPermInfo
}
