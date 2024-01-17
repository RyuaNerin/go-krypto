//go:build exclude
// +build exclude

package aria

//go:noescape
func __EncKeySetup(rk *byte, mk *byte, keyBits uint64)

//go:noescape
func __DecKeySetup(rk *byte, rounds uint64)

//go:noescape
func __Crypt(dst, src, rk *byte, round uint64)
