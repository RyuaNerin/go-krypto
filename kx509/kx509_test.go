package kx509

import (
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

var (
	curveList = []elliptic.Curve{
		elliptic.P256(),
		elliptic.P224(),
		elliptic.P384(),
		elliptic.P521(),
	}
	sizeList = []kcdsa.ParameterSizes{
		kcdsa.L2048N224SHA256,

		kcdsa.L2048N224SHA224,
		kcdsa.L2048N256SHA256,
		kcdsa.L3072N256SHA256,
	}

	testCases []struct {
		PrivateKey interface{}
		Marshaled  []byte
	}
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

			_ = kcdsa.GenerateParameters(&p1.Parameters, rand.Reader, size)
			_ = kcdsa.GenerateKey(&p1, rand.Reader)

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
}

func TestMarshalAndParseMarshalECPrivateKey(t *testing.T) {
	for _, curve := range curveList {
		p1, _ := eckcdsa.GenerateKey(curve, rand.Reader)

		der, err := MarshalECPrivateKey(p1)
		if err != nil {
			t.Error(err)
			return
		}

		p2, err := ParseECPrivateKey(der)
		if err != nil {
			t.Error(err)
			return
		}

		if !p1.Equal(p2) {
			t.Error("not equals!")
			return
		}
	}
}
