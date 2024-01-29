package internal

import (
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

func Y(Q, P, G, X *big.Int) *big.Int {
	// x의 역원 생성
	xInv := internal.FermatInverse(X, Q)

	// 전자서명 검증키 y 생성(Y = G^{X^{-1} mod Q} mod P)
	return new(big.Int).Exp(G, xInv, P)
}
