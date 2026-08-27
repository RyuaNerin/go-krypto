//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh256 && amd64

package lsh256

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

	case hasAVX2 && enableAVX2: // avx2 is slower than ssse3
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

func (s simdSSSE2) update(c, d unsafe.Pointer, n int64) { __lsh256_sse2_update(c, d, n) }
func (s simdSSSE2) final(c *lsh256ContextData) (sum [Size]byte) {
	__lsh256_sse2_final(unsafe.Pointer(c), unsafe.Pointer(&sum[0]))
	return sum
}

func (s simdSSSE3) update(c, d unsafe.Pointer, n int64) { __lsh256_ssse3_update(c, d, n) }
func (s simdSSSE3) final(c *lsh256ContextData) (sum [Size]byte) {
	__lsh256_ssse3_final(unsafe.Pointer(c), unsafe.Pointer(&sum[0]))
	return sum
}

func (s simdAVX2) update(c, d unsafe.Pointer, n int64) { __lsh256_avx2_update(c, d, n) }
func (s simdAVX2) final(c *lsh256ContextData) (sum [Size]byte) {
	__lsh256_avx2_final(unsafe.Pointer(c), unsafe.Pointer(&sum[0]))
	return sum
}

func (s simdAVX512) update(c, d unsafe.Pointer, n int64) { __lsh256_avx512_update(c, d, n) }
func (s simdAVX512) final(c *lsh256ContextData) (sum [Size]byte) {
	__lsh256_avx512_final(unsafe.Pointer(c), unsafe.Pointer(&sum[0]))
	return sum
}

// constants for SIMD implementation

var (
	bytePermInfo = [...]uint32{
		0x03020100, 0x06050407, 0x09080b0a, 0x0c0f0e0d, // ssse3,
		0x10131211, 0x15141716, 0x1a19181b, 0x1f1e1d1c, // avx2
	}
	bytePermInfoR = [...]uint32{
		0x00030201, 0x05040706, 0x0a09080b, 0x0f0e0d0c, // ssse3,
	}
	msgWordPermInfo = [...]uint32{
		// ssse3, avx2
		0x0f0e0d0c, 0x0b0a0908, 0x03020100, 0x07060504, 0x1f1e1d1c, 0x13121110, 0x17161514, 0x1b1a1918,
	}

	// sse2
	gammaMask0 = [...]uint32{0x00000000, 0xffffffff, 0xffffffff, 0xffffffff}
	gammaMask1 = [...]uint32{0xffffffff, 0x00000000, 0x00000000, 0x00000000}
	gammaMask2 = [...]uint32{0x00000000, 0x00000000, 0xffffffff, 0xffffffff}
	gammaMask3 = [...]uint32{0xffffffff, 0xffffffff, 0x00000000, 0x00000000}
	gammaMask4 = [...]uint32{0x00000000, 0x00000000, 0x00000000, 0xffffffff}
	gammaMask5 = [...]uint32{0xffffffff, 0xffffffff, 0xffffffff, 0x00000000}
)

func init() {
	_ = bytePermInfo
	_ = bytePermInfoR
	_ = msgWordPermInfo
	_ = gammaMask0
	_ = gammaMask1
	_ = gammaMask2
	_ = gammaMask3
	_ = gammaMask4
	_ = gammaMask5
}
