//go:build !noasm && !noasm_krypto && !noasm_krypto_lea && arm64

package lea

import "crypto/cipher"

func newCipherAsm(key []byte) (cipher.Block, error) {
	ctx := new(leaContextAsm)
	ctx.round = setKeyGo(&ctx.rk, key)

	return ctx, nil
}
