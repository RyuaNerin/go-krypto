package eckcdsa_test

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"io"
	"math/big"
	"slices"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/cavp"
	"github.com/RyuaNerin/go-krypto/internal/randutil"

	"github.com/RyuaNerin/elliptic2/nist"
	"github.com/stretchr/testify/require"
)

const (
	svtP224SHA224 = "15. EC-KCDSA/EC-KCDSA_(P-224)(SHA-224)_SVT"
	svtP224SHA256 = "" // "15. EC-KCDSA/EC-KCDSA_(P-224)(SHA-256)_SVT" // invalid test vector file
	svtP256SHA224 = "15. EC-KCDSA/EC-KCDSA_(P-256)(SHA-224)_SVT"
	svtP256SHA256 = "15. EC-KCDSA/EC-KCDSA_(P-256)(SHA-256)_SVT"

	svtB233SHA224 = "15. EC-KCDSA/EC-KCDSA_(B-233)(SHA-224)_SVT"
	svtB233SHA256 = "" //  "15. EC-KCDSA/EC-KCDSA_(B-233)(SHA-256)_SVT" // invalid test vector file
	svtB283SHA224 = "15. EC-KCDSA/EC-KCDSA_(B-283)(SHA-224)_SVT"
	svtB283SHA256 = "15. EC-KCDSA/EC-KCDSA_(B-283)(SHA-256)_SVT"

	svtK233SHA224 = "15. EC-KCDSA/EC-KCDSA_(K-233)(SHA-224)_SVT"
	svtK233SHA256 = "" // "15. EC-KCDSA/EC-KCDSA_(K-233)(SHA-256)_SVT" // invalid test vector file
	svtK283SHA224 = "15. EC-KCDSA/EC-KCDSA_(K-283)(SHA-224)_SVT"
	svtK283SHA256 = "15. EC-KCDSA/EC-KCDSA_(K-283)(SHA-256)_SVT"
)

type curveAndHashPair struct {
	name        string
	curve       elliptic.Curve
	hash        hash.Hash
	svtFilename string
}

var (
	curves = []elliptic.Curve{
		elliptic.P224(),
		elliptic.P256(),
		elliptic.P384(),
		elliptic.P521(),
		nist.B233(),
		nist.B283(),
		nist.K233(),
		nist.K283(),
	}
	curveAndHashPairs = []curveAndHashPair{
		{"P224_SHA224", elliptic.P224(), sha256.New224(), svtP224SHA224},
		{"P224_SHA256", elliptic.P224(), sha256.New(), svtP224SHA256},
		{"P256_SHA224", elliptic.P256(), sha256.New224(), svtP256SHA224},
		{"P256_SHA256", elliptic.P256(), sha256.New(), svtP256SHA256},

		{"B233_SHA224", nist.B233(), sha256.New224(), svtB233SHA224},
		{"B233_SHA256", nist.B233(), sha256.New(), svtB233SHA256},
		{"B283_SHA224", nist.B283(), sha256.New224(), svtB283SHA224},
		{"B283_SHA256", nist.B283(), sha256.New(), svtB283SHA256},

		{"K233_SHA224", nist.K233(), sha256.New224(), svtK233SHA224},
		{"K233_SHA256", nist.K233(), sha256.New(), svtK233SHA256},
		{"K283_SHA224", nist.K283(), sha256.New224(), svtK283SHA224},
		{"K283_SHA256", nist.K283(), sha256.New(), svtK283SHA256},
	}
)

func Test_ECKCDSA_SVT(t *testing.T) {
	test := func(pair curveAndHashPair) func(t *testing.T) {
		return func(t *testing.T) {
			if pair.svtFilename == "" {
				t.Skipf("no SVT file for %s", pair.name)
			}
			cavp.TestECKCDSASvt(t, pair.svtFilename, pair.curve, pair.hash)
		}
	}

	for _, pair := range curveAndHashPairs {
		t.Run(pair.name, test(pair))
	}
}

func Test_SignVerify_With_BadPublicKey(t *testing.T) {
	for idx, tc := range testCase_TTAK {
		tc2 := testCase_TTAK[(idx+1)%len(testCase_TTAK)]

		key := eckcdsa.PublicKey{
			Curve: tc2.curve,
			X:     tc2.Qx,
			Y:     tc2.Qy,
		}

		ok := eckcdsa.Verify(&key, tc.hash, tc.M, tc.R, tc.S)
		require.Falsef(t, ok, "%d: Verify unexpected success with bad public key", idx)
	}
}

func Test_ECKCDSA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
	}

	for _, pair := range curveAndHashPairs {
		t.Run(pair.name, testKCDSA(pair.curve, pair.hash))
	}
}

func Test_Signing_With_DegenerateKeys(t *testing.T) {
	// Signing with degenerate private keys should not cause an infinite
	// loop.
	badKeys := []struct {
		d, y, x string
	}{
		{"0000", "0001", "0101"},
		{"0100", "0f0f", "1010"},
	}

	for i, test := range badKeys {
		priv := eckcdsa.PrivateKey{
			PublicKey: eckcdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     internal.HI(test.x),
				Y:     internal.HI(test.y),
			},
			D: internal.HI(test.d),
		}

		data := []byte("testing")
		_, _, err := eckcdsa.Sign(rand.Reader, &priv, sha256.New(), data)
		require.Error(t, err, "#%d: expected error signing with degenerate key", i)
	}
}

func Test_Signing_With_ZeroPrivateKey(t *testing.T) {
	invalidD := []*big.Int{
		new(big.Int),
		new(big.Int).Set(elliptic.P256().Params().N),
	}

	priv, err := eckcdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	for _, d := range invalidD {
		priv.D = d

		_, _, err = eckcdsa.Sign(rand.Reader, priv, sha256.New(), []byte("input"))
		require.ErrorIs(t, err, eckcdsa.ErrInvalidPrivateKey)
	}
}

func Test_SignVerify_HashLongerThanTwoCoordinates(t *testing.T) {
	// SHA-512 output (64 bytes) is longer than xQ ‖ yQ of P-224 (56 bytes) and
	// B-233/K-233 (60 bytes), so cQ is zero-padded.
	for _, curve := range []elliptic.Curve{elliptic.P224(), nist.B233(), nist.K233()} {
		priv, err := eckcdsa.GenerateKey(curve, rand.Reader)
		require.NoError(t, err)

		data := []byte("testing")
		r, s, err := eckcdsa.Sign(rand.Reader, priv, sha512.New(), data)
		require.NoError(t, err)
		require.True(t, eckcdsa.Verify(&priv.PublicKey, sha512.New(), data, r, s), curve.Params().Name)
	}
}

func testKCDSA(
	curve elliptic.Curve,
	h hash.Hash,
) func(t *testing.T) {
	return func(t *testing.T) {
		rnd, _ := randutil.NewTestRand()

		priv, err := eckcdsa.GenerateKey(curve, rnd)
		require.NoError(t, err)

		testSignAndVerify(t, priv, h, rnd)
		testSignAndVerifyASN1(t, priv, h, rnd)
	}
}

func testSignAndVerify(
	t *testing.T,
	priv *eckcdsa.PrivateKey,
	h hash.Hash,
	rnd io.Reader,
) {
	data := []byte("testing")
	dataClone := slices.Clone(data)
	privClone := priv.Clone()
	r, s, err := eckcdsa.Sign(rnd, priv, h, data)
	require.NoError(t, err)
	require.Equal(t, dataClone, data, "Sign modified the input data")
	require.True(t, privClone.Equal(priv), "Sign modified the private key")

	rClone := new(big.Int).Set(r)
	sClone := new(big.Int).Set(s)
	ok := eckcdsa.Verify(&priv.PublicKey, h, data, r, s)
	require.True(t, ok, "Verify failed")
	require.Equal(t, rClone.Text(16), r.Text(16), "Verify modified the signature R")
	require.Equal(t, sClone.Text(16), s.Text(16), "Verify modified the signature S")
	require.True(t, privClone.PublicKey.Equal(&priv.PublicKey), "Verify modified the public key")

	data[0] ^= 0xff
	ok = eckcdsa.Verify(&priv.PublicKey, h, data, r, s)
	require.False(t, ok, "Verify always works!")
}

func testSignAndVerifyASN1(t *testing.T, priv *eckcdsa.PrivateKey, h hash.Hash, rnd io.Reader) {
	data := []byte("testing")

	sig, err := eckcdsa.SignASN1(rnd, priv, h, data)
	require.NoError(t, err, "error signing")

	ok := eckcdsa.VerifyASN1(&priv.PublicKey, h, data, sig)
	require.True(t, ok, "VerifyASN1 failed")

	data[0] ^= 0xff
	ok = eckcdsa.VerifyASN1(&priv.PublicKey, h, data, sig)
	require.False(t, ok, "VerifyASN1 always works!")
}
