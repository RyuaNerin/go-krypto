//go:build !(arm64 || amd64 || amd64p32) || purego
// +build !arm64,!amd64,!amd64p32 purego

package lsh256

var (
	newContext = newContextGo
	sum        = sumGo
)
