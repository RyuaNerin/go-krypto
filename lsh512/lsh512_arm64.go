//go:build arm64 && !purego && (!gccgo || go1.18)
// +build arm64
// +build !purego
// +build !gccgo go1.18

package lsh512

var (
	newContext = simdSetNEON.NewContext
	sum        = simdSetNEON.Sum
)
