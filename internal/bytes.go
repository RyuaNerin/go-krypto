package internal

import "github.com/RyuaNerin/go-krypto/internal/kryptoutil"

func SliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

// without guarantee of data
func ResizeBuffer(buf []byte, bytes int) []byte {
	if bytes < cap(buf) {
		return buf[:bytes]
	} else {
		return make([]byte, bytes)
	}
}

// keep data
func ResizeSlice(arr []byte, bytes int) []byte {
	arrLen := len(arr)
	arrCap := cap(arr)

	switch {
	case arrLen == bytes:
	case arrLen < bytes:
		arr = arr[:bytes]

	case bytes <= arrCap:
		arr = arr[:bytes]
		kryptoutil.MemsetByte(arr[arrLen:], 0)

	default:
		arr2 := make([]byte, bytes)
		copy(arr2, arr)
		arr = arr2
	}
	return arr
}

// 0 0[0 0 0 0 0 0]
func RightMost(b []byte, bits int) []byte {
	bytes := Bytes(bits)
	b = b[len(b)-bytes:]

	remain := bits % 8
	if remain > 0 {
		b[0] = b[0] & ((1 << remain) - 1)
	}

	return b
}

// [0 0 0 0 0 0]0 0
func LeftMost(b []byte, bits int) []byte {
	bytes := Bytes(bits)
	b = b[:bytes]

	remain := bits % 8
	if remain > 0 {
		b[0] = b[0] & byte(0b_11111111<<(8-remain))
	}

	return b
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
