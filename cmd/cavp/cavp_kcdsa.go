package main

import (
	"log"
	"strconv"
	"strings"

	kcdsainternal "github.com/RyuaNerin/go-krypto/internal/kcdsa"
	"github.com/RyuaNerin/go-krypto/kcdsa"
)

func processKCDSA(path, filename string) {
	sz := getKCDSA(filename)
	if sz == -1 {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	switch {
	case strings.HasSuffix(filename, "KPG.REQ"): // 키 생성
		processKCDSA_KPG(cavp, sz)

	case strings.HasSuffix(filename, "SGT.REQ"): // Sign
		processKCDSA_SGT(cavp, sz)

	case strings.HasSuffix(filename, "SVT.REQ"): // Verify
		processKCDSA_SVT(cavp, sz)
	}
}

// 키 쌍 생성
func processKCDSA_KPG(cavp *cavpProcessor, sz kcdsa.ParameterSizes) {
	/**
	req
		|P| = 2048
		|Q| = 224
		10개의 키 및 관련 생성값 (P, Q, G, X, Y, Seed, Count, XKEY, OUPRI, h, J) 생성 후 시험기관에 제출바랍니다.
		rsp 파일 형식은 sam 파일 형식을 유지하여 작성바랍니다.

	sam
		|P| = 2048
		|Q| = 224

		P = ?
		Q = ?
		G = ?
		X = ?
		Y = ?
		Seed = ?
		Count = ?
		XKEY = ?
		OUPRI = ?
		h = ?
		J = ?

		P = ?
		Q = ?
		G = ?
		X = ?
		Y = ?
		Seed = ?
		Count = ?
		XKEY = ?
		OUPRI = ?
		h = ?
		J = ?

		...
	*/
	domain, _ := kcdsainternal.GetDomain(int(sz))

	bitsP := domain.A
	bitsQ := domain.B
	bitsG := domain.A
	bitsX := domain.B
	bitsY := domain.A
	bitsH := domain.A
	bitsJ := domain.A - domain.B

	cavp.WriteLine(cavp.ReadLine())
	cavp.WriteLine(cavp.ReadLine())
	cavp.WriteLine("")

	res := cavpSection{
		{"P", "", false},
		{"Q", "", false},
		{"G", "", false},
		{"X", "", false},
		{"Y", "", false},
		{"Seed", "", false},
		{"Count", "", false},
		{"XKEY", "", false},
		{"OUPRI", "", false},
		{"h", "", false},
		{"J", "", false},
	}

	for i := 0; i < 10; i++ {
		domain, _ := kcdsainternal.GetDomain(int(sz))

		generated, err := kcdsainternal.GenerateParameters(rnd, domain)
		if err != nil {
			panic(err)
		}
		privKey := kcdsa.PrivateKey{
			PublicKey: kcdsa.PublicKey{
				Parameters: kcdsa.Parameters{
					P: generated.P,
					Q: generated.Q,
					G: generated.G,
					TTAKParams: kcdsa.TTAKParameters{
						J:     generated.J,
						Seed:  generated.Seed,
						Count: generated.Count,
					},
				},
			},
		}

		XKEY, OUPRI, err := kcdsa.GenerateKeyWithSeed(&privKey, rnd, nil, nil, sz)
		if err != nil {
			panic(err)
		}

		res[0x0].Value = hexInt(privKey.P, bitsP)               // P
		res[0x1].Value = hexInt(privKey.Q, bitsQ)               // Q
		res[0x2].Value = hexInt(privKey.G, bitsG)               // G
		res[0x3].Value = hexInt(privKey.X, bitsX)               // X
		res[0x4].Value = hexInt(privKey.Y, bitsY)               // Y
		res[0x5].Value = hexStr(privKey.TTAKParams.Seed)        // Seed
		res[0x6].Value = strconv.Itoa(privKey.TTAKParams.Count) // Count
		res[0x7].Value = hexStr(XKEY)                           // XKEY
		res[0x8].Value = hexStr(OUPRI)                          // OUPRI
		res[0x9].Value = hexInt(generated.H, bitsH)             // h
		res[0xA].Value = hexInt(privKey.TTAKParams.J, bitsJ)    // J

		cavp.WriteValues(res)
	}
}

// Sign
func processKCDSA_SGT(cavp *cavpProcessor, sz kcdsa.ParameterSizes) {
	/**
	req
		|P| = 2048
		|Q| = 224
		SHAAlg = SHA-224

		M = ...

		M = ...

		...

	sam
		|P| = 2048
		|Q| = 224
		SHAAlg = SHA-224

		P = ?
		Q = ?
		G = ?
		M = ...
		X = ?
		Y = ?
		R = ?
		S = ?

		P = ?
		Q = ?
		G = ?
		M = ...
		X = ?
		Y = ?
		R = ?
		S = ?

		...
	*/

	domain, _ := kcdsainternal.GetDomain(int(sz))

	bitsP := domain.A
	bitsQ := domain.B
	bitsG := domain.A
	bitsX := domain.B
	bitsY := domain.A

	bitsR := domain.LH * 8
	bitsS := domain.B
	if bitsR > domain.B {
		bitsR = domain.B
	}

	res := cavpSection{
		{"P", "", false},
		{"Q", "", false},
		{"G", "", false},
		{"M", "", false},
		{"X", "", false},
		{"Y", "", false},
		{"R", "", false},
		{"S", "", false},
	}

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("M") {
			var priv kcdsa.PrivateKey

			if err := kcdsa.GenerateParameters(&priv.Parameters, rnd, sz); err != nil {
				panic(err)
			}
			if err := kcdsa.GenerateKey(&priv, rnd); err != nil {
				panic(err)
			}
			M := cs.Hex("M")

			r, s, err := kcdsa.Sign(rnd, &priv, sz, M)
			if err != nil {
				panic(err)
			}

			res[0].Value = hexInt(priv.P, bitsP) // P
			res[1].Value = hexInt(priv.Q, bitsQ) // Q
			res[2].Value = hexInt(priv.G, bitsG) // G
			res[3].Value = hexStr(M)             // M
			res[4].Value = hexInt(priv.X, bitsX) // X
			res[5].Value = hexInt(priv.Y, bitsY) // Y
			res[6].Value = hexInt(r, bitsR)      // R
			res[7].Value = hexInt(s, bitsS)      // S

			cs = res
		}

		cavp.WriteValues(cs)
	}
}

// Verify
func processKCDSA_SVT(cavp *cavpProcessor, sz kcdsa.ParameterSizes) {
	/**
	req
		|P| = 2048, |Q| = 224
		SHAAlg = SHA-224

		P = ...
		Q = ...
		G = ...
		M = ...
		X = ...
		Y = ...
		R = ...
		S = ...

		P = ...
		Q = ...
		G = ...
		M = ...
		X = ...
		Y = ...
		R = ...
		S = ...

		...

	sam
		|P| = 2048, |Q| = 224
		SHAAlg = SHA-224

		P = ...
		Q = ...
		G = ...
		M = ...
		X = ...
		Y = ...
		R = ...
		S = ...
		Result = P or F

		P = ...
		Q = ...
		G = ...
		M = ...
		X = ...
		Y = ...
		R = ...
		S = ...
		Result = P or F

		...
	*/
	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("M") {
			pub := kcdsa.PublicKey{
				Parameters: kcdsa.Parameters{
					P: cs.BigInt("P"),
					Q: cs.BigInt("Q"),
					G: cs.BigInt("G"),
				},
				Y: cs.BigInt("Y"),
			}

			M := cs.Hex("M")
			r := cs.BigInt("R")
			s := cs.BigInt("S")

			if kcdsa.Verify(&pub, sz, M, r, s) {
				cs = append(cs, cavpRow{"Result", "P", false})
			} else {
				cs = append(cs, cavpRow{"Result", "F", false})
			}
		}

		cavp.WriteValues(cs)
	}
}
