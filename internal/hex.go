package internal

import (
	"encoding/hex"
	"math/big"
	"strings"
)

func H(s string) string {
	var sb strings.Builder
	sb.Grow(len(s))
	for _, c := range s {
		if '0' <= c && c <= '9' {
			sb.WriteRune(c)
		} else if 'a' <= c && c <= 'f' {
			sb.WriteRune(c)
		} else if 'A' <= c && c <= 'F' {
			sb.WriteRune(c)
		}
	}

	return sb.String()
}

// hex to *big.Int
func HI(s string) *big.Int {
	s = H(s)
	if len(s)%2 != 0 {
		panic("len(s) must be a multiple of 2")
	}
	result, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic(s)
	}
	return result
}

// hex to byte
func HB(s string) []byte {
	s = H(s)
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(s)
	}
	return b
}

func Reverse(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
