package main

import (
	"strings"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
)

func processGCM(path, filename string) error {
	newCipher := getBlock(filename)
	if newCipher == nil {
		return unknownAlgorithm(filename)
	}

	output := responsePath(path)
	switch {
	case strings.HasSuffix(filename, "AE.REQ"):
		return cavp.ProcessGCMAe(path, output, newCipher)
	case strings.HasSuffix(filename, "AD.REQ"):
		return cavp.ProcessGCMAd(path, output, newCipher)
	default:
		return unknownAlgorithm(filename)
	}
}
