//go:build (amd64 || amd64p32) && !purego
// +build amd64 amd64p32
// +build !purego

package lsh512

var (
	newContext = simdSetSSE2.NewContext
	sum        = simdSetSSE2.Sum
)

func init() {
	if hasAVX2 {
		newContext = simdSetAVX2.NewContext
		sum = simdSetAVX2.Sum
	} else if hasSSSE3 {
		newContext = simdSetSSSE3.NewContext
		sum = simdSetSSSE3.Sum
	}
}
