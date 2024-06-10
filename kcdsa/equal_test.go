// https://github.com/golang/go/blob/go1.22.4/src/crypto/ecdsa/equal_test.go

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package kcdsa_test

import (
	"crypto"
	"crypto/rand"
	"io"
	"testing"

	"github.com/RyuaNerin/go-krypto/kcdsa"
	"github.com/RyuaNerin/go-krypto/kx509"
)

func generateKey(rand io.Reader, sizes kcdsa.ParameterSizes) (*kcdsa.PrivateKey, error) {
	var key kcdsa.PrivateKey
	if err := kcdsa.GenerateParameters(&key.Parameters, rand, sizes); err != nil {
		return nil, err
	}
	if err := kcdsa.GenerateKey(&key, rand); err != nil {
		return nil, err
	}
	return &key, nil
}

func testEqual(t *testing.T, sizes kcdsa.ParameterSizes) {
	private, _ := generateKey(rand.Reader, sizes)
	public := &private.PublicKey

	if !public.Equal(public) { //nolint:gocritic
		t.Errorf("public key is not equal to itself: %v", public)
	}
	if !public.Equal(crypto.Signer(private).Public().(*kcdsa.PublicKey)) {
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

	other, _ := generateKey(rand.Reader, sizes)
	if public.Equal(other.Public()) {
		t.Errorf("different public keys are Equal")
	}
	if private.Equal(other) {
		t.Errorf("different private keys are Equal")
	}

	// Ensure that keys with the same coordinates but on different parameters
	// aren't considered Equal.
	differentParameters := &kcdsa.PublicKey{}
	*differentParameters = *public // make a copy of the public key
	differentParameters.Parameters = other.Parameters
	if public.Equal(differentParameters) {
		t.Errorf("public keys with different parameters are Equal")
	}
}

func TestEqual(t *testing.T) {
	t.Run("A2048B224SHA224", func(t *testing.T) { testEqual(t, kcdsa.A2048B224SHA224) })
	if testing.Short() {
		return
	}
	t.Run("A2048B224SHA256", func(t *testing.T) { testEqual(t, kcdsa.A2048B224SHA256) })
	t.Run("A2048B256SHA256", func(t *testing.T) { testEqual(t, kcdsa.A2048B256SHA256) })
	t.Run("A3072B256SHA256", func(t *testing.T) { testEqual(t, kcdsa.A3072B256SHA256) })
	t.Run("A1024B160HAS160", func(t *testing.T) { testEqual(t, kcdsa.A1024B160HAS160) })
}
