//go:build !(gc && (arm64 || amd64) && !purego && go1.18)
// +build !gc !arm64,!amd64 purego !go1.18

package aria

import "crypto/cipher"

func newCipher(key []byte) (cipher.Block, error) {
	return newCipherGo(key)
}
