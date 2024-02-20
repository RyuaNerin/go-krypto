package main

import (
	"log"
	"strings"

	"github.com/RyuaNerin/go-krypto/kdf/kbkdf"
)

type funcKBKDF func(prf kbkdf.PRF, key, label, context, iv []byte, r, length int) []byte

func processKBKDF(path, filename string) {
	cavp := NewCavp(path, filename)
	defer cavp.Close()

	counterMode := func(prf kbkdf.PRF, key, label, context, _ []byte, r, length int) []byte {
		return kbkdf.CounterMode(prf, key, label, context, r, length)
	}
	doublePipeMode := func(prf kbkdf.PRF, key, label, context, _ []byte, r, length int) []byte {
		return kbkdf.DoublePipeMode(prf, key, label, context, r, length)
	}
	feedbackMode := kbkdf.FeedbackMode

	switch {
	case strings.Contains(filename, "HMAC"):
		hashInfo := getHash(cavp.filename)
		if hashInfo == nil {
			log.Println("Unknown algorithm: ", cavp.filename)
			return
		}

		switch {
		case strings.Contains(filename, "(CTR)") || strings.Contains(filename, "-CTR"):
			processKBKDF_HMAC(cavp, hashInfo, counterMode)

		case strings.Contains(filename, "(DP)") || strings.Contains(filename, "-DP"):
			processKBKDF_HMAC(cavp, hashInfo, doublePipeMode)

		case strings.Contains(filename, "(FB)") || strings.Contains(filename, "-FB"):
			processKBKDF_HMAC(cavp, hashInfo, feedbackMode)

		default:
			log.Println("Unknown algorithm: ", filename)
		}

	case strings.Contains(filename, "CMAC"):
		newCipher := getBlock(cavp.filename)
		if newCipher == nil {
			log.Println("Unknown algorithm: ", cavp.filename)
			return
		}

		switch {
		case strings.Contains(filename, "(CTR)") || strings.Contains(filename, "-CTR"):
			processKBKDF_CMAC(cavp, newCipher, counterMode)

		case strings.Contains(filename, "(DP)") || strings.Contains(filename, "-DP"):
			processKBKDF_CMAC(cavp, newCipher, doublePipeMode)

		case strings.Contains(filename, "(FB)") || strings.Contains(filename, "-FB"):
			processKBKDF_CMAC(cavp, newCipher, feedbackMode)

		default:
			log.Println("Unknown algorithm: ", filename)
		}

	default:
		log.Println("Unknown algorithm: ", filename)
	}
}

func processKBKDF_HMAC(cavp *cavpProcessor, hashInfo *HashInfo, fn funcKBKDF) {
	prf := kbkdf.NewHMACPRF(hashInfo.New)

	var rlen int
	var dst []byte

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("PRF") {
			if cs.ContainsValue("NO COUNTER") {
				rlen = 0
			} else {
				rlen = cs.Int("RLEN")
			}
		} else if cs.ContainsKey("COUNT") {
			L := cs.Int("L")
			KI := cs.Hex("KI")
			Label := cs.Hex("Label")
			Context := cs.Hex("Context")

			var IV []byte
			if cs.ContainsKey("IV") {
				IV = cs.Hex("IV")
			}

			dst = fn(prf, KI, Label, Context, IV, rlen/8, L/8)
			cs = append(cs, cavpRow{"KO", hexStr(dst), false})
		}

		cavp.WriteValues(cs)
	}
}

func processKBKDF_CMAC(cavp *cavpProcessor, newCipher funcNewBlockCipher, fn funcKBKDF) {
	var rlen int
	var dst []byte

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("PRF") {
			if cs.ContainsValue("NO COUNTER") {
				rlen = 0
			} else {
				rlen = cs.Int("RLEN")
			}
		} else if cs.ContainsKey("COUNT") {
			L := cs.Int("L")
			KI := cs.Hex("KI")
			Label := cs.Hex("Label")
			Context := cs.Hex("Context")

			var IV []byte
			if cs.ContainsKey("IV") {
				IV = cs.Hex("IV")
			}

			prf := kbkdf.NewCMACPRF(newCipher)
			dst = fn(prf, KI, Label, Context, IV, rlen/8, L/8)
			cs = append(cs, cavpRow{"KO", hexStr(dst), false})
		}

		cavp.WriteValues(cs)
	}
}
