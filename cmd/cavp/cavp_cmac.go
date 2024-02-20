package main

import (
	"log"
	"strings"

	"github.com/RyuaNerin/go-krypto/kipher"
)

func processCMAC(path, filename string) {
	newCipher := getBlock(filename)
	if newCipher == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	switch {
	case strings.HasSuffix(filename, "GEN.REQ"):
		processCMAC_Gen(cavp, newCipher)

	case strings.HasSuffix(filename, "VER.REQ"):
		processCMAC_Ver(cavp, newCipher)

	default:
		log.Println("Unknown algorithm: ", filename)
	}
}

func processCMAC_Gen(cavp *cavpProcessor, newCipher funcNewBlockCipher) {
	buf := make([]byte, 0, 64)

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("K") {
			K := cs.Hex("K")
			M := cs.Hex("M")
			TLen := cs.Int("Tlen")

			b, err := newCipher(K)
			if err != nil {
				panic(err)
			}

			cmac := kipher.NewCMAC(b)
			cmac.Write(M)
			buf = cmac.Sum(buf[:0])[:TLen/8]

			cs = append(cs, cavpRow{"T", hexStr(buf), false})
		}

		cavp.WriteValues(cs)
	}
}

func processCMAC_Ver(cavp *cavpProcessor, newCipher funcNewBlockCipher) {
	buf := make([]byte, 0, 64)

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("K") {
			K := cs.Hex("K")
			M := cs.Hex("M")
			TLen := cs.Int("Tlen")
			T := cs.Hex("T")

			b, err := newCipher(K)
			if err != nil {
				panic(err)
			}

			cmac := kipher.NewCMAC(b)
			cmac.Write(M)
			buf = cmac.Sum(buf[:0])[:TLen/8]

			if kipher.Equal(buf, T) {
				cs = append(cs, cavpRow{"", "VALID", false})
			} else {
				cs = append(cs, cavpRow{"", "INVALID", false})
			}
		}

		cavp.WriteValues(cs)
	}
}
