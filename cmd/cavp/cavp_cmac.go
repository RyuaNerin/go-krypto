package main

import (
	"log"
	"strings"

	"github.com/RyuaNerin/go-krypto/kipher"
)

func processCMAC(path, filename string) {
	var fnCipher newCipher
	for substr, v := range prefixBlocks {
		if strings.Contains(filename, substr) {
			fnCipher = v
			break
		}
	}
	if fnCipher == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	switch {
	case strings.HasSuffix(filename, "Gen.req"):
		processCMAC_Gen(cavp, fnCipher)
	case strings.HasSuffix(filename, "Ver.req"):
		processCMAC_Ver(cavp, fnCipher)

	default:
		log.Println("Unknown algorithm: ", filename)
	}
}

func processCMAC_Gen(cavp *cavpProcessor, fnCipher newCipher) {
	buf := make([]byte, 0, 64)

	for cavp.Next() {
		cvl := cavp.ReadValues()

		if cvl.Contains("K") {
			K := cvl.Hex("K")
			M := cvl.Hex("M")
			TLen := cvl.Int("Tlen")

			b, err := fnCipher(K)
			if err != nil {
				panic(err)
			}

			cmac := kipher.NewCMAC(b, TLen/8)
			cmac.Write(M)
			buf = cmac.Sum(buf[:0])

			cvl = append(cvl, cavpValue{"T", hexStr(buf)})
		}

		cavp.WriteValues(cvl)
	}
}

func processCMAC_Ver(cavp *cavpProcessor, fnCipher newCipher) {
	buf := make([]byte, 0, 64)

	for cavp.Next() {
		cvl := cavp.ReadValues()

		if cvl.Contains("K") {
			K := cvl.Hex("K")
			M := cvl.Hex("M")
			TLen := cvl.Int("Tlen")
			T := cvl.Hex("T")

			b, err := fnCipher(K)
			if err != nil {
				panic(err)
			}

			cmac := kipher.NewCMAC(b, TLen/8)
			cmac.Write(M)
			buf = cmac.Sum(buf[:0])

			if kipher.Equal(buf, T) {
				cvl = append(cvl, cavpValue{"", "VALID"})
			} else {
				cvl = append(cvl, cavpValue{"", "INVALID"})
			}
		}

		cavp.WriteValues(cvl)
	}
}
