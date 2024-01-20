package kx509

import (
	"crypto/rand"
	"testing"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
)

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
