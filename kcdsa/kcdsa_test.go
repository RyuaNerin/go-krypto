package kcdsa

import (
	"bufio"
	"crypto/rand"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	. "github.com/RyuaNerin/testingutil"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

var (
	as = []CipherSize{
		{Name: "L2048 N224 SHA224", Size: int(L2048N224SHA224)},
		{Name: "L2048 N224 SHA256", Size: int(L2048N224SHA256)},
		{Name: "L2048 N256 SHA256", Size: int(L2048N256SHA256)},
		{Name: "L3072 N256 SHA256", Size: int(L3072N256SHA256)},
	}
)

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

		if _, _, err := Sign(rand.Reader, &priv, L2048N224SHA224.Hash(), msg); err == nil {
			t.Errorf("#%d: unexpected success", i)
			return
		}
	}
}
func Test_KCDSA(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping parameter generation test in short mode")
	}

	gp := func(params *Parameters, sizes ParameterSizes) error {
		return GenerateParameters(params, rand.Reader, sizes)
	}
	gk := func(priv *PrivateKey, sizes ParameterSizes) error {
		return GenerateKey(priv, rand.Reader)
	}

	testKCDSA(t, L2048N224SHA224, 2048, 224, gp, gk)
	testKCDSA(t, L2048N224SHA256, 2048, 224, gp, gk)
	testKCDSA(t, L2048N256SHA256, 2048, 256, gp, gk)
	testKCDSA(t, L3072N256SHA256, 3072, 256, gp, gk)
}

func testKCDSA(t *testing.T, sizes ParameterSizes, L, N int, gp func(params *Parameters, sizes ParameterSizes) error, gk func(priv *PrivateKey, sizes ParameterSizes) error) {
	var priv PrivateKey
	params := &priv.Parameters

	err := gp(params, sizes)
	if err != nil {
		t.Errorf("%d: %s", int(sizes), err)
		return
	}

	if params.P.BitLen() > L {
		t.Errorf("%d: params.BitLen got:%d want:%d", int(sizes), params.P.BitLen(), L)
		return
	}

	if params.Q.BitLen() > N {
		t.Errorf("%d: q.BitLen got:%d want:%d", int(sizes), params.Q.BitLen(), L)
		return
	}

	err = gk(&priv, sizes)
	if err != nil {
		t.Errorf("error generating key: %s", err)
		return
	}

	testSignAndVerify(t, int(sizes), &priv, sizes)
}

func testSignAndVerify(t *testing.T, i int, priv *PrivateKey, sizes ParameterSizes) {
	data := []byte("testing")
	r, s, err := Sign(rand.Reader, priv, sizes.Hash(), data)
	if err != nil {
		t.Errorf("%d: error signing: %s", i, err)
		return
	}

	ok := Verify(&priv.PublicKey, sizes.Hash(), data, r, s)
	if !ok {
		t.Errorf("%d: Verify failed", i)
		return
	}
}
