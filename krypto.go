// Package krypto collects cryptographic algorithms designed in Republic of Korea
package krypto

import (
	"crypto"
	"crypto/cipher"
	"hash"
)

type (
	Hash interface {
		hash.Hash
		Zeroizable
	}
)

// MaxSupportedBlocks is the maximum number of blocks some functions in the supported
// algorithms can process in a single call.
const MaxSupportedBlocks = 16

type (
	Block interface {
		cipher.Block
		Zeroizable
	}
	Block4 interface {
		Block

		// Encrypt4 encrypts 4 blocks of data (ECB mode).
		//
		// It may be faster than calling Encrypt four times.
		Encrypt4(dst, src []byte)
		// Decrypt4 decrypts 4 blocks of data (ECB mode).
		//
		// It may be faster than calling Decrypt four times.
		Decrypt4(dst, src []byte)
	}
	Block8 interface {
		Block4

		// Encrypt8 encrypts 8 blocks of data (ECB mode).
		//
		// It may be faster than calling Encrypt eight times.
		Encrypt8(dst, src []byte)
		// Decrypt8 decrypts 8 blocks of data (ECB mode).
		//
		// It may be faster than calling Decrypt eight times.
		Decrypt8(dst, src []byte)
	}
	Block16 interface {
		Block8

		// Encrypt16 encrypts 16 blocks of data (ECB mode).
		//
		// It may be faster than calling Encrypt sixteen times.
		Encrypt16(dst, src []byte)
		// Decrypt16 decrypts 16 blocks of data (ECB mode).
		//
		// It may be faster than calling Decrypt sixteen times.
		Decrypt16(dst, src []byte)
	}

	BlockMode interface {
		cipher.BlockMode
		Zeroizable
	}
	Stream interface {
		cipher.Stream
		Zeroizable
	}
	AEAD interface {
		cipher.AEAD
		Zeroizable
	}
)

type (
	CBCEncAble interface {
		NewCBCEncrypter(iv []byte) cipher.BlockMode
	}

	CBCDecAble interface {
		NewCBCDecrypter(iv []byte) cipher.BlockMode
	}

	GCMAble interface {
		NewGCM(nonceSize, tagSize int) (cipher.AEAD, error)
	}

	CTRAble interface {
		NewCTR(iv []byte) cipher.Stream
	}
)

type (
	DRBG interface {
		// Read reads random bytes from the DRBG.
		Read(p []byte) (n int, err error)

		// Generate generates random bytes from the DRBG with additional input.
		Generate(dst []byte, additionalInput []byte) (n int, err error)

		// Reseed reseeds the DRBG with additional input.
		//
		// uses entropy source to reseed the DRBG.
		Reseed(additionalInput []byte) error

		// ReseedWithEntropy reseeds the DRBG with additional input and entropy input.
		ReseedWithEntropy(additionalInput, entropyInput []byte) error

		Zeroizable
	}
)

type (
	PrivateKey interface {
		crypto.PrivateKey

		Equal(x crypto.PrivateKey) bool

		Public() crypto.PublicKey

		Zeroizable
	}
	PublicKey interface {
		crypto.PublicKey

		Equal(x crypto.PublicKey) bool
		Zeroizable
	}
)
