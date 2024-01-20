package kx509

import (
	"crypto/rand"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

func TestMarshalAndParsePKIXPublicKey(t *testing.T) {
	t.Run("EC-KCDSA", func(t *testing.T) {
		for _, curve := range curveList {
			p1p, _ := eckcdsa.GenerateKey(curve, rand.Reader)
			p1 := &p1p.PublicKey

			der, err := MarshalPKIXPublicKey(p1)
			if err != nil {
				t.Error(err)
				return
			}

			p2, err := ParsePKIXPublicKey(der)
			if err != nil {
				t.Error(err)
				return
			}

			if !p1.Equal(p2) {
				t.Error("not equals!")
				return
			}
		}
	})
	t.Run("KCDSA", func(t *testing.T) {
		for _, size := range sizeList {
			var p1p kcdsa.PrivateKey

			_ = kcdsa.GenerateParameters(&p1p.Parameters, rand.Reader, size)
			_ = kcdsa.GenerateKey(&p1p, rand.Reader)

			p1 := &p1p.PublicKey

			der, err := MarshalPKIXPublicKey(p1)
			if err != nil {
				t.Error(err)
				return
			}

			p2, err := ParsePKIXPublicKey(der)
			if err != nil {
				t.Error(err)
				return
			}

			if !p1.Equal(p2) {
				t.Error("not equals!")
				return
			}
		}
	})
}
