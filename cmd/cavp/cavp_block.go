package main

import (
	"log"
	"strings"

	"github.com/RyuaNerin/go-krypto/internal"
)

func processBlock(path, filename string) {
	newBlockCipher := getBlock(filename)
	if newBlockCipher == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	////////////////////////////////////////////////////////////

	newBlockMode := getBlockMode(filename)
	if newBlockMode == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	////////////////////////////////////////////////////////////

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	switch {
	case strings.HasSuffix(filename, "KAT.REQ"): // 기지 답안 검사
		processBlock_KAT_MMT(cavp, newBlockCipher, newBlockMode)

	case strings.HasSuffix(filename, "MCT.REQ"): // Monte
		processBlock_MCT(cavp, newBlockCipher)

	case strings.HasSuffix(filename, "MMT.REQ"): // 다중 블록 메시지 검사
		processBlock_KAT_MMT(cavp, newBlockCipher, newBlockMode)
	}
}

func processBlock_KAT_MMT(cavp *cavpProcessor, fnCipher funcNewBlockCipher, fnNewBlock funcNewBlockMode) {
	dst := make([]byte, 32)

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("KEY") {
			key := cs.Hex("KEY")
			var iv []byte
			if cs.ContainsKey("IV") {
				iv = cs.Hex("IV")
			} else if cs.ContainsKey("CTR") {
				iv = cs.Hex("CTR")
			}

			c, err := fnCipher(key)
			if err != nil {
				panic(err)
			}
			fn := fnNewBlock(c, iv)

			if cs.ContainsKey("PT") {
				src := cs.Hex("PT")

				dst = internal.ResizeBuffer(dst, len(src))
				fn(dst, src, true)

				cs = append(cs, cavpRow{"CT", hexStr(dst), false})
			} else {
				src := cs.Hex("CT")

				dst = internal.ResizeBuffer(dst, len(src))
				fn(dst, src, false)

				cs = append(cs, cavpRow{"PT", hexStr(dst), false})
			}
		}

		cavp.WriteValues(cs)
	}
}

func processBlock_MCT(cavp *cavpProcessor, fnCipher funcNewBlockCipher) {
	var Key, Pt, IV []byte

	for cavp.Next() {
		cs := cavp.ReadValues()
		if cs.ContainsKey("KEY") {
			Key = cs.Hex("KEY")
			Pt = cs.Hex("PT")
			if cs.ContainsKey("IV") {
				IV = cs.Hex("IV")
			} else if cs.ContainsKey("CTR") {
				IV = cs.Hex("CTR")
			}

			break
		}
	}

	c, err := fnCipher(Key)
	if err != nil {
		panic(err)
	}
	blockSize := c.BlockSize()

	const Repeat = 1001

	key := make([][]byte, Repeat)
	pt := make([][]byte, Repeat)
	ct := make([][]byte, Repeat)
	iv := make([][]byte, Repeat)

	for idx := 0; idx < Repeat; idx++ {
		key[idx] = make([]byte, len(Key))
		pt[idx] = make([]byte, len(Pt))
		ct[idx] = make([]byte, len(Pt))
		iv[idx] = make([]byte, len(IV))
	}
	copy(key[0], Key)
	copy(pt[0], Pt)
	copy(iv[0], IV)

	fn := getBlockModeMCT(cavp.filename)
	if fn != nil {
		fn(cavp, fnCipher, key, pt, ct, iv, blockSize)
	} else {
		log.Println("Unknown algorithm: ", cavp.filename)
	}
}
