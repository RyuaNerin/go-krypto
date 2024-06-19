//go:build !((arm64 || amd64) && !purego && (!gccgo || go1.18))
// +build !arm64,!amd64 purego gccgo,!go1.18

package lsh256

import "hash"

func newContext(size int) hash.Hash {
	return newContextGo(size)
}

func sum(size int, data []byte) [Size]byte {
	return sumGo(size, data)
}
