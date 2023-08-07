package elliptic2m

import (
	"crypto/elliptic"
	"math/big"
	"testing"
)

type testCase struct {
	Qx, Qy *big.Int
	Fail   bool
}

func testPoint(t *testing.T, testCases []testCase, curve elliptic.Curve) {
	for idx, tc := range testCases {
		ok := curve.IsOnCurve(tc.Qx, tc.Qy)
		if ok == tc.Fail {
			t.Errorf("%d: Verify failed, got:%v want:%v", idx, ok, !tc.Fail)
			return
		}
	}
}
