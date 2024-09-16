// https://github.com/golang/go/blob/release-branch.go1.21/src/crypto/subtle/xor_amd64.go

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build gc && !go1.20 && go1.18 && !purego && amd64
// +build gc,!go1.20,go1.18,!purego,amd64

package subtle

//go:noescape
func xorBytes(dst, a, b *byte, n int)
