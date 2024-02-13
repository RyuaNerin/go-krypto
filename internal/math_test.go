package internal

import (
	"math"
	"testing"
)

func TestBytes(t *testing.T) {
	for bits := 0; bits < 0xFF_FFFF; bits++ {
		answer := int(math.Ceil(float64(bits) / 8))

		if Bytes(bits) != answer {
			t.Fail()
			return
		}
	}
}
