package kcdsattak

import (
	"crypto/rand"
	"testing"

	"github.com/RyuaNerin/go-krypto/kcdsa"
)

func Test_TTAK_GenerateJ(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
		return
	}

	for _, tc := range testCase_TestVector {
		J, err := GenerateJ(tc.Seed_, tc.Sizes)
		if err != nil {
			t.Error(err)
			return
		}
		if J.Cmp(tc.J) != 0 {
			t.Errorf("GenerateTTAKJ failed")
			return
		}
	}
}

func Test_TTAK_GeneratePQ(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
		return
	}

	for _, tc := range testCase_TestVector {
		P, Q, count, err := GeneratePQ(tc.J, tc.Seed_, tc.Sizes)
		if err != nil {
			t.Error(err)
			return
		}
		if P.Cmp(tc.P) != 0 || Q.Cmp(tc.Q) != 0 || count != tc.Count {
			t.Errorf("GenerateTTAKPQ failed")
			return
		}
	}
}

func Test_TTAK_GenerateHG(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
		return
	}

	for _, tc := range testCase_TestVector {
		_, _, err := GenerateHG(rand.Reader, tc.P, tc.J)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func Test_TTAK_GenerateG(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
		return
	}

	for _, tc := range testCase_TestVector {
		G, err := GenerateG(tc.P, tc.J, tc.H)
		if err != nil {
			t.Error(err)
			return
		}
		if G.Cmp(tc.G) != 0 {
			t.Errorf("GenerateTTAKG failed")
			return
		}
	}
}

func Test_TTAK_GenerateXYZ(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
		return
	}

	for _, tc := range testCase_TestVector {
		X, Y, Z, _, err := GenerateXYZ(tc.P, tc.Q, tc.G, UserProvidedRandomInput, tc.XKEY, tc.Sizes)
		if err != nil {
			t.Error(err)
			return
		}
		if X.Cmp(tc.X) != 0 || Y.Cmp(tc.Y) != 0 || Z.Cmp(tc.Z) != 0 {
			t.Errorf("GenerateTTAKX failed")
			return
		}
	}
}

func Test_kcdsa_GenerateParametersTTAK(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
		return
	}

	gp := func(params *Parameters, sizes kcdsa.ParameterSizes) error {
		_, _, err := GenerateParameters(params, rand.Reader, sizes)
		return err
	}
	gk := func(priv *PrivateKey, sizes kcdsa.ParameterSizes) error {
		return GenerateKey(priv, rand.Reader, UserProvidedRandomInput, sizes)
	}

	testKCDSA(t, kcdsa.L2048N224SHA224, 2048, 224, gp, gk)
	testKCDSA(t, kcdsa.L2048N224SHA256, 2048, 224, gp, gk)
	testKCDSA(t, kcdsa.L2048N256SHA256, 2048, 256, gp, gk)
	testKCDSA(t, kcdsa.L3072N256SHA256, 3072, 256, gp, gk)
}
