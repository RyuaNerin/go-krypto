//go:build amd64 && gc && !purego

package lea

import "unsafe"

//go:noescape
func __lea_encrypt_8block(ct, pt, rk unsafe.Pointer, round uint64)

//go:noescape
func __lea_decrypt_8block(pt, ct, rk unsafe.Pointer, round uint64)
