package main

import (
	"log"
	"strings"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/drbg/ctrdrbg"
)

func processDRBGCTR(path, filename string) {
	var err error

	newBlockCipher := getBlock(filename)
	if newBlockCipher == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	keySize := getBlockKeySize(filename)

	usePR := strings.Contains(filename, "(USE PR)")
	useDF := strings.Contains(filename, "(USE DF)")

	dst := make([]byte, 0, 1024)

	var returnedBitsLen int

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("ReturnedBitsLen") {
			returnedBitsLen = cs.Int("ReturnedBitsLen")
			dst = internal.Grow(dst, internal.BitsToBytes(returnedBitsLen))
		} else if cs.ContainsKey("COUNT") {
			EntropyInput := cs.Hex("EntropyInput")
			var Nonce []byte
			if cs.ContainsKey("Nonce") {
				Nonce = cs.Hex("Nonce")
			}
			PersonalizationString := cs.Hex("PersonalizationString")

			drbg := ctrdrbg.Instantiate(
				newBlockCipher,
				keySize,
				0,
				0,
				EntropyInput,
				Nonce,
				PersonalizationString,
				useDF,
				usePR,
			)
			if usePR {
				EntropyInputPR := cs.HexList("EntropyInputPR")
				AdditionalInput := cs.HexList("AdditionalInput")

				err = drbg.Generate(dst, ret(EntropyInputPR[0]), AdditionalInput[0])
				if err != nil {
					panic(err)
				}
				err = drbg.Generate(dst, ret(EntropyInputPR[1]), AdditionalInput[1])
				if err != nil {
					panic(err)
				}
			} else {
				EntropyInputReseed := cs.Hex("EntropyInputReseed")
				AdditionalInputReseed := cs.Hex("AdditionalInputReseed")
				AdditionalInput := cs.HexList("AdditionalInput")

				drbg.Reseed(EntropyInputReseed, AdditionalInputReseed)

				err = drbg.Generate(dst, nil, AdditionalInput[0])
				if err != nil {
					panic(err)
				}
				err = drbg.Generate(dst, nil, AdditionalInput[1])
				if err != nil {
					panic(err)
				}
			}

			cs = append(cs, cavpRow{"ReturnedBits", hexStr(dst), false})
		}

		cavp.WriteValues(cs)
	}
}
