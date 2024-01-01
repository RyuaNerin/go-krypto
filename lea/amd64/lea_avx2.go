//go:build exclude

package lea

import "unsafe"

//go:noescape
func __lea_encrypt_8block(ctx unsafe.Pointer, dst, src []byte)

//go:noescape
func __lea_decrypt_8block(ctx unsafe.Pointer, dst, src []byte)
