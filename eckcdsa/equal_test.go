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

	"github.com/stretchr/testify/require"
)

func testEqual(t *testing.T, c elliptic.Curve) {
	private, _ := eckcdsa.GenerateKey(c, rand.Reader)
	public := &private.PublicKey

	require.True(t, public.Equal(public), "public key is not equal to itself")
	require.True(t, crypto.Signer(private).Public().(*eckcdsa.PublicKey).Equal(public), "private.Public() is not Equal to public")
	require.True(t, private.Equal(private), "private key is not equal to itself")

	enc, err := kx509.MarshalPKCS8PrivateKey(private)
	require.NoError(t, err)

	decoded, err := kx509.ParsePKCS8PrivateKey(enc)
	require.NoError(t, err)
	require.IsType(t, &eckcdsa.PrivateKey{}, decoded)
	require.IsType(t, &eckcdsa.PublicKey{}, decoded.(crypto.Signer).Public())
	require.True(t, public.Equal(decoded.(crypto.Signer).Public()), "public key is not equal to itself after decoding")
	require.True(t, private.Equal(decoded), "private key is not equal to itself after decoding")

	other, _ := eckcdsa.GenerateKey(c, rand.Reader)
	require.False(t, public.Equal(other.Public()), "different public keys are Equal")
	require.False(t, private.Equal(other), "different private keys are Equal")

	// Ensure that keys with the same coordinates but on different curves
	// aren't considered Equal.
	differentCurve := &eckcdsa.PublicKey{}
	*differentCurve = *public // make a copy of the public key
	if differentCurve.Curve == elliptic.P256() {
		differentCurve.Curve = elliptic.P224()
	} else {
		differentCurve.Curve = elliptic.P256()
	}
	require.False(t, public.Equal(differentCurve), "public keys with different curves are Equal")
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
