package kcdsa

import (
	"crypto/dsa"
	"crypto/rand"
	"testing"
)

var (
	sizeDSA2KCDSA = map[dsa.ParameterSizes]ParameterSizes{
		dsa.L2048N224: L2048N224SHA224,
		dsa.L2048N256: L2048N256SHA256,
		dsa.L3072N256: L3072N256SHA256,
	}
)

func Test_DSA_TO_KCDSA(t *testing.T) {
	for sz := range sizeDSA2KCDSA {
		for {
			var expect dsa.PrivateKey
			dsa.GenerateParameters(&expect.Parameters, rand.Reader, sz)
			dsa.GenerateKey(&expect, rand.Reader)

			cvt, err := FromDSA(&expect)
			if err != nil {
				continue
			}

			answer := cvt.ToDSA()

			equals := true &&
				expect.X.Cmp(answer.X) == 0 &&
				expect.Y.Cmp(answer.Y) == 0 &&
				expect.P.Cmp(answer.P) == 0 &&
				expect.Q.Cmp(answer.Q) == 0 &&
				expect.G.Cmp(answer.G) == 0

			if !equals {
				t.Fail()
				return
			}

			break
		}
	}
}

func Test_KCDSA_TO_DSA(t *testing.T) {
	for _, sz := range sizeDSA2KCDSA {
		var expect PrivateKey
		GenerateParameters(&expect.Parameters, rand.Reader, sz)
		GenerateKey(&expect, rand.Reader)

		cvt := expect.ToDSA()

		answer, err := FromDSA(cvt)
		if err != nil {
			t.Error(err)
			return
		}

		if !expect.Equal(answer) {
			t.Fail()
			return
		}
	}
}
