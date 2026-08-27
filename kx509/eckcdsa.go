package kx509

import (
	"crypto/elliptic"
	"math/big"

	"github.com/RyuaNerin/elliptic2"
)

func getP(curve elliptic2.CurveExtended) *big.Int {
	params2 := curve.Params2()
	p := params2.P
	if p == nil {
		p = params2.P
	}
	return p
}

func ellipticUnmarshal(curve elliptic.Curve, data []byte) (x, y *big.Int) {
	curve2, ok := curve.(elliptic2.CurveExtended)
	if !ok {
		return elliptic.Unmarshal(curve, data)
	}

	byteLen := (curve.Params().BitSize + 7) / 8
	if len(data) != 1+2*byteLen {
		return nil, nil
	}
	if data[0] != 4 { // uncompressed form
		return nil, nil
	}
	p := getP(curve2)
	x = new(big.Int).SetBytes(data[1 : 1+byteLen])
	y = new(big.Int).SetBytes(data[1+byteLen:])
	if x.Cmp(p) >= 0 || y.Cmp(p) >= 0 {
		return nil, nil
	}
	if !curve.IsOnCurve(x, y) {
		return nil, nil
	}
	return
}

func ellipticUnmarshalCompressed(curve elliptic.Curve, data []byte) (x, y *big.Int) {
	curve2, ok := curve.(elliptic2.CurveExtended)
	if !ok {
		return elliptic.UnmarshalCompressed(curve, data)
	}

	byteLen := (curve.Params().BitSize + 7) / 8
	if len(data) != 1+byteLen {
		return nil, nil
	}
	if data[0] != 2 && data[0] != 3 { // compressed form
		return nil, nil
	}
	p := getP(curve2)
	x = new(big.Int).SetBytes(data[1:])
	if x.Cmp(p) >= 0 {
		return nil, nil
	}
	y = curve2.ComputeY(x, false)
	if y == nil {
		return nil, nil
	}
	if byte(y.Bit(0)) != data[0]&1 {
		y = curve2.ComputeY(x, true)
	}

	if !curve.IsOnCurve(x, y) {
		return nil, nil
	}

	return x, y
}
