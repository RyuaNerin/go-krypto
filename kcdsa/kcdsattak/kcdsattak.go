// Package kcdsattak implements functions what generate the KCDSA parameters as defined in TTAK.KO-12.0001/R4
package kcdsattak

import (
	"encoding/binary"
	"errors"
	"io"
	"math/big"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/kcdsa"
	kcdsainternal "github.com/RyuaNerin/go-krypto/kcdsa/internal"
)

var (
	ErrInvalidPublicKey      = kcdsa.ErrInvalidPublicKey      // Same with kcdsa.ErrInvalidPublicKey
	ErrInvalidParameterSizes = kcdsa.ErrInvalidParameterSizes // Same with kcdsa.ErrInvalidParameterSizes
	ErrParametersNotSetUp    = kcdsa.ErrParametersNotSetUp    // Same with kcdsa.ErrParametersNotSetUp

	ErrUseAnotherSeed = errors.New("krypto/kcdsa/kcdsattak: use another seed")
	ErrUseAnotherH    = errors.New("krypto/kcdsa/kcdsattak: use another H")
	ErrWrongSeed      = errors.New("krypto/kcdsa/kcdsattak: wrong seed length")
	ErrShortXKey      = errors.New("krypto/kcdsa/kcdsattak: XKEY is too short")

	one = big.NewInt(1)
	two = big.NewInt(2)
)

// Pre-prime generating function
func PPGF(seed []byte, nBits int, sizes kcdsa.ParameterSizes) ([]byte, error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, kcdsa.ErrInvalidParameterSizes
	}

	return ppgf(nil, seed, nBits, domain), nil
}
func ppgf(U []byte, seed []byte, nBits int, domain kcdsainternal.Domain) []byte {
	// p.12
	// from java
	i := internal.Bytes(nBits)
	iBuf := make([]byte, 1)

	if len(U) < i {
		U = U[:cap(U)]
		for len(U) < i {
			U = append(U, make([]byte, i)...)
		}
	}

	h := domain.NewHash()

	var hbuf []byte
	count := 0
	for {
		iBuf[0] = byte(count)

		h.Reset()
		h.Write(seed)
		h.Write(iBuf)
		hbuf = h.Sum(hbuf[:0])

		if i >= domain.LH {
			i -= domain.LH
			for j := 0; j < domain.LH; j++ {
				U[j+i] = hbuf[j]
			}
			if i == 0 {
				break
			}
		} else {
			for j := 0; j < i; j++ {
				U[j] = hbuf[j+domain.LH-i]
			}
			break
		}

		count++
	}

	i = nBits & 0x07
	if i != 0 {
		U[0] &= byte((1 << i) - 1)
	}

	return U
}

// Generate J, defined in TTAK.KO-12.0001/R4
// bits of seed > domain.B
func GenerateJ(seed []byte, sizes kcdsa.ParameterSizes) (J *big.Int, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, kcdsa.ErrInvalidParameterSizes
	}
	return generateJ(seed, domain)
}
func generateJ(seed []byte, domain kcdsainternal.Domain) (J *big.Int, err error) {
	// p.14
	if len(seed) != internal.Bytes(domain.B) {
		return nil, ErrWrongSeed
	}

	// 2: Seed를 일방향 함수 PPGF의 입력으로 하여 비트 길이가 n = (α - β - 4)인 난수 U를 생성한다.
	// (U ← PPGF(Seed, n))
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("U ← PPGF(Seed, n)")
	U := new(big.Int).SetBytes(ppgf(nil, seed, domain.A-domain.B-4, domain))
	//fmt.Println(U.BitLen())
	//fmt.Println("U = 0x" + hex.EncodeToString(U.Bytes()))

	// 3: U의 상위에 4 비트 '1000'을 붙이고 최하위 비트는 1로 만들어 이를 J로 둔다.
	// (J ← 2^(α-β-1) ∨ U ∨ 1)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("J ← 2^(α-β-1) ∨ U ∨ 1")
	J = new(big.Int).SetInt64(0b1)
	J.Lsh(J, uint(domain.A-domain.B-1))
	J.Or(J, U)
	J.Or(J, one)
	//fmt.Println("J = 0x" + hex.EncodeToString(J.Bytes()))

	// 4: 강한 소수 판정 알고리즘으로 J를 판정하여 소수가 아니면 단계 1로 간다.
	if !J.ProbablyPrime(internal.NumMRTests) {
		return nil, ErrUseAnotherSeed
	}
	return J, nil
}

// Generate P, Q, defined in TTAK.KO-12.0001/R4
func GeneratePQ(J *big.Int, seed []byte, sizes kcdsa.ParameterSizes) (p, q *big.Int, count int, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, nil, 0, kcdsa.ErrInvalidParameterSizes
	}
	return generatePQ(J, seed, domain)
}
func generatePQ(J *big.Int, seed []byte, domain kcdsainternal.Domain) (p, q *big.Int, count int, err error) {
	// p.14
	if len(seed) != internal.Bytes(domain.B) {
		return nil, nil, 0, ErrWrongSeed
	}

	// 5: Count를 0으로 둔다. (Count ← 0)
	count = 0

	seedCount := make([]byte, len(seed)+4)
	copy(seedCount, seed)

	q = new(big.Int)
	p = new(big.Int)

	uBuf := make([]byte, internal.Bytes(domain.B))

	// 7: Count > 2^24이면 단계 1로 간다.
	for count <= (1 << 24) {
		// 6: Count를 1 증가시킨다. (Count ← (Count + 1))
		count += 1
		binary.BigEndian.PutUint32(seedCount[len(seedCount)-4:], uint32(count))

		// 8: Seed에 Count를 연접한 것을 일방향 함수 PPGF의 입력으로 하여 비트 길이가
		// β인 난수 U를 생성한다. (U ← PPGF(Seed ‖ Count, β))
		U := ppgf(uBuf[:0], seedCount, domain.B, domain)

		// 9: U의 최상위 및 최하위 비트를 1로 만들어 이를 q로 둔다.
		// (q ← 2^(β-1) ∨ U ∨ 1)
		U[0] |= 0b1000_0000
		U[len(U)-1] |= 1
		q.SetBytes(U)

		// 10: p ← (2Jq + 1)의 비트 길이가 α보다 길면 단계 6으로 간다.
		p.Add(p.Lsh(p.Mul(J, q), 1), one)
		if p.BitLen() > domain.A {
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

		return
	}

	return nil, nil, 0, ErrUseAnotherSeed
}

// Generate H, G, defined in TTAK.KO-12.0001/R4
func GenerateHG(rand io.Reader, P, J *big.Int) (H []byte, G *big.Int, err error) {
	pm1 := new(big.Int).Set(P)
	pm1.Sub(pm1, one)

	hInt := new(big.Int)
	for {
		H, err = internal.ReadBits(H, rand, P.BitLen())
		if err != nil {
			return nil, nil, err
		}
		hInt.Mod(hInt.Add(hInt.SetBytes(H), two), pm1)

		G, err := generateG(P, J, H, pm1)
		if err != nil {
			continue
		}

		return H, G, nil
	}
}

// Generate P, defined in TTAK.KO-12.0001/R4
// 1 < H < (p-1)
func GenerateG(P, J *big.Int, H []byte) (G *big.Int, err error) {
	pm1 := new(big.Int).Set(P)
	pm1.Sub(pm1, one)

	return generateG(P, J, H, pm1)
}
func generateG(P, J *big.Int, H []byte, pm1 *big.Int) (G *big.Int, err error) {
	h := new(big.Int).SetBytes(H)

	// 1: p보다 작은 임의의 수 h를 생성한다.
	// 1 < h < (p - 1)
	if h.Cmp(one) != 1 || h.Cmp(pm1) != -1 {
		return nil, errors.New("H must be 1 < H < p-1")
	}

	// 2: g ← h^(2J) mod p를 계산한다.
	g := new(big.Int).Set(J)
	g.Lsh(g, 1)
	g.Exp(h, g, P)

	// 3: g = 1이면 단계 1로 간다.
	if g.Cmp(one) == 0 {
		return nil, ErrUseAnotherH
	}

	return g, nil
}

// Generate X, Y, Z, defined in TTAK.KO-12.0001/R4
// bits of xkey > B
func GenerateXYZ(P, Q, G *big.Int, userProvidedRandomInput []byte, xkey []byte, sizes kcdsa.ParameterSizes) (X, Y, Z *big.Int, xkeyNext []byte, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, nil, nil, nil, kcdsa.ErrInvalidParameterSizes
	}
	return generateXYZ(P, Q, G, userProvidedRandomInput, xkey, domain)
}
func generateXYZ(P, Q, G *big.Int, userProvidedRandomInput []byte, xkey []byte, domain kcdsainternal.Domain) (X, Y, Z *big.Int, xkeyNext []byte, err error) {
	// p.16
	if len(xkey) < internal.Bytes(domain.B) {
		return nil, nil, nil, nil, ErrShortXKey
	}

	i2b := new(big.Int).Lsh(one, uint(domain.B))
	i2l := new(big.Int).Lsh(one, uint(domain.L))

	ppgfBuf := make([]byte, 0, internal.Bytes(domain.B))

	// 3: XSEEDj ← PPGF(user_provided_random_input, b)
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("3: XSEEDj ← PPGF(user_provided_random_input, b)")
	xseed := new(big.Int).SetBytes(ppgf(ppgfBuf, userProvidedRandomInput, domain.B, domain))
	//fmt.Println("xseed = 0x" + hex.EncodeToString(xseed.Bytes()))

	// 4: XVAL ← (XKEY + XSEEDj) mod 2^b
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("4: XVAL ← (XKEY + XSEEDj) mod 2^b")
	xval := new(big.Int).SetBytes(xkey)
	xval.Add(xval, xseed)
	xval.Mod(xval, i2b)
	//fmt.Println("xval = 0x" + hex.EncodeToString(xval.Bytes()))

	// 5: xj ← PPGF(XVAL, b) mod q
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("5: xj ← PPGF(XVAL, b) mod q")
	X = new(big.Int).SetBytes(ppgf(ppgfBuf, xval.Bytes(), domain.B, domain))
	X.Mod(X, Q)
	//fmt.Println("X = 0x" + hex.EncodeToString(X.Bytes()))

	// 6: XKEY ← (XKEY + PPGF((xj + XSEEDj) mod 2^b, b)) mod 2^b
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("5: XKEY ← (XKEY + PPGF((xj + XSEEDj) mod 2^b, b)) mod 2^b")
	xkeyNextInt := new(big.Int).Set(X)
	xkeyNextInt.Mod(xkeyNextInt.Add(xkeyNextInt, xseed), i2b)

	xkeyNextInt.SetBytes(ppgf(ppgfBuf, xkeyNextInt.Bytes(), domain.B, domain))
	xkeyNextInt.Mod(xkeyNextInt.Add(xkeyNextInt, new(big.Int).SetBytes(xkey)), i2b)
	//fmt.Println("XKEY = 0x" + hex.EncodeToString(xkeyNext.Bytes()))
	xkeyNext = xkeyNextInt.FillBytes(xkey)

	// x′ = x^(-1) mod q
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("x′ = x^(-1) mod q")
	Xinv := internal.FermatInverse(X, Q)
	//fmt.Println("Xinv = 0x" + hex.EncodeToString(Xinv.Bytes()))

	// y = g^x′ mod p
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("y = g^x′ mod p")
	Y = new(big.Int).Exp(G, Xinv, P)
	//fmt.Println("Y = 0x" + hex.EncodeToString(Y.Bytes()))

	// z = y mod 2^ℓ
	//fmt.Println("--------------------------------------------------")
	//fmt.Println("z = y mod 2^ℓ")
	Z = new(big.Int).Mod(Y, i2l)
	//fmt.Println("Z = 0x" + hex.EncodeToString(Z.Bytes()))

	return
}

// Generate K defined in TTAK.KO-12.0001/R4
func GenerateK(randReader io.Reader, Q *big.Int, machine_generated_random_input []byte, sizes kcdsa.ParameterSizes) (K *big.Int, err error) {
	domain, ok := kcdsainternal.GetDomain(int(sizes))
	if !ok {
		return nil, kcdsa.ErrInvalidParameterSizes
	}
	return generateK(randReader, Q, machine_generated_random_input, domain)
}
func generateK(randReader io.Reader, Q *big.Int, machine_generated_random_input []byte, domain kcdsainternal.Domain) (K *big.Int, err error) {
	kkeyBytes := make([]byte, internal.Bytes(domain.B))
	if _, err = randReader.Read(kkeyBytes); err != nil {
		return nil, err
	}
	kkey := new(big.Int).SetBytes(kkeyBytes)

	kseedBytes := ppgf(nil, machine_generated_random_input, domain.B, domain)
	kseed := new(big.Int).SetBytes(kseedBytes)

	kval := kkey.Add(kkey, kseed).Bytes()
	kval = internal.TruncateLeft(kval, domain.B)

	KBytes := ppgf(nil, kval, domain.B, domain)
	K = new(big.Int).SetBytes(KBytes)
	K.Mod(K, Q)

	return
}
