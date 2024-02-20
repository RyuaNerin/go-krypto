//go:build !(arm64 || amd64 || amd64p32) || purego
// +build !arm64,!amd64,!amd64p32 purego

package lsh512

var (
	newContext = newContextGo
	sum        = sumGo
)
