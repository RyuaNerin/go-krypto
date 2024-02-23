package eckcdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"testing"

	"github.com/RyuaNerin/go-krypto/internal"

	"github.com/RyuaNerin/elliptic2/nist"
)

var convertTestCases = []struct {
	curve elliptic.Curve
	D     *big.Int
	X     *big.Int
	Y     *big.Int
	Xk    *big.Int
	Yk    *big.Int
}{
	{
		curve: elliptic.P224(),
		D:     internal.HI(`99b4d4bc70626ab185206ab890ffcb7a5319a2a88b325f4f2992861d`),
		X:     internal.HI(`19ae5855ada42406160e8a714efd9082ec7e540ee0819cf6b44271d3`),
		Y:     internal.HI(`87ff0ab3faf76add576013f918f7cb7a97c1dfc439daa8bc28a7253c`),
		Xk:    internal.HI(`c5cdf38f05476f090f112613517fa83c1f6b7956d7d417a6b9b590d9`),
		Yk:    internal.HI(`e0c6767e07c080128a7551e0133db88acb7abe841acb316ceb8297b8`),
	},
	{
		curve: elliptic.P256(),
		D:     internal.HI(`c2cdc894031d20e65d38fab2e091f864de4bf2630357cc3b5a76808a7d972fb6`),
		X:     internal.HI(`3bfb4a2a1025a0987a02e7b6dc26ffea7b87ca3c9196bfe7353b07aded59deb0`),
		Y:     internal.HI(`662e530d79f6e197f52abdf1e92f718a010a993da369e759bbd233ecb432d8c5`),
		Xk:    internal.HI(`eb3ee504db3e764a3118d59fae4c2a683844d5dc155755217d371df2ed5c5838`),
		Yk:    internal.HI(`469acd87e995ae3e8ddb45fe12031ec7f84370dbe9257db3f1109cfa06b702b1`),
	},
	{
		curve: elliptic.P384(),
		D:     internal.HI(`5f1692086d16eef76d50fc499209332f4a897bc8d288ebed5701d95d131025bb038aee8701b8d6425d1892c563e4ee30`),
		X:     internal.HI(`ed186320ecb0dde24e69ae185b94f5a75c79d34b9c19445baacef7663e957a6894327ef01adf94906629df933332524d`),
		Y:     internal.HI(`58b475e227da277e1858bc92681e1d5e742f1054a20331e04910aefb3eef7a47276490d1a616cdd54b7f987fdf469be5`),
		Xk:    internal.HI(`b9bfe8a8fce4e3ed58ac07ec4d3169c8a9bde5016ece04607ab2490c45ae4045277bda17a3bcd861bab90e9577d811c2`),
		Yk:    internal.HI(`b9a1dbb5478c1c16ecee0457c843e5ce811e09882f798c4f3a6df6206eb53f14ecd536e88dc963e28a7218f535aa7164`),
	},
	{
		curve: elliptic.P521(),
		D:     internal.HI(`1d5f7709fc3b079020667e0ed0c96a7cec12e866e80f9a7d5cdd58ca2bf14ad759cd95827c4ee4ad27ec9f519e8c7f8c32fc7ec49d0c4d24ec3862021a4add0dce3`),
		X:     internal.HI(`14a66787e624eb330de45b5dde2d208c220f7c187dad49959209642b51b4885318b0e53e7358ca7edbd9b82536f8813c0f406fbc890ce73ceadc80129bb018dd234`),
		Y:     internal.HI(`0d64e2bce7bf8484766550f03cbf1fdc54f946bb04407f170cd49111148a91ba2935b0b279b78d5e01cd9115099c3684ca70c1fcb3931f806e57856cde6b01bf84a`),
		Xk:    internal.HI(`060e5c33466a13cb6f9a9f557c74be327483f4b2f55ed33a3f5f19a7dedb88f83e558dce6b18d8d757698162a930c7502e05b44a03216340e4510031887d2478269`),
		Yk:    internal.HI(`1f57d63a9ca8fec19cfe6529e378c5496e93b0714051bd52d95ceb11ebe55ab1f37de703337b3a245a6de43006fe494aa5648a1275159b67c3144bc59c088b1887e`),
	},
	{
		curve: nist.K283(),
		D:     internal.HI(`10d64434603b497e976d9bf192c021dfe2a5a77e5f217cfccf748fde2af6d39abedcbef`),
		X:     internal.HI(`1d70ffdaabab4d16c99c1fef5824f522675008d3e9e27b83719935f7410301c7b85d6bc`),
		Y:     internal.HI(`6ff7a26bec1a9a00a1dce7ba4f80fcacfede8c1a2e0c1e83f74c8ac77448ad66949e2f6`),
		Xk:    internal.HI(`358de9bc5d65cbc011f47f3da730513967491228056ff8c7780f283df9ebaaeaeef5ec9`),
		Yk:    internal.HI(`2bc5b17bca09a7851ed77c662152c597978135e8ccd085595db1835de956c847d9f4e8a`),
	},
}

func Test_ECDSA_TO_ECKCDSA(t *testing.T) {
	for _, tc := range convertTestCases {
		input := &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: tc.curve,
				X:     tc.X,
				Y:     tc.Y,
			},
			D: tc.D,
		}
		expect := &PrivateKey{
			PublicKey: PublicKey{
				Curve: tc.curve,
				X:     tc.Xk,
				Y:     tc.Yk,
			},
			D: tc.D,
		}

		answer := FromECDSA(input)

		if !answer.Equal(expect) {
			t.Fail()
		}
	}
}

func Test_ECKCDSA_TO_ECDSA(t *testing.T) {
	for _, tc := range convertTestCases {
		input := &PrivateKey{
			PublicKey: PublicKey{
				Curve: tc.curve,
				X:     tc.Xk,
				Y:     tc.Yk,
			},
			D: tc.D,
		}
		expect := &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: tc.curve,
				X:     tc.X,
				Y:     tc.Y,
			},
			D: tc.D,
		}

		answer := input.ToECDSA()

		if !answer.Equal(expect) {
			t.Fail()
		}
	}
}
