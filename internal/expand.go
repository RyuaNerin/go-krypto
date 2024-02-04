package internal

import "github.com/RyuaNerin/go-krypto/internal/kryptoutil"

func Expand(arr []byte, bytes int) []byte {
	if bytes < cap(arr) {
		return arr[:bytes]
	} else {
		return make([]byte, bytes)
	}
}

func FitSize(arr []byte, bytes int) []byte {
	arrLen := len(arr)
	arrCap := cap(arr)

	switch {
	case arrLen == bytes:
	case arrLen < bytes:
		arr = arr[:bytes]

	case bytes <= arrCap:
		arr = arr[:bytes]
		kryptoutil.MemsetByte(arr[bytes:], 0)

	default:
		arr2 := make([]byte, bytes)
		copy(arr2, arr)
		arr = arr2
	}
	return arr
}
