package kcdsa_test

// https://github.com/golang/go/blob/go1.22.4/src/crypto/ecdsa/equal_test.go

// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"crypto"
	"io"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/randutil"
	"github.com/RyuaNerin/go-krypto/kcdsa"
	"github.com/RyuaNerin/go-krypto/kx509"

	"github.com/stretchr/testify/require"
)

func TestEqual(t *testing.T) {
	for _, tc := range testConfigList {
		t.Run(tc.Name, func(t *testing.T) {
			testEqual(t, tc.Sizes)
		})
	}
}

func generateKey(t testing.TB, rand io.Reader, sizes kcdsa.ParameterSizes) (*kcdsa.PrivateKey, error) {
	t.Helper()

	var key kcdsa.PrivateKey
	err := kcdsa.GenerateParameters(&key.Parameters, rand, sizes)
	require.NoError(t, err, "GenerateParameters failed")

	err = kcdsa.GenerateKey(&key, rand)
	require.NoError(t, err, "GenerateKey failed")
	return &key, nil
}

func testEqual(t *testing.T, sizes kcdsa.ParameterSizes) {
	rnd, _ := randutil.NewTestRand()

	private, _ := generateKey(t, rnd, sizes)
	public := &private.PublicKey

	require.True(t, public.Equal(public), "public key is not equal to itself")
	require.True(t, public.Equal(crypto.Signer(private).Public().(*kcdsa.PublicKey)), "public key is not equal to itself via crypto.Signer")
	require.True(t, private.Equal(private), "private key is not equal to itself")

	enc, err := kx509.MarshalPKCS8PrivateKey(private)
	require.NoError(t, err, "MarshalPKCS8PrivateKey failed")

	decoded, err := kx509.ParsePKCS8PrivateKey(enc)
	require.NoError(t, err, "ParsePKCS8PrivateKey failed")
	require.True(t, public.Equal(decoded.(crypto.Signer).Public()), "public key is not equal to itself after decoding")
	require.True(t, private.Equal(decoded), "private key is not equal to itself after decoding")

	other, _ := generateKey(t, rnd, sizes)
	require.False(t, public.Equal(other.Public()), "different public keys are Equal")
	require.False(t, private.Equal(other), "different private keys are Equal")

	// Ensure that keys with the same coordinates but on different parameters
	// aren't considered Equal.
	differentParameters := &kcdsa.PublicKey{}
	*differentParameters = *public // make a copy of the public key
	differentParameters.Parameters = other.Parameters
	require.False(t, public.Equal(differentParameters), "public keys with different parameters are Equal")
}
