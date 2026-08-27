package main

import (
	"strings"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
)

func processBlock(path, filename string) error {
	newCipher := getBlock(filename)
	mode, ok := getBlockMode(filename)
	if newCipher == nil || !ok {
		return unknownAlgorithm(filename)
	}

	output := responsePath(path)
	switch {
	case strings.HasSuffix(filename, "KAT.REQ"):
		return cavp.ProcessBlockKAT(path, output, mode, newCipher)
	case strings.HasSuffix(filename, "MCT.REQ"):
		return cavp.ProcessBlockMCT(path, output, mode, newCipher, 100)
	case strings.HasSuffix(filename, "MMT.REQ"):
		return cavp.ProcessBlockMMT(path, output, mode, newCipher)
	default:
		return unknownAlgorithm(filename)
	}
}
