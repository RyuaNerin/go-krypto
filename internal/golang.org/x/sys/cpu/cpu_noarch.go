//go:build !(arm || 386 || amd64 || amd64p32)
// +build !arm,!386,!amd64,!amd64p32

package cpu

const cacheLineSize = 1

func archInit() {
}

func initOptions() {
}
