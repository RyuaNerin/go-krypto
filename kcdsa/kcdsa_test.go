package kcdsa

import (
	"bufio"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"

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

	testKCDSA(t, A2048B224SHA224, 2048, 224)
	testKCDSA(t, A2048B224SHA256, 2048, 224)
	testKCDSA(t, A2048B256SHA256, 2048, 256)
	testKCDSA(t, A3072B256SHA256, 3072, 256)
}

func testKCDSA(t *testing.T, sizes ParameterSizes, l, n int) {
	var priv PrivateKey
	params := &priv.Parameters

	err := GenerateParameters(params, rand.Reader, sizes)
	if err != nil {
		t.Errorf("%d: %s", int(sizes), err)
		return
	}

	if params.P.BitLen() > l {
		t.Errorf("%d: params.BitLen got:%d want:%d", int(sizes), params.P.BitLen(), l)
		return
	}

	if params.Q.BitLen() > n {
		t.Errorf("%d: q.BitLen got:%d want:%d", int(sizes), params.Q.BitLen(), l)
		return
	}

	err = GenerateKey(&priv, rand.Reader)
	if err != nil {
		t.Errorf("error generating key: %s", err)
		return
	}

	testSignAndVerify(t, int(sizes), &priv, sizes)
}

func testSignAndVerify(t *testing.T, i int, priv *PrivateKey, sizes ParameterSizes) {
	data := []byte("testing")
	r, s, err := Sign(rand.Reader, priv, sizes, data)
	if err != nil {
		t.Errorf("%d: error signing: %s", i, err)
		return
	}

	ok := Verify(&priv.PublicKey, sizes, data, r, s)
	if !ok {
		t.Errorf("%d: Verify failed", i)
		return
	}
}

func testVerify(t *testing.T, testCases []testCase) {
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
