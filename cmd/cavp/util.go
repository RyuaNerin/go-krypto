package main

import (
	"log"
	"math/big"
	"strings"
)

func hexStr(src []byte) string {
	const hextable = "0123456789ABCDEF"

	dst := make([]byte, len(src)*2)

	j := 0
	for _, v := range src {
		dst[j] = hextable[v>>4]
		dst[j+1] = hextable[v&0x0f]
		j += 2
	}
	return string(dst)
}

func hexInt(b *big.Int, bits int) string {
	targetLength := (bits + 3) / 4

	str := hexStr(b.Bytes())
	if len(str) < targetLength {
		str = strings.Repeat("0", targetLength-len(str)) + str
	} else if len(str) > targetLength { // 보통 앞에 자르는거 함.
		log.Println("cut....")
		log.Println("bits", bits)
		log.Println("targetLength", targetLength)
		log.Println("len", len(str))
		log.Println(str)
		str = str[len(str)-targetLength:]
		log.Println(str)
	}
	return str
}

func ret(b []byte) func() ([]byte, error) {
	return func() ([]byte, error) {
		return b, nil
	}
}
