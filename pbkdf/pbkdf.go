// Package pbkdf implements HMAC-based Key Derivation Functions, as defined in TTAK.KO-12.0334-Part1
package pbkdf

import (
	"crypto/hmac"
	"crypto/subtle"
	"hash"

	"github.com/RyuaNerin/go-krypto/internal"
)

// Generate a key from the password, salt and iteration count,
// then returns a []byte of length keylen.
func Generate(dst, password, salt []byte, iteration, keyLen int, h func() hash.Hash) []byte {
	hh := hmac.New(h, password)
	hLen := hh.Size()

	U := make([]byte, hLen, hLen*2)
	T := U[hLen : hLen*2]

	dst, out := internal.SliceForAppend(dst, keyLen)

	var i [4]byte
	for off := 0; off < keyLen; off += hLen {
		internal.IncCtr(i[:])

		hh.Reset()
		hh.Write(salt)
		hh.Write(i[:])
		U = hh.Sum(U[:0])

		copy(T, U)

		for iter := 1; iter < iteration; iter++ {
			hh.Reset()
			hh.Write(U)
			U = hh.Sum(U[:0])

			subtle.XORBytes(T, T, U)
		}

		copy(out[off:], T)
	}

	return dst
}
