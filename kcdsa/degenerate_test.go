package kcdsa_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"
	kcdsa_internal "github.com/RyuaNerin/go-krypto/internal/kcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"

	"github.com/stretchr/testify/require"
)

// With Y = 0 or G = P, W' = Y^S' G^E' mod P is 0 for every message, so R' =
// trunc(Hash(0…0)) would verify for anything.
func Test_Verify_RejectsDegenerateKeys(t *testing.T) {
	tc := testCaseTTAK[0]
	domain, ok := kcdsa_internal.GetDomain(int(tc.Sizes))
	require.True(t, ok)

	h := domain.NewHash()
	h.Write(make([]byte, internal.BitsToBytes(tc.P.BitLen())))
	forgedR := new(big.Int).SetBytes(internal.RightMost(h.Sum(nil), domain.B))
	forgedS := big.NewInt(1)

	newPub := func(g, y *big.Int) *kcdsa.PublicKey {
		return &kcdsa.PublicKey{
			Parameters: kcdsa.Parameters{P: tc.P, Q: tc.Q, G: g},
			Y:          y,
		}
	}
	pMinusOne := new(big.Int).Sub(tc.P, big.NewInt(1))

	for name, pub := range map[string]*kcdsa.PublicKey{
		"Y=0":   newPub(tc.G, big.NewInt(0)),
		"Y=1":   newPub(tc.G, big.NewInt(1)),
		"Y=P-1": newPub(tc.G, pMinusOne),
		"G=P":   newPub(tc.P, tc.Y),
		"G=1":   newPub(big.NewInt(1), tc.Y),
	} {
		require.False(t, kcdsa.Verify(pub, tc.Sizes, []byte("any message"), forgedR, forgedS), name)
	}
}

func Test_Verify_RejectsMismatchedSizes(t *testing.T) {
	tc := testCaseTTAK[0]
	pub := &kcdsa.PublicKey{
		Parameters: kcdsa.Parameters{P: tc.P, Q: tc.Q, G: tc.G},
		Y:          tc.Y,
	}
	require.True(t, kcdsa.Verify(pub, tc.Sizes, tc.M, tc.R, tc.S))

	for _, sizes := range []kcdsa.ParameterSizes{kcdsa.A2048B224SHA224, kcdsa.A2048B256SHA256, kcdsa.A3072B256SHA256} {
		if sizes == tc.Sizes {
			continue
		}
		require.False(t, kcdsa.Verify(pub, sizes, tc.M, tc.R, tc.S))
	}
}

func Test_Sign_RejectsMismatchedSizes(t *testing.T) {
	// A key whose P is shorter or longer than the domain used to panic or produce
	// signatures that never verify.
	tc := testCaseTTAK[0]
	priv := &kcdsa.PrivateKey{
		PublicKey: kcdsa.PublicKey{
			Parameters: kcdsa.Parameters{P: tc.P, Q: tc.Q, G: tc.G},
			Y:          tc.Y,
		},
		X: tc.X,
	}

	sizes := kcdsa.A3072B256SHA256
	if tc.Sizes == sizes {
		sizes = kcdsa.A2048B224SHA224
	}
	_, _, err := kcdsa.Sign(rand.Reader, priv, sizes, tc.M)
	require.ErrorIs(t, err, kcdsa.ErrInvalidPublicKey)
}
