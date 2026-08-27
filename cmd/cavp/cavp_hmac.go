package main

import "github.com/RyuaNerin/go-krypto/internal/cavp"

func processHMAC(path, filename string) error {
	hashInfo := getHash(filename)
	if hashInfo == nil {
		return unknownAlgorithm(filename)
	}
	return cavp.ProcessHMAC(path, responsePath(path), hashInfo.New)
}
