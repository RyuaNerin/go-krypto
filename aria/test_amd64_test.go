//go:build amd64 && gc && !purego

package aria

import (
	"testing"

	"github.com/RyuaNerin/go-krypto/testingutil"
)

func Test_processFinSSE2(t *testing.T) {
	testingutil.BTTC(
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
	)
}
