package main

import (
	"strings"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
)

func processCCM(path, filename string) error {
	newCipher := getBlock(filename)
	if newCipher == nil {
		return unknownAlgorithm(filename)
	}

	output := responsePath(path)
	switch {
	case strings.HasSuffix(filename, "GE.REQ"):
		return cavp.ProcessCCMGe(path, output, newCipher)
	case strings.HasSuffix(filename, "DV.REQ"):
		return cavp.ProcessCCMDv(path, output, newCipher)
	default:
		return unknownAlgorithm(filename)
	}
}
