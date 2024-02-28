package main

import (
	"log"

	"github.com/RyuaNerin/go-krypto/internal"
	"github.com/RyuaNerin/go-krypto/internal/drbg/hashdrbg"
)

func processDRBGHash(path, filename string) {
	var err error

	hashInfo := getHash(filename)
	if hashInfo == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	dst := make([]byte, 0, 1024)

	var returnedBitsLen int
	var usePR bool

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("ReturnedBitsLen") {
			returnedBitsLen = cs.Int("ReturnedBitsLen")
			dst = internal.ResizeBuffer(dst, internal.Bytes(returnedBitsLen))
		}
		if cs.ContainsKey("PredictionResistance") {
			usePR = cs.Bool("PredictionResistance")
		}

		if cs.ContainsKey("COUNT") {
			EntropyInput := cs.Hex("EntropyInput")
			Nonce := cs.Hex("Nonce")
			PersonalizationString := cs.Hex("PersonalizationString")

			drbg := hashdrbg.Instantiate_Hash_DRBG(
				hashInfo.New(),
				returnedBitsLen,
				EntropyInput,
				Nonce,
				PersonalizationString,
				usePR,
			)
			if usePR {
				EntropyInputPR := cs.HexList("EntropyInputPR")
				AdditionalInput := cs.HexList("AdditionalInput")

				err = drbg.Generate_Hash_DRBG(dst, ret(EntropyInputPR[0]), AdditionalInput[0])
				if err != nil {
					panic(err)
				}
				err = drbg.Generate_Hash_DRBG(dst, ret(EntropyInputPR[1]), AdditionalInput[1])
				if err != nil {
					panic(err)
				}
			} else {
				EntropyInputReseed := cs.Hex("EntropyInputReseed")
				AdditionalInputReseed := cs.Hex("AdditionalInputReseed")
				AdditionalInput := cs.HexList("AdditionalInput")

				drbg.Reseed_Hash_DRBG(EntropyInputReseed, AdditionalInputReseed)

				err = drbg.Generate_Hash_DRBG(dst, nil, AdditionalInput[0])
				if err != nil {
					panic(err)
				}
				err = drbg.Generate_Hash_DRBG(dst, nil, AdditionalInput[1])
				if err != nil {
					panic(err)
				}
			}

			cs = append(cs, cavpRow{"ReturnedBits", hexStr(dst), false})
		}

		cavp.WriteValues(cs)
	}
}
