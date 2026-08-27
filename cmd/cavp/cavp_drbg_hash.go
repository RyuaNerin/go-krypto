package main

import "github.com/RyuaNerin/go-krypto/internal/cavp"

func processDRBGHash(path, filename string) error {
	hashInfo := getHash(filename)
	if hashInfo == nil {
		return unknownAlgorithm(filename)
	}
	return cavp.ProcessDrbgHashKAT(path, responsePath(path), hashInfo.New)
}
