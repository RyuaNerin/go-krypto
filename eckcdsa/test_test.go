package eckcdsa

import (
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
)

type testCase struct {
	M []byte

	curve elliptic.Curve
	hash  hash.Hash

	D  *big.Int
	Qx *big.Int
	Qy *big.Int

	K *big.Int

	R *big.Int
	S *big.Int

	Fail bool
}

var (
	p224     = elliptic.P224()
	p256     = elliptic.P256()
	secp224r = elliptic.P224() // Also known as: P-224, wap-wsg-idm-ecid-wtls12, ansip224r1
	secp256r = elliptic.P256() // Also known as: P-256, prime256v1
	/**
	sect233r = elliptic.Curve(nil) // Also known as: B-233, wap-wsg-idm-ecid-wtls11, ansit233r1
	sect233k = elliptic.Curve(nil) // Also known as: K-233, wap-wsg-idm-ecid-wtls10, ansit233k1
	sect283r = elliptic.Curve(nil) // Also known as: B-283, ansit283r1
	*/

	hashSHA256     = sha256.New()
	hashSHA256_224 = sha256.New224()
)

func testVerify(t *testing.T, testCases []testCase) {
	for idx, tc := range testCases {
		key := PublicKey{
			Curve: tc.curve,
			X:     tc.Qx,
			Y:     tc.Qy,
		}

		ok := Verify(&key, tc.hash, tc.M, tc.R, tc.S)
		if ok == tc.Fail {
			t.Errorf("%d: Verify failed, got:%v want:%v\nM=%s", idx, ok, !tc.Fail, hex.EncodeToString(tc.M))
		}
	}
}

func testSignVerify(t *testing.T, testCases []testCase) {
	for idx, tc := range testCases {
		key := PrivateKey{
			PublicKey: PublicKey{
				Curve: tc.curve,
				X:     tc.Qx,
				Y:     tc.Qy,
			},
			D: tc.D,
		}

		R, S, err := SignWithK(tc.K, &key, tc.hash, tc.M)
		if err != nil {
			t.Errorf("%d: error signing: %s", idx, err)
		}

		if R.Cmp(tc.R) != 0 || S.Cmp(tc.S) != 0 {
			t.Errorf("%d: error signing: (r, s)", idx)
		}

		ok := Verify(&key.PublicKey, tc.hash, tc.M, tc.R, tc.S)
		if ok == tc.Fail {
			t.Errorf("%d: Verify failed, got:%v want:%v\nM=%s", idx, ok, !tc.Fail, hex.EncodeToString(tc.M))
		}
	}
}

func Test_SignVerify_With_BadPublicKey(t *testing.T) {
	for idx, tc := range testCase_TTAK {
		tc2 := testCase_TTAK[(idx+1)%len(testCase_TTAK)]

		key := PublicKey{
			Curve: tc2.curve,
			X:     tc2.Qx,
			Y:     tc2.Qy,
		}

		ok := Verify(&key, tc.hash, tc.M, tc.R, tc.S)
		if ok {
			t.Errorf("Verify unexpected success with non-existent mod inverse of Q")
		}
	}
}

func Test_ECKCDSA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
	}

	testKCDSA(t, "P224_SHA256_224", p224, hashSHA256_224)
	testKCDSA(t, "P224_SHA256_256", p224, hashSHA256)

	testKCDSA(t, "P256_SHA256_224", p256, hashSHA256_224)
	testKCDSA(t, "P256_SHA256_256", p256, hashSHA256)
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
		priv := PrivateKey{
			PublicKey: PublicKey{
				Curve: secp224r,
				X:     internal.HI(test.x),
				Y:     internal.HI(test.y),
			},
			D: internal.HI(test.d),
		}

		data := []byte("testing")
		if _, _, err := Sign(rand.Reader, &priv, sha256.New(), data); err == nil {
			t.Errorf("#%d: unexpected success", i)
		}
	}
}

func testKCDSA(
	t *testing.T,
	name string,
	curve elliptic.Curve,
	h hash.Hash,
) {
	priv, err := GenerateKey(curve, rand.Reader)
	if err != nil {
		t.Errorf("%s: error generating key: %s", name, err)
		return
	}

	testSignAndVerify(t, name, priv, h)
}

func testSignAndVerify(
	t *testing.T,
	name string,
	priv *PrivateKey,
	h hash.Hash,
) {
	data := []byte("testing")
	r, s, err := Sign(rand.Reader, priv, h, data)
	if err != nil {
		t.Errorf("%s: error signing: %s", name, err)
		return
	}

	ok := Verify(&priv.PublicKey, h, data, r, s)
	if !ok {
		t.Errorf("%s: Verify failed", name)
	}
}
