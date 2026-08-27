package main

import (
	"strings"

	"github.com/RyuaNerin/go-krypto/internal/cavp"
)

func processDRBGCTR(path, filename string) error {
	newCipher := getBlock(filename)
	keySize := getBlockKeySize(filename)
	if newCipher == nil || keySize == 0 {
		return unknownAlgorithm(filename)
	}
	useDF := strings.Contains(filename, "(USE DF)")
	return cavp.ProcessDrbgCtrKAT(path, responsePath(path), newCipher, keySize, useDF)
}
