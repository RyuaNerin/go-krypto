package elliptic2m

import (
	"crypto/elliptic"
	"testing"
)

func benchmarkAllCurves(b *testing.B, f func(*testing.B, elliptic.Curve)) {
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
		b.Run(test.name, func(B *testing.B) {
			f(b, test.curve)
		})
	}
}

func BenchmarkScalarBaseMult(b *testing.B) {
	benchmarkAllCurves(b, func(b *testing.B, curve elliptic.Curve) {
		priv := getK(curve)

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x, _ := curve.ScalarBaseMult(priv)
			priv[0] ^= byte(x.Bits()[0])
		}
	})
}

func BenchmarkScalarMult(b *testing.B) {
	benchmarkAllCurves(b, func(b *testing.B, curve elliptic.Curve) {
		priv := getK(curve)
		x, y := curve.ScalarBaseMult(priv)

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x, y = curve.ScalarMult(x, y, priv)
			priv[0] ^= byte(x.Bits()[0])
		}
	})
}

func BenchmarkDouble(b *testing.B) {
	benchmarkAllCurves(b, func(b *testing.B, curve elliptic.Curve) {
		x, y := curve.ScalarBaseMult(getK(curve))

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x, y = curve.Double(x, y)
		}
	})
}

func BenchmarkAdd(b *testing.B) {
	benchmarkAllCurves(b, func(b *testing.B, curve elliptic.Curve) {
		x1, y1 := curve.ScalarBaseMult(getK(curve))
		x2, y2 := curve.ScalarBaseMult(getK(curve))

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			x, y := curve.Add(x1, y1, x2, y2)
			x2, y2 = x1, y1
			x1, y1 = x, y
		}
	})
}
