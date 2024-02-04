package kcdsa

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"hash"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	kcdsainternal "github.com/RyuaNerin/go-krypto/internal/kcdsa"
)

var (
	ErrTTAKParametersNotSetUp = errors.New("krypto/kcdsa: ttakparameters not set up before generating key")
	ErrShortXKEY              = errors.New("krypto/kcdsa: XKEY is too small.")

	two   = big.NewInt(2)
	three = big.NewInt(3)
)

// Generate the paramters
// using the prime number generator defined in TTAK.KO12.0001/R4
func GenerateParametersTTAK(params *Parameters, rand io.Reader, sizes ParameterSizes) (H *big.Int, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, ErrInvalidParameterSizes
	}

	h := domain.NewHash()

	// p. 13
	var seed []byte
	var ubuf []byte
	for {
		seed, err = internal.ReadBits(seed, rand, domain.B)
		if err != nil {
			return nil, err
		}

		// 2 ~ 4
		J, ubuf2, ok := generateJ(seed, ubuf, h, domain)
		if !ok {
			ubuf = ubuf2
			continue
		}
		/**
		J, ubuf2, ok, err := generateJAlt(rand, seed, ubuf[:0], domain)
		if err != nil {
			return err
		}
		if !ok {
			ubuf = ubuf2
			continue
		}
		*/

		// 5 ~ 12
		P, Q, count, ok := generatePQ(J, seed, h, domain)
		if !ok {
			continue
		}

		H, G, err := generateHG(rand, P, J)
		if err != nil {
			return nil, err
		}

		params.TTAKParams = TTAKParameters{
			J:     J,
			Seed:  seed,
			Count: count,
		}

		params.P = P
		params.Q = Q
		params.G = G

		return H, nil
	}
}

// TTAKParameters -> P, Q, G(randomly)
func RegenerateParametersTTAK(params *Parameters, rand io.Reader, sizes ParameterSizes) error {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return ErrInvalidParameterSizes
	}

	if params.TTAKParams.Count == 0 || params.TTAKParams.J == nil || params.TTAKParams.Seed == nil || params.TTAKParams.J.Sign() <= 0 {
		return ErrInvalidTTAKParameters
	}
	if params.TTAKParams.J.Sign() <= 0 {
		return ErrInvalidTTAKParameters
	}

	if len(params.TTAKParams.Seed) != internal.Bytes(domain.B) {
		return ErrInvalidTTAKParameters
	}

	q := new(big.Int)
	p := new(big.Int)

	seedCount := make([]byte, len(params.TTAKParams.Seed)+4)
	copy(seedCount, params.TTAKParams.Seed)
	binary.BigEndian.PutUint32(seedCount[len(params.TTAKParams.Seed):], uint32(params.TTAKParams.Count))

	uBuf := make([]byte, internal.Bytes(domain.B))

	// 8: Seed에 Count를 연접한 것을 일방향 함수 PPGF의 입력으로 하여 비트 길이가
	// β인 난수 U를 생성한다. (U ← PPGF(Seed ‖ Count, β))
	U := ppgf(uBuf[:0], seedCount, domain.B, domain.NewHash())

	// 9: U의 최상위 및 최하위 비트를 1로 만들어 이를 q로 둔다.
	// (q ← 2^(β-1) ∨ U ∨ 1)
	U[0] |= 0b1000_0000
	U[len(U)-1] |= 1
	q.SetBytes(U)

	// 10: p ← (2Jq + 1)의 비트 길이가 α보다 길면 단계 6으로 간다.
	p.Add(p.Lsh(p.Mul(params.TTAKParams.J, q), 1), one)
	if p.BitLen() > domain.A {
		return ErrInvalidTTAKParameters
	}

	// 11: 강한 소수 판정 알고리즘으로 q를 판정하여 소수가 아니면 단계 6으로 간다.
	if !q.ProbablyPrime(internal.NumMRTests) {
		return ErrInvalidTTAKParameters
	}

	// 12: 강한 소수 판정 알고리즘으로 p를 판정하여 소수가 아니면 단계 6으로 간다
	if !p.ProbablyPrime(internal.NumMRTests) {
		return ErrInvalidTTAKParameters
	}

	_, g, err := generateHG(rand, p, params.TTAKParams.J)
	if err != nil {
		return err
	}

	params.P = p
	params.Q = q
	params.G = g
	return nil
}

func GenerateKeyTTAK(priv *PrivateKey, rand io.Reader, xkey, upri []byte, sizes ParameterSizes) (xkeyOut, upriOut []byte, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, nil, ErrInvalidParameterSizes
	}

	if priv.P == nil || priv.Q == nil || priv.G == nil {
		return nil, nil, ErrParametersNotSetUp
	}

	if len(xkey) == 0 {
		xkey = internal.Expand(xkey, internal.Bytes(domain.B))
		xkey, err = internal.ReadBits(xkey, rand, domain.B)
		if err != nil {
			return nil, nil, err
		}
	} else if len(xkey) < internal.Bytes(domain.B) {
		return nil, nil, ErrShortXKEY
	}
	if len(upri) == 0 {
		upri = internal.Expand(upri, 64)
		_, err = io.ReadFull(rand, upri)
		if err != nil {
			return nil, nil, err
		}
	}

	h := domain.NewHash()

	priv.X = generateX(priv.Q, upri, xkey, h, domain)
	priv.Y = kcdsainternal.Y(priv.P, priv.Q, priv.G, priv.X)

	return xkey, upri, nil
}

func ppgf(in []byte, seed []byte, nBits int, h hash.Hash) []byte {
	// p.12
	// from java
	i := internal.Bytes(nBits)
	in = internal.Expand(in, i)

	LH := h.Size()

	count := 0
	iBuf := make([]byte, 1)
	hbuf := make([]byte, 0, LH)

	for {
		iBuf[0] = byte(count)

		h.Reset()
		h.Write(seed)
		h.Write(iBuf)
		hbuf = h.Sum(hbuf[:0])

		if i >= LH {
			i -= LH
			copy(in[i:], hbuf)
			if i == 0 {
				break
			}
		} else {
			copy(in, hbuf[len(hbuf)-i:])
			break
		}

		count++
	}

	return internal.TruncateLeft(in, nBits)
}

// performance issue of ppgf...
func generateJAlt(rand io.Reader, seed []byte, ubuf []byte, d kcdsainternal.Domain) (J *big.Int, UBytes []byte, ok bool, err error) {
	UBytes, err = internal.ReadBits(ubuf[:0], rand, d.A-d.B-4)
	if err != nil {
		return
	}

	U := new(big.Int).SetBytes(UBytes)

	// 3: U의 상위에 4 비트 '1000'을 붙이고 최하위 비트는 1로 만들어 이를 J로 둔다.
	// (J ← 2^(α-β-1) ∨ U ∨ 1)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("J ← 2^(α-β-1) ∨ U ∨ 1")
	J = big.NewInt(0b1)
	J.Lsh(J, uint(d.A-d.B-1))
	J.Or(J, U)
	J.Or(J, one)
	//fmt.Println("J = 0x" + hex.EncodeToString(J.Bytes()))

	// 4: 강한 소수 판정 알고리즘으로 J를 판정하여 소수가 아니면 단계 1로 간다.
	if !J.ProbablyPrime(internal.NumMRTests) {
		return
	}

	ok = true
	return
}

func generateJ(seed, UBytes []byte, h hash.Hash, d kcdsainternal.Domain) (J *big.Int, UBytes2 []byte, ok bool) {
	// 2: Seed를 일방향 함수 PPGF의 입력으로 하여 비트 길이가 n = (α - β - 4)인 난수 U를 생성한다.
	// (U ← PPGF(Seed, n))
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("U ← PPGF(Seed, n)")
	U := new(big.Int).SetBytes(ppgf(UBytes[:0], seed, d.A-d.B-4, h))
	//fmt.Println(U.BitLen())
	//fmt.Println("U = 0x" + hex.EncodeToString(U.Bytes()))

	// 3: U의 상위에 4 비트 '1000'을 붙이고 최하위 비트는 1로 만들어 이를 J로 둔다.
	// (J ← 2^(α-β-1) ∨ U ∨ 1)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("J ← 2^(α-β-1) ∨ U ∨ 1")
	J = big.NewInt(0b1)
	J.Lsh(J, uint(d.A-d.B-1))
	J.Or(J, U)
	J.Or(J, one)
	//fmt.Println("J = 0x" + hex.EncodeToString(J.Bytes()))

	// 4: 강한 소수 판정 알고리즘으로 J를 판정하여 소수가 아니면 단계 1로 간다.
	if !J.ProbablyPrime(internal.NumMRTests) {
		return
	}

	ok = true
	return
}

func generatePQ(J *big.Int, seed []byte, h hash.Hash, d kcdsainternal.Domain) (p, q *big.Int, count int, ok bool) {
	// 5: Count를 0으로 둔다. (Count ← 0)
	count = 0

	seedCount := make([]byte, len(seed)+4)
	copy(seedCount, seed)

	q = new(big.Int)
	p = new(big.Int)

	uBuf := make([]byte, internal.Bytes(d.B))

	// 7: Count > 2^24이면 단계 1로 간다.
	for count <= (1 << 24) {
		// 6: Count를 1 증가시킨다. (Count ← (Count + 1))
		count += 1
		binary.BigEndian.PutUint32(seedCount[len(seedCount)-4:], uint32(count))

		// 8: Seed에 Count를 연접한 것을 일방향 함수 PPGF의 입력으로 하여 비트 길이가
		// β인 난수 U를 생성한다. (U ← PPGF(Seed ‖ Count, β))
		U := ppgf(uBuf[:0], seedCount, d.B, h)

		// 9: U의 최상위 및 최하위 비트를 1로 만들어 이를 q로 둔다.
		// (q ← 2^(β-1) ∨ U ∨ 1)
		U[0] |= 0b1000_0000
		U[len(U)-1] |= 1
		q.SetBytes(U)

		// 10: p ← (2Jq + 1)의 비트 길이가 α보다 길면 단계 6으로 간다.
		p.Add(p.Lsh(p.Mul(J, q), 1), one)
		if p.BitLen() > d.A {
			continue
		}

		// 11: 강한 소수 판정 알고리즘으로 q를 판정하여 소수가 아니면 단계 6으로 간다.
		if !q.ProbablyPrime(internal.NumMRTests) {
			continue
		}

		// 12: 강한 소수 판정 알고리즘으로 p를 판정하여 소수가 아니면 단계 6으로 간다
		if !p.ProbablyPrime(internal.NumMRTests) {
			continue
		}

		ok = true
		return
	}

	return
}

func generateHG(randReader io.Reader, P, J *big.Int) (H *big.Int, G *big.Int, err error) {
	pSub3 := new(big.Int).Sub(P, three)

	for {
		// 1: p보다 작은 임의의 수 h를 생성한다.
		// 1 < h < (p - 1)
		//     -1 < h < p - 3
		//         is same with 0 <= h < p-3
		//     than, h + 2
		H, err = rand.Int(randReader, pSub3)
		if err != nil {
			return
		}
		H.Add(H, two)

		G, ok := generateG(P, J, H)
		if !ok {
			continue
		}

		return H, G, nil
	}
}
func generateG(P, J *big.Int, H *big.Int) (G *big.Int, ok bool) {
	// 2: g ← h^(2J) mod p를 계산한다.
	g := new(big.Int).Set(J)
	g.Lsh(g, 1)
	g.Exp(H, g, P)

	// 3: g = 1이면 단계 1로 간다.
	if g.Cmp(one) == 0 {
		return nil, false
	}

	return g, true
}

func generateX(Q *big.Int, upri, xkey []byte, h hash.Hash, d kcdsainternal.Domain) (X *big.Int) {
	// 3: XSEEDj ← PPGF(user_provided_random_input, b)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("3: XSEEDj ← PPGF(user_provided_random_input, b)")
	xseed := ppgf(nil, upri, d.B, h)
	//fmt.Println("xseed = 0x" + hex.EncodeToString(xseed))

	// 4: XVAL ← (XKEY + XSEEDj) mod 2^b
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("4: XVAL ← (XKEY + XSEEDj) mod 2^b")
	var carry int
	xval := make([]byte, internal.Bytes(d.B))
	for i := 0; i < len(xseed); i++ {
		idx := len(xseed) - i - 1
		sum := int(xseed[idx]) + carry

		if i < len(xkey) {
			sum += int(xkey[len(xkey)-i-1])
		}

		xval[idx] = byte(sum)
		carry = sum >> 8
	}
	xval = internal.TruncateLeft(xval, d.B)
	//fmt.Println("xval = 0x" + hex.EncodeToString(xval.Bytes()))

	// 5: xj ← PPGF(XVAL, b) mod q
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("5: xj ← PPGF(XVAL, b) mod q")
	tmp := xseed
	X = new(big.Int).SetBytes(ppgf(tmp, xval, d.B, h))
	X.Mod(X, Q)
	//fmt.Println("X = 0x" + hex.EncodeToString(X.Bytes()))

	return X
}
