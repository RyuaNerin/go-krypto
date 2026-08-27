package main

import (
	"strings"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
)

func processCMAC(path, filename string) error {
	newCipher := getBlock(filename)
	if newCipher == nil {
		return unknownAlgorithm(filename)
	}

	output := responsePath(path)
	switch {
	case strings.HasSuffix(filename, "GEN.REQ"):
		return cavp.ProcessCMACGen(path, output, newCipher)
	case strings.HasSuffix(filename, "VER.REQ"):
		return cavp.ProcessCMACVer(path, output, newCipher)
	default:
		return unknownAlgorithm(filename)
	}
}
