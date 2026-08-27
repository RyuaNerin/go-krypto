package kcdsa_test

import (
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/tzeroize"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

func Test_Zeroize(t *testing.T) {
	getKey := func() *kcdsa.PrivateKey {
		return &kcdsa.PrivateKey{
			PublicKey: kcdsa.PublicKey{
				Parameters: kcdsa.Parameters{
					P: big.NewInt(1),
					Q: big.NewInt(2),
					G: big.NewInt(3),
					GenParameters: kcdsa.GenerationParameters{
						J:     big.NewInt(4),
						Seed:  []byte{1, 2, 3},
						Count: 5,
					},
				},
				Y: big.NewInt(6),
			},
			X: big.NewInt(7),
		}
	}

	t.Run("PrivateKey", func(t *testing.T) {
		tzeroize.TestZeroize(t, getKey())
	})
	t.Run("PublicKey", func(t *testing.T) {
		tzeroize.TestZeroize(t, &getKey().PublicKey)
	})
}
