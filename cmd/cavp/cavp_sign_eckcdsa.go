package main

import (
	"crypto/elliptic"
	"crypto/rand"
	"hash"
	"log"
	"strings"

	"github.com/RyuaNerin/go-krypto/eckcdsa"
)

func processECKCDSA(path, filename string) {
	hashInfo := getHash(filename)
	if hashInfo == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}
	h := hashInfo.New()

	curve := getECKCDSACurve(filename)
	if curve == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	switch {
	case strings.HasSuffix(filename, "KPG.REQ"):
		processECKCDSA_KPG(cavp, curve)

	case strings.HasSuffix(filename, "PKV.REQ"):
		processECKCDSA_PKV(cavp, curve)

	case strings.HasSuffix(filename, "SGT.REQ"):
		processECKCDSA_SGT(cavp, h, curve)

	case strings.HasSuffix(filename, "SVT.REQ"):
		processECKCDSA_SVT(cavp, h, curve)
	}
}

// 키 쌍 생성
//

func processECKCDSA_KPG(cavp *cavpProcessor, curve elliptic.Curve) {
	/**
	.req
		[B-233]
		d = ?
		10개의 키(d)를 생성 후 시험기관에 제출바랍니다.
		rsp 파일 형식은 sam 파일 형식을 유지하여 작성바랍니다.

	.sam
		[B-233]
		d = ?
		Qx = ?
		Qy = ?

		d = ?
		Qx = ?
		Qy = ?

		...
	*/
	cavp.WriteLine(cavp.ReadLine())

	xyBits := curve.Params().BitSize

	res := cavpSection{
		{"d", "", false},
		{"Qx", "", false},
		{"Qy", "", false},
	}
	for i := 0; i < 10; i++ {
		pk, err := eckcdsa.GenerateKey(curve, rnd)
		if err != nil {
			panic(err)
		}

		res[0].Value = hexInt(pk.D, xyBits)
		res[1].Value = hexInt(pk.X, xyBits)
		res[2].Value = hexInt(pk.Y, xyBits)

		cavp.WriteValues(res)
	}
}

// 타원곡선 위에 있는지 확인
//

func processECKCDSA_PKV(cavp *cavpProcessor, curve elliptic.Curve) {
	/**
	#  ECKCDSA

	Qx = ...
	Qy = ...

	Qx = ...
	Qy = ...

	...
	*/
	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("Qx") {
			if curve.IsOnCurve(cs.BigInt("Qx"), cs.BigInt("Qy")) {
				cs = append(cs, cavpRow{"Result", "P", false})
			} else {
				cs = append(cs, cavpRow{"Result", "F", false})
			}
		}
		cavp.WriteValues(cs)
	}
}

// Sign
//

func processECKCDSA_SGT(cavp *cavpProcessor, h hash.Hash, curve elliptic.Curve) {
	/**
	.req
		[B-233, SHA-224]
		M = ...

		M = ...

		...

	.sam
		[B-233, SHA-224]
		M = ...
		Qx = ?
		Qy = ?
		R = ?
		S = ?

		M = ...
		Qx = ?
		Qy = ?
		R = ?
		S = ?

		...
	*/
	w := curve.Params().BitSize

	bitsXY := w
	bitsR := h.Size() * 8
	bitsS := w

	if bitsR > bitsS {
		bitsR = bitsS
	}

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("M") {
			privKey, err := eckcdsa.GenerateKey(curve, rand.Reader)
			if err != nil {
				panic(err)
			}

			r, s, err := eckcdsa.Sign(rnd, privKey, h, cs.Hex("M"))
			if err != nil {
				panic(err)
			}

			cs = append(
				cs,
				cavpRow{"Qx", hexInt(privKey.X, bitsXY), false},
				cavpRow{"Qy", hexInt(privKey.Y, bitsXY), false},
				cavpRow{"R", hexInt(r, bitsR), false},
				cavpRow{"S", hexInt(s, bitsS), false},
			)
		}

		cavp.WriteValues(cs)
	}
}

// Verify
//

func processECKCDSA_SVT(cavp *cavpProcessor, h hash.Hash, curve elliptic.Curve) {
	/**
	req
		#  ECKCDSA

		[B-233, SHA-224]
		M = ...
		Qx = ...
		Qy = ...
		R = ...
		S = ...

		M = ...
		Qx = ...
		Qy = ..
		R = ...
		S = ...

		...

	sam
		#  ECKCDSA

		[B-233, SHA-224]
		M = ...
		Qx = ...
		Qy = ...
		R = ...
		S = ...
		Result = P or F

		M = ...
		Qx = ...
		Qy = ...
		R = ...
		S = ...
		Result = P or F

		...
	*/
	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("M") {
			pubKey := eckcdsa.PublicKey{
				Curve: curve,
				X:     cs.BigInt("Qx"),
				Y:     cs.BigInt("Qy"),
			}

			M := cs.Hex("M")
			r := cs.BigInt("R")
			s := cs.BigInt("S")

			if eckcdsa.Verify(&pubKey, h, M, r, s) {
				cs = append(cs, cavpRow{"Result", "P", false})
			} else {
				cs = append(cs, cavpRow{"Result", "F", false})
			}
		}

		cavp.WriteValues(cs)
	}
}
