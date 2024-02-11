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

			kcdsa.GenerateParameters(&p1p.Parameters, rand.Reader, size)
			kcdsa.GenerateKey(&p1p, rand.Reader)

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
	t.Run("KCDSA-TTAK", func(t *testing.T) {
		for _, size := range sizeList {
			var p1p kcdsa.PrivateKey

			kcdsa.GenerateParameters(&p1p.Parameters, rand.Reader, size)
			kcdsa.GenerateKey(&p1p, rand.Reader)

			p1 := &p1p.PublicKey

			der, err := MarshalPKIXPublicKey(p1)
			if err != nil {
				t.Error(err)
				return
			}

			p2r, err := ParsePKIXPublicKey(der)
			if err != nil {
				t.Error(err)
				return
			}

			p2, ok := p2r.(*kcdsa.PublicKey)
			if !ok {
				t.Error("type error")
				return
			}

			if !p1.Equal(p2) || !p1.TTAKParams.Equal(p2.TTAKParams) {
				t.Error("not equals!")
				return
			}
		}
	})
}
