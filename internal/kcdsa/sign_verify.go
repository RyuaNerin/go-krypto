package kcdsa

import (
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

func Sign(
	P, Q, G, Y, X *big.Int,
	domain Domain,
	K *big.Int,
	data []byte,
) (r, s *big.Int, ok bool) {
	h := domain.NewHash()

	B := domain.B
	l := h.BlockSize()

	tmp := make([]byte, internal.Bytes(domain.A))

	// step 2. w = g^k mod p
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 2. w = g^k mod p")
	// fmt.Println("G = 0x" + hex.EncodeToString(G.Bytes()))
	// fmt.Println("K = 0x" + hex.EncodeToString(K.Bytes()))
	// fmt.Println("P = 0x" + hex.EncodeToString(P.Bytes()))
	W := new(big.Int).Exp(G, K, P)
	WBytes := tmp[:internal.Bytes(domain.A)]
	W.FillBytes(WBytes)
	// fmt.Println("W = 0x" + hex.EncodeToString(WBytes))
	// fmt.Println("W = 0x" + hex.EncodeToString(W.Bytes()))

	//	step 3. R = h(W) mod 2^β (w를 바이트 열로 변환 후 해시한 결과의 바이트 열에서 	β 비트만큼 절삭):
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 3. R = h(W) mod 2^β")
	h.Reset()
	h.Write(WBytes)
	RBytes := internal.RightMost(h.Sum(tmp[:0]), B)
	R := new(big.Int).SetBytes(RBytes)
	// fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	// step 4. Z = Y mod 2^l
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 4. Z = Y mod 2^l")
	// fmt.Println("Y = 0x" + hex.EncodeToString(Y.Bytes()))
	ZBytesLen := internal.Bytes(Y.BitLen())
	if ZBytesLen < l {
		ZBytesLen = l
	}
	ZBytes := tmp[:ZBytesLen]
	Y.FillBytes(ZBytes)
	ZBytes = internal.RightMost(ZBytes, l*8)
	// fmt.Println("Z = 0x" + hex.EncodeToString(ZBytes))

	// step 5. h = trunc(Hash(Z||M), β)
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 5. h = trunc(Hash(Z||M), β)")
	h.Reset()
	h.Write(ZBytes)
	h.Write(data)
	HBytes := internal.RightMost(h.Sum(tmp[:0]), B)
	H := new(big.Int).SetBytes(HBytes)
	// fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))

	// step 6. E = (R xor H) mod Q
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 6. E = (R xor H) mod Q")
	// fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))
	// fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))
	// fmt.Println("Q = 0x" + hex.EncodeToString(Q.Bytes()))
	E := new(big.Int).Xor(R, H)
	E.Mod(E, Q)
	// fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))

	// step 7. S = X(K-E) mod Q
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 7. S = X(K-E) mod Q")
	// fmt.Println("X = 0x" + hex.EncodeToString(X.Bytes()))
	// fmt.Println("K = 0x" + hex.EncodeToString(K.Bytes()))
	// fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))
	// fmt.Println("Q = 0x" + hex.EncodeToString(Q.Bytes()))
	K.Mod(K.Sub(K, E), Q)
	S := new(big.Int).Mul(X, K)
	S.Mod(S, Q)
	// fmt.Println("S = 0x" + hex.EncodeToString(S.Bytes()))

	r = R
	s = S

	ok = r.Sign() != 0 && s.Sign() != 0
	return
}

func Verify(
	P, Q, G, Y *big.Int,
	domain Domain,
	data []byte,
	R, S *big.Int,
) bool {
	h := domain.NewHash()

	B := domain.B
	l := h.BlockSize()

	tmpSize := l
	YBytesLen := internal.Bytes(Y.BitLen())
	PBytesLen := internal.Bytes(P.BitLen())
	if tmpSize < YBytesLen {
		tmpSize = YBytesLen
	}
	if tmpSize < PBytesLen {
		tmpSize = PBytesLen
	}

	tmp := make([]byte, tmpSize)

	// step 2. Z = Y mod 2^l
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 2. Z = Y mod 2^l")
	// fmt.Println("Y = 0x" + hex.EncodeToString(Y.Bytes()))
	if YBytesLen < l {
		YBytesLen = l
	}
	ZBytes := tmp[:YBytesLen]
	Y.FillBytes(ZBytes)
	ZBytes = internal.RightMost(ZBytes, l*8)
	// fmt.Println("Z = 0x" + hex.EncodeToString(ZBytes))

	// step 3. h = trunc(Hash(Z||M), β)
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 3. h = trunc(Hash(Z||M), β)")
	h.Reset()
	h.Write(ZBytes)
	h.Write(data)
	HBytes := internal.RightMost(h.Sum(tmp[:0]), B)
	H := new(big.Int).SetBytes(HBytes)
	// fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))

	// step 4. E' = (R' xor H') mod Q
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 4. E' = (R' xor H') mod Q")
	// fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))
	// fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))
	// fmt.Println("Q = 0x" + hex.EncodeToString(Q.Bytes()))
	E := new(big.Int).Xor(R, H)
	E.Mod(E, Q)
	// fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))

	// step 5. W' = Y ^ {S'} G ^ {E'} mod P
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 5. W' = Y ^ {S'} G ^ {E'} mod P")
	// fmt.Println("Y = 0x" + hex.EncodeToString(Y.Bytes()))
	// fmt.Println("G = 0x" + hex.EncodeToString(G.Bytes()))
	// fmt.Println("S = 0x" + hex.EncodeToString(S.Bytes()))
	// fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))
	// fmt.Println("P = 0x" + hex.EncodeToString(P.Bytes()))
	W := new(big.Int).Exp(Y, S, P)
	E.Exp(G, E, P)
	W.Mul(W, E)
	W.Mod(W, P)

	WBytes := tmp[:PBytesLen]
	W.FillBytes(WBytes)
	// fmt.Println("W = 0x" + hex.EncodeToString(W.Bytes()))

	// step 6. trunc(Hash(W'), β) = R'이 성립하는지 확인한다.
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 6. trunc(Hash(W'), β) = R'")
	h.Reset()
	h.Write(WBytes)
	rBytes := internal.RightMost(h.Sum(tmp[:0]), B)
	r := new(big.Int).SetBytes(rBytes)
	// fmt.Println("r = 0x" + hex.EncodeToString(r.Bytes()))
	// fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	return internal.BigIntEqual(R, r)
}
