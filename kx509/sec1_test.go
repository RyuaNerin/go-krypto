package kx509

import (
	"testing"
)

func TestSEC1ASN1DER(t *testing.T) {
	for _, tc := range eckcdsaTestCases {
		b, err := MarshalECKCPrivateKey(&tc.key)
		if err != nil {
			t.Error(err)
			return
		}

		key, err := ParseECKCPrivateKey(b)
		if err != nil {
			t.Error(err)
			return
		}

		if !tc.key.Equal(key) {
			t.Error("not equals!")
			return
		}
	}
}
