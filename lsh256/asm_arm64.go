//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh256 && arm64

package lsh256

import "unsafe"

func init() {
	defaultSimd = simdNEON{}
}

type simdNEON struct{}

var _ simdSet = simdNEON{}

func (s simdNEON) update(c, d unsafe.Pointer, n int64) { __lsh256_neon_update(c, d, n) }
func (s simdNEON) final(c *lsh256ContextData) (sum [Size]byte) {
	__lsh256_neon_final(unsafe.Pointer(c), unsafe.Pointer(&sum[0]))
	return sum
}

// constants for SIMD implementation

const uint32Max uint32 = ^uint32(0)

var (
	gammaNEON    = [...]int32{0, 8, 16, 24, 24, 16, 8, 0}
	gammaNEONInv = [...]uint32{
		uint32Max - 32 + 1, uint32Max - 24 + 1, uint32Max - 16 + 1, uint32Max - 8 + 1,
		uint32Max - 8 + 1, uint32Max - 16 + 1, uint32Max - 24 + 1, uint32Max - 32 + 1,
	}
)

func init() {
	_ = gammaNEON
	_ = gammaNEONInv
}
