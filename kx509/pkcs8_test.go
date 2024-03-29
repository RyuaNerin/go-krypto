package kx509

import (
	"crypto/rand"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

func TestMarshalAndParsePKCS8PrivateKey(t *testing.T) {
	t.Run("EC-KCDSA", func(t *testing.T) {
		for _, curve := range curveList {
			p1, _ := eckcdsa.GenerateKey(curve, rand.Reader)

			der, err := MarshalPKCS8PrivateKey(p1)
			if err != nil {
				t.Error(err)
				return
			}

			p2, err := ParsePKCS8PrivateKey(der)
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
			var p1 kcdsa.PrivateKey

			kcdsa.GenerateParameters(&p1.Parameters, rand.Reader, size)
			kcdsa.GenerateKey(&p1, rand.Reader)

			der, err := MarshalPKCS8PrivateKey(&p1)
			if err != nil {
				t.Error(err)
				return
			}

			p2, err := ParsePKCS8PrivateKey(der)
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
			var p1 kcdsa.PrivateKey

			kcdsa.GenerateParameters(&p1.Parameters, rand.Reader, size)
			kcdsa.GenerateKey(&p1, rand.Reader)

			der, err := MarshalPKCS8PrivateKey(&p1)
			if err != nil {
				t.Error(err)
				return
			}

			p2r, err := ParsePKCS8PrivateKey(der)
			if err != nil {
				t.Error(err)
				return
			}

			p2, ok := p2r.(*kcdsa.PrivateKey)
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
