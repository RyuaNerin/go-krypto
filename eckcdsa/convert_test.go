package eckcdsa

import (
	"crypto/ecdsa"
	"crypto/rand"
	"testing"
)

func Test_ECDSA_TO_ECKCDSA(t *testing.T) {
	for _, curve := range curveList {
		for {
			expect, _ := ecdsa.GenerateKey(curve, rand.Reader)

			cvt, err := FromECDSA(expect)
			if err != nil {
				continue
			}

			answer := cvt.ToECDSA()

			if !expect.Equal(answer) {
				t.Fail()
				return
			}

			break
		}
	}
}

func Test_ECKCDSA_TO_ECDSA(t *testing.T) {
	for _, curve := range curveList {
		expect, _ := GenerateKey(curve, rand.Reader)

		cvt := expect.ToECDSA()

		answer, err := FromECDSA(cvt)
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