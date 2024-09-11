//go:build gc && (arm64 || amd64) && !purego && go1.18
// +build gc
// +build arm64 amd64
// +build !purego
// +build go1.18

package lsh512

import "hash"

func newContext(size int) hash.Hash {
	return newContextAsm(size)
}

func sum(size int, data []byte) [Size]byte {
	return sumAsm(size, data)
}
