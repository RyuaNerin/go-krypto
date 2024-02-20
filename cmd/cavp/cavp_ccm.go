package main

import (
	"log"
	"strings"

	"github.com/RyuaNerin/go-krypto/kipher"
)

func processCCM(path, filename string) {
	newCipher := getBlock(filename)
	if newCipher == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	switch {
	case strings.HasSuffix(filename, "GE.REQ"):
		processCCM_GE(cavp, newCipher)

	case strings.HasSuffix(filename, "DV.REQ"):
		processCCM_DV(cavp, newCipher)

	default:
		log.Println("Unknown algorithm: ", filename)
	}
}

func processCCM_GE(cavp *cavpProcessor, fnCipher funcNewBlockCipher) {
	buf := make([]byte, 0, 64)

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("K") {
			key := cs.Hex("K")
			nonce := cs.Hex("N")
			pt := cs.Hex("P")
			adata := cs.Hex("A")
			tlen := cs.Int("Tlen")

			b, err := fnCipher(key)
			if err != nil {
				panic(err)
			}

			aead, err := kipher.NewCCM(b, len(nonce), tlen/8)
			if err != nil {
				panic(err)
			}
			buf = aead.Seal(buf[:0], nonce, pt, adata)

			cs = append(
				cs,
				cavpRow{"C", hexStr(buf), false},
			)
		}

		cavp.WriteValues(cs)
	}
}

func processCCM_DV(cavp *cavpProcessor, fnCipher funcNewBlockCipher) {
	pt := make([]byte, 0, 64)

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("K") {
			key := cs.Hex("K")
			nonce := cs.Hex("N")
			adata := cs.Hex("A")
			c := cs.Hex("C")
			tlen := cs.Int("Tlen")

			b, err := fnCipher(key)
			if err != nil {
				panic(err)
			}

			aead, err := kipher.NewCCM(b, len(nonce), tlen/8)
			if err != nil {
				panic(err)
			}

			pt, err = aead.Open(pt[:0], nonce, c, adata)
			if err == nil {
				cs = append(cs, cavpRow{"P", hexStr(pt), false})
			} else {
				cs = append(cs, cavpRow{"", "INVALID", false})
			}
		}

		cavp.WriteValues(cs)
	}
}
