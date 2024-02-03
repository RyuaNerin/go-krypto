//go:build (amd64 || arm64) && !purego
// +build amd64 arm64
// +build !purego

package lea

import (
	"crypto/cipher"

	"github.com/RyuaNerin/go-krypto/kipher"
)

type cbcDecAble interface {
	NewCBCDecrypter(iv []byte) cipher.BlockMode
}
type ctrAble interface {
	NewCTR(iv []byte) cipher.Stream
}

var (
	_ cbcDecAble = (*leaContext)(nil)
	_ ctrAble    = (*leaContext)(nil)
)

func (ctx *leaContext) NewCBCDecrypter(iv []byte) cipher.BlockMode {
	return kipher.NewCBCDecrypter(ctx, iv)
}

func (leaCtx *leaContext) NewCTR(iv []byte) cipher.Stream {
	return kipher.NewCTR(leaCtx, iv)
}
