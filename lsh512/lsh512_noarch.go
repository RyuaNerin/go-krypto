//go:build !(arm64 || amd64) || purego || (gccgo && !go1.18)
// +build !arm64,!amd64 purego gccgo,!go1.18

package lsh512

var (
	newContext = newContextGo
	sum        = sumGo
)
