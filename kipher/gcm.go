// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kipher

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/gcm.go

import (
	"crypto/cipher"
	"errors"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/alias"
	"github.com/RyuaNerin/go-krypto/internal/kipher"
	"github.com/RyuaNerin/go-krypto/internal/memory"
	"github.com/RyuaNerin/go-krypto/internal/subtle"
)

type gcm struct {
	gcm       kipher.GCM
	cipher    kipher.Block
	nonceSize int
	tagSize   int
}

// NewGCM returns the given 128-bit, block cipher wrapped in Galois Counter Mode
// with the standard nonce length.
//
// In general, the GHASH operation performed by this implementation of GCM is not constant-time.
// An exception is when the underlying Block was created by aes.NewCipher
// on systems with hardware support for AES. See the crypto/aes package documentation for details.
func NewGCM(b cipher.Block) (cipher.AEAD, error) {
	return NewGCMWithSize(b, kipher.GCMStandardNonceSize, kipher.GCMTagSize)
}

// NewGCMWithNonceSize returns the given 128-bit, block cipher wrapped in Galois
// Counter Mode, which accepts nonces of the given length. The length must not
// be zero.
//
// Only use this function if you require compatibility with an existing
// cryptosystem that uses non-standard nonce lengths. All other users should use
// NewGCM, which is faster and more resistant to misuse.
func NewGCMWithNonceSize(b cipher.Block, size int) (cipher.AEAD, error) {
	return NewGCMWithSize(b, size, kipher.GCMTagSize)
}

// NewGCMWithTagSize returns the given 128-bit, block cipher wrapped in Galois
// Counter Mode, which generates tags with the given length.
//
// Tag sizes between 1 and 16 bytes are allowed.
//
// Only use this function if you require compatibility with an existing
// cryptosystem that uses non-standard tag lengths. All other users should use
// NewGCM, which is more resistant to misuse.
func NewGCMWithTagSize(b cipher.Block, tagSize int) (cipher.AEAD, error) {
	return NewGCMWithSize(b, kipher.GCMStandardNonceSize, tagSize)
}

// NewGCM returns the given 128-bit, block cipher wrapped in Galois Counter Mode
// with the standard nonce length.
//
// In general, the GHASH operation performed by this implementation of GCM is not constant-time.
// An exception is when the underlying Block was created by aes.NewCipher
// on systems with hardware support for AES. See the crypto/aes package documentation for details.
//
// if nonceSize = 0, nonceSize = 12
// if tagSize = 0, tagSize = 16
func NewGCMWithSize(b cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	if nonceSize == 0 {
		nonceSize = kipher.GCMStandardNonceSize
	}
	if tagSize == 0 {
		tagSize = kipher.GCMTagSize
	}

	if nonceSize <= 0 {
		return nil, errors.New(msgInvalidNonceZero)
	}

	kb, ok := b.(kipher.Block)
	if !ok {
		if cipher, ok := b.(kipher.GCMAble); ok {
			if tagSize < kipher.GCMMinimumTagSize || tagSize > kipher.GCMBlockSize {
				return nil, errors.New(msgInvalidTagSizeGCM)
			}
			return cipher.NewGCM(nonceSize, tagSize)
		}
		kb = kipher.WrapKipher(b)
	}

	if b.BlockSize() != kipher.GCMBlockSize {
		return nil, errors.New(msgRequire128Bits)
	}

	return newGCMWithNonceAndTagSize(kb, nonceSize, tagSize)
}

func newGCMWithNonceAndTagSize(b kipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	g := &gcm{
		cipher:    b,
		nonceSize: nonceSize,
		tagSize:   tagSize,
	}
	g.gcm.Init(b)

	return g, nil
}

func (g *gcm) NonceSize() int {
	return g.nonceSize
}

func (g *gcm) Overhead() int {
	return g.tagSize
}

func (g *gcm) Seal(dst, nonce, plaintext, data []byte) []byte {
	if len(nonce) != g.nonceSize {
		panic(msgInvalidNonce)
	}
	if uint64(len(plaintext)) > ((1<<32)-2)*uint64(kipher.GCMBlockSize) {
		panic(msgDataTooLarge)
	}

	ret, out := internal.SliceForAppend(dst, len(plaintext)+g.tagSize)
	if alias.InexactOverlap(out, plaintext) {
		panic(msgBufferOverlap)
	}

	var counter, tagMask [kipher.GCMBlockSize]byte
	g.gcm.DeriveCounter(&counter, nonce)

	g.cipher.Encrypt(tagMask[:], counter[:])
	internal.IncCtr(counter[kipher.GCMBlockSize-4:])

	kipher.GCMCounterCrypt(out, plaintext, g.cipher, &counter)

	var tag [kipher.GCMTagSize]byte
	g.gcm.Auth(tag[:], out[:len(plaintext)], data, &tagMask)
	copy(out[len(plaintext):], tag[:])

	return ret
}

func (g *gcm) Open(dst, nonce, ciphertext, data []byte) ([]byte, error) {
	if len(nonce) != g.nonceSize {
		panic(msgInvalidNonce)
	}

	if len(ciphertext) < g.tagSize {
		return nil, errors.New(msgOpenFailed)
	}
	if uint64(len(ciphertext)) > ((1<<32)-2)*kipher.GCMBlockSize+uint64(g.tagSize) {
		return nil, errors.New(msgOpenFailed)
	}

	tag := ciphertext[len(ciphertext)-g.tagSize:]
	ciphertext = ciphertext[:len(ciphertext)-g.tagSize]

	var counter, tagMask [kipher.GCMBlockSize]byte
	g.gcm.DeriveCounter(&counter, nonce)

	g.cipher.Encrypt(tagMask[:], counter[:])
	internal.IncCtr(counter[kipher.GCMBlockSize-4:])

	var expectedTag [kipher.GCMTagSize]byte
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
		memory.Memclr(out)
		return nil, errors.New(msgOpenFailed)
	}

	kipher.GCMCounterCrypt(out, ciphertext, g.cipher, &counter)

	return ret, nil
}
