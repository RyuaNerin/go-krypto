//go:build amd64 && !purego && (!gccgo || go1.18)
// +build amd64
// +build !purego
// +build !gccgo go1.18

package aria

var newCipher = newCipherGo

func init() {
	if hasSSSE3 {
		newCipher = newCipherAsm
	}
}
