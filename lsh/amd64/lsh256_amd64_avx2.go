//go:build exclude
// +build exclude

package lsh256

import "unsafe"

//go:noescape
func __lsh256_avx2_init(ctx unsafe.Pointer, algtype uint32)

//go:noescape
func __lsh256_avx2_update(ctx, data unsafe.Pointer, databytelen uint64)

//go:noescape
func __lsh256_avx2_final(ctx, hashval unsafe.Pointer)
