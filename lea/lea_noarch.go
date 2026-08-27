//go:build !(!noasm && !noasm_krypto && !noasm_krypto_lea && (amd64 || arm64))

package lea

import "crypto/cipher"

func newCipher(key []byte) (cipher.Block, error) {
	return newCipherGo(key)
}
