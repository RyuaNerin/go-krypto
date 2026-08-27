package main

import (
	"strings"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
)

func processKCDSA(path, filename string) error {
	parameterSizes, ok := getKCDSA(filename)
	if !ok {
		return unknownAlgorithm(filename)
	}

	output := responsePath(path)
	switch {
	case strings.HasSuffix(filename, "KPG.REQ"):
		return cavp.ProcessKCDSAKpg(path, output, parameterSizes, 10)
	case strings.HasSuffix(filename, "SGT.REQ"):
		return cavp.ProcessKCDSASgt(path, output, parameterSizes, 1)
	case strings.HasSuffix(filename, "SVT.REQ"):
		return cavp.ProcessKCDSASvt(path, output, parameterSizes)
	default:
		return unknownAlgorithm(filename)
	}
}
