package eckcdsa

import (
	"bufio"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"

	"github.com/RyuaNerin/elliptic2/nist"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

type testCase struct {
	curve elliptic.Curve
	hash  hash.Hash

	D  *big.Int
	Qx *big.Int
	Qy *big.Int

	K *big.Int

	M []byte
	R *big.Int
	S *big.Int

	Fail bool
}

var (
	p224     = elliptic.P224()
	p256     = elliptic.P256()
	secp224r = elliptic.P224()
	secp256r = elliptic.P256()

	b233     = nist.B233()
	k233     = nist.K233()
	b283     = nist.B283()
	k283     = nist.K283()
	sect233r = nist.B233()
	sect233k = nist.K233()
	sect283r = nist.B283()
	sect283k = nist.K283()

	hashSHA256     = sha256.New()
	hashSHA256_224 = sha256.New224()
)

func testVerify(t *testing.T, testCases []testCase, curve elliptic.Curve, hash hash.Hash) {
	for idx, tc := range testCases {
		key := PublicKey{
			Curve: curve,
			X:     tc.Qx,
			Y:     tc.Qy,
		}

		ok := Verify(&key, hash, tc.M, tc.R, tc.S)
		if ok == tc.Fail {
			t.Errorf("%d: Verify failed, got:%v want:%v\nM=%s", idx, ok, !tc.Fail, hex.EncodeToString(tc.M))
			return
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

		R, S, err := signUsingK(tc.K, &key, tc.hash, tc.M)
		if err != nil {
			t.Errorf("%d: error signing: %s", idx, err)
			return
		}

		if R.Cmp(tc.R) != 0 || S.Cmp(tc.S) != 0 {
			t.Errorf("%d: error signing: (r, s)", idx)
			return
		}

		ok := Verify(&key.PublicKey, tc.hash, tc.M, tc.R, tc.S)
		if ok == tc.Fail {
			t.Errorf("%d: Verify failed, got:%v want:%v\nM=%s", idx, ok, !tc.Fail, hex.EncodeToString(tc.M))
			return
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
			t.Errorf("%d: Verify unexpected success with non-existent mod inverse of Q", idx)
			return
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

	testKCDSA(t, "B233_SHA256_224", b233, hashSHA256_224)
	testKCDSA(t, "B233_SHA256_256", b233, hashSHA256)

	testKCDSA(t, "B283_SHA256_224", b283, hashSHA256_224)
	testKCDSA(t, "B283_SHA256_256", b283, hashSHA256)

	testKCDSA(t, "K233_SHA256_224", k233, hashSHA256_224)
	testKCDSA(t, "K233_SHA256_256", k233, hashSHA256)

	testKCDSA(t, "K283_SHA256_224", k283, hashSHA256_224)
	testKCDSA(t, "K283_SHA256_256", k283, hashSHA256)
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
			return
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
		return
	}
}
