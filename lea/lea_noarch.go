//go:build !((arm64 || amd64) && !purego && (!gccgo || go1.18))
// +build !arm64,!amd64 purego gccgo,!go1.18

package lea

import "crypto/cipher"

func newCipher(key []byte) (cipher.Block, error) {
	return newCipherGo(key)
}
