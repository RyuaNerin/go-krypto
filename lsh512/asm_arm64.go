//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh512 && arm64

package lsh512

import "unsafe"

func init() {
	defaultSimd = simdNEON{}
}

type simdNEON struct{}

var _ simdSet = simdNEON{}

func (s simdNEON) update(c, d unsafe.Pointer, n int64) { __lsh512_neon_update(c, d, n) }
func (s simdNEON) final(c *lsh512ContextData) (sum [Size]byte) {
	__lsh512_neon_final(unsafe.Pointer(c), unsafe.Pointer(&sum[0]))
	return sum
}

// constants for SIMD implementation

const uint64Max uint64 = ^uint64(0)

var (
	gammaNEON    = [...]uint64{0, 16, 32, 48, 8, 24, 40, 56}
	gammaNEONInv = [...]uint64{
		uint64Max - 64 + 1, uint64Max - 48 + 1, uint64Max - 32 + 1, uint64Max - 16 + 1,
		uint64Max - 56 + 1, uint64Max - 40 + 1, uint64Max - 24 + 1, uint64Max - 8 + 1,
	}
)

func init() {
	_ = gammaNEON
	_ = gammaNEONInv
}
