package randutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"io"
)

// https://github.com/golang/go/blob/master/src/crypto/ecdsa/ecdsa.go#L409-L464

func MixedCSPRNG(rand io.Reader, D []byte, hash []byte) (io.Reader, error) {
	// This implementation derives the nonce from an AES-CTR CSPRNG keyed by:
	//
	//    SHA2-512(priv.D || entropy || hash)[:32]
	//
	// The CSPRNG key is indifferentiable from a random oracle as shown in
	// [Coron], the AES-CTR stream is indifferentiable from a random oracle
	// under standard cryptographic assumptions (see [Larsson] for examples).
	//
	// [Coron]: https://cs.nyu.edu/~dodis/ps/merkle.pdf
	// [Larsson]: https://web.archive.org/web/20040719170906/https://www.nada.kth.se/kurser/kth/2D1441/semteo03/lecturenotes/assump.pdf

	// Get 256 bits of entropy from rand.
	entropy := make([]byte, 32)
	if _, err := io.ReadFull(rand, entropy); err != nil {
		return nil, err
	}

	// Initialize an SHA-512 hash context; digest...
	md := sha512.New()
	md.Write(D)             // the private key,
	md.Write(entropy)       // the entropy,
	md.Write(hash)          // and the input hash;
	key := md.Sum(nil)[:32] // and compute ChopMD-256(SHA-512),
	// which is an indifferentiable MAC.

	// Create an AES-CTR instance to use as a CSPRNG.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a CSPRNG that xors a stream of zeros with
	// the output of the AES-CTR instance.
	const aesIV = "IV for go-krypto"
	return &cipher.StreamReader{
		R: zeroReader,
		S: cipher.NewCTR(block, []byte(aesIV)),
	}, nil
}

type zr struct{}

var zeroReader = zr{}

// Read replaces the contents of dst with zeros. It is safe for concurrent use.
func (zr) Read(dst []byte) (n int, err error) {
	for i := range dst {
		dst[i] = 0
	}
	return len(dst), nil
}
