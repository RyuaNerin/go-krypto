package kcdsa_test

import (
	"crypto/sha256"
	"io"
	"slices"
	"testing"

	"github.com/RyuaNerin/go-krypto/has160"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/cavp"
	kcdsa_internal "github.com/RyuaNerin/go-krypto/internal/kcdsa"
	"github.com/RyuaNerin/go-krypto/internal/randutil"
	"github.com/RyuaNerin/go-krypto/kcdsa"

	"github.com/stretchr/testify/require"
)

func init() {
	_ = sha256.Size
	_ = has160.Size
}

type KcdsaTestConfig struct {
	Name        string
	Sizes       kcdsa.ParameterSizes
	KPGFilename string
	SVTFilename string
}

var testConfigList = []KcdsaTestConfig{
	{
		Name:        "A2048 B224 SHA224",
		Sizes:       kcdsa.A2048B224SHA224,
		KPGFilename: "14. KCDSA/KCDSA_(2048)(224)(SHA-224)_KPG",
		SVTFilename: "14. KCDSA/KCDSA_(2048)(224)(SHA-224)_SVT",
	},
	{
		Name:  "A2048 B224 SHA256",
		Sizes: kcdsa.A2048B224SHA256,
	},
	{
		Name:        "A2048 B256 SHA256",
		Sizes:       kcdsa.A2048B256SHA256,
		KPGFilename: "14. KCDSA/KCDSA_(2048)(256)(SHA-256)_KPG",
		SVTFilename: "14. KCDSA/KCDSA_(2048)(256)(SHA-256)_SVT",
	},
	{
		Name:  "A3072 B256 SHA256",
		Sizes: kcdsa.A3072B256SHA256,
	},
}

func Test_KCDSA_KPG(t *testing.T) {
	for _, tc := range testConfigList {
		if tc.KPGFilename == "" {
			continue
		}
		t.Run(tc.Name, func(t *testing.T) {
			cavp.TestKCDSAKpg(t, tc.KPGFilename, tc.Sizes)
		})
	}
}

func Test_KCDSA_SVT(t *testing.T) {
	for _, tc := range testConfigList {
		if tc.SVTFilename == "" {
			continue
		}
		t.Run(tc.Name, func(t *testing.T) {
			cavp.TestKCDSASvt(t, tc.SVTFilename, tc.Sizes)
		})
	}
}

func Test_SignVerify_With_BadPublicKey(t *testing.T) {
	for idx, tc := range testCaseTTAK {
		tc2 := testCaseTTAK[(idx+1)%len(testCaseTTAK)]

		pub := kcdsa.PublicKey{
			Parameters: kcdsa.Parameters{
				P: tc2.P,
				Q: tc2.Q,
				G: tc2.G,
			},
			Y: tc2.Y,
		}

		ok := kcdsa.Verify(&pub, tc.Sizes, tc.M, tc.R, tc.S)
		require.False(t, ok, "Verify unexpected success with bad public key")
	}
}

func Test_Signing_With_DegenerateKeys(t *testing.T) {
	badKeys := []struct {
		p, q, g, y, x string
	}{
		{"00", "01", "00", "00", "00"},
		{"01", "ff", "00", "00", "00"},
	}

	rnd, _ := randutil.NewTestRand()

	msg := []byte("testing")
	for i, test := range badKeys {
		priv := kcdsa.PrivateKey{
			PublicKey: kcdsa.PublicKey{
				Parameters: kcdsa.Parameters{
					P: internal.HI(test.p),
					Q: internal.HI(test.q),
					G: internal.HI(test.g),
				},
				Y: internal.HI(test.y),
			},
			X: internal.HI(test.x),
		}

		_, _, err := kcdsa.Sign(rnd, &priv, kcdsa.A2048B224SHA224, msg)
		require.Error(t, err, "Sign should fail with degenerate key(%d)", i)
	}
}

func Test_KCDSA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
	}

	t.Run("A2048B224SHA224", testKCDSA(kcdsa.A2048B224SHA224))
	t.Run("A2048B224SHA256", testKCDSA(kcdsa.A2048B224SHA256))
	t.Run("A2048B256SHA256", testKCDSA(kcdsa.A2048B256SHA256))
	t.Run("A3072B256SHA256", testKCDSA(kcdsa.A3072B256SHA256))
	t.Run("A1024B160HAS160", testKCDSA(kcdsa.A1024B160HAS160))
}

func testKCDSA(sizes kcdsa.ParameterSizes) func(*testing.T) {
	return func(t *testing.T) {
		rnd, _ := randutil.NewTestRand()

		d, ok := kcdsa_internal.GetDomain(int(sizes))
		require.True(t, ok, "domain not found")

		var priv kcdsa.PrivateKey
		params := &priv.Parameters

		err := kcdsa.GenerateParameters(params, rnd, sizes)
		require.NoError(t, err, "GenerateParameters failed")

		require.LessOrEqual(t, params.P.BitLen(), d.A, "P is too large")
		require.LessOrEqual(t, params.Q.BitLen(), d.B, "Q is too large")

		err = kcdsa.GenerateKey(&priv, rnd)
		require.NoError(t, err, "GenerateKey failed")

		testSignAndVerify(t, &priv, sizes, rnd)
		testSignAndVerifyASN1(t, &priv, sizes, rnd)
	}
}

func testSignAndVerify(t *testing.T, priv *kcdsa.PrivateKey, sizes kcdsa.ParameterSizes, rnd io.Reader) {
	data := []byte("testing")
	r, s, err := kcdsa.Sign(rnd, priv, sizes, data)
	require.NoError(t, err, "Sign failed")

	ok := kcdsa.Verify(&priv.PublicKey, sizes, data, r, s)
	require.True(t, ok, "Verify failed")

	data[0] ^= 0xff
	ok = kcdsa.Verify(&priv.PublicKey, sizes, data, r, s)
	require.False(t, ok, "Verify should fail with modified data")
}

func testSignAndVerifyASN1(t *testing.T, priv *kcdsa.PrivateKey, sizes kcdsa.ParameterSizes, rnd io.Reader) {
	data := []byte("testing")
	dataClone := slices.Clone(data)
	privClone := priv.Clone()
	sig, err := kcdsa.SignASN1(rnd, priv, sizes, data)
	require.NoError(t, err, "SignASN1 failed")
	require.Equal(t, dataClone, data, "SignASN1 modified data")
	require.True(t, privClone.Equal(priv), "SignASN1 modified private key")

	sigClone := slices.Clone(sig)
	ok := kcdsa.VerifyASN1(&priv.PublicKey, sizes, data, sig)
	require.True(t, ok, "VerifyASN1 failed")
	require.True(t, privClone.PublicKey.Equal(&priv.PublicKey), "VerifyASN1 modified public key")
	require.Equal(t, dataClone, data, "VerifyASN1 modified data")
	require.Equal(t, sigClone, sig, "VerifyASN1 modified signature")

	data[0] ^= 0xff
	ok = kcdsa.VerifyASN1(&priv.PublicKey, sizes, data, sig)
	require.False(t, ok, "VerifyASN1 should fail with modified data")
}

func verifyTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		pub := kcdsa.PublicKey{
			Parameters: kcdsa.Parameters{
				P: tc.P,
				Q: tc.Q,
				G: tc.G,
			},
			Y: tc.Y,
		}

		ok := kcdsa.Verify(&pub, tc.Sizes, tc.M, tc.R, tc.S)
		if tc.Fail {
			require.False(t, ok, "Verify should fail for test case with Fail=true")
		} else {
			require.True(t, ok, "Verify failed for test case with Fail=false")
		}
	}
}
