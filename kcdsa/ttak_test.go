package kcdsa

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func Test_TTAK_GenerateJ(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ttak parameter generation test in short mode")
		return
	}

	for _, tc := range testCase_TestVector {
		d, _ := tc.Sizes.domain()
		J, _, ok := generateJ(tc.Seed_, nil, d.NewHash(), d)
		if !ok {
			t.Fail()
			return
		}
		if J.Cmp(tc.J) != 0 {
			t.Errorf("GenerateTTAKJ failed")
			return
		}
	}
}

func Test_TTAK_GeneratePQ(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ttak parameter generation test in short mode")
		return
	}

	for _, tc := range testCase_TestVector {
		d, _ := tc.Sizes.domain()
		P, Q, count, ok := generatePQ(tc.J, tc.Seed_, d.NewHash(), d)
		if !ok {
			t.Fail()
			return
		}
		if P.Cmp(tc.P) != 0 || Q.Cmp(tc.Q) != 0 || count != tc.Count {
			t.Errorf("GenerateTTAKPQ failed")
			return
		}
	}
}

func Test_TTAK_GenerateHG(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ttak parameter generation test in short mode")
		return
	}

	for _, tc := range testCase_TestVector {
		_, _, err := generateHG(rand.Reader, tc.P, tc.J)
		if err != nil {
			t.Error(err)
			return
		}
	}
}

func Test_TTAK_GenerateG(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ttak parameter generation test in short mode")
		return
	}

	for _, tc := range testCase_TestVector {
		G, ok := generateG(tc.P, tc.J, new(big.Int).SetBytes(tc.H))
		if !ok {
			t.Fail()
			return
		}
		if G.Cmp(tc.G) != 0 {
			t.Errorf("GenerateTTAKG failed")
			return
		}
	}
}

func Test_RegenerateParametersTTAK(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ttak parameter generation test in short mode")
		return
	}

	for _, tc := range testCase_TestVector {
		params := Parameters{
			TTAKParams: TTAKParameters{
				J:     tc.J,
				Seed:  tc.Seed_,
				Count: tc.Count,
			},
		}
		err := RegenerateParametersTTAK(&params, rnd, tc.Sizes)
		if err != nil {
			t.Error(err)
			return
		}

		if params.P.Cmp(tc.P) != 0 || params.Q.Cmp(tc.Q) != 0 {
			t.Errorf("GenerateTTAKG failed")
			return
		}
	}
}