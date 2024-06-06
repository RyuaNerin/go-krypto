package kx509

import (
	"bytes"
	"encoding/hex"
	"encoding/pem"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

func TestPKIXPublicKey(t *testing.T) {
	t.Run("EC-KCDSA", func(t *testing.T) {
		for _, tc := range eckcdsaTestCases {
			pb, _ := pem.Decode([]byte(tc.pkixPublicKey))
			expectedDER := pb.Bytes

			actualDER, err := MarshalPKIXPublicKey(&tc.key.PublicKey)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Equal(expectedDER, actualDER) {
				t.Errorf("Not equal:\nexpected: %s\nactual  : %s", hex.EncodeToString(expectedDER), hex.EncodeToString(actualDER))
				return
			}

			key, err := ParsePKIXPublicKey(expectedDER)
			if err != nil {
				t.Error(err)
				return
			}
			if !eqPubECKCDSA(t, &tc.key.PublicKey, key.(*eckcdsa.PublicKey)) {
				return
			}
		}
	})

	t.Run("KCDSA", func(t *testing.T) {
		for _, tc := range kcdsaTestCases {
			der, err := MarshalPKIXPublicKey(&tc.key.PublicKey)
			if err != nil {
				t.Error(err)
				return
			}

			key, err := ParsePKIXPublicKey(der)
			if err != nil {
				t.Error(err)
				return
			}
			if !eqPubKCDSA(t, &tc.key.PublicKey, key.(*kcdsa.PublicKey)) {
				return
			}
		}
	})
}
