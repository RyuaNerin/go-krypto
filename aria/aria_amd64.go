//go:build amd64 && !purego
// +build amd64,!purego

package aria

var newCipher = newCipherGo

func init() {
	if hasSSSE3 {
		newCipher = newCipherAsm
	}
}
