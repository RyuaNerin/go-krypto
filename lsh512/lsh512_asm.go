//go:build !noasm && !noasm_krypto && !noasm_krypto_lsh512 && (arm64 || amd64)

package lsh512

import "hash"

func newContext(size int) hash.Hash {
	return newContextAsm(size)
}

func sum(size int, data []byte) [Size]byte {
	return sumAsm(size, data)
}
