package eckcdsa

import (
	"crypto/elliptic"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal/testing/tzeroize"
)

func Test_Zeroize(t *testing.T) {
	t.Run("PrivateKey", func(t *testing.T) {
		privKey := PrivateKey{
			PublicKey: PublicKey{
				Curve: elliptic.P256(),
				X:     big.NewInt(1),
				Y:     big.NewInt(2),
			},
			D: big.NewInt(3),
		}
		tzeroize.TestZeroize(t, &privKey)
	})
	t.Run("PublicKey", func(t *testing.T) {
		pubKey := PublicKey{
			Curve: elliptic.P256(),
			X:     big.NewInt(1),
			Y:     big.NewInt(2),
		}
		tzeroize.TestZeroize(t, &pubKey)
	})
}
