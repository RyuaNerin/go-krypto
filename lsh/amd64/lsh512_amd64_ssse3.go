//go:build exclude

package lsh512

import "unsafe"

//go:noescape
func __lsh512_ssse3_init(ctx unsafe.Pointer, algtype uint64)

//go:noescape
func __lsh512_ssse3_update(ctx, data unsafe.Pointer, databytelen uint64)

//go:noescape
func __lsh512_ssse3_final(ctx, hashval unsafe.Pointer)
