//go:build amd64

package lea

import "golang.org/x/sys/cpu"

var (
	hasAVX2 = cpu.X86.HasAVX2
)

func init() {
	leaEnc4 = leaEnc4SSE2
	leaDec4 = leaDec4SSE2

	leaEnc8 = leaEnc8SSE2
	leaDec8 = leaDec8SSE2

	if hasAVX2 {
		leaEnc8 = leaEnc8AVX2
		leaDec8 = leaDec8AVX2
	}
}

func leaEnc8SSE2(round int, rk []uint32, dst, src []byte) {
	leaEnc4SSE2(round, rk, dst[0x00:], src[0x00:])
	leaEnc4SSE2(round, rk, dst[0x40:], src[0x40:])
}
func leaDec8SSE2(round int, rk []uint32, dst, src []byte) {
	leaDec4SSE2(round, rk, dst[0x00:], src[0x00:])
	leaDec4SSE2(round, rk, dst[0x40:], src[0x40:])
}

func leaEnc4SSE2(round int, rk []uint32, dst []byte, src []byte)

func leaDec4SSE2(round int, rk []uint32, dst []byte, src []byte)

func leaEnc8AVX2(round int, rk []uint32, dst []byte, src []byte)

func leaDec8AVX2(round int, rk []uint32, dst []byte, src []byte)
