//go:build (amd64 || arm64) && !purego && (!gccgo || go1.18)
// +build amd64 arm64
// +build !purego
// +build !gccgo go1.18

package aria

import "crypto/cipher"

func newCipher(key []byte) (cipher.Block, error) {
	return newCipherAsm(key)
}
