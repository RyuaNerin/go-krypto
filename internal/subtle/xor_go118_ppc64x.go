// https://github.com/golang/go/blob/go1.23.1/src/crypto/subtle/xor_ppc64x.go

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build gc && !go1.20 && go1.18 && !purego && (ppc64 || ppc64le)
// +build gc
// +build !go1.20
// +build go1.18
// +build !purego
// +build ppc64 ppc64le

package subtle

//go:noescape
func xorBytes(dst, a, b *byte, n int)
