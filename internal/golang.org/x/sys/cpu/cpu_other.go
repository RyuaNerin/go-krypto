//go:build !386 && !amd64 && !amd64p32 && !arm && !arm64
// +build !386,!amd64,!amd64p32,!arm,!arm64

package cpu

const cacheLineSize = 1

var Initialized bool

func initOptions() {}
