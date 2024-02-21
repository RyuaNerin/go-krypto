package eckcdsa

import (
	"errors"
	"hash"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal/golang.org/x/crypto/cryptobyte"
	"github.com/RyuaNerin/go-krypto/internal/golang.org/x/crypto/cryptobyte/asn1"
)

// Sign data using K generated randomly like in crypto/ecdsa packages.
// returns the ASN.1 encoded signature.
func SignASN1(randReader io.Reader, priv *PrivateKey, h hash.Hash, data []byte) (sig []byte, err error) {
	r, s, err := Sign(randReader, priv, h, data)
	if err != nil {
		return nil, err
	}

	return encodeSignature(r.Bytes(), s.Bytes())
}

// VerifyASN1 verifies the ASN.1 encoded signature, sig, M, of hash using the
// public key, pub. Its return value records whether the signature is valid.
func VerifyASN1(pub *PublicKey, h hash.Hash, data, sig []byte) bool {
	r, s, err := parseSignature(sig)
	if err != nil {
		return false
	}

	return Verify(
		pub,
		h,
		data,
		new(big.Int).SetBytes(r),
		new(big.Int).SetBytes(s),
	)
}

// https://github.com/golang/go/blob/go1.21.6/src/crypto/ecdsa/ecdsa.go#L338-L345
func encodeSignature(r, s []byte) ([]byte, error) {
	var b cryptobyte.Builder
	b.AddASN1(asn1.SEQUENCE, func(b *cryptobyte.Builder) {
		addASN1IntBytes(b, r)
		addASN1IntBytes(b, s)
	})
	return b.Bytes()
}

// https://github.com/golang/go/blob/go1.21.6/src/crypto/ecdsa/ecdsa.go#L349-L363
func addASN1IntBytes(b *cryptobyte.Builder, bytes []byte) {
	for len(bytes) > 0 && bytes[0] == 0 {
		bytes = bytes[1:]
	}
	if len(bytes) == 0 {
		b.SetError(errors.New("invalid integer"))
		return
	}
	b.AddASN1(asn1.INTEGER, func(c *cryptobyte.Builder) {
		if bytes[0]&0x80 != 0 {
			c.AddUint8(0)
		}
		c.AddBytes(bytes)
	})
}

// https://github.com/golang/go/blob/master/src/crypto/ecdsa/ecdsa.go#L549-L560
func parseSignature(sig []byte) (r, s []byte, err error) {
	var inner cryptobyte.String
	input := cryptobyte.String(sig)
	if !input.ReadASN1(&inner, asn1.SEQUENCE) ||
		!input.Empty() ||
		!inner.ReadASN1Integer(&r) ||
		!inner.ReadASN1Integer(&s) ||
		!inner.Empty() {
		return nil, nil, errors.New("invalid ASN.1")
	}
	return r, s, nil
}
