//go:build exclude
// +build exclude

package aria

import "unsafe"

//go:noescape
func __initEncKey_SSSE3(mk, rk unsafe.Pointer, keyBits uint64)

//go:noescape
func __initDecKey_SSSE3(mk, rk unsafe.Pointer, keyBits uint64)

//go:noescape
func __process_SSSE3(dst, src, rk unsafe.Pointer, rounds uint64)
