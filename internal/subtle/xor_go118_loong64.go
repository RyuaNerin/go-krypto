// https://github.com/golang/go/blob/69827b5c8ddf93be65bfc8b17d331bb3ff7b704c/src/crypto/subtle/xor_loong64.go

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (gc && !go1.20 && go1.18 && !purego && loong64) || (gc && go1.20 && !purego && loong64)
// +build gc,!go1.20,go1.18,!purego,loong64 gc,go1.20,!purego,loong64

package subtle

//go:noescape
func xorBytes(dst, a, b *byte, n int)
