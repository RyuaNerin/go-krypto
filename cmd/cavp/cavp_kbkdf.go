package main

import (
	"strings"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
	"github.com/RyuaNerin/go-krypto/kbkdf"
)

func processKBKDF(path, filename string) error {
	mode, ok := getKBKDFMode(filename)
	if !ok {
		return unknownAlgorithm(filename)
	}

	var prf kbkdf.PRF
	switch {
	case strings.Contains(filename, "HMAC"):
		hashInfo := getHash(filename)
		if hashInfo == nil {
			return unknownAlgorithm(filename)
		}
		prf = kbkdf.NewHMACPRF(hashInfo.New)
	case strings.Contains(filename, "CMAC"):
		newCipher := getBlock(filename)
		if newCipher == nil {
			return unknownAlgorithm(filename)
		}
		prf = kbkdf.NewCMACPRF(newCipher)
	default:
		return unknownAlgorithm(filename)
	}

	return cavp.ProcessKbkdfKat(path, responsePath(path), mode, prf)
}

func getKBKDFMode(filename string) (cavp.KBKDFMode, bool) {
	switch {
	case strings.Contains(filename, "(CTR)"), strings.Contains(filename, "-CTR"):
		return cavp.KBKDFModeCounter, true
	case strings.Contains(filename, "(DP)"), strings.Contains(filename, "-DP"):
		return cavp.KBKDFModeDoublePipeline, true
	case strings.Contains(filename, "(FB)"), strings.Contains(filename, "-FB"):
		return cavp.KBKDFModeFeedback, true
	default:
		return 0, false
	}
}
