//go:build (!amd64 && !arm64) || purego
// +build !amd64,!arm64 purego

package aria

import "crypto/cipher"

func newCipher(key []byte) (cipher.Block, error) {
	return newCipherGo(key)
}
