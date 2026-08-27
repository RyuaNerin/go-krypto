package main

import "github.com/RyuaNerin/go-krypto/internal/cavp"

func processDRBGHMAC(path, filename string) error {
	hashInfo := getHash(filename)
	if hashInfo == nil {
		return unknownAlgorithm(filename)
	}
	return cavp.ProcessDrbgHmacKAT(path, responsePath(path), hashInfo.New)
}
