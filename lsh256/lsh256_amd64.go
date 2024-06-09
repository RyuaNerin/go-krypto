//go:build amd64 && !purego && (!gccgo || go1.18)
// +build amd64
// +build !purego
// +build !gccgo go1.18

package lsh256

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
