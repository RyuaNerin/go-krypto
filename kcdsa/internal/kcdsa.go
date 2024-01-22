package internal

import (
	"crypto/subtle"
	"hash"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
)

var (
	one = big.NewInt(1)
)

func Sign(P, Q, G, Y, X, K *big.Int, h hash.Hash, data []byte) (r, s *big.Int, err error) {
	// Q 생성할 때, Q 사이즈를 doamin.B 사이즈랑 동일하게 생성한다.
	B := Q.BitLen()

	buf := make([]byte, 0, h.Size())

	// step 2. w = g^k mod p
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 2. w = g^k mod p")
	//fmt.Println("G = 0x" + hex.EncodeToString(G.Bytes()))
	//fmt.Println("K = 0x" + hex.EncodeToString(K.Bytes()))
	//fmt.Println("P = 0x" + hex.EncodeToString(P.Bytes()))
	W := new(big.Int).Exp(G, K, P)
	//fmt.Println("W = 0x" + hex.EncodeToString(W.Bytes()))

	//	step 3. R = h(W) mod 2^β (w를 바이트 열로 변환 후 해시한 결과의 바이트 열에서 	β 비트만큼 절삭):
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 3. R = h(W) mod 2^β")
	h.Reset()
	h.Write(W.Bytes())
	RBytes := internal.TruncateLeft(h.Sum(buf[:0]), B)
	R := new(big.Int).SetBytes(RBytes)
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	// step 4. Z = Y mod 2^l
	i2l := new(big.Int).Lsh(one, uint(h.BlockSize())*8)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 4. Z = Y mod 2^l")
	//fmt.Println("Y = 0x" + hex.EncodeToString(Y.Bytes()))
	//fmt.Println("2l = 0x" + hex.EncodeToString(i2l.Bytes()))
	Z := new(big.Int).Mod(Y, i2l)
	ZBytes := Z.Bytes()
	//fmt.Println("Z = 0x" + hex.EncodeToString(ZBytes))

	// step 5. h = trunc(Hash(Z||M), β)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 5. h = trunc(Hash(Z||M), β)")
	h.Reset()
	h.Write(ZBytes)
	h.Write(data)
	HBytes := internal.TruncateLeft(h.Sum(buf[:0]), B)
	H := new(big.Int).SetBytes(HBytes)
	//fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))

	// step 6. E = (R xor H) mod Q
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 6. E = (R xor H) mod Q")
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))
	//fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))
	//fmt.Println("Q = 0x" + hex.EncodeToString(Q.Bytes()))
	E := new(big.Int).Xor(R, H)
	E.Mod(E, Q)
	//fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))

	//step 7. S = X(K-E) mod Q
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 7. S = X(K-E) mod Q")
	//fmt.Println("X = 0x" + hex.EncodeToString(X.Bytes()))
	//fmt.Println("K = 0x" + hex.EncodeToString(K.Bytes()))
	//fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))
	//fmt.Println("Q = 0x" + hex.EncodeToString(Q.Bytes()))
	K.Mod(K.Sub(K, E), Q)
	S := new(big.Int).Mul(X, K)
	S.Mod(S, Q)
	//fmt.Println("S = 0x" + hex.EncodeToString(S.Bytes()))

	r = R
	s = S

	return
}

func Verify(P, Q, G, Y *big.Int, h hash.Hash, data []byte, R, S *big.Int) bool {
	// Q 생성할 때, Q 사이즈를 doamin.B 사이즈랑 동일하게 생성한다.
	B := Q.BitLen()

	// step 1. 수신된 서명 {R', S'}에 대해 |R'|=LH, 0 < S' < Q 임을 확인한다.
	if P.Sign() <= 0 {
		return false
	}

	if S.Sign() < 1 || S.Cmp(Q) >= 0 {
		return false
	}

	buf := make([]byte, h.Size())

	// step 2. Z = Y mod 2^l
	i2l := new(big.Int).Lsh(one, uint(h.BlockSize())*8)

	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 2. Z = Y mod 2^l")
	//fmt.Println("Y = 0x" + hex.EncodeToString(Y.Bytes()))
	//fmt.Println("2l = 0x" + hex.EncodeToString(i2l.Bytes()))
	Z := new(big.Int).Mod(Y, i2l)
	ZBytes := Z.Bytes()
	//fmt.Println("Z = 0x" + hex.EncodeToString(ZBytes))

	// step 3. h = trunc(Hash(Z||M), β)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 3. h = trunc(Hash(Z||M), β)")
	h.Reset()
	h.Write(ZBytes)
	h.Write(data)
	HBytes := internal.TruncateLeft(h.Sum(buf[:0]), B)
	H := new(big.Int).SetBytes(HBytes)
	//fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))

	// step 4. E' = (R' xor H') mod Q
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 4. E' = (R' xor H') mod Q")
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))
	//fmt.Println("H = 0x" + hex.EncodeToString(H.Bytes()))
	//fmt.Println("Q = 0x" + hex.EncodeToString(Q.Bytes()))
	E := new(big.Int).Xor(R, H)
	E.Mod(E, Q)
	//fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))

	// step 5. W' = Y ^ {S'} G ^ {E'} mod P
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 5. W' = Y ^ {S'} G ^ {E'} mod P")
	//fmt.Println("Y = 0x" + hex.EncodeToString(Y.Bytes()))
	//fmt.Println("G = 0x" + hex.EncodeToString(G.Bytes()))
	//fmt.Println("S = 0x" + hex.EncodeToString(S.Bytes()))
	//fmt.Println("E = 0x" + hex.EncodeToString(E.Bytes()))
	//fmt.Println("P = 0x" + hex.EncodeToString(P.Bytes()))
	W := new(big.Int).Exp(Y, S, P)
	E.Exp(G, E, P)
	W.Mul(W, E)
	W.Mod(W, P)
	//fmt.Println("W = 0x" + hex.EncodeToString(W.Bytes()))

	// step 6. trunc(Hash(W'), β) = R'이 성립하는지 확인한다.
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("step 6. trunc(Hash(W'), β) = R'")
	h.Reset()
	h.Write(W.Bytes())
	rBytes := internal.TruncateLeft(h.Sum(buf[:0]), B)
	r := new(big.Int).SetBytes(rBytes)
	//fmt.Println("r = 0x" + hex.EncodeToString(r.Bytes()))
	//fmt.Println("R = 0x" + hex.EncodeToString(R.Bytes()))

	return bigIntEqual(R, r)
}

func bigIntEqual(a, b *big.Int) bool {
	return subtle.ConstantTimeCompare(a.Bytes(), b.Bytes()) == 1
}
