package kx509

import (
	"crypto/elliptic"
	"encoding/asn1"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"

	"github.com/stretchr/testify/require"
)

func TestSEC1ASN1DER(t *testing.T) {
	for _, tc := range eckcdsaTestCases {
		b, err := MarshalECKCPrivateKey(&tc.key)
		require.NoError(t, err, "failed to marshal EC-KCDSA private key to SEC1 ASN.1 DER")

		key, err := ParseECKCPrivateKey(b)
		require.NoError(t, err, "failed to parse EC-KCDSA private key from SEC1 ASN.1 DER")

		require.True(t, tc.key.Equal(key), "EC-KCDSA private key did not round trip through SEC1 ASN.1 DER")
	}
}

func TestECPrivateKeyUsesFixedWidthScalar(t *testing.T) {
	curve := elliptic.P256()
	x, y := curve.ScalarBaseMult([]byte{1})
	key := &eckcdsa.PrivateKey{
		D: big.NewInt(1),
		PublicKey: eckcdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
	}

	der, err := MarshalECKCPrivateKey(key)
	require.NoError(t, err, "failed to marshal EC-KCDSA private key to SEC1 ASN.1 DER")

	var encoded eckcdsaPrivateKey
	rest, err := asn1.Unmarshal(der, &encoded)
	require.NoError(t, err, "failed to unmarshal ECPrivateKey")
	require.Empty(t, rest, "unexpected trailing data after ECPrivateKey")

	require.Len(t, encoded.PrivateKey, 32, "unexpected EC private key length")

	want := make([]byte, 32)
	want[31] = 1
	require.Equal(t, want, encoded.PrivateKey, "EC private key scalar is not fixed-width")

	parsed, err := ParseECKCPrivateKey(der)
	require.NoError(t, err, "failed to parse EC-KCDSA private key from SEC1 ASN.1 DER")
	require.True(t, key.Equal(parsed), "fixed-width EC-KCDSA private key did not round trip")
}

func TestECPrivateKeyRejectsZeroScalar(t *testing.T) {
	der, err := asn1.Marshal(eckcdsaPrivateKey{
		Version:       ecPrivKeyVersion,
		PrivateKey:    make([]byte, 32),
		NamedCurveOID: oidNamedCurveP256,
	})
	require.NoError(t, err, "failed to marshal ECPrivateKey with zero private scalar")

	_, err = ParseECKCPrivateKey(der)
	require.Error(t, err, "zero private scalar was accepted")
}
