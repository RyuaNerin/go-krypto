package internal

import (
	"io"

	"github.com/RyuaNerin/go-krypto/internal/kryptoutil"
)

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
func Grow(buf []byte, bytes int) []byte {
	if bytes < cap(buf) {
		return buf[:bytes]
	} else {
		return make([]byte, bytes)
	}
}

// keep data
func Resize(arr []byte, bytes int) []byte {
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
	bytes := BitsToBytes(bits)
	b = b[len(b)-bytes:]

	remain := bits % 8
	if remain > 0 {
		b[0] &= ((1 << remain) - 1)
	}

	return b
}

// [0 0 0 0 0 0]0 0
func LeftMost(b []byte, bits int) []byte {
	bytes := BitsToBytes(bits)
	b = b[:bytes]

	remain := bits % 8
	if remain > 0 {
		b[0] &= byte(0b_11111111 << (8 - remain))
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
		value >>= 8
	}
	kryptoutil.MemsetByte(dst[:len(dst)-n], 0)
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

// resize dst, ReadFull, cut from right
func ReadBits(dst []byte, rand io.Reader, bits int) ([]byte, error) {
	bytes := BitsToBytes(bits)

	dst = Grow(dst, bytes)

	if _, err := io.ReadFull(rand, dst); err != nil {
		return dst, err
	}

	bytes = bits & 0x07
	if bytes != 0 {
		dst[0] &= byte((1 << bytes) - 1)
	}

	return dst, nil
}

// resize dst, ReadFull, cut from right
func ReadBytes(dst []byte, rand io.Reader, bytes int) ([]byte, error) {
	dst = Grow(dst, bytes)

	if _, err := io.ReadFull(rand, dst); err != nil {
		return dst, err
	}

	return dst, nil
}
