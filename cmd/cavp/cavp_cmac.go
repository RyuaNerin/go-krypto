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
		cvl := cavp.ReadValues()

		if cvl.ContainsKey("K") {
			K := cvl.Hex("K")
			M := cvl.Hex("M")
			TLen := cvl.Int("Tlen")

			b, err := newCipher(K)
			if err != nil {
				panic(err)
			}

			cmac := kipher.NewCMAC(b)
			cmac.Write(M)
			buf = cmac.Sum(buf[:0])[:TLen/8]

			cvl = append(cvl, cavpRow{"T", hexStr(buf), false})
		}

		cavp.WriteValues(cvl)
	}
}

func processCMAC_Ver(cavp *cavpProcessor, newCipher funcNewBlockCipher) {
	buf := make([]byte, 0, 64)

	for cavp.Next() {
		cvl := cavp.ReadValues()

		if cvl.ContainsKey("K") {
			K := cvl.Hex("K")
			M := cvl.Hex("M")
			TLen := cvl.Int("Tlen")
			T := cvl.Hex("T")

			b, err := newCipher(K)
			if err != nil {
				panic(err)
			}

			cmac := kipher.NewCMAC(b)
			cmac.Write(M)
			buf = cmac.Sum(buf[:0])[:TLen/8]

			if kipher.Equal(buf, T) {
				cvl = append(cvl, cavpRow{"", "VALID", false})
			} else {
				cvl = append(cvl, cavpRow{"", "INVALID", false})
			}
		}

		cavp.WriteValues(cvl)
	}
}
