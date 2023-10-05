//go:build amd64 && gc && !purego

package aria

import (
	"testing"

	. "github.com/RyuaNerin/go-krypto/testingutil"
)

func Test_processFin_SSE2(t *testing.T) {
	BTTC(
		t,
		0,
		0,
		BlockSize*2,
		0,
		func(key, additional []byte) (interface{}, error) {
			return nil, nil
		},
		func(c interface{}, dst, src []byte) {
			processFinGo(dst, src[:BlockSize], src[BlockSize:])
		},
		func(c interface{}, dst, src []byte) {
			processFinSSE2(dst, src[:BlockSize], src[BlockSize:])
		},
		false,
	)
}

////////////////////////////////////////////////////////////////////////////////////////////////////

func Benchmark_processFin_Go(b *testing.B)   { benchmarkProcessFin(b, processFinGo) }
func Benchmark_processFin_SSE2(b *testing.B) { benchmarkProcessFin(b, processFinSSE2) }

func benchmarkProcessFin(b *testing.B, f func(dst, rk, t []byte)) {
	BBD(
		b,
		0,
		0,
		BlockSize*2,
		func(key, additional []byte) (interface{}, error) {
			return nil, nil
		},
		func(c interface{}, dst, src []byte) {
			processFinGo(dst, src[:BlockSize], src[BlockSize:])
		},
		false,
	)
}
