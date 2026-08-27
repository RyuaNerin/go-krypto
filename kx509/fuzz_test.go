package kx509

import "testing"

const maxFuzzKeyDER = 4096

func FuzzParsePKIXPublicKey(f *testing.F) {
	for _, key := range []any{&eckcdsaTestCases[0].key.PublicKey, &kcdsaTestCases[0].key.PublicKey} {
		der, err := MarshalPKIXPublicKey(key)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(der)
	}
	f.Add([]byte{0x30, 0})
	f.Fuzz(func(t *testing.T, der []byte) {
		if len(der) > maxFuzzKeyDER {
			return
		}
		_, _ = ParsePKIXPublicKey(der)
	})
}

func FuzzParsePKCS8PrivateKey(f *testing.F) {
	for _, key := range []any{&eckcdsaTestCases[0].key, &kcdsaTestCases[0].key} {
		der, err := MarshalPKCS8PrivateKey(key)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(der)
	}
	f.Add([]byte{0x30, 0})
	f.Fuzz(func(t *testing.T, der []byte) {
		if len(der) > maxFuzzKeyDER {
			return
		}
		_, _ = ParsePKCS8PrivateKey(der)
	})
}

func FuzzParseECKCPrivateKey(f *testing.F) {
	der, err := MarshalECKCPrivateKey(&eckcdsaTestCases[0].key)
	if err != nil {
		f.Fatal(err)
	}
	f.Add(der)
	f.Add([]byte{0x30, 0})
	f.Fuzz(func(t *testing.T, der []byte) {
		if len(der) > maxFuzzKeyDER {
			return
		}
		_, _ = ParseECKCPrivateKey(der)
	})
}
