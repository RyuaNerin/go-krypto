package main

import (
	"log"
	"strings"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/kipher"
)

func processGCM(path, filename string) {
	newCipher := getBlock(filename)
	if newCipher == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	switch {
	case strings.HasSuffix(filename, "AE.REQ"):
		processGCM_AE(cavp, newCipher)
	case strings.HasSuffix(filename, "AD.REQ"):
		processGCM_AD(cavp, newCipher)

	default:
		log.Println("Unknown algorithm: ", filename)
	}
}

func processGCM_AE(cavp *cavpProcessor, newCipher funcNewBlockCipher) {
	buf := make([]byte, 0, 64)

	var keyLen, ivLen, ptLen, aadLen, tagLen int

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("Key") {
			key := cs.Hex("Key")
			iv := cs.Hex("IV")
			pt := cs.Hex("PT")
			adata := cs.Hex("Adata")

			key = internal.Grow(key, keyLen/8)
			iv = internal.Grow(iv, ivLen/8)
			pt = internal.Grow(pt, ptLen/8)
			adata = internal.Grow(adata, aadLen/8)

			b, err := newCipher(key)
			if err != nil {
				panic(err)
			}

			aead, err := kipher.NewGCMWithSize(b, ivLen/8, tagLen/8)
			if err != nil {
				panic(err)
			}
			buf = aead.Seal(buf[:0], iv, pt, adata)

			c := buf[:len(pt)]
			t := buf[len(pt):]

			cs = append(
				cs,
				cavpRow{"C", hexStr(c), false},
				cavpRow{"T", hexStr(t), false},
			)
		} else {
			keyLen = cs.Int("KeyLen")
			ivLen = cs.Int("IVLen")
			ptLen = cs.Int("PTLen")
			aadLen = cs.Int("AADLen")
			tagLen = cs.Int("TagLen")
		}

		cavp.WriteValues(cs)
	}
}

func processGCM_AD(cavp *cavpProcessor, newCipher funcNewBlockCipher) {
	ct := make([]byte, 0, 64)
	pt := make([]byte, 0, 64)

	var keyLen, ivLen, ptLen, aadLen, tagLen int

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("Key") {
			key := cs.Hex("Key")
			iv := cs.Hex("IV")
			adata := cs.Hex("Adata")
			c := cs.Hex("C")
			t := cs.Hex("T")

			key = internal.Resize(key, keyLen/8)
			iv = internal.Resize(iv, ivLen/8)
			adata = internal.Resize(adata, aadLen/8)
			c = internal.Resize(c, ptLen/8)
			t = internal.Resize(t, tagLen/8)

			ct = internal.Grow(ct, len(c)+len(t))
			copy(ct[:len(c)], c)
			copy(ct[len(c):], t)

			b, err := newCipher(key)
			if err != nil {
				panic(err)
			}

			aead, err := kipher.NewGCMWithSize(b, ivLen/8, tagLen/8)
			if err != nil {
				panic(err)
			}

			pt, err = aead.Open(pt[:0], iv, ct, adata)
			if err == nil {
				cs = append(cs, cavpRow{"PT", hexStr(pt), false})
			} else {
				cs = append(cs, cavpRow{"", "INVALID", false})
			}
		} else {
			keyLen = cs.Int("keylen")
			ivLen = cs.Int("ivlen")
			ptLen = cs.Int("ptlen")
			aadLen = cs.Int("aadlen")
			tagLen = cs.Int("taglen")
		}

		cavp.WriteValues(cs)
	}
}
