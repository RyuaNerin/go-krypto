package internal

import (
	"bufio"
	"crypto/rand"
	"errors"
	"math/big"
	"testing"
)

var rnd = bufio.NewReaderSize(rand.Reader, 1<<15)

type testCase struct {
	Sizes _ParameterSizes

	M []byte

	Seed_ []byte
	J     *big.Int
	Count int
	P, Q  *big.Int

	H []byte
	G *big.Int

	XKEY []byte
	X    *big.Int
	Y, Z *big.Int

	KKEY *big.Int
	R    *big.Int
	S    *big.Int

	Fail bool
}

func Test_SignVerify_With_BadPublicKey(t *testing.T) {
	for idx, tc := range testCase_TestVector {
		tc2 := testCase_TestVector[(idx+1)%len(testCase_TestVector)]

		domain, _ := GetDomain(int(tc.Sizes))

		ok := Verify(tc2.P, tc2.Q, tc2.G, tc2.Y, domain.NewHash(), tc.M, tc.R, tc.S)
		if ok {
			t.Errorf("Verify unexpected success with non-existent mod inverse of Q")
			return
		}
	}
}

func generateK(Q *big.Int) (K *big.Int, err error) {
	if Q.Sign() <= 0 || Q.BitLen()%8 != 0 {
		return nil, errors.New("invalid public key")
	}

	privQMinus1 := new(big.Int).Sub(Q, one)

	// step 1. 난수 k를 [1, Q-1]에서 임의로 선택한다.
	for {
		// K = [0 ~ q-2]
		K, err = rand.Int(rnd, privQMinus1)
		if err != nil {
			return
		}
		// k =  K + 1 -> [1 ~ q-1]
		K.Add(K, one)

		if K.Sign() > 0 && K.Cmp(Q) < 0 {
			break
		}
	}

	return
}

func testVerify(t *testing.T, testCases []testCase) {
	for _, tc := range testCases {
		domain, _ := GetDomain(int(tc.Sizes))

		ok := Verify(tc.P, tc.Q, tc.G, tc.Y, domain.NewHash(), tc.M, tc.R, tc.S)
		if ok == tc.Fail {
			t.Errorf("verify failed")
			return
		}
	}
}
