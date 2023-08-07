package elliptic2m

import (
	"bufio"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"
)

var (
	rnd = bufio.NewReaderSize(rand.Reader, 1<<15)
)

type internalTestcase struct {
	x1, y1 *big.Int
	x2, y2 *big.Int
	x, y   *big.Int
}

func testAllCurves(t *testing.T, f func(*testing.T, elliptic.Curve)) {
	tests := []struct {
		name  string
		curve elliptic.Curve
	}{
		{"B-233", B233()},
		{"K-233", K233()},
		{"B-283", B283()},
		{"K-283", K283()},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			f(t, test.curve)
		})
	}
}

func getK(c elliptic.Curve) []byte {
	k, _ := rand.Int(rnd, c.Params().N)
	return k.Bytes()
}

func Test_ScalarBaseMult(t *testing.T) {
	testAllCurves(t, func(t *testing.T, c elliptic.Curve) {
		x, y := c.ScalarBaseMult(getK(c))
		if !c.IsOnCurve(x, y) {
			t.Fail()
		}
	})
}

func Test_ScalarMult(t *testing.T) {
	testAllCurves(t, func(t *testing.T, c elliptic.Curve) {
		x1, y1 := c.ScalarBaseMult(getK(c))
		if !c.IsOnCurve(x1, y1) {
			t.Fail()
		}

		x, y := c.ScalarMult(x1, y1, getK(c))
		if !c.IsOnCurve(x, y) {
			t.Fail()
		}
	})
}

func Test_Double(t *testing.T) {
	testAllCurves(t, func(t *testing.T, c elliptic.Curve) {
		x1, y1 := c.ScalarBaseMult(getK(c))
		if !c.IsOnCurve(x1, y1) {
			t.Fail()
		}

		x, y := c.Double(x1, y1)
		if !c.IsOnCurve(x, y) {
			t.Fail()
		}
	})
}

func Test_Add(t *testing.T) {
	testAllCurves(t, func(t *testing.T, c elliptic.Curve) {
		x1, y1 := c.ScalarBaseMult(getK(c))
		if !c.IsOnCurve(x1, y1) {
			t.Fail()
		}
		x2, y2 := c.ScalarBaseMult(getK(c))
		if !c.IsOnCurve(x2, y2) {
			t.Fail()
		}

		x, y := c.Add(x1, y1, x2, y2)
		if !c.IsOnCurve(x, y) {
			t.Fail()
		}
	})
}
