//go:build !(gc && (arm64 || amd64) && !purego && go1.18)
// +build !gc !arm64,!amd64 purego !go1.18

package lsh512

import "hash"

func newContext(size int) hash.Hash {
	return newContextGo(size)
}

func sum(size int, data []byte) [Size]byte {
	return sumGo(size, data)
}
