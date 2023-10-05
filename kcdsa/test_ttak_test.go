package kcdsa

import (
	"crypto/rand"
	"testing"

	"github.com/RyuaNerin/go-krypto/kcdsa/kcdsattak"
)

func Test_TTAK_GenerateJ(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
		return
	}

	for _, tc := range testCase_TTAK {
		domain, _ := tc.Sizes.domain()
		J, err := kcdsattak.GenerateJ(tc.Seed_, domain.Domain)
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

	for _, tc := range testCase_TTAK {
		domain, _ := tc.Sizes.domain()
		P, Q, count, err := kcdsattak.GeneratePQ(tc.J, tc.Seed_, domain.Domain)
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

	for _, tc := range testCase_TTAK {
		_, _, err := kcdsattak.GenerateHG(rand.Reader, tc.P, tc.J)
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

	for _, tc := range testCase_TTAK {
		G, err := kcdsattak.GenerateG(tc.P, tc.J, tc.H)
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

	for _, tc := range testCase_TTAK {
		domain, _ := tc.Sizes.domain()
		X, Y, Z, _, err := kcdsattak.GenerateXYZ(tc.P, tc.Q, tc.G, UserProvidedRandomInput, tc.XKEY, domain.Domain)
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

	gp := func(params *Parameters, sizes ParameterSizes) error {
		_, _, err := GenerateParametersTTAK(params, rand.Reader, sizes)
		return err
	}
	gk := func(priv *PrivateKey) error {
		return GenerateKeyTTAK(priv, rand.Reader, UserProvidedRandomInput)
	}

	testKCDSA(t, L2048N224SHA224, 2048, 224, gp, gk)
	testKCDSA(t, L2048N224SHA256, 2048, 224, gp, gk)
	testKCDSA(t, L2048N256SHA256, 2048, 256, gp, gk)
	testKCDSA(t, L3072N256SHA256, 3072, 256, gp, gk)
}
