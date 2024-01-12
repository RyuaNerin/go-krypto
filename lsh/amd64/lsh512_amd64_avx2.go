//go:build exclude
// +build exclude

package lsh512

import "unsafe"

//go:noescape
func __lsh512_avx2_init(ctx unsafe.Pointer, algtype uint64)

//go:noescape
func __lsh512_avx2_update(ctx, data unsafe.Pointer, databytelen uint64)

//go:noescape
func __lsh512_avx2_final(ctx, hashval unsafe.Pointer)
