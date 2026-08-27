package eckcdsa_test

import (
	"crypto/elliptic"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/stretchr/testify/require"
)

func TestVerifyRejectsNilSignature(t *testing.T) {
	curve := elliptic.P256()
	pub := &eckcdsa.PublicKey{Curve: curve, X: curve.Params().Gx, Y: curve.Params().Gy}

	testCases := []struct {
		Name string
		R    *big.Int
		S    *big.Int
	}{
		{"nil nil", nil, nil},
		{"nil r", nil, big.NewInt(1)},
		{"nil s", big.NewInt(1), nil},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ok := eckcdsa.Verify(pub, sha256.New(), []byte("message"), tc.R, tc.S)
			require.True(t, !ok)
		})
	}
}
