package kx509

import (
	"bytes"
	"encoding/hex"
	"encoding/pem"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

func TestPKCS8PrivateKey(t *testing.T) {
	t.Run("EC-KCDSA", func(t *testing.T) {
		for _, tc := range eckcdsaTestCases {
			pb, _ := pem.Decode([]byte(tc.pkcs8PrivateKey))
			expectedDER := pb.Bytes

			actualDER, err := MarshalPKCS8PrivateKey(&tc.key)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(expectedDER, actualDER) {
				t.Errorf("Not equal:\nexpected: %s\nactual  : %s", hex.EncodeToString(expectedDER), hex.EncodeToString(actualDER))
				return
			}

			key, err := ParsePKCS8PrivateKey(expectedDER)
			if err != nil {
				t.Error(err)
				return
			}
			if !eqPrivECKCDSA(t, &tc.key, key.(*eckcdsa.PrivateKey)) {
				return
			}
		}
	})

	t.Run("KCDSA", func(t *testing.T) {
		for _, tc := range kcdsaTestCases {
			der, err := MarshalPKCS8PrivateKey(tc.key)
			if err != nil {
				t.Error(err)
				return
			}

			psdKeyRaw, err := ParsePKCS8PrivateKey(der)
			if err != nil {
				t.Error(err)
				return
			}
			psdKey := psdKeyRaw.(*kcdsa.PrivateKey)

			if !eqPrivKCDSA(t, &tc.key, psdKey) {
				return
			}
		}
	})
}
