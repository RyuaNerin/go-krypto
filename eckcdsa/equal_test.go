// https://github.com/golang/go/blob/go1.22.4/src/crypto/ecdsa/equal_test.go

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eckcdsa_test

import (
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kx509"
)

func testEqual(t *testing.T, c elliptic.Curve) {
	private, _ := eckcdsa.GenerateKey(c, rand.Reader)
	public := &private.PublicKey

	if !public.Equal(public) { //nolint:gocritic
		t.Errorf("public key is not equal to itself: %v", public)
	}
	if !public.Equal(crypto.Signer(private).Public().(*eckcdsa.PublicKey)) {
		t.Errorf("private.Public() is not Equal to public: %q", public)
	}
	if !private.Equal(private) { //nolint:gocritic
		t.Errorf("private key is not equal to itself: %v", private)
	}

	enc, err := kx509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := kx509.ParsePKCS8PrivateKey(enc)
	if err != nil {
		t.Fatal(err)
	}
	if !public.Equal(decoded.(crypto.Signer).Public()) {
		t.Errorf("public key is not equal to itself after decoding: %v", public)
	}
	if !private.Equal(decoded) {
		t.Errorf("private key is not equal to itself after decoding: %v", private)
	}

	other, _ := eckcdsa.GenerateKey(c, rand.Reader)
	if public.Equal(other.Public()) {
		t.Errorf("different public keys are Equal")
	}
	if private.Equal(other) {
		t.Errorf("different private keys are Equal")
	}

	// Ensure that keys with the same coordinates but on different curves
	// aren't considered Equal.
	differentCurve := &eckcdsa.PublicKey{}
	*differentCurve = *public // make a copy of the public key
	if differentCurve.Curve == elliptic.P256() {
		differentCurve.Curve = elliptic.P224()
	} else {
		differentCurve.Curve = elliptic.P256()
	}
	if public.Equal(differentCurve) {
		t.Errorf("public keys with different curves are Equal")
	}
}

func TestEqual(t *testing.T) {
	t.Run("P224", func(t *testing.T) { testEqual(t, elliptic.P224()) })
	if testing.Short() {
		return
	}
	t.Run("P256", func(t *testing.T) { testEqual(t, elliptic.P256()) })
	t.Run("P384", func(t *testing.T) { testEqual(t, elliptic.P384()) })
	t.Run("P521", func(t *testing.T) { testEqual(t, elliptic.P521()) })
}
