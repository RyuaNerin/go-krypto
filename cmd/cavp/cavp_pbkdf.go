package main

import (
	"log"

	"github.com/RyuaNerin/go-krypto/pbkdf2"
)

func processPBKDF(path, filename string) {
	hashInfo := getHash(filename)
	if hashInfo == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	var iteration int

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	var dst []byte

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("Iteration") {
			iteration = cs.Int("Iteration")
		} else if cs.ContainsKey("COUNT") {
			Password := cs.Value("Password")
			Salt := cs.Hex("Salt")
			KLen := cs.Int("KLen")

			dst = pbkdf2.Generate(dst[:0], []byte(Password), Salt, iteration, KLen/8, hashInfo.New)

			cs = append(cs, cavpRow{"MK", hexStr(dst), false})
		}

		cavp.WriteValues(cs)
	}
}
