package kx509

import (
	"crypto/elliptic"

	"github.com/RyuaNerin/go-krypto/kcdsa"
)

var (
	curveList = []elliptic.Curve{
		elliptic.P256(),
		elliptic.P224(),
		elliptic.P384(),
		elliptic.P521(),
	}
	sizeList = []kcdsa.ParameterSizes{
		kcdsa.L2048N224SHA224,
		kcdsa.L2048N224SHA256,
		kcdsa.L2048N256SHA256,
		kcdsa.L3072N256SHA256,
	}
)
