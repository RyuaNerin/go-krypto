package main

import (
	"crypto/hmac"
	"log"
)

func processHMAC(path, filename string) {
	hashInfo := getHash(filename)
	if hashInfo == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	mac := make([]byte, 64)

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("Key") {
			tlen := cs.Int("tlen")
			key := cs.Hex("Key")
			msg := cs.Hex("Msg")

			hh := hmac.New(hashInfo.New, key)
			if _, err := hh.Write(msg); err != nil {
				panic(err)
			}
			mac = hh.Sum(mac[:0])

			cs = append(cs, cavpRow{"Mac", hexStr(mac[:tlen]), false})
		}

		cavp.WriteValues(cs)
	}
}
