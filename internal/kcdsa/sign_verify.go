package kcdsa

import (
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

func Sign(
	R, S *big.Int,
	P, Q, G, Y, X *big.Int,
	domain Domain,
	K *big.Int,
	data []byte,
	tmpInt *big.Int,
	tmpBuf []byte,
) (tmpBufOut []byte, ok bool) {
	h := domain.NewHash()

	B := domain.B
	l := h.BlockSize()

	tmpBuf = internal.Grow(tmpBuf, internal.BitsToBytes(domain.A))

	// step 2. w = g^k mod p
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 2. w = g^k mod p")
	// fmt.Println("G = 0x" + hex.EncodeToString(G.Bytes()))
	// fmt.Println("K = 0x" + hex.EncodeToString(K.Bytes()))
	// fmt.Println("P = 0x" + hex.EncodeToString(P.Bytes()))
	W := tmpInt.Exp(G, K, P)
	WBytes := tmpBuf[:internal.BitsToBytes(domain.A)]
	W.FillBytes(WBytes)
	// fmt.Println("W = 0x" + hex.EncodeToString(WBytes))
	// fmt.Println("W = 0x" + hex.EncodeToString(W.Bytes()))

	//	step 3. R = h(W) mod 2^β (w를 바이트 열로 변환 후 해시한 결과의 바이트 열에서 	β 비트만큼 절삭):
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 3. R = h(W) mod 2^β")
	h.Reset()
	h.Write(WBytes)
	RBytes := internal.RightMost(h.Sum(tmpBuf[:0]), B)
	R.SetBytes(RBytes)
	// fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	// step 4. Z = Y mod 2^l
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 4. Z = Y mod 2^l")
	// fmt.Println("Y = 0x" + hex.EncodeToString(Y.Bytes()))
	ZBytesLen := internal.BitsToBytes(Y.BitLen())
	if ZBytesLen < l {
		ZBytesLen = l
	}
	ZBytes := tmpBuf[:ZBytesLen]
	Y.FillBytes(ZBytes)
	ZBytes = internal.RightMost(ZBytes, l*8)
	// fmt.Println("Z = 0x" + hex.EncodeToString(ZBytes))

	// step 5. h = trunc(Hash(Z||M), β)
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 5. h = trunc(Hash(Z||M), β)")
	h.Reset()
	h.Write(ZBytes)
	h.Write(data)
	HBytes := internal.RightMost(h.Sum(tmpBuf[:0]), B)
	H := tmpInt.SetBytes(HBytes)
	// fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))

	// step 6. E = (R xor H) mod Q
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 6. E = (R xor H) mod Q")
	// fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))
	// fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))
	// fmt.Println("Q = 0x" + hex.EncodeToString(Q.Bytes()))
	E := tmpInt.Xor(R, H)
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
	S.Mul(X, K)
	S.Mod(S, Q)
	// fmt.Println("S = 0x" + hex.EncodeToString(S.Bytes()))

	return tmpBuf, R.Sign() != 0 && S.Sign() != 0
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
	YBytesLen := internal.BitsToBytes(Y.BitLen())
	PBytesLen := internal.BitsToBytes(P.BitLen())
	if tmpSize < YBytesLen {
		tmpSize = YBytesLen
	}
	if tmpSize < PBytesLen {
		tmpSize = PBytesLen
	}

	tmp := make([]byte, tmpSize)

	tmpInt1 := new(big.Int)
	tmpInt2 := new(big.Int)

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
	H := tmpInt1.SetBytes(HBytes)
	// fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))

	// step 4. E' = (R' xor H') mod Q
	// fmt.Println("--------------------------------------------------")
	// fmt.Println("step 4. E' = (R' xor H') mod Q")
	// fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))
	// fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))
	// fmt.Println("Q = 0x" + hex.EncodeToString(Q.Bytes()))
	E := tmpInt1.Xor(R, H)
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
	W := tmpInt2.Exp(Y, S, P)
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
	r := tmpInt1.SetBytes(rBytes)
	// fmt.Println("r = 0x" + hex.EncodeToString(r.Bytes()))
	// fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	return internal.BigEqual(R, r)
}
