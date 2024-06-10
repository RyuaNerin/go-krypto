package kcdsa

import (
	"bufio"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/kcdsa"

	. "github.com/RyuaNerin/testingutil"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

type testCase struct {
	Sizes ParameterSizes

	M []byte

	Seedb []byte
	J     *big.Int
	Count int
	P, Q  *big.Int

	H []byte
	G *big.Int

	XKEY []byte
	X    *big.Int
	Y, Z *big.Int

	KKEY *big.Int
	R    *big.Int
	S    *big.Int

	Fail bool
}

var as = []CipherSize{
	{Name: "A2048 B224 SHA224", Size: int(A2048B224SHA224)},
	{Name: "A2048 B224 SHA256", Size: int(A2048B224SHA256)},
	{Name: "A2048 B256 SHA256", Size: int(A2048B256SHA256)},
	{Name: "A3072 B256 SHA256", Size: int(A3072B256SHA256)},
}

func Test_SignVerify_With_BadPublicKey(t *testing.T) {
	for idx, tc := range testCaseTTAK {
		tc2 := testCaseTTAK[(idx+1)%len(testCaseTTAK)]

		pub := PublicKey{
			Parameters: Parameters{
				P: tc2.P,
				Q: tc2.Q,
				G: tc2.G,
			},
			Y: tc2.Y,
		}

		ok := Verify(&pub, tc.Sizes, tc.M, tc.R, tc.S)
		if ok {
			t.Errorf("Verify unexpected success with non-existent mod inverse of Q")
			return
		}
	}
}

func Test_Signing_With_DegenerateKeys(t *testing.T) {
	badKeys := []struct {
		p, q, g, y, x string
	}{
		{"00", "01", "00", "00", "00"},
		{"01", "ff", "00", "00", "00"},
	}

	msg := []byte("testing")
	for i, test := range badKeys {
		priv := PrivateKey{
			PublicKey: PublicKey{
				Parameters: Parameters{
					P: internal.HI(test.p),
					Q: internal.HI(test.q),
					G: internal.HI(test.g),
				},
				Y: internal.HI(test.y),
			},
			X: internal.HI(test.x),
		}

		if _, _, err := Sign(rand.Reader, &priv, A2048B224SHA224, msg); err == nil {
			t.Errorf("#%d: unexpected success", i)
			return
		}
	}
}

func Test_KCDSA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
	}

	t.Run("A2048B224SHA224", testKCDSA(A2048B224SHA224))
	t.Run("A2048B224SHA256", testKCDSA(A2048B224SHA256))
	t.Run("A2048B256SHA256", testKCDSA(A2048B256SHA256))
	t.Run("A3072B256SHA256", testKCDSA(A3072B256SHA256))
	t.Run("A1024B160HAS160", testKCDSA(A1024B160HAS160))
}

func testKCDSA(sizes ParameterSizes) func(*testing.T) {
	return func(t *testing.T) {
		d, ok := kcdsa.GetDomain(int(sizes))
		if !ok {
			t.Errorf("domain not found")
			return
		}

		var priv PrivateKey
		params := &priv.Parameters

		err := GenerateParameters(params, rand.Reader, sizes)
		if err != nil {
			t.Error(err)
			return
		}

		if params.P.BitLen() > d.A {
			t.Errorf("params.BitLen got:%d want:%d", params.P.BitLen(), d.A)
			return
		}

		if params.Q.BitLen() > d.B {
			t.Errorf("q.BitLen got:%d want:%d", params.Q.BitLen(), d.B)
			return
		}

		err = GenerateKey(&priv, rand.Reader)
		if err != nil {
			t.Errorf("error generating key: %s", err)
			return
		}

		testSignAndVerify(t, &priv, sizes)
		testSignAndVerifyASN1(t, &priv, sizes)
	}
}

func testSignAndVerify(t *testing.T, priv *PrivateKey, sizes ParameterSizes) {
	data := []byte("testing")
	r, s, err := Sign(rand.Reader, priv, sizes, data)
	if err != nil {
		t.Errorf("error signing: %s", err)
		return
	}

	ok := Verify(&priv.PublicKey, sizes, data, r, s)
	if !ok {
		t.Error("Verify failed")
		return
	}

	data[0] ^= 0xff
	if Verify(&priv.PublicKey, sizes, data, r, s) {
		t.Errorf("Verify always works!")
	}
}

func testSignAndVerifyASN1(t *testing.T, priv *PrivateKey, sizes ParameterSizes) {
	data := []byte("testing")
	sig, err := SignASN1(rand.Reader, priv, sizes, data)
	if err != nil {
		t.Errorf("error signing: %s", err)
		return
	}

	if !VerifyASN1(&priv.PublicKey, sizes, data, sig) {
		t.Errorf("VerifyASN1 failed")
	}

	data[0] ^= 0xff
	if VerifyASN1(&priv.PublicKey, sizes, data, sig) {
		t.Errorf("VerifyASN1 always works!")
	}
}

func verifyTestCases(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		pub := PublicKey{
			Parameters: Parameters{
				P: tc.P,
				Q: tc.Q,
				G: tc.G,
			},
			Y: tc.Y,
		}

		ok := Verify(&pub, tc.Sizes, tc.M, tc.R, tc.S)
		if ok == tc.Fail {
			t.Errorf("verify failed")
			return
		}
	}
}
