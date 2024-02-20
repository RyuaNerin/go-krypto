//go:build arm64 && !purego
// +build arm64,!purego

package lsh256

var (
	newContext = simdSetNEON.NewContext
	sum        = simdSetNEON.Sum
)
