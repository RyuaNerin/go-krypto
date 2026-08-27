package kipher

// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"crypto/subtle"

	"github.com/RyuaNerin/go-krypto"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/kipher/kipherimpl"
)

func NewGCM(b krypto.Block8, nonceSize, tagSize int) *GCM {
	g := &GCM{
		cipher:    b,
		nonceSize: nonceSize,
		tagSize:   tagSize,
	}
	g.gcm.Init(b)

	return g
}

type GCM struct {
	gcm       kipherimpl.GCM
	cipher    krypto.Block8
	nonceSize int
	tagSize   int
}

var _ krypto.AEAD = (*GCM)(nil)

func (g *GCM) Zeroize() {
	g.gcm.Zeroize()
	g.cipher = nil
	g.nonceSize = 0
	g.tagSize = 0
}

func (g *GCM) NonceSize() int {
	return g.nonceSize
}

func (g *GCM) Overhead() int {
	return g.tagSize
}

func (g *GCM) Seal(dst, nonce, plaintext, data []byte) []byte {
	if len(nonce) != g.nonceSize {
		panic(msgInvalidNonce)
	}
	if uint64(len(plaintext)) > ((1<<32)-2)*uint64(kipherimpl.GCMBlockSize) {
		panic(msgDataTooLarge)
	}

	ret, out := internal.SliceForAppend(dst, len(plaintext)+g.tagSize)
	if alias.InexactOverlap(out, plaintext) {
		panic(msgBufferOverlap)
	}

	var counter, tagMask [kipherimpl.GCMBlockSize]byte
	defer clear(counter[:])
	defer clear(tagMask[:])
	g.gcm.DeriveCounter(&counter, nonce)

	g.cipher.Encrypt(tagMask[:], counter[:])
	internal.IncCtr(counter[kipherimpl.GCMBlockSize-4:])

	kipherimpl.GCMCounterCrypt(out, plaintext, g.cipher, &counter)

	var tag [kipherimpl.GCMTagSize]byte
	defer clear(tag[:])
	g.gcm.Auth(tag[:], out[:len(plaintext)], data, &tagMask)
	copy(out[len(plaintext):], tag[:])

	return ret
}

func (g *GCM) Open(dst, nonce, ciphertext, data []byte) ([]byte, error) {
	if len(nonce) != g.nonceSize {
		panic(msgInvalidNonce)
	}

	if len(ciphertext) < g.tagSize {
		return nil, ErrOpenFailed
	}
	if uint64(len(ciphertext)) > ((1<<32)-2)*kipherimpl.GCMBlockSize+uint64(g.tagSize) {
		return nil, ErrOpenFailed
	}

	tag := ciphertext[len(ciphertext)-g.tagSize:]
	ciphertext = ciphertext[:len(ciphertext)-g.tagSize]

	var counter, tagMask [kipherimpl.GCMBlockSize]byte
	defer clear(counter[:])
	defer clear(tagMask[:])
	g.gcm.DeriveCounter(&counter, nonce)

	g.cipher.Encrypt(tagMask[:], counter[:])
	internal.IncCtr(counter[kipherimpl.GCMBlockSize-4:])

	var expectedTag [kipherimpl.GCMTagSize]byte
	defer clear(expectedTag[:])
	g.gcm.Auth(expectedTag[:], ciphertext, data, &tagMask)

	ret, out := internal.SliceForAppend(dst, len(ciphertext))
	if alias.InexactOverlap(out, ciphertext) {
		panic(msgBufferOverlap)
	}

	if subtle.ConstantTimeCompare(expectedTag[:g.tagSize], tag) != 1 {
		// The AESNI code decrypts and authenticates concurrently, and
		// so overwrites dst in the event of a tag mismatch. That
		// behavior is mimicked here in order to be consistent across
		// platforms.
		clear(out)
		return nil, ErrOpenFailed
	}

	kipherimpl.GCMCounterCrypt(out, ciphertext, g.cipher, &counter)

	return ret, nil
}
