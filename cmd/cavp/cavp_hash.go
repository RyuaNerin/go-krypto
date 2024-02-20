package main

import (
	"log"
	"math"
	"strconv"
	"strings"
)

func processHash(path, filename string) {
	hashInfo := getHash(filename)
	if hashInfo == nil {
		log.Println("Unknown algorithm: ", filename)
		return
	}

	cavp := NewCavp(path, filename)
	defer cavp.Close()

	switch {
	case strings.HasSuffix(filename, "MCT.REQ"): // Monte
		processHash_MCT(cavp, hashInfo)

	case strings.HasSuffix(filename, "SMT.REQ"): // Short Message
		processHash_LT(cavp, hashInfo)

	case strings.HasSuffix(filename, "LMT.REQ"): // Long Message
		processHash_LT(cavp, hashInfo)

	default:
		log.Println("Unknown algorithm: ", filename)
	}
}

func processHash_MCT(cavp *cavpProcessor, hashInfo *HashInfo) {
	var seed []byte

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("Seed") {
			seed = cs.Hex("Seed")
		}

		cavp.WriteValues(cs)
	}

	h := hashInfo.New()

	res := cavpSection{
		{"COUNT", "0", false},
		{"MD", "", false},
	}

	w := float64(hashInfo.w)
	n := float64(len(seed) * 8)
	N := int(math.Floor(32*w/n)) + 1

	log.Println(cavp.filename, w, n, N)

	md := make([][]byte, 1000+N)
	for idx := range md {
		md[idx] = make([]byte, h.Size())
	}

	for i := 0; i < 100; i++ {
		for j := 0; j < N; j++ {
			copy(md[j], seed)
		}
		for k := N; k < 1000+N; k++ {
			h.Reset()
			for kk := k - N; kk <= k-1; kk++ {
				h.Write(md[kk])
			}
			md[k] = h.Sum(md[k][:0])
		}

		copy(seed, md[1000+N-1])
		res[0].Value = strconv.Itoa(i)
		res[1].Value = hexStr(seed)
		cavp.WriteValues(res)
	}
}

func processHash_LT(cavp *cavpProcessor, hashInfo *HashInfo) {
	h := hashInfo.New()

	digest := make([]byte, h.Size())

	for cavp.Next() {
		cs := cavp.ReadValues()

		if cs.ContainsKey("Len") {
			h.Reset()
			h.Write(cs.Hex("Msg")[:cs.Int("Len")/8])
			digest = h.Sum(digest[:0])

			cs = append(cs, cavpRow{"MD", hexStr(digest), false})
		}

		cavp.WriteValues(cs)
	}
}
