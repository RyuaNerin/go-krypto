// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kipher

// Based on https://github.com/golang/go/blob/go1.21.6/src/crypto/cipher/gcm.go

import (
	"crypto/cipher"
	"crypto/subtle"
	"encoding/binary"
	"errors"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/alias"
	igcm "github.com/RyuaNerin/go-krypto/internal/gcm"
)

type gcm struct {
	gcm igcm.GCM

	nonceSize int
	tagSize   int
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
func NewGCM(b cipher.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	if nonceSize <= 0 {
		return nil, errors.New(msgInvalidNonceZero)
	}

	kb, ok := b.(internal.Block)
	if !ok {
		if cipher, ok := b.(internal.GCMAble); ok {
			if tagSize < igcm.GCMMinimumTagSize || tagSize > igcm.GCMBlockSize {
				return nil, errors.New(msgInvalidTagSizeGCM)
			}
			return cipher.NewGCM(nonceSize, tagSize)
		}
		kb = internal.WrapBlock(b)
	}

	if b.BlockSize() != igcm.GCMBlockSize {
		return nil, errors.New(msgRequire128Bits)
	}

	if nonceSize == 0 {
		nonceSize = igcm.GCMStandardNonceSize
	}
	if tagSize == 0 {
		tagSize = igcm.GCMTagSize
	}

	return newGCMWithNonceAndTagSize(kb, nonceSize, tagSize)
}

// NewGCM returns the given 128-bit, block cipher wrapped in Galois Counter Mode
// with the standard nonce length.
//
// In general, the GHASH operation performed by this implementation of GCM is not constant-time.
// An exception is when the underlying Block was created by aes.NewCipher
// on systems with hardware support for AES. See the crypto/aes package documentation for details.
func NewGCMWithDefaultSize(b cipher.Block) (cipher.AEAD, error) {
	return NewGCM(b, igcm.GCMStandardNonceSize, igcm.GCMTagSize)
}

// NewGCMWithNonceSize returns the given 128-bit, block cipher wrapped in Galois
// Counter Mode, which accepts nonces of the given length. The length must not
// be zero.
//
// Only use this function if you require compatibility with an existing
// cryptosystem that uses non-standard nonce lengths. All other users should use
// NewGCM, which is faster and more resistant to misuse.
func NewGCMWithNonceSize(b cipher.Block, size int) (cipher.AEAD, error) {
	return NewGCM(b, size, igcm.GCMTagSize)
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
	return NewGCM(b, igcm.GCMStandardNonceSize, tagSize)
}

func newGCMWithNonceAndTagSize(b internal.Block, nonceSize, tagSize int) (cipher.AEAD, error) {
	g := &gcm{
		nonceSize: nonceSize,
		tagSize:   tagSize,
	}
	igcm.Init(&g.gcm, b)

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
	if uint64(len(plaintext)) > ((1<<32)-2)*uint64(g.gcm.Cipher.BlockSize()) {
		panic(msgDataTooLarge)
	}

	ret, out := internal.SliceForAppend(dst, len(plaintext)+g.tagSize)
	if alias.InexactOverlap(out, plaintext) {
		panic(msgBufferOverlap)
	}

	var counter, tagMask [igcm.GCMBlockSize]byte
	g.gcm.DeriveCounter(&counter, nonce)

	g.gcm.Cipher.Encrypt(tagMask[:], counter[:])
	igcm.GCMInc32(&counter)

	g.gcm.CounterCrypt(out, plaintext, &counter)

	var tag [igcm.GCMTagSize]byte
	g.auth(tag[:], out[:len(plaintext)], data, &tagMask)
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
	if uint64(len(ciphertext)) > ((1<<32)-2)*uint64(g.gcm.Cipher.BlockSize())+uint64(g.tagSize) {
		return nil, errors.New(msgOpenFailed)
	}

	tag := ciphertext[len(ciphertext)-g.tagSize:]
	ciphertext = ciphertext[:len(ciphertext)-g.tagSize]

	var counter, tagMask [igcm.GCMBlockSize]byte
	g.gcm.DeriveCounter(&counter, nonce)

	g.gcm.Cipher.Encrypt(tagMask[:], counter[:])
	igcm.GCMInc32(&counter)

	var expectedTag [igcm.GCMTagSize]byte
	g.auth(expectedTag[:], ciphertext, data, &tagMask)

	ret, out := internal.SliceForAppend(dst, len(ciphertext))
	if alias.InexactOverlap(out, ciphertext) {
		panic(msgBufferOverlap)
	}

	if subtle.ConstantTimeCompare(expectedTag[:g.tagSize], tag) != 1 {
		// The AESNI code decrypts and authenticates concurrently, and
		// so overwrites dst in the event of a tag mismatch. That
		// behavior is mimicked here in order to be consistent across
		// platforms.
		for i := range out {
			out[i] = 0
		}
		return nil, errors.New(msgOpenFailed)
	}

	g.gcm.CounterCrypt(out, ciphertext, &counter)

	return ret, nil
}

// auth calculates GHASH(ciphertext, additionalData), masks the result with
// tagMask and writes the result to out.
func (g *gcm) auth(out, ciphertext, additionalData []byte, tagMask *[igcm.GCMTagSize]byte) {
	var y igcm.GCMFieldElement
	g.gcm.Update(&y, additionalData)
	g.gcm.Update(&y, ciphertext)

	y.Low ^= uint64(len(additionalData)) * 8
	y.High ^= uint64(len(ciphertext)) * 8

	g.gcm.Mul(&y)

	binary.BigEndian.PutUint64(out, y.Low)
	binary.BigEndian.PutUint64(out[8:], y.High)

	subtle.XORBytes(out, out, tagMask[:])
}
