package kcdsattak

import (
	"bufio"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/kcdsa"

	. "github.com/RyuaNerin/testingutil"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

var (
	as = []CipherSize{
		{Name: "L2048 N224 SHA224", Size: int(kcdsa.L2048N224SHA224)},
		{Name: "L2048 N224 SHA256", Size: int(kcdsa.L2048N224SHA256)},
		{Name: "L2048 N256 SHA256", Size: int(kcdsa.L2048N256SHA256)},
		{Name: "L3072 N256 SHA256", Size: int(kcdsa.L3072N256SHA256)},
	}
)

type testCase struct {
	Sizes kcdsa.ParameterSizes

	M []byte

	Seed_ []byte
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

func testKCDSA(
	t *testing.T,
	sizes kcdsa.ParameterSizes,
	L, N int,
	gp func(params *Parameters, sizes kcdsa.ParameterSizes) error,
	gk func(priv *PrivateKey, sizes kcdsa.ParameterSizes) error,
) {
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

func testSignAndVerify(t *testing.T, i int, priv *PrivateKey, sizes kcdsa.ParameterSizes) {
	data := []byte("testing")
	r, s, err := Sign(rand.Reader, priv, data, sizes)
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
