//go:build amd64 && gc && !purego

package lea

import "unsafe"

//go:noescape
func __lea_encrypt_4block(ct, pt, rk unsafe.Pointer, round uint64)

//go:noescape
func __lea_decrypt_4block(pt, ct, rk unsafe.Pointer, round uint64)
