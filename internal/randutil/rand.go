package randutil

import (
	"bufio"
	"crypto/elliptic"
	"io"
	"math/big"
	"math/rand/v2"
)

func NewTestRand() (io.Reader, *rand.Rand) {
	r := rand.NewChaCha8([32]byte{})
	rnd := rand.New(r)
	return r, rnd
}

func NewBenchRand() (io.Reader, *rand.Rand) {
	r := rand.NewChaCha8([32]byte{})
	rnd := rand.New(r)
	return bufio.NewReaderSize(r, 1<<20), rnd
}

// https://cs.opensource.google/go/go/+/refs/tags/go1.20.7:src/crypto/ecdsa/ecdsa_legacy.go;l=168-188
// randFieldElement returns a random element of the order of the given
// curve using the procedure given in FIPS 186-4, Appendix B.5.2.
func RandFieldElement(k *big.Int, c elliptic.Curve, rand io.Reader) (err error) {
	// See randomPoint for notes on the algorithm. This has to match, or s390x
	// signatures will come out different from other architectures, which will
	// break TLS recorded tests.

	N := c.Params().N                   // 1. N = len(n)
	b := make([]byte, (N.BitLen()+7)/8) // 2. If N is invalid, then return an ERROR indication, Invalid_d, and Invalid_Q.
	defer clear(b)

	for {
		if _, err = io.ReadFull(rand, b); err != nil { // 3
			return err
		}
		if excess := len(b)*8 - N.BitLen(); excess > 0 {
			b[0] >>= excess
		}
		// 6. If (c > n-2), then go to step 4.
		// 7. d = c + 1.
		//
		// d > n-1 ===> d >= n
		k.SetBytes(b)
		if k.Sign() > 0 && k.Cmp(N) < 0 {
			return nil
		}
	}

	/**
	B.4.2 Key Pair Generation by Testing Candidates

	In this method, a random number is obtained and tested to determine that it will produce a value
		of d in the correct range. If d is out-of-range, another random number is obtained (i.e., the
		process is iterated until an acceptable value of d is obtained.
	The following process or its equivalent may be used to generate an ECC key pair.
	Input:
		1. (q, FR, a, b {, domain_parameter_seed}, G, n, h)
			The domain parameters that are used for this process. n is a prime number,
			and G is a point on the elliptic curve.
	Output:
		1. status The status returned from the key pair generation procedure. The status will
			indicate SUCCESS or an ERROR.
		2. (d, Q) The generated private and public keys. If an error is encountered during
			the generation process, invalid values for d and Q should be returned, as
			represented by Invalid_d and Invalid_Q in the following specification. d is
			an integer, and Q is an elliptic curve point. The generated private key d is
			in the range [1, n-1].

	Process:
		1. N = len(n)
		2. If N is invalid, then return an ERROR indication, Invalid_d, and Invalid_Q.
		3. requested_security_strength = the security strength associated with N; see SP 800-57, Part 1.
		4. Obtain a string of N returned_bits from an RBG with a security strength of
			requested_security_strength or more. If an ERROR indication is returned, then
			return an ERROR indication, Invalid_d, and Invalid_Q.
		5. Convert returned_bits to the (non-negative) integer c (see Appendix C.2.1)
		6. If (c > n-2), then go to step 4.
		7. d = c + 1.
	*/
}
