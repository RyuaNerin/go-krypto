package internal

import "github.com/RyuaNerin/go-krypto/internal/kryptoutil"

// same with `int(math.Ceil(float64(a) / float64(b)))`
func CeilDiv(a, b int) int {
	if b == 0 {
		return 0
	}
	return (a + b - 1) / b
}

func IncCtr(b []byte) {
	for i := len(b) - 1; i >= 0; i-- {
		c := b[i]
		c++
		b[i] = c
		if c > 0 {
			return
		}
	}
}

func Add(dst []byte, src ...[]byte) {
	n := len(dst)

	var value uint64
	for idx := 0; idx < n; idx++ {
		for _, v := range src {
			if idx < len(v) {
				value += uint64(v[len(v)-idx-1])
			}
		}

		dst[len(dst)-idx-1] = byte(value & 0xFF)
		value = value >> 8
	}
	kryptoutil.MemsetByte(dst[:len(dst)-n], 0)
}
