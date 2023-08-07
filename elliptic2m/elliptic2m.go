// Package elliptic2m implements the B-233, B-283, K-233, K-283
// based on https://gist.github.com/kurtbrose/4423605

package elliptic2m

import (
	"crypto/elliptic"
	"math/big"
	"sync"

	"github.com/RyuaNerin/go-krypto/internal"
)

var initonce sync.Once

var (
	b233 curve
	b283 curve
	k233 curve
	k283 curve
)

// Also known as: sect233r1, wap-wsg-idm-ecid-wtls11, ansit233r1
func B233() elliptic.Curve {
	initonce.Do(initAll)
	return &b233
}
func Sect233r1() elliptic.Curve  { return B233() }
func Ansit233r1() elliptic.Curve { return B233() }

// Also known as: sect283r1, ansit283r1
func B283() elliptic.Curve {
	initonce.Do(initAll)
	return &b283
}
func Sect283r1() elliptic.Curve  { return B283() }
func Ansit283r1() elliptic.Curve { return B283() }

// Also known as: sect233k1, wap-wsg-idm-ecid-wtls10, ansit233k1
func K233() elliptic.Curve {
	initonce.Do(initAll)
	return &k233
}
func Sect233k1() elliptic.Curve  { return K233() }
func Ansit233k1() elliptic.Curve { return K233() }

// Also known as: sect283k1, ansit283k1
func K283() elliptic.Curve {
	initonce.Do(initAll)
	return &k283
}
func Sect283k1() elliptic.Curve  { return K283() }
func Ansit283k1() elliptic.Curve { return K283() }

func initAll() {
	initB233()
	initK233()
	initB283()
	initK283()
}

func calcF(f ...int) *big.Int {
	ret := new(big.Int)
	for _, v := range f {
		tmp := big.NewInt(1)
		tmp.Lsh(tmp, uint(v))

		ret.Add(ret, tmp)
	}
	ret.Add(ret, big.NewInt(1))

	return ret
}

func initB233() {
	// https://neuromancer.sk/std/nist/B-233
	// sect233r1
	b233.params = &curveParams{
		A: internal.HI(`0x000000000000000000000000000000000000000000000000000000000001`),
		CurveParams: elliptic.CurveParams{
			Name:    "B-233",
			BitSize: 233,
			P:       calcF(233, 74),
			B:       internal.HI(`0x0066647ede6c332c7f8c0923bb58213b333b20e9ce4281fe115f7d8f90ad`),
			Gx:      internal.HI(`0x00fac9dfcbac8313bb2139f1bb755fef65bc391f8b36f8f8eb7371fd558b`),
			Gy:      internal.HI(`0x01006a08a41903350678e58528bebf8a0beff867a7ca36716f7e01f81052`),
			N:       internal.HI(`0x1000000000000000000000000000013e974e72f8a6922031d2603cfe0d7`),
		},
	}
}
func initK233() {
	// https://neuromancer.sk/std/nist/K-233
	// sect233k1
	k233.params = &curveParams{
		A: internal.HI(`0x000000000000000000000000000000000000000000000000000000000000`),
		CurveParams: elliptic.CurveParams{
			Name:    "K-233",
			BitSize: 233,
			P:       calcF(233, 74),
			B:       internal.HI(`0x000000000000000000000000000000000000000000000000000000000001`),
			Gx:      internal.HI(`0x017232ba853a7e731af129f22ff4149563a419c26bf50a4c9d6eefad6126`),
			Gy:      internal.HI(`0x01db537dece819b7f70f555a67c427a8cd9bf18aeb9b56e0c11056fae6a3`),
			N:       internal.HI(`0x8000000000000000000000000000069d5bb915bcd46efb1ad5f173abdf`),
		},
	}
}
func initB283() {
	// https://neuromancer.sk/std/nist/B-283
	b283.params = &curveParams{
		A: internal.HI(`0x00000000000000000000000000000000000000000000000000000000000000000000001`),
		CurveParams: elliptic.CurveParams{
			Name:    "B-283",
			BitSize: 283,
			P:       calcF(283, 12, 7, 5),
			B:       internal.HI(`0x27b680ac8b8596da5a4af8a19a0303fca97fd7645309fa2a581485af6263e313b79a2f5`),
			Gx:      internal.HI(`0x5f939258db7dd90e1934f8c70b0dfec2eed25b8557eac9c80e2e198f8cdbecd86b12053`),
			Gy:      internal.HI(`0x3676854fe24141cb98fe6d4b20d02b4516ff702350eddb0826779c813f0df45be8112f4`),
			N:       internal.HI(`0x3ffffffffffffffffffffffffffffffffffef90399660fc938a90165b042a7cefadb307`),
		},
	}
}
func initK283() {
	// https://neuromancer.sk/std/nist/K-283
	k283.params = &curveParams{
		A: internal.HI(`0x00000000000000000000000000000000000000000000000000000000000000000000000`),
		CurveParams: elliptic.CurveParams{
			Name:    "K-283",
			BitSize: 283,
			P:       calcF(283, 12, 7, 5),
			B:       internal.HI(`0x00000000000000000000000000000000000000000000000000000000000000000000001`),
			Gx:      internal.HI(`0x503213f78ca44883f1a3b8162f188e553cd265f23c1567a16876913b0c2ac2458492836`),
			Gy:      internal.HI(`0x1ccda380f1c9e318d90f95d07e5426fe87e45c0e8184698e45962364e34116177dd2259`),
			N:       internal.HI(`0x1ffffffffffffffffffffffffffffffffffe9ae2ed07577265dff7f94451e061e163c61`),
		},
	}
}
