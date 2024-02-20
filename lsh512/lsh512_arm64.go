//go:build arm64 && !purego
// +build arm64,!purego

package lsh512

var (
	newContext = simdSetNEON.NewContext
	sum        = simdSetNEON.Sum
)
