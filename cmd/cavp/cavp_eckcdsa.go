package main

import (
	"strings"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
)

func processECKCDSA(path, filename string) error {
	curve := getECKCDSACurve(filename)
	if curve == nil {
		return unknownAlgorithm(filename)
	}

	output := responsePath(path)
	switch {
	case strings.HasSuffix(filename, "KPG.REQ"):
		return cavp.ProcessECKCDSAKpg(path, output, curve, 10)
	case strings.HasSuffix(filename, "PKV.REQ"):
		return cavp.ProcessECKCDSAPkv(path, output, curve)
	case strings.HasSuffix(filename, "SGT.REQ"), strings.HasSuffix(filename, "SVT.REQ"):
		hashInfo := getHash(filename)
		if hashInfo == nil {
			return unknownAlgorithm(filename)
		}
		if strings.HasSuffix(filename, "SGT.REQ") {
			return cavp.ProcessECKCDSASgt(path, output, curve, hashInfo.New())
		}
		return cavp.ProcessECKCDSASvt(path, output, curve, hashInfo.New())
	default:
		return unknownAlgorithm(filename)
	}
}
