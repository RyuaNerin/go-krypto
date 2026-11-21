package kx509

import (
	"encoding/asn1"
	"encoding/pem"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"

	"github.com/stretchr/testify/require"
)

func TestPKCS8PrivateKey(t *testing.T) {
	t.Run("EC-KCDSA", func(t *testing.T) {
		for _, tc := range eckcdsaTestCases {
			pb, _ := pem.Decode([]byte(tc.pkcs8PrivateKey))
			require.NotNil(t, pb, "failed to decode test PKCS#8 PEM")

			// Existing Botan keys, including a P-521 key with a shortened d,
			// remain parseable for interoperability.
			key, err := ParsePKCS8PrivateKey(pb.Bytes)
			require.NoError(t, err, "failed to parse test PKCS#8 private key")
			equalECKCDSAPrivateKey(t, &tc.key, key.(*eckcdsa.PrivateKey))

			actualDER, err := MarshalPKCS8PrivateKey(&tc.key)
			require.NoError(t, err, "failed to marshal EC-KCDSA private key to PKCS#8")

			var info pkcds8PrivateKey
			rest, err := asn1.Unmarshal(actualDER, &info)
			require.NoError(t, err)
			require.Empty(t, rest, "unexpected trailing data after PrivateKeyInfo")

			require.Equal(t, 0, info.Version, "unexpected PrivateKeyInfo version")
			require.Equal(t, oidPublicKeyECKCDSA, info.Algo.Algorithm, "unexpected PrivateKeyInfo algorithm OID")

			var ecKey eckcdsaPrivateKey
			rest, err = asn1.Unmarshal(info.PrivateKey, &ecKey)
			require.NoError(t, err)
			require.Empty(t, rest, "unexpected trailing data after ECPrivateKey")

			wantPrivateKeyLen := (tc.key.Curve.Params().N.BitLen() + 7) / 8
			require.Len(t, ecKey.PrivateKey, wantPrivateKeyLen, "unexpected EC private key length")

			key, err = ParsePKCS8PrivateKey(actualDER)
			require.NoError(t, err)
			equalECKCDSAPrivateKey(t, &tc.key, key.(*eckcdsa.PrivateKey))
		}
	})

	t.Run("KCDSA", func(t *testing.T) {
		for _, tc := range kcdsaTestCases {
			der, err := MarshalPKCS8PrivateKey(tc.key)
			require.NoError(t, err, "failed to marshal KCDSA private key to PKCS#8")

			psdKeyRaw, err := ParsePKCS8PrivateKey(der)
			require.NoError(t, err, "failed to parse KCDSA private key from PKCS#8")
			psdKey := psdKeyRaw.(*kcdsa.PrivateKey)

			equalCDSAPrivateKey(t, &tc.key, psdKey)

			var info pkcds8PrivateKey
			rest, err := asn1.Unmarshal(der, &info)
			require.NoError(t, err, "failed to unmarshal PKCS#8 private key")
			require.Empty(t, rest, "unexpected trailing data after PrivateKeyInfo")

			require.Equal(t, 0, info.Version, "unexpected PrivateKeyInfo version")
			require.Equal(t, oidPublicKeyKCDSA, info.Algo.Algorithm, "unexpected PrivateKeyInfo algorithm OID")

			x := new(big.Int)
			rest, err = asn1.Unmarshal(info.PrivateKey, &x)
			require.NoError(t, err, "failed to unmarshal KCDSA private key INTEGER x")
			require.Empty(t, rest, "unexpected trailing data after KCDSA private key INTEGER x")
			require.Equal(t, tc.key.X.Text(16), x.Text(16), "KCDSA private key INTEGER x does not match expected value")

			params, err := parseKCDSAParameters(info.Algo.Parameters.FullBytes)
			require.NoError(t, err, "failed to parse KCDSA generation parameters")
			require.NotNil(t, params, "KCDSA generation parameters are nil")
			require.NotNil(t, params.J, "KCDSA generation parameters J is nil")
			require.Equal(t, tc.key.GenParameters.J.Text(16), params.J.Text(16), "KCDSA generation parameters J does not match expected value")
		}
	})
}

func TestPKCS8PrivateKeyAlgorithmAliases(t *testing.T) {
	t.Run("EC-KCDSA", func(t *testing.T) {
		der, err := MarshalPKCS8PrivateKey(&eckcdsaTestCases[0].key)
		require.NoError(t, err, "failed to marshal EC-KCDSA private key to PKCS#8")
		require.NotNil(t, der, "failed to marshal EC-KCDSA private key to PKCS#8")

		var info pkcds8PrivateKey
		_, err = asn1.Unmarshal(der, &info)
		require.NoError(t, err, "failed to unmarshal PKCS#8 private key")

		for _, oid := range []asn1.ObjectIdentifier{oidPublicKeyECKCDSAAlteGOV} {
			info.Algo.Algorithm = oid
			aliasDER, err := asn1.Marshal(info)
			require.NoError(t, err, "failed to marshal PKCS#8 private key with alias OID")

			key, err := ParsePKCS8PrivateKey(aliasDER)
			require.NoError(t, err, "failed to parse EC-KCDSA private key from PKCS#8 with alias OID")
			require.IsType(t, &eckcdsa.PrivateKey{}, key, "parsed key is not of type *eckcdsa.PrivateKey")
		}
	})

	t.Run("KCDSA", func(t *testing.T) {
		der, err := MarshalPKCS8PrivateKey(&kcdsaTestCases[0].key)
		require.NoError(t, err, "failed to marshal KCDSA private key to PKCS#8")

		var info pkcds8PrivateKey
		_, err = asn1.Unmarshal(der, &info)
		require.NoError(t, err, "failed to unmarshal PKCS#8 private key")

		for _, oid := range []asn1.ObjectIdentifier{oidPublicKeyKCDSAAlteGOV, oidPublicKeyKCDSA1AlteGOV} {
			info.Algo.Algorithm = oid
			aliasDER, err := asn1.Marshal(info)
			require.NoError(t, err, "failed to marshal PKCS#8 private key with alias OID")

			key, err := ParsePKCS8PrivateKey(aliasDER)
			require.NoError(t, err, "failed to parse KCDSA private key from PKCS#8 with alias OID")
			require.IsType(t, &kcdsa.PrivateKey{}, key, "parsed key is not of type *kcdsa.PrivateKey")
		}
	})
}
