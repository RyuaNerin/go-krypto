package main

import (
	"strings"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
)

func processHash(path, filename string) error {
	hashInfo := getHash(filename)
	if hashInfo == nil {
		return unknownAlgorithm(filename)
	}

	output := responsePath(path)
	switch {
	case strings.HasSuffix(filename, "MCT.REQ"):
		return cavp.ProcessHashMCT(path, output, hashInfo.New, 100)
	case strings.HasSuffix(filename, "SMT.REQ"):
		return cavp.ProcessHashSMT(path, output, hashInfo.New)
	case strings.HasSuffix(filename, "LMT.REQ"):
		return cavp.ProcessHashLMT(path, output, hashInfo.New)
	default:
		return unknownAlgorithm(filename)
	}
}
