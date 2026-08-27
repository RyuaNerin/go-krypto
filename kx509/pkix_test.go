package kx509

import (
	"crypto/elliptic"
	"encoding/asn1"
	"encoding/pem"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"

	"github.com/stretchr/testify/require"
)

func TestPKIXPublicKey(t *testing.T) {
	t.Run("EC-KCDSA", func(t *testing.T) {
		for _, tc := range eckcdsaTestCases {
			pb, _ := pem.Decode([]byte(tc.pkixPublicKey))
			require.NotNil(t, pb, "failed to decode test PKIX PEM")
			expectedDER := pb.Bytes

			actualDER, err := MarshalPKIXPublicKey(&tc.key.PublicKey)
			require.NoError(t, err, "failed to marshal EC-KCDSA public key to PKIX")

			require.Equal(t, expectedDER, actualDER, "marshaled PKIX public key does not match expected DER")

			key, err := ParsePKIXPublicKey(expectedDER)
			require.NoError(t, err, "failed to parse EC-KCDSA public key from PKIX")
			equalECKCDSAPublicKey(t, &tc.key.PublicKey, key.(*eckcdsa.PublicKey))
		}
	})

	t.Run("KCDSA", func(t *testing.T) {
		for _, tc := range kcdsaTestCases {
			der, err := MarshalPKIXPublicKey(&tc.key.PublicKey)
			require.NoError(t, err, "failed to marshal KCDSA public key to PKIX")

			key, err := ParsePKIXPublicKey(der)
			require.NoError(t, err, "failed to parse KCDSA public key from PKIX")
			equalCDSAPublicKey(t, &tc.key.PublicKey, key.(*kcdsa.PublicKey))
		}
	})
}

func TestPKIXPublicKeyAlgorithmAliases(t *testing.T) {
	t.Run("EC-KCDSA", func(t *testing.T) {
		der, err := MarshalPKIXPublicKey(&eckcdsaTestCases[0].key.PublicKey)
		require.NoError(t, err, "failed to marshal EC-KCDSA public key to PKIX")

		var info pkixPublicKey
		_, err = asn1.Unmarshal(der, &info)
		require.NoError(t, err, "failed to unmarshal PKIX public key")

		for _, oid := range []asn1.ObjectIdentifier{oidPublicKeyECKCDSAAlteGOV} {
			info.Algo.Algorithm = oid
			info.Raw = nil
			aliasDER, err := asn1.Marshal(info)
			require.NoError(t, err, "failed to marshal PKIX public key with alias OID")

			key, err := ParsePKIXPublicKey(aliasDER)
			require.NoError(t, err, "failed to parse EC-KCDSA public key from PKIX with alias OID")
			require.IsType(t, &eckcdsa.PublicKey{}, key, "parsed key is not of type *eckcdsa.PublicKey")
		}
	})

	t.Run("KCDSA", func(t *testing.T) {
		der, err := MarshalPKIXPublicKey(&kcdsaTestCases[0].key.PublicKey)
		require.NoError(t, err, "failed to marshal KCDSA public key to PKIX")

		var info pkixPublicKey
		_, err = asn1.Unmarshal(der, &info)
		require.NoError(t, err, "failed to unmarshal PKIX public key")

		for _, oid := range []asn1.ObjectIdentifier{oidPublicKeyKCDSAAlteGOV, oidPublicKeyKCDSA1AlteGOV} {
			info.Algo.Algorithm = oid
			info.Raw = nil
			aliasDER, err := asn1.Marshal(info)
			require.NoError(t, err, "failed to marshal PKIX public key with alias OID")

			key, err := ParsePKIXPublicKey(aliasDER)
			require.NoError(t, err, "failed to parse KCDSA public key from PKIX with alias OID")
			require.IsType(t, &kcdsa.PublicKey{}, key, "parsed key is not of type *kcdsa.PublicKey")
		}
	})
}

func TestParseCompressedECKCDSAPublicKey(t *testing.T) {
	tc := eckcdsaTestCases[1]
	der, err := MarshalPKIXPublicKey(&tc.key.PublicKey)
	require.NoError(t, err, "failed to marshal EC-KCDSA public key to PKIX")

	var info pkixPublicKey
	_, err = asn1.Unmarshal(der, &info)
	require.NoError(t, err, "failed to unmarshal PKIX public key")

	point := elliptic.MarshalCompressed(tc.key.Curve, tc.key.X, tc.key.Y)
	info.BitString = asn1.BitString{Bytes: point, BitLength: len(point) * 8}
	info.Raw = nil
	der, err = asn1.Marshal(info)
	require.NoError(t, err, "failed to marshal PKIX public key with compressed point")

	key, err := ParsePKIXPublicKey(der)
	require.NoError(t, err, "failed to parse compressed EC-KCDSA public key from PKIX")
	equalECKCDSAPublicKey(t, &tc.key.PublicKey, key.(*eckcdsa.PublicKey))
}

func TestKCDSAOptionalGenerationParameters(t *testing.T) {
	key := kcdsaTestCases[0].key.PublicKey
	key.GenParameters = kcdsa.GenerationParameters{J: key.GenParameters.J}

	der, err := MarshalPKIXPublicKey(&key)
	require.NoError(t, err, "failed to marshal KCDSA public key to PKIX")

	parsedKey, err := ParsePKIXPublicKey(der)
	require.NoError(t, err, "failed to parse KCDSA public key from PKIX")

	parsed := parsedKey.(*kcdsa.PublicKey)
	require.NotNil(t, parsed.GenParameters, "KCDSA generation parameters were not preserved")
	require.NotNil(t, parsed.GenParameters.J, "KCDSA J parameter was not preserved")
	require.Equal(t, key.GenParameters.J.Text(16), parsed.GenParameters.J.Text(16), "KCDSA J parameter was not preserved")

	require.Empty(t, parsed.GenParameters.Seed, "KCDSA Seed parameter was not preserved")
	require.Zero(t, parsed.GenParameters.Count, "KCDSA Count parameter was not preserved")
}
