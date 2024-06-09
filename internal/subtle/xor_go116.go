// based on https://github.com/golang/go/blob/go1.15.15/src/crypto/cipher/xor_generic.go

//go:build !go1.17 && (purego || gccgo)
// +build !go1.17
// +build purego gccgo

package subtle

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"reflect"
	"runtime"
	"unsafe"
)

const wordSize = int(unsafe.Sizeof(uintptr(0)))

const supportsUnaligned = runtime.GOARCH == "386" ||
	runtime.GOARCH == "amd64" ||
	runtime.GOARCH == "ppc64" ||
	runtime.GOARCH == "ppc64le" ||
	runtime.GOARCH == "s390x"

func XORBytes(dst, x, y []byte) int {
	n := len(x)
	if len(y) < n {
		n = len(y)
	}
	if n == 0 {
		return 0
	}
	if n > len(dst) {
		panic("subtle.XORBytes: dst too short")
	}

	dst = dst[:n]
	x = x[:n]
	y = y[:n]

	// xorBytes assembly is written using pointers and n. Back to slices.
	if supportsUnaligned || aligned(&dst[0], &x[0], &y[0]) {
		xorLoopUintptr(words(dst), words(x), words(y))
		if n%wordSize == 0 {
			return n
		}
		done := n &^ (wordSize - 1)
		dst = dst[done:]
		x = x[done:]
		y = y[done:]
	}
	xorLoopByte(dst, x, y)
	return n
}

func words(b []byte) []uintptr {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b)) //nolint:govet
	header.Len /= wordSize
	header.Cap /= wordSize
	return *(*[]uintptr)(unsafe.Pointer(&header)) //nolint:govet
}

// aligned reports whether dst, x, and y are all word-aligned pointers.
func aligned(dst, x, y *byte) bool {
	return (uintptr(unsafe.Pointer(dst))|uintptr(unsafe.Pointer(x))|uintptr(unsafe.Pointer(y)))&uintptr(wordSize-1) == 0
}

func xorLoopUintptr(dst, x, y []uintptr) {
	x = x[:len(dst)] // remove bounds check in loop
	y = y[:len(dst)] // remove bounds check in loop
	for i := range dst {
		dst[i] = x[i] ^ y[i]
	}
}

func xorLoopByte(dst, x, y []byte) {
	x = x[:len(dst)] // remove bounds check in loop
	y = y[:len(dst)] // remove bounds check in loop
	for i := range dst {
		dst[i] = x[i] ^ y[i]
	}
}
